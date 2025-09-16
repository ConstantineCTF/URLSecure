package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8" // Redis client library
)

// NewRedisClient creates and returns a new Redis client instance.
func NewRedisClient(host, port string) *redis.Client {
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port), // Redis server address
		PoolSize:     20,                               // Connection pool size for concurrency
		MinIdleConns: 5,                                // Minimum idle connections to keep
		DialTimeout:  5 * time.Second,                  // Dial timeout duration
		ReadTimeout:  3 * time.Second,                  // Read timeout duration
		WriteTimeout: 3 * time.Second,                  // Write timeout duration
	}

	client := redis.NewClient(options)

	// Ping Redis to verify connection on startup, panic on failure to alert immediately
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	return client
}

// Ctx is a globally available context for Redis commands
var Ctx = context.Background()
