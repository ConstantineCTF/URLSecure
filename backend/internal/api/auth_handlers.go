package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"

	authpkg "github.com/ConstantineCTF/URLSecure/backend/pkg/auth"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB, rdb *redis.Client) {
	api := r.Group("/api")
	{
		api.POST("/register", registerHandler(db))
		api.POST("/login", loginHandler(db))
	}
	secured := api.Group("/")
	secured.Use(authpkg.AuthMiddleware())
	{
		secured.POST("/shorten", shortenHandler(db, rdb))
		secured.GET("/stats/:code", statsHandler(db))
		secured.GET("/links", listLinksHandler(db))
	}
}

func registerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash password
		hash, err := authpkg.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}

		// Insert user with username, email, and password_hash
		_, err = db.Exec(
			"INSERT INTO users (username,email,password_hash) VALUES (?,?,?)",
			req.Username, req.Email, hash,
		)
		if err != nil {
			// Check for MySQL duplicate-entry error (1062)
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				msg := "username or email already in use"
				// More specific messaging
				if strings.Contains(mysqlErr.Message, "username") {
					msg = "username already taken"
				} else if strings.Contains(mysqlErr.Message, "email") {
					msg = "email already registered"
				}
				c.JSON(http.StatusConflict, gin.H{"error": msg})
			} else {
				// Log unexpected DB error for debugging
				c.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
			}
			return
		}

		// Retrieve new user ID
		var userID uint64
		row := db.QueryRow("SELECT id FROM users WHERE username = ?", req.Username)
		if err := row.Scan(&userID); err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve user ID"})
			return
		}

		// Create JWT
		token, err := authpkg.CreateJWT(userID)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"token": token})
	}
}

func loginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Identifier string `json:"identifier"`
			Password   string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Lookup by email or username
		var id uint64
		var hash string
		query := "SELECT id,password_hash FROM users WHERE email = ? OR username = ?"
		if err := db.QueryRow(query, req.Identifier, req.Identifier).Scan(&id, &hash); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// Verify password
		if err := authpkg.CheckPassword(hash, req.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// Issue JWT
		token, err := authpkg.CreateJWT(id)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
