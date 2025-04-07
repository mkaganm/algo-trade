package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	BinanceWSURL       string
	MaxConnectionRetry int
	RetryDelay         time.Duration
	MongoURI           string
	DatabaseName       string
	CollectionName     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	maxConnectionRetry, err := strconv.Atoi(os.Getenv("MAX_CONNECTION_RETRY"))
	if err != nil {
		return nil, fmt.Errorf("failed to convert MAX_CONNECTION_RETRY to int: %w", err)
	}

	retryDelay, err := time.ParseDuration(os.Getenv("RETRY_DELAY"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse RETRY_DELAY: %w", err)
	}

	return &Config{
		BinanceWSURL:       os.Getenv("BINANCE_WS_URL"),
		MaxConnectionRetry: maxConnectionRetry,
		RetryDelay:         retryDelay,
		MongoURI:           os.Getenv("MONGO_URI"),
		DatabaseName:       os.Getenv("DATABASE_NAME"),
		CollectionName:     os.Getenv("COLLECTION_NAME"),
	}, nil
}
