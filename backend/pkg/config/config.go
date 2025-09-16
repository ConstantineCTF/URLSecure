package config

import (
	"github.com/spf13/viper" // Configuration library for reading env and config files
)

// Config holds all configuration values loaded from environment or config file
type Config struct {
	AppEnv             string // Application environment (development, production, etc.)
	HTTPPort           string // Port for HTTP server
	DBHost             string // Database host address
	DBPort             string // Database port
	DBUser             string // Database user name
	DBPass             string // Database password
	DBName             string // Database name
	RedisHost          string // Redis host address
	RedisPort          string // Redis port
	JWTSecret          string // JWT secret key
	RateLimitRequests  int    // Number of requests allowed in rate limit window
	RateLimitWindowSec int    // Duration of rate limit window in seconds
}

// Load reads configuration from .env file and environment variables
func Load() (*Config, error) {
	viper.SetConfigFile(".env") // Load config values from .env file
	viper.AutomaticEnv()        // Override with environment variables if set

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Populate Config struct using Viper getters
	return &Config{
		AppEnv:             viper.GetString("APP_ENV"),
		HTTPPort:           viper.GetString("HTTP_PORT"),
		DBHost:             viper.GetString("DB_HOST"),
		DBPort:             viper.GetString("DB_PORT"),
		DBUser:             viper.GetString("DB_USER"),
		DBPass:             viper.GetString("DB_PASS"),
		DBName:             viper.GetString("DB_NAME"),
		RedisHost:          viper.GetString("REDIS_HOST"),
		RedisPort:          viper.GetString("REDIS_PORT"),
		JWTSecret:          viper.GetString("JWT_SECRET"),
		RateLimitRequests:  viper.GetInt("RATE_LIMIT_REQUESTS"),
		RateLimitWindowSec: viper.GetInt("RATE_LIMIT_WINDOW"),
	}, nil
}
