package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI       = "mongodb://admin:admin@localhost:27017"
	databaseName   = "btc_data"
	collectionName = "depth"
	signalsColName = "trade_signals"
	shortPeriod    = 50
	longPeriod     = 200
	redisStream    = "trade_signals_stream"
)

type OrderBookRecord struct {
	Data      OrderBookData `bson:"data"`
	Timestamp time.Time     `bson:"timestamp"`
}

type OrderBookData struct {
	EventType     string     `json:"e"`
	EventTime     int64      `json:"E"`
	Symbol        string     `json:"s"`
	FirstUpdateID int64      `json:"U"`
	FinalUpdateID int64      `json:"u"`
	BidUpdates    [][]string `json:"b"`
	AskUpdates    [][]string `json:"a"`
}

func calculateSMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	sma := make([]float64, len(prices)-period+1)
	for i := 0; i <= len(prices)-period; i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices[i+j]
		}
		sma[i] = sum / float64(period)
	}
	return sma
}

func processSignals() {
	// MongoDB connection
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(databaseName).Collection(collectionName)
	signalsCollection := client.Database(databaseName).Collection(signalsColName)

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Fetch records
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", -1}}) // Sort by timestamp

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Fatalf("Failed to fetch records: %v", err)
	}
	defer cursor.Close(ctx)

	var rawRecords []bson.M
	if err = cursor.All(ctx, &rawRecords); err != nil {
		log.Fatalf("Failed to decode raw records: %v", err)
	}

	var records []OrderBookRecord
	for _, rawRecord := range rawRecords {
		var record OrderBookRecord
		rawData, ok := rawRecord["data"].(string)
		if !ok {
			log.Printf("Unexpected data format: %+v\n", rawRecord["data"])
			continue
		}

		if err := json.Unmarshal([]byte(rawData), &record.Data); err != nil {
			log.Printf("Failed to unmarshal data: %v", err)
			continue
		}

		record.Timestamp, _ = rawRecord["timestamp"].(time.Time)
		records = append(records, record)
	}

	// Check if the data length is less than 200
	if len(records) < 200 {
		log.Println("Not enough data to process signals. Skipping...")
		return
	}

	// Extract prices
	var prices []float64
	for _, record := range records {
		// Example: Use the first bid price
		if len(record.Data.BidUpdates) > 0 {
			price := record.Data.BidUpdates[0][0]
			var p float64
			fmt.Sscanf(price, "%f", &p)
			prices = append(prices, p)
		}
	}

	// Calculate SMA
	shortSMA := calculateSMA(prices, shortPeriod)
	longSMA := calculateSMA(prices, longPeriod)

	// Generate a single BUY/SELL signal
	if len(shortSMA) > 0 && len(longSMA) > 0 {
		lastShortSMA := shortSMA[len(shortSMA)-1]
		lastLongSMA := longSMA[len(longSMA)-1]

		var signal string
		if lastShortSMA > lastLongSMA {
			signal = "BUY"
		} else if lastShortSMA < lastLongSMA {
			signal = "SELL"
		} else {
			signal = "NEUTRAL"
		}

		// Add signal to Redis Stream with detailed information
		_, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: redisStream,
			Values: map[string]interface{}{
				"signal":   signal,
				"shortSMA": lastShortSMA,
				"longSMA":  lastLongSMA,
				"time":     time.Now().Format(time.RFC3339),
			},
		}).Result()
		if err != nil {
			log.Printf("Failed to add signal to Redis Stream: %v", err)
		} else {
			log.Printf("Added signal to Redis Stream: %s (shortSMA: %.2f, longSMA: %.2f)", signal, lastShortSMA, lastLongSMA)
		}

		// Log detailed signal to MongoDB
		_, err = signalsCollection.InsertOne(ctx, bson.M{
			"signal":    signal,
			"shortSMA":  lastShortSMA,
			"longSMA":   lastLongSMA,
			"timestamp": time.Now(),
		})
		if err != nil {
			log.Printf("Failed to log signal to MongoDB: %v", err)
		} else {
			log.Printf("Logged detailed signal to MongoDB: %s", signal)
		}
	} else {
		fmt.Println("Not enough data to calculate SMA.")
	}
}

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Health check endpoint
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		// MongoDB connection check
		client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "DOWN",
				"message": "Failed to create MongoDB client",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = client.Connect(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "DOWN",
				"message": "Failed to connect to MongoDB",
			})
		}
		defer client.Disconnect(ctx)

		// Redis connection check
		rdb := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		defer rdb.Close()

		_, err = rdb.Ping(ctx).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "DOWN",
				"message": "Failed to connect to Redis",
			})
		}

		// If both checks pass
		return c.JSON(fiber.Map{
			"status":  "UP",
			"message": "Service is running",
		})
	})

	// Schedule signal processing
	c := cron.New()
	_, err := c.AddFunc("*/5 * * * *", processSignals) // Every 5 minutes
	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v", err)
	}

	log.Println("Cron job scheduled to run every 5 minutes.")
	c.Start()

	// Start Fiber app
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("Failed to start Fiber app: %v", err)
		}
	}()

	// Keep the application running
	select {}
}
