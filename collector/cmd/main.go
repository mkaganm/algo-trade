package main

import (
	"context"
	"log"

	"github.com/mkaganm/algo-trade/collector/internal/adapters/binance"
	"github.com/mkaganm/algo-trade/collector/internal/adapters/mongodb"
	"github.com/mkaganm/algo-trade/collector/internal/config"
	"github.com/mkaganm/algo-trade/collector/internal/core"
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

	log.Println("Application shutdown complete")
}
