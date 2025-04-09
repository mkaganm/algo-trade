package config

import (
	"os"
	"strconv"
)

type Config struct {
	MongoURI       string
	DatabaseName   string
	CollectionName string
	SignalsColName string
	RedisAddr      string
	RedisStream    string
	ShortPeriod    int
	LongPeriod     int
	ServerPort     string
}

func Load() *Config {
	shortPeriod, err := strconv.Atoi(getEnv("SHORT_PERIOD", "50"))
	if err != nil {
		shortPeriod = 50
	}

	longPeriod, err := strconv.Atoi(getEnv("LONG_PERIOD", "200"))
	if err != nil {
		longPeriod = 200
	}

	return &Config{
		MongoURI:       getEnv("MONGO_URI", "mongodb://admin:admin@localhost:27017"),
		DatabaseName:   getEnv("DATABASE_NAME", "btc_data"),
		CollectionName: getEnv("COLLECTION_NAME", "depth"),
		SignalsColName: getEnv("SIGNALS_COL_NAME", "trade_signals"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisStream:    getEnv("REDIS_STREAM", "trade_signals_stream"),
		ShortPeriod:    shortPeriod,
		LongPeriod:     longPeriod,
		ServerPort:     getEnv("SERVER_PORT", ":8082"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
