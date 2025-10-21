package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/ChrisChantszto/futurefrontier/internal/config"
	"github.com/ChrisChantszto/futurefrontier/internal/db"
	"github.com/ChrisChantszto/futurefrontier/internal/middleware"
	"github.com/ChrisChantszto/futurefrontier/internal/models"
	"github.com/ChrisChantszto/futurefrontier/internal/service"
	"github.com/ChrisChantszto/futurefrontier/internal/transport/http"
)

func connectMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	// Logger
	zlogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer zlogger.Sync()

	// DB (optional for hello world)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := connectMongo(ctx, cfg.MongoURI)
	if err != nil {
		zlogger.Sugar().Fatalf("mongo connect failed: %v", err)
	}
	var database *mongo.Database
	if mongoClient == nil {
		zlogger.Sugar().Fatal("mongo client is nil; aborting startup")
	}
	database = mongoClient.Database(cfg.DBName)

	// Initialize database indexes
	if err := db.InitializeDatabase(ctx, database, zlogger); err != nil {
		zlogger.Sugar().Warnf("failed to initialize database indexes: %v", err)
	}

	// Initialize logger service
	loggerService, err := service.NewLoggerService(cfg.Logging, database, zlogger)
	if err != nil {
		zlogger.Sugar().Fatalf("failed to initialize logger service: %v", err)
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName: "onetake-corpsite-backend",
		ErrorHandler: middleware.ErrorHandlerMiddleware(loggerService, models.LogConfig{
			ProjectID:         cfg.Logging.ProjectID,
			EnableRequestBody: cfg.Logging.EnableRequestBody,
			MaxBodySize:       cfg.Logging.MaxBodySize,
		}),
	})

	// Middlewares
	app.Use(middleware.RequestIDMiddleware())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization, X-Locale, Accept-Language, X-Request-ID",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie, X-Request-ID",
		MaxAge:           3600,
	}))

	// Custom logging middleware
	app.Use(middleware.LoggerMiddleware(loggerService, models.LogConfig{
		ProjectID:          cfg.Logging.ProjectID,
		EnableRequestBody:  cfg.Logging.EnableRequestBody,
		EnableResponseBody: cfg.Logging.EnableResponseBody,
		EnableHeaders:      cfg.Logging.EnableHeaders,
		MaxBodySize:        cfg.Logging.MaxBodySize,
		SensitiveHeaders:   []string{"authorization", "cookie", "x-api-key"},
		SensitivePaths:     []string{}, // Log all paths
	}, zlogger))

	// Keep fiberzap for development visibility
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zlogger,
	}))

	// Add logging stats endpoint - public access
	app.Get("/health/logging", middleware.LoggingStatsHandler(loggerService))

	// Routes
	http.SetupRoutes(app, database, cfg, zlogger)

	// Start server in goroutine
	go func() {
		port := cfg.Port
		if port == "" {
			port = "8080"
		}
		zlogger.Sugar().Infof("listening on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			zlogger.Sugar().Error("Server failed to start:", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	zlogger.Info("Shutting down server...")

	// Graceful shutdown
	if err := app.Shutdown(); err != nil {
		zlogger.Sugar().Error("Server shutdown error:", err)
	}

	// Shutdown logger service
	loggerService.Shutdown()
	zlogger.Info("Server shutdown complete")
}
