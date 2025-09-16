package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ConstantineCTF/URLSecure/backend/internal/middleware"  // Custom middleware (RateLimit, Auth)
	"github.com/ConstantineCTF/URLSecure/backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// NewRouter constructs the Gin engine and sets up routes and middleware
func NewRouter(cfg *config.Config, db *sql.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	// Trust only localhost (loopback) for proxy IPs, enhancing security
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	// Serve static assets from ./public/assets
	r.Static("/assets", "./public")

	// Serve SPA index.html on root
	r.GET("/", func(c *gin.Context) { c.File("./public/index.html") })

	// Health check endpoint (public)
	r.GET("/api/health", healthHandler)

	// Public authentication endpoints (register + login)
	public := r.Group("/api")
	{
		public.POST("/register", registerHandler(db))
		public.POST("/login", loginHandler(db))
	}

	// Protected endpoints - require rate limit and JWT auth middleware
	protected := r.Group("/api")
	protected.Use(
		middleware.RateLimitMiddleware(),
		middleware.AuthMiddleware(),
	)
	{
		protected.POST("/shorten", shortenHandler(db, rdb))        // Create short URL
		protected.GET("/stats/:code", statsHandler(db))            // Get stats for code
		protected.GET("/links", listLinksHandler(db))              // List all user links
	}

	// Redirect endpoint for short URLs (public)
	r.GET("/r/:code", redirectHandler(db, rdb))

	return r
}

// shortenHandler stores a new URL in DB and caches it in Redis asynchronously
func shortenHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			URL string `json:"url" binding:"required,url"` // URL must be valid
		}

		// Validate JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Retrieve authenticated user ID from context
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, ok := userIDVal.(uint64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
			return
		}

		// Generate random 6-character short code for URL
		code := generateCode(6)

		// Insert link record into DB synchronously before responding
		if _, err := db.Exec(
			"INSERT INTO links (user_id, code, target) VALUES (?, ?, ?)",
			userID, code, req.URL,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		// Cache short URL target asynchronously; doesn't block response
		go func() {
			ctx := context.Background()
			rdb.Set(ctx, "url:"+code, req.URL, 24*time.Hour)
		}()

		// Return code of new shortened URL
		c.JSON(http.StatusCreated, gin.H{"code": code})
	}
}

// statsHandler returns statistics (click count, creation date) for a short code
func statsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")

		var clicks int
		var created time.Time

		// Query DB for click count and creation date of the short URL
		if err := db.QueryRow(
			"SELECT click_count, created_at FROM links WHERE code = ?", code,
		).Scan(&clicks, &created); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// Return stats as JSON
		c.JSON(http.StatusOK, gin.H{"code": code, "clicks": clicks, "createdAt": created})
	}
}

// redirectHandler resolves short URL from cache or DB, increments click, redirects user
func redirectHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		ctx := context.Background()

		log.Printf("Redirect handler for code: %s", code)

		// Try Redis cache first
		target, err := rdb.Get(ctx, "url:"+code).Result()
		if err == redis.Nil {
			log.Println("Cache miss—query DB")

			// Cache miss, query DB for target URL
			if err := db.QueryRow(
				"SELECT target FROM links WHERE code = ?", code,
			).Scan(&target); err != nil {
				log.Printf("DB lookup failed for code %s: %v", code, err)
				c.String(http.StatusNotFound, "Not found")
				return
			}

			// Cache result asynchronously
			rdb.Set(ctx, "url:"+code, target, 24*time.Hour)
		} else if err != nil {
			// Redis failure
			log.Printf("Redis error for code %s: %v", code, err)
			c.String(http.StatusInternalServerError, "Internal error")
			return
		}

		// Increment the click count asynchronously in DB, no need to await
		go db.Exec(
			"UPDATE links SET clicks = clicks + 1 WHERE code = ?", code,
		)

		log.Printf("Redirecting code %s to target: %s", code, target)
		// Redirect client to target URL
		c.Redirect(http.StatusFound, target)
	}
}

// healthHandler returns basic health check JSON
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
