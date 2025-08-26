package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/config"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/service"
)

// Public: GET /api/settings
func RegisterSettingsPublic(r fiber.Router, db *mongo.Database, cfg config.Config, log *zap.Logger) {
	svc := service.NewSettingsService(db)

	r.Get("/settings", func(c *fiber.Ctx) error {
		set, err := svc.Get(c.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "settings error")
		}
		return c.JSON(set)
	})
}

// Protected: PATCH /api/settings
func RegisterSettingsProtected(r fiber.Router, db *mongo.Database, cfg config.Config, log *zap.Logger) {
	svc := service.NewSettingsService(db)

	r.Patch("/settings", func(c *fiber.Ctx) error {
		var body map[string]any
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		set, err := svc.Update(c.Context(), bson.M(body))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "update error")
		}
		return c.JSON(set)
	})
}