package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mkaganm/algo-trade/trader/internal/adapters/http"
	"github.com/mkaganm/algo-trade/trader/internal/adapters/redisdapter"
	"github.com/mkaganm/algo-trade/trader/internal/app"
	"github.com/mkaganm/algo-trade/trader/internal/config"
	"github.com/robfig/cron/v3"
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

	// Initialize cron job
	c := cron.New()

	_, err := c.AddFunc("@every 5s", func() {
		messageProcessor.ProcessMessages(context.Background())
	})
	if err != nil {
		log.Printf("Failed to add cron job: %v", err)
	}

	c.Start()

	defer c.Stop()

	// Initialize Fiber app
	app := fiber.New()

	// Register health check handler
	healthHandler := http.NewHealthHandler(redisRepo)
	healthHandler.RegisterRoutes(app)

	// Start the server
	log.Printf("Starting server on port %s...", cfg.AppPort)
	log.Println(app.Listen(":" + cfg.AppPort))
}
