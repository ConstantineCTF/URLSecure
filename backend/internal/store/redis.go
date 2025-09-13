package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(host, port string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
}

var Ctx = context.Background()
