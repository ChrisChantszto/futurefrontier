package http

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/ChrisChantszto/futurefrontier/internal/config"
	"github.com/ChrisChantszto/futurefrontier/internal/service"
	"github.com/ChrisChantszto/futurefrontier/internal/transport/http/handlers"
	"github.com/ChrisChantszto/futurefrontier/internal/transport/http/middleware"
)

func SetupRoutes(app *fiber.App, db *mongo.Database, cfg config.Config, log *zap.Logger) {
    // Serve local test page from same origin to avoid CORS when opening the file directly
    app.Get("/test-ai", func(c *fiber.Ctx) error {
        return c.SendFile("./test-ai.html")
    })

	api := app.Group("/api")

	// Public routes
	handlers.RegisterAuth(api, db, cfg, log)
	handlers.RegisterSettingsPublic(api, db, cfg, log)
	handlers.RegisterLocalesPublic(api, db, cfg, log)
	// Public pages (read-only)
	handlers.RegisterPagesPublic(api, db, cfg, log)

	api.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	// Protected routes
	protected := api.Group("", middleware.RequireAuth(cfg))
	handlers.RegisterSettingsProtected(protected, db, cfg, log)
	handlers.RegisterLocalesProtected(protected, db, cfg, log)
	handlers.RegisterPagesProtected(protected, db, cfg, log)

	// Elasticsearch service for logs and demo data
	esService, err := service.NewElasticsearchService(cfg.Logging, log)
	if err != nil {
		log.Error("Failed to initialize Elasticsearch service", zap.Error(err))
	} else {
		// Demo data generation endpoint
		demoGen := service.NewDemoDataGenerator(esService, log)
		demoHandler := handlers.NewDemoHandler(demoGen, log)
		protected.Post("/demo/generate", demoHandler.GenerateDemoData)
		log.Info("Demo data endpoint registered")

		// Logs API endpoints for frontend
		logsHandler := handlers.NewLogsHandler(esService, log)
		protected.Get("/logs", logsHandler.GetLogs)
		protected.Get("/logs/stats", logsHandler.GetStats)
		log.Info("Logs API endpoints registered")
	}

	// AI routes (protected) - only initialize if Google Cloud is configured
	if cfg.GoogleCloud.ProjectID != "" {
		RegisterAIRoutes(protected, cfg, log)
	} else {
		log.Warn("Google Cloud not configured - AI routes disabled")
	}
}

// RegisterAIRoutes sets up AI-powered endpoints
func RegisterAIRoutes(router fiber.Router, cfg config.Config, log *zap.Logger) {
	// Initialize Vertex AI service
	vertexAI, err := service.NewVertexAIService(cfg.GoogleCloud, log)
	if err != nil {
		log.Error("Failed to initialize Vertex AI service", zap.Error(err))
		return
	}

	// Initialize Elasticsearch service for log retrieval
	esService, err := service.NewElasticsearchService(cfg.Logging, log)
	if err != nil {
		log.Error("Failed to initialize Elasticsearch service for AI", zap.Error(err))
		return
	}

	// Create AI handler
	aiHandler := handlers.NewAIHandler(vertexAI, esService, log)

	// AI endpoints
	ai := router.Group("/ai")
	ai.Post("/generate", aiHandler.GenerateContent)
	ai.Post("/analyze-logs", aiHandler.AnalyzeLogs)
	ai.Post("/detect-anomalies", aiHandler.DetectAnomalies)
	ai.Post("/suggest-optimizations", aiHandler.SuggestOptimizations)
	ai.Post("/chat", aiHandler.ChatWithLogs)
	
	// Enhanced AI endpoints for pattern detection
	ai.Post("/detect-patterns", aiHandler.DetectPatterns)
	ai.Post("/analyze-failures", aiHandler.AnalyzeFailures)
	ai.Post("/root-cause-analysis", aiHandler.RootCauseAnalysis)
	ai.Post("/analyze-traffic-failures", aiHandler.AnalyzeTrafficFailures)

	log.Info("AI routes registered successfully")
}