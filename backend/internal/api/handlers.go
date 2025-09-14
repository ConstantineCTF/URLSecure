package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/ConstantineCTF/URLSecure/backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func NewRouter(cfg *config.Config, db *sql.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Serve static assets under /assets
	r.Static("/assets", "./public")

	// SPA entry point
	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	// Health check
	r.GET("/api/health", healthHandler)

	// API routes
	api := r.Group("/api")
	{
		api.POST("/shorten", shortenHandler(db, rdb))
		api.GET("/stats/:code", statsHandler(db))
	}

	// Redirect handler
	r.GET("/r/:code", redirectHandler(db, rdb))

	return r
}

func shortenHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			URL string `json:"url" binding:"required,url"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		code := generateCode(6)
		if _, err := db.Exec(
			"INSERT INTO urls (code,target) VALUES (?,?)",
			code, req.URL,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		rdb.Set(context.Background(), "url:"+code, req.URL, 24*time.Hour)
		c.JSON(http.StatusCreated, gin.H{"code": code})
	}
}

func statsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		var clicks int
		var created time.Time
		if err := db.QueryRow(
			"SELECT clicks,created_at FROM urls WHERE code=?", code,
		).Scan(&clicks, &created); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": code, "clicks": clicks, "createdAt": created})
	}
}

func redirectHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		ctx := context.Background()
		target, err := rdb.Get(ctx, "url:"+code).Result()
		if err == redis.Nil {
			row := db.QueryRow("SELECT target FROM urls WHERE code=?", code)
			if err := row.Scan(&target); err != nil {
				c.String(http.StatusNotFound, "Not found")
				return
			}
			rdb.Set(ctx, "url:"+code, target, 24*time.Hour)
		} else if err != nil {
			c.String(http.StatusInternalServerError, "Internal error")
			return
		}
		go db.Exec("UPDATE urls SET clicks=clicks+1 WHERE code=?", code)
		c.Redirect(http.StatusFound, target)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
