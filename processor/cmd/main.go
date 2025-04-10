package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mkaganm/algo-trade/processor/internal/application"
	"github.com/mkaganm/algo-trade/processor/internal/config"
	"github.com/mkaganm/algo-trade/processor/internal/helpers"
	"github.com/mkaganm/algo-trade/processor/internal/infrastructure/api"
	"github.com/mkaganm/algo-trade/processor/internal/infrastructure/persistence"
	"github.com/mkaganm/algo-trade/processor/internal/infrastructure/scheduler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	signalProcessingTimeout  = 30 * time.Second
	mongoDBConnectionTimeout = 10 * time.Second
)

//nolint:funlen
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB
	mongoRepo, err := persistence.NewMongoOrderBookRepository(
		cfg.MongoURI,
		cfg.DatabaseName,
		cfg.CollectionName,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB repository: %v", err)
	}

	// Initialize Redis
	redisPublisher := persistence.NewRedisSignalPublisher(cfg.RedisAddr, cfg.RedisStream)

	// Initialize application services
	signalProcessor := application.NewSignalProcessor(
		mongoRepo,
		mongoRepo, // Assuming MongoOrderBookRepository also implements SignalRepository
		redisPublisher,
	)

	// Initialize scheduler
	cronScheduler := scheduler.NewCronScheduler()

	_, err = cronScheduler.Schedule("*/5 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), signalProcessingTimeout)
		defer cancel()

		signal, err := signalProcessor.GenerateSignal(ctx, cfg.ShortPeriod, cfg.LongPeriod)
		if err != nil {
			log.Printf("Failed to generate signal: %v", err)

			return
		}

		log.Printf("Generated signal: %s (shortSMA: %.2f, longSMA: %.2f)", signal.Signal, signal.ShortSMA, signal.LongSMA)
	})
	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v", err)
	}

	cronScheduler.Start()

	defer cronScheduler.Stop()

	// Initialize HTTP server
	app := fiber.New()

	// Create MongoDB client for health check
	ctx, cancel := context.WithTimeout(context.Background(), mongoDBConnectionTimeout)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect MongoDB client: %v", err)
		}
	}()

	// Create Redis client for health check
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	// Setup health check handler
	healthHandler := api.NewHealthHandler(mongoClient, redisClient)
	app.Get("/healthcheck", healthHandler.Check)

	// Start server
	go startServer(app, cfg.ServerPort)

	// Keep the application running
	select {}
}

func startServer(app *fiber.App, serverPort string) {
	defer helpers.RecoverRoutine(make(chan error))

	if err := app.Listen(serverPort); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
