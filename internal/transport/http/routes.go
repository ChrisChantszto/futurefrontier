package http

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/config"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/transport/http/handlers"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/transport/http/middleware"
)

func SetupRoutes(app *fiber.App, db *mongo.Database, cfg config.Config, log *zap.Logger) {
	api := app.Group("/api")

	// Public routes
	handlers.RegisterAuth(api, db, cfg, log)
	handlers.RegisterSettingsPublic(api, db, cfg, log)
	// Public pages (read-only)
	handlers.RegisterPagesPublic(api, db, cfg, log)

	api.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	// Protected routes
	protected := api.Group("", middleware.RequireAuth(cfg))
	handlers.RegisterSettingsProtected(protected, db, cfg, log)
	handlers.RegisterPagesProtected(protected, db, cfg, log)
}