package handlers

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/ChrisChantszto/futurefrontier/internal/config"
	"github.com/ChrisChantszto/futurefrontier/internal/models"
	"github.com/ChrisChantszto/futurefrontier/internal/service"
)

var idRE = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Public
func RegisterPagesPublic(r fiber.Router, db *mongo.Database, cfg config.Config, log *zap.Logger) {
	svc := service.NewPagesService(db)

	r.Get("/pages", func(c *fiber.Ctx) error {
		limit := int64(cfg.PageLimit)
		if q := c.QueryInt("limit", 0); q > 0 {
			limit = int64(q)
		}
		pages, err := svc.List(c.Context(), limit)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "list error")
		}
		return c.JSON(pages)
	})

	r.Get("/pages/:identifier", func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")
		if !idRE.MatchString(identifier) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid identifier")
		}
		p, err := svc.GetByIdentifier(c.Context(), identifier)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, "get error")
		}
		return c.JSON(p)
	})

	r.Get("/pages/:identifier/export", func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")
		if !idRE.MatchString(identifier) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid identifier")
		}
		p, err := svc.GetByIdentifier(c.Context(), identifier)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, "get error")
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		c.Set(fiber.HeaderContentDisposition, "attachment; filename=\""+identifier+".json\"")
		return c.JSON(p)
	})
}

// Protected (write operations)
func RegisterPagesProtected(r fiber.Router, db *mongo.Database, cfg config.Config, log *zap.Logger) {
	svc := service.NewPagesService(db)

	r.Post("/pages", func(c *fiber.Ctx) error {
		var p models.Page
		if err := c.BodyParser(&p); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if !idRE.MatchString(p.Identifier) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid identifier")
		}
		// Optional: check exists to avoid duplicates
		if existing, err := svc.GetByIdentifier(c.Context(), p.Identifier); err == nil && existing != nil {
			return fiber.NewError(fiber.StatusConflict, "identifier exists")
		}
		created, err := svc.Create(c.Context(), &p)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "create error")
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	})

	r.Patch("/pages/:identifier", func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")
		if !idRE.MatchString(identifier) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid identifier")
		}
		var body map[string]any
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		p, err := svc.UpdateByIdentifier(c.Context(), identifier, bson.M(body))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, "update error")
		}
		return c.JSON(p)
	})

	r.Delete("/pages/:identifier", func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")
		if !idRE.MatchString(identifier) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid identifier")
		}
		if err := svc.DeleteByIdentifier(c.Context(), identifier); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "delete error")
		}
		return c.JSON(fiber.Map{"ok": true})
	})

	r.Post("/pages/import", func(c *fiber.Ctx) error {
		var payload struct{ Pages []models.Page `json:"pages"` }
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		for _, p := range payload.Pages {
			if !idRE.MatchString(p.Identifier) {
				return fiber.NewError(fiber.StatusBadRequest, "invalid identifier in import")
			}
		}
		count, err := svc.Import(c.Context(), payload.Pages)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "import error")
		}
		return c.JSON(fiber.Map{"ok": true, "count": count})
	})

	r.Post("/pages/:identifier/generate-template", func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")
		if !idRE.MatchString(identifier) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid identifier")
		}
		tpl := svc.GenerateTemplate(c.Context(), identifier)
		return c.JSON(tpl)
	})
}
