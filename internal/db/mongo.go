package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func Connect(ctx context.Context, uri, dbName string) (*Mongo, error) {
	clientOpts := options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(10 * time.Second).
		SetConnectTimeout(10 * time.Second)
	
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &Mongo{Client: client, DB: client.Database(dbName)}, nil
}

func WithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}