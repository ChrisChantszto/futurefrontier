package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ChrisChantszto/futurefrontier/internal/service"
	"go.uber.org/zap"
)

type DemoHandler struct {
	demoGen *service.DemoDataGenerator
	logger  *zap.Logger
}

func NewDemoHandler(demoGen *service.DemoDataGenerator, logger *zap.Logger) *DemoHandler {
	return &DemoHandler{
		demoGen: demoGen,
		logger:  logger,
	}
}

// GenerateDemoData generates demo log data for testing
func (h *DemoHandler) GenerateDemoData(c *fiber.Ctx) error {
	countStr := c.Query("count", "1000")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 || count > 10000 {
		return JSONError(c, fiber.StatusBadRequest, "Invalid count parameter (must be between 1 and 10000)", 0, nil)
	}

	h.logger.Info("Starting demo data generation", zap.Int("count", count))

	ctx := context.Background()
	if err := h.demoGen.GenerateRealisticLogs(ctx, count); err != nil {
		h.logger.Error("Failed to generate demo data", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to generate demo data", 0, err.Error())
	}

	return JSONSuccess(c, "Demo data generated successfully", fiber.Map{
		"logs_generated": count,
		"message":        "Demo logs have been written to Elasticsearch",
	})
}
