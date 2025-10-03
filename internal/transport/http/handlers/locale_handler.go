package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/config"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/models"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/service"
)

// Public: GET /api/locales - Get all locales
// Public: GET /api/locales/config - Get frontend config
// Public: GET /api/locales/:id - Get single locale
func RegisterLocalesPublic(r fiber.Router, db *mongo.Database, cfg config.Config, log *zap.Logger) {
	svc := service.NewLocaleService(db)

	// Initialize default locales on first call
	r.Use(func(c *fiber.Ctx) error {
		if err := svc.InitializeDefaultLocales(c.Context()); err != nil {
			log.Warn("Failed to initialize default locales", zap.Error(err))
		}
		return c.Next()
	})

	// GET /api/locales - List all locales
	r.Get("/locales", func(c *fiber.Ctx) error {
		enabledOnly := c.Query("enabled") == "true"
		
		locales, err := svc.GetAll(c.Context(), enabledOnly)
		if err != nil {
			log.Error("Failed to get locales", zap.Error(err))
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve locales")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    locales,
		})
	})

	// GET /api/locales/config - Get frontend routing config
	r.Get("/locales/config", func(c *fiber.Ctx) error {
		config, err := svc.GetConfig(c.Context())
		if err != nil {
			log.Error("Failed to get locale config", zap.Error(err))
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve locale config")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    config,
		})
	})

	// GET /api/locales/:id - Get single locale
	r.Get("/locales/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		
		locale, err := svc.GetByID(c.Context(), id)
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "Locale not found")
		}
		if err != nil {
			log.Error("Failed to get locale", zap.Error(err), zap.String("id", id))
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve locale")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    locale,
		})
	})
}

// Protected: POST /api/locales - Create new locale
// Protected: PATCH /api/locales/:id - Update locale
// Protected: DELETE /api/locales/:id - Delete locale
// Protected: POST /api/locales/:id/set-default - Set as default
func RegisterLocalesProtected(r fiber.Router, db *mongo.Database, cfg config.Config, log *zap.Logger) {
	svc := service.NewLocaleService(db)

	// POST /api/locales - Create new locale
	r.Post("/locales", func(c *fiber.Ctx) error {
		var req models.CreateLocaleRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		// Basic validation
		if req.Code == "" || req.Name == "" || req.NativeName == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Code, name, and nativeName are required")
		}

		locale, err := svc.Create(c.Context(), req)
		if err != nil {
			log.Error("Failed to create locale", zap.Error(err), zap.String("code", req.Code))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		log.Info("Locale created", zap.String("code", locale.Code), zap.String("name", locale.Name))

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    locale,
			"message": "Locale created successfully",
		})
	})

	// PATCH /api/locales/:id - Update locale
	r.Patch("/locales/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		
		var req models.UpdateLocaleRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		locale, err := svc.Update(c.Context(), id, req)
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "Locale not found")
		}
		if err != nil {
			log.Error("Failed to update locale", zap.Error(err), zap.String("id", id))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		log.Info("Locale updated", zap.String("id", id))

		return c.JSON(fiber.Map{
			"success": true,
			"data":    locale,
			"message": "Locale updated successfully",
		})
	})

	// DELETE /api/locales/:id - Delete locale
	r.Delete("/locales/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := svc.Delete(c.Context(), id); err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "Locale not found")
			}
			log.Error("Failed to delete locale", zap.Error(err), zap.String("id", id))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		log.Info("Locale deleted", zap.String("id", id))

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Locale deleted successfully",
		})
	})

	// POST /api/locales/:id/set-default - Set locale as default
	r.Post("/locales/:id/set-default", func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := svc.SetDefault(c.Context(), id); err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "Locale not found")
			}
			log.Error("Failed to set default locale", zap.Error(err), zap.String("id", id))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		log.Info("Default locale set", zap.String("id", id))

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Default locale set successfully",
		})
	})
}
