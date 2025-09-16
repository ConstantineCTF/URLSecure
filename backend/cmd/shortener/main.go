package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"   // For graceful shutdown on OS signals
	"syscall"     // For signal constants like SIGINT, SIGTERM
	"time"
	
	"github.com/ConstantineCTF/URLSecure/backend/internal/api"    // HTTP router and handlers
	"github.com/ConstantineCTF/URLSecure/backend/internal/store"  // Database and redis clients
	"github.com/ConstantineCTF/URLSecure/backend/pkg/config"     // Config loading from env
	"github.com/gin-gonic/gin"                                  // HTTP web framework
	"github.com/joho/godotenv"                                  // Load .env file for env vars
)

func main() {
	// Load environment variables from .env file if present; safe to ignore errors if missing
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it:", err)
	}

	// Set Gin to release mode for production level logging/performance
	gin.SetMode(gin.ReleaseMode)

	// Load all configuration from env vars into cfg struct
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to MySQL with config credentials
	db, err := store.ConnectMySQL(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}
	// Set MySQL connection pool parameters for efficiency
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)
	defer db.Close() // Close DB connection on program exit

	// Initialize Redis client
	redisClient := store.NewRedisClient(cfg.RedisHost, cfg.RedisPort)
	defer redisClient.Close() // Close Redis client on exit

	// Create HTTP router with all routes and middleware
	router := api.NewRouter(cfg, db, redisClient)

	// Setup HTTP server with address from config and our router
	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	// Run server asynchronously for graceful shutdown handling
	go func() {
		log.Printf("starting server on port %s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Fatal if server crashes unexpectedly
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Channel to listen for interrupt signals (CTRL+C or kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Notify channel on SIGINT or SIGTERM

	<-quit // Block here until quit signal detected

	log.Println("shutting down server gracefully...")
	// Context with timeout of 5 seconds for server shutdown to cleanup ongoing requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down server cleanly or force exit on error
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	// Confirm server shutdown and exit
	log.Println("server exiting properly")
}
