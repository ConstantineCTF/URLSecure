package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-redis/redis/v8"            // Redis client import for side effects
	"github.com/go-sql-driver/mysql"             // MySQL driver specific errors
	authpkg "github.com/ConstantineCTF/URLSecure/backend/pkg/auth"  // Auth utilities (hashing, JWT)
)

/* RegisterRoutes wires auth and protected endpoints.
func RegisterRoutes(r *gin.Engine, db *sql.DB, rdb *redis.Client) {
	api := r.Group("/api")
	{
		api.POST("/register", registerHandler(db))               // Public registration endpoint
		api.POST("/login", loginHandler(db))                     // Public login endpoint
	}
	secured := api.Group("/")
	secured.Use(authpkg.AuthMiddleware())                       // JWT middleware to protect endpoints
	{
		secured.POST("/shorten", shortenHandler(db, rdb))       // Shorten URL endpoint
		secured.GET("/stats/:code", statsHandler(db))           // Stats endpoint by code
		secured.GET("/links", listLinksHandler(db))              // List user links endpoint
	}
}
*/

func registerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		// Bind incoming JSON to req, return 400 if invalid
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash password for secure storage
		hash, err := authpkg.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}

		// Insert new user into MySQL users table
		_, err = db.Exec(
			"INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
			req.Username, req.Email, hash,
		)
		if err != nil {
			// Handle duplicate username/email error
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				msg := "username or email already in use"
				if strings.Contains(mysqlErr.Message, "username") {
					msg = "username already taken"
				} else if strings.Contains(mysqlErr.Message, "email") {
					msg = "email already registered"
				}
				c.JSON(http.StatusConflict, gin.H{"error": msg})
			} else {
				c.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
			}
			return
		}

		// Query new user's ID for token creation
		var userID uint64
		db.QueryRow("SELECT id FROM users WHERE username = ?", req.Username).Scan(&userID)

		// Issue JWT token for newly registered user
		token, err := authpkg.CreateJWT(userID)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}

		// Return JWT token with 201 Created
		c.JSON(http.StatusCreated, gin.H{"token": token})
	}
}

func loginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Identifier string `json:"identifier"` // Can be email or username
			Password   string `json:"password"`
		}
		// Bind login request JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var id uint64
		var hash string

		// Query user by email or username to get ID and hashed password
		query := "SELECT id, password_hash FROM users WHERE email = ? OR username = ?"
		if err := db.QueryRow(query, req.Identifier, req.Identifier).Scan(&id, &hash); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// Check password correctness against stored hash
		if err := authpkg.CheckPassword(hash, req.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// Create JWT token after successful auth
		token, err := authpkg.CreateJWT(id)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}

		// Return JWT token on success
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
