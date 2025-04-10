package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mkaganm/algo-trade/trader/internal/adapters/http"
	"github.com/mkaganm/algo-trade/trader/internal/adapters/redisdapter"
	"github.com/mkaganm/algo-trade/trader/internal/app"
	"github.com/mkaganm/algo-trade/trader/internal/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer rdb.Close()

	// Initialize repository and use case
	redisRepo := redisdapter.NewRedisRepository(rdb)
	messageProcessor := app.NewMessageProcessor(redisRepo)

	// Start message processing in a separate goroutine
	go func() {
		for {
			messageProcessor.ProcessMessages(context.Background())
			time.Sleep(2 * time.Second)
		}
	}()

	// Initialize Fiber app
	app := fiber.New()

	// Register health check handler
	healthHandler := http.NewHealthHandler(redisRepo)
	healthHandler.RegisterRoutes(app)

	// Start the server
	log.Printf("Starting server on port %s...", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
