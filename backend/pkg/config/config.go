package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv             string
	HTTPPort           string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPass             string
	DBName             string
	RedisHost          string
	RedisPort          string
	JWTSecret          string
	RateLimitRequests  int
	RateLimitWindowSec int
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

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
