package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mkaganm/algo-trade/collector/internal/adapters/binance"
	"github.com/mkaganm/algo-trade/collector/internal/adapters/healthcheck"
	"github.com/mkaganm/algo-trade/collector/internal/adapters/mongodb"
	"github.com/mkaganm/algo-trade/collector/internal/config"
	"github.com/mkaganm/algo-trade/collector/internal/core"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 15 * time.Second
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize MongoDB repository
	repo, err := mongodb.NewMongoOrderBookRepository(cfg.MongoURI, cfg.DatabaseName, cfg.CollectionName)
	if err != nil {
		log.Fatalf("Failed to create MongoDB repository: %v", err)
	}
	defer repo.Close()

	// Initialize Binance WebSocket client
	wsClient := binance.NewBinanceWebSocket(cfg.BinanceWSURL, cfg.MaxConnectionRetry, cfg.RetryDelay)

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	})

	// Register health check endpoint
	app.Get("/healthcheck", func(_ *fiber.Ctx) error {
		healthcheck.CheckHandler(repo.Client)

		return nil
	})

	// fixme : add recover go routine and spearete roitine function
	// Start health check endpoint
	go func() {
		log.Println("Starting health check endpoint at :8080")

		if err := app.Listen(":8080"); err != nil {
			log.Printf("Failed to start health check endpoint: %v", err)
		}
	}()

	// Create and run service
	service := core.NewDataCollectorService(wsClient, repo)
	if err := service.Run(ctx); err != nil {
		log.Printf("Service failed: %v", err)

		return
	}

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	log.Println("Application shutdown complete")
}
