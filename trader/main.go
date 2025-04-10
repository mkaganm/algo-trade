package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

const (
	redisAddr       = "localhost:6379"       // Redis address
	inputStream     = "trade_signals_stream" // Input stream name
	outputStream    = "processed_signals"    // Output stream name
	consumerGroup   = "signal_group"         // Consumer group name
	consumerName    = "signal_consumer"      // Consumer name
	startingID      = "0"                    // Starting ID
	pollingInterval = 2 * time.Second        // Polling interval
)

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Create consumer group (ignore error if it already exists)
	err := rdb.XGroupCreateMkStream(ctx, inputStream, consumerGroup, startingID).Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Failed to create consumer group: %v", err)
	}

	fmt.Println("Connected to Redis and consumer group created.")

	// Initialize Fiber app
	app := fiber.New()

	// Health check endpoint
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		// Ping Redis to check connection
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "unhealthy",
				"message": "Redis connection failed",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "healthy",
			"message": "Redis connection is active",
		})
	})

	// Start the Fiber app
	go func() {
		if err := app.Listen(":8083"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Main processing loop
	for {
		// Read messages from the stream
		streams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: consumerName,
			Streams:  []string{inputStream, ">"},
			Count:    10,
			Block:    pollingInterval,
		}).Result()

		if err != nil && !errors.Is(err, redis.Nil) {
			log.Printf("Error reading messages: %v", err)

			continue
		}

		// Process messages
		for _, stream := range streams {
			for _, message := range stream.Messages {
				fmt.Printf("Message received: %v\n", message.Values)

				fmt.Println(reflect.TypeOf(message.Values))
				fmt.Println(message.Values["signal"])
				fmt.Println(message.Values["shortSMA"])
				fmt.Println(message.Values["longSMA"])
				fmt.Println(message.Values["time"])

				processedData, _ := json.Marshal(message.Values)

				// Add detailed information to the processed stream
				_, err := rdb.XAdd(ctx, &redis.XAddArgs{
					Stream: outputStream,
					Values: map[string]interface{}{
						"original_id":   message.ID,
						"original_data": processedData,
						"processed":     true,
						"processed_at":  time.Now().Format(time.RFC3339),
					},
				}).Result()
				if err != nil {
					log.Printf("Error writing processed message: %v", err)
				}

				// Acknowledge the message
				err = rdb.XAck(ctx, inputStream, consumerGroup, message.ID).Err()
				if err != nil {
					log.Printf("Error acknowledging message: %v", err)
				} else {
					fmt.Printf("Message acknowledged: %s\n", message.ID)
				}
			}
		}
	}
}
