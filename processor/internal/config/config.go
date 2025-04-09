package config

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
	return &Config{
		MongoURI:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:   "btc_data",
		CollectionName: "depth",
		SignalsColName: "trade_signals",
		RedisAddr:      "localhost:6379",
		RedisStream:    "trade_signals_stream",
		ShortPeriod:    50,
		LongPeriod:     200,
		ServerPort:     ":8082",
	}
}
