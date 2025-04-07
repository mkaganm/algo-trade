package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

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

	// Create and run service
	service := core.NewDataCollectorService(wsClient, repo)
	if err := service.Run(ctx); err != nil {
		log.Printf("Service failed: %v", err)

		return
	}

	// Start health check endpoint
	server := &http.Server{
		Addr:         ":8080",
		Handler:      healthcheck.CheckHandler(repo.Client),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	log.Println("Starting health check endpoint at :8080")

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Failed to start health check endpoint: %v", err)

		return
	}

	log.Println("Application shutdown complete")
}
