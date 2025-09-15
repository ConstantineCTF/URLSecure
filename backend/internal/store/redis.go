package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(host, port string) *redis.Client {
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		PoolSize:     20, // Increase pool size for concurrency
		MinIdleConns: 5,  // Maintain some idle connections
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	client := redis.NewClient(options)

	// Ping Redis to verify connection on startup
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	return client
}

var Ctx = context.Background()
