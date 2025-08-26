package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/config"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/transport/http"
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
	var db *mongo.Database
	if mongoClient == nil {
		zlogger.Sugar().Fatal("mongo client is nil; aborting startup")
	}
	db = mongoClient.Database(cfg.DBName)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName: "onetake-corpsite-backend",
	})

	// Middlewares
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zlogger,
	}))

	// Routes
	http.SetupRoutes(app, db, cfg, zlogger)

	// Start
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	zlogger.Sugar().Infof("listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		zlogger.Sugar().Fatal(err)
	}
}