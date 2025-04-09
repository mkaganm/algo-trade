package api

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthHandler struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
}

func NewHealthHandler(mongoClient *mongo.Client, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		mongoClient: mongoClient,
		redisClient: redisClient,
	}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// MongoDB health check
	err := h.mongoClient.Ping(ctx, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "DOWN",
			"message": "MongoDB connection failed",
		})
	}

	// Redis health check
	_, err = h.redisClient.Ping(ctx).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "DOWN",
			"message": "Redis connection failed",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "UP",
		"message": "Service is healthy",
	})
}
