package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ConstantineCTF/URLSecure/backend/internal/middleware"
	"github.com/ConstantineCTF/URLSecure/backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// NewRouter constructs the Gin engine with all routes and middleware.
func NewRouter(cfg *config.Config, db *sql.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	// Trust only loopback (localhost) or whatever proxy IPs you use
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}
	// Static assets & SPA
	r.Static("/assets", "./public")
	r.GET("/", func(c *gin.Context) { c.File("./public/index.html") })

	// Public health
	r.GET("/api/health", healthHandler)

	// Public auth endpoints
	public := r.Group("/api")
	{
		public.POST("/register", registerHandler(db))
		public.POST("/login", loginHandler(db))
	}

	// Protected endpoints: rate limiting + JWT auth
	protected := r.Group("/api")
	protected.Use(
		middleware.RateLimitMiddleware(),
		middleware.AuthMiddleware(),
	)
	{
		protected.POST("/shorten", shortenHandler(db, rdb))
		protected.GET("/stats/:code", statsHandler(db))
		protected.GET("/links", listLinksHandler(db))
	}

	// Redirect (public)
	r.GET("/r/:code", redirectHandler(db, rdb))

	return r
}

// shortenHandler stores a new URL and caches it.
func shortenHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			URL string `json:"url" binding:"required,url"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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

		code := generateCode(6)

		// Insert into DB synchronously (must finish before responding)
		if _, err := db.Exec(
			"INSERT INTO links (user_id, code, target) VALUES (?, ?, ?)",
			userID, code, req.URL,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		// Cache asynchronously (won't delay response)
		go func() {
			ctx := context.Background()
			rdb.Set(ctx, "url:"+code, req.URL, 24*time.Hour)
		}()

		c.JSON(http.StatusCreated, gin.H{"code": code})
	}
}

// statsHandler returns click stats.
func statsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		var clicks int
		var created time.Time
		if err := db.QueryRow(
			"SELECT click_count, created_at FROM links WHERE code = ?", code,
		).Scan(&clicks, &created); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": code, "clicks": clicks, "createdAt": created})
	}
}

// redirectHandler retrieves URL from cache or DB, tracks click.
func redirectHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		ctx := context.Background()

		log.Printf("Redirect handler for code: %s", code)
		target, err := rdb.Get(ctx, "url:"+code).Result()
		if err == redis.Nil {
			log.Println("Cache miss—query DB")
			if err := db.QueryRow(
				"SELECT target FROM links WHERE code = ?", code,
			).Scan(&target); err != nil {
				log.Printf("DB lookup failed for code %s: %v", code, err)
				c.String(http.StatusNotFound, "Not found")
				return
			}
			rdb.Set(ctx, "url:"+code, target, 24*time.Hour)
		} else if err != nil {
			log.Printf("Redis error for code %s: %v", code, err)
			c.String(http.StatusInternalServerError, "Internal error")
			return
		}

		go db.Exec(
			"UPDATE links SET clicks = clicks + 1 WHERE code = ?", code,
		)

		log.Printf("Redirecting code %s to target: %s", code, target)
		c.Redirect(http.StatusFound, target)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
