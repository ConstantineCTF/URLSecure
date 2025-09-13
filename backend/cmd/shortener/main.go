package main

import (
	"log"

	"github.com/ConstantineCTF/URLSecure/backend/internal/api"
	"github.com/ConstantineCTF/URLSecure/backend/internal/store"
	"github.com/ConstantineCTF/URLSecure/backend/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to MySQL
	dsn := store.NewMySQLDSN(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := store.ConnectMySQL(dsn)
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redisClient := store.NewRedisClient(cfg.RedisHost, cfg.RedisPort)
	defer redisClient.Close()

	// Create and start router, passing db and redisClient
	router := api.NewRouter(cfg, db, redisClient)
	log.Printf("starting server on port %s", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
