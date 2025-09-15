package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ConstantineCTF/URLSecure/backend/internal/api"
	"github.com/ConstantineCTF/URLSecure/backend/internal/store"
	"github.com/ConstantineCTF/URLSecure/backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // already here
)

func main() {
	// This will run before auth.init in most cases, but it's safe to keep
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it:", err)
	}

	gin.SetMode(gin.ReleaseMode)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := store.ConnectMySQL(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)
	defer db.Close()

	redisClient := store.NewRedisClient(cfg.RedisHost, cfg.RedisPort)
	defer redisClient.Close()

	router := api.NewRouter(cfg, db, redisClient)
	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("starting server on port %s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exiting properly")
}
