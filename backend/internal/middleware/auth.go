package middleware

import (
	"net/http"
	"strings"

	authpkg "github.com/ConstantineCTF/URLSecure/backend/pkg/auth" // Auth utilities (JWT parsing)
	"github.com/gin-gonic/gin"
)

// AuthMiddleware enforces JWT authentication on protected routes
// It validates the Authorization header and extracts the user ID from the token,
// then sets userID in Gin's context for handlers to use.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Check if the Authorization header is present and properly formatted as Bearer token
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}

		// Extract the JWT token string after "Bearer "
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the JWT token, obtain claims containing user ID
		claims, err := authpkg.ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Set the extracted user ID in the request context for downstream handlers
		c.Set("userID", claims.UserID)

		// Continue processing request
		c.Next()
	}
}
