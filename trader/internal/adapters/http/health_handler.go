package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mkaganm/algo-trade/trader/internal/ports"
)

type HealthHandler struct {
	healthService ports.HealthService
}

func NewHealthHandler(healthService ports.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

func (h *HealthHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/healthcheck", h.HealthCheck)
}

func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	err := h.healthService.CheckHealth(context.Background())
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
}
