package main

import (
	"github.com/ConstantineCTF/URLSecure/backend/internal/api"
	"github.com/ConstantineCTF/URLSecure/backend/pkg/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	router := api.NewRouter(cfg)
	log.Printf("starting server on port %s", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
