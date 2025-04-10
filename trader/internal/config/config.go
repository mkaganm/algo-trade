package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr string
	AppPort   string
}

func LoadConfig() *Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve configuration values
	return &Config{
		RedisAddr: os.Getenv("REDIS_ADDR"),
		AppPort:   os.Getenv("APP_PORT"),
	}
}
