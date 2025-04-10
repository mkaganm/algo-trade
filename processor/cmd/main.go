package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mkaganm/algo-trade/processor/internal/application"
	"github.com/mkaganm/algo-trade/processor/internal/config"
	"github.com/mkaganm/algo-trade/processor/internal/infrastructure/api"
	"github.com/mkaganm/algo-trade/processor/internal/infrastructure/persistence"
	"github.com/mkaganm/algo-trade/processor/internal/infrastructure/scheduler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	cronTimeout    = 30 * time.Second
	contextTimeout = 10 * time.Second
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize dependencies
	_, signalProcessor, cronScheduler := initializeDependencies(cfg)

	// Schedule tasks
	scheduleTasks(cronScheduler, signalProcessor, cfg)

	// Start HTTP server
	startHTTPServer(cfg)

	// Keep the application running
	select {}
}

func initializeDependencies(
	cfg *config.Config,
) (
	*persistence.MongoOrderBookRepository,
	*application.SignalProcessor,
	*scheduler.CronScheduler,
) {
	// Initialize MongoDB
	mongoRepo, err := persistence.NewMongoOrderBookRepository(cfg.MongoURI, cfg.DatabaseName, cfg.CollectionName)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB repository: %v", err)
	}

	// Initialize Redis
	redisPublisher := persistence.NewRedisSignalPublisher(cfg.RedisAddr, cfg.RedisStream)

	// Initialize application services
	signalProcessor := application.NewSignalProcessor(mongoRepo, mongoRepo, redisPublisher)

	// Initialize scheduler
	cronScheduler := scheduler.NewCronScheduler()

	return mongoRepo, signalProcessor, cronScheduler
}

func scheduleTasks(
	cronScheduler *scheduler.CronScheduler,
	signalProcessor *application.SignalProcessor,
	cfg *config.Config,
) {
	_, err := cronScheduler.Schedule("*/1 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), cronTimeout)
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
}

func startHTTPServer(cfg *config.Config) {
	app := fiber.New()

	// Create MongoDB client for health check
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Printf("Failed to create and connect MongoDB client: %v", err)

		return
	}

	defer mongoClient.Disconnect(ctx) //nolint

	// Create Redis client for health check
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	// Setup health check handler
	healthHandler := api.NewHealthHandler(mongoClient, redisClient)
	app.Get("/healthcheck", healthHandler.Check)

	// Start server
	go func() {
		if err := app.Listen(cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}
