package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	// Map to hold rate limiters per IP address
	visitors = make(map[string]*rate.Limiter)
	// Mutex to synchronize access to the visitors map
	mu sync.Mutex
)

// getVisitor returns the rate limiter for the given IP, creating one if not existing
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// Create a new rate limiter allowing 20 requests per second with bursts allowed
		limiter = rate.NewLimiter(rate.Every(time.Second), 20)
		visitors[ip] = limiter
	}
	return limiter
}

// RateLimitMiddleware applies rate limiting per IP address to incoming requests.
// It aborts with HTTP 429 if the IP exceeds the allowed rate.
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Reject requests with invalid IP addresses
		if ip == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid IP"})
			return
		}

		limiter := getVisitor(ip)

		// If request exceeds limiter allowance, respond with rate limit error
		if !limiter.Allow() {
			c.JSON(429, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		// Continue processing request if within rate limit
		c.Next()
	}
}
