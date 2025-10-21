package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/ChrisChantszto/futurefrontier/internal/models"
)

const (
	APILogsCollection   = "api_logs_backup"
	ErrorLogsCollection = "error_logs_backup"
)

// MongoBackupService handles MongoDB backup operations for failed logs
type MongoBackupService struct {
	database *mongo.Database
	logger   *zap.Logger
}

// NewMongoBackupService creates a new MongoDB backup service
func NewMongoBackupService(database *mongo.Database, logger *zap.Logger) (*MongoBackupService, error) {
	service := &MongoBackupService{
		database: database,
		logger:   logger,
	}

	// Create indexes for better performance
	if err := service.createIndexes(); err != nil {
		logger.Warn("Failed to create MongoDB backup indexes", zap.Error(err))
	}

	return service, nil
}

// createIndexes creates necessary indexes for the backup collections
func (mbs *MongoBackupService) createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// API logs indexes
	apiCollection := mbs.database.Collection(APILogsCollection)
	apiIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "es_written", Value: 1},
				{Key: "created_at", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "request_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "retry_count", Value: 1},
				{Key: "last_retry_at", Value: 1},
			},
		},
		// TTL index to automatically delete old backup logs after 30 days
		{
			Keys: bson.D{
				{Key: "created_at", Value: 1},
			},
			Options: options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60), // 30 days
		},
	}

	if _, err := apiCollection.Indexes().CreateMany(ctx, apiIndexes); err != nil {
		return err
	}

	// Error logs indexes
	errorCollection := mbs.database.Collection(ErrorLogsCollection)
	errorIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "es_written", Value: 1},
				{Key: "created_at", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "request_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "retry_count", Value: 1},
				{Key: "last_retry_at", Value: 1},
			},
		},
		// TTL index to automatically delete old backup logs after 30 days
		{
			Keys: bson.D{
				{Key: "created_at", Value: 1},
			},
			Options: options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60), // 30 days
		},
	}

	if _, err := errorCollection.Indexes().CreateMany(ctx, errorIndexes); err != nil {
		return err
	}

	mbs.logger.Info("Created MongoDB backup indexes")
	return nil
}

// SaveAPILog saves an API log entry to MongoDB backup
func (mbs *MongoBackupService) SaveAPILog(ctx context.Context, logEntry *models.APILogEntry) error {
	collection := mbs.database.Collection(APILogsCollection)
	
	_, err := collection.InsertOne(ctx, logEntry)
	if err != nil {
		return err
	}

	return nil
}

// SaveErrorLog saves an error log entry to MongoDB backup
func (mbs *MongoBackupService) SaveErrorLog(ctx context.Context, logEntry *models.ErrorLogEntry) error {
	collection := mbs.database.Collection(ErrorLogsCollection)
	
	_, err := collection.InsertOne(ctx, logEntry)
	if err != nil {
		return err
	}

	return nil
}

// GetFailedAPILogs retrieves API logs that failed to write to Elasticsearch
func (mbs *MongoBackupService) GetFailedAPILogs(ctx context.Context, limit int) ([]models.APILogEntry, error) {
	collection := mbs.database.Collection(APILogsCollection)
	
	filter := bson.M{
		"es_written": false,
		"$or": []bson.M{
			{"retry_count": bson.M{"$lt": 3}}, // Less than 3 retries
			{"retry_count": bson.M{"$exists": false}}, // No retry count set
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: 1}}) // Oldest first

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.APILogEntry
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetFailedErrorLogs retrieves error logs that failed to write to Elasticsearch
func (mbs *MongoBackupService) GetFailedErrorLogs(ctx context.Context, limit int) ([]models.ErrorLogEntry, error) {
	collection := mbs.database.Collection(ErrorLogsCollection)
	
	filter := bson.M{
		"es_written": false,
		"$or": []bson.M{
			{"retry_count": bson.M{"$lt": 3}}, // Less than 3 retries
			{"retry_count": bson.M{"$exists": false}}, // No retry count set
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: 1}}) // Oldest first

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.ErrorLogEntry
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// DeleteAPILog removes an API log entry from MongoDB backup
func (mbs *MongoBackupService) DeleteAPILog(ctx context.Context, id primitive.ObjectID) error {
	collection := mbs.database.Collection(APILogsCollection)
	
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteErrorLog removes an error log entry from MongoDB backup
func (mbs *MongoBackupService) DeleteErrorLog(ctx context.Context, id primitive.ObjectID) error {
	collection := mbs.database.Collection(ErrorLogsCollection)
	
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// UpdateAPILogRetry increments the retry count for an API log entry
func (mbs *MongoBackupService) UpdateAPILogRetry(ctx context.Context, id primitive.ObjectID) error {
	collection := mbs.database.Collection(APILogsCollection)
	
	now := time.Now()
	update := bson.M{
		"$inc": bson.M{"retry_count": 1},
		"$set": bson.M{"last_retry_at": now},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateErrorLogRetry increments the retry count for an error log entry
func (mbs *MongoBackupService) UpdateErrorLogRetry(ctx context.Context, id primitive.ObjectID) error {
	collection := mbs.database.Collection(ErrorLogsCollection)
	
	now := time.Now()
	update := bson.M{
		"$inc": bson.M{"retry_count": 1},
		"$set": bson.M{"last_retry_at": now},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// GetBackupStats returns statistics about backup logs
func (mbs *MongoBackupService) GetBackupStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// API logs stats
	apiCollection := mbs.database.Collection(APILogsCollection)
	
	totalAPI, err := apiCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats["total_api_logs"] = totalAPI

	failedAPI, err := apiCollection.CountDocuments(ctx, bson.M{"es_written": false})
	if err != nil {
		return nil, err
	}
	stats["failed_api_logs"] = failedAPI

	// Error logs stats
	errorCollection := mbs.database.Collection(ErrorLogsCollection)
	
	totalError, err := errorCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats["total_error_logs"] = totalError

	failedError, err := errorCollection.CountDocuments(ctx, bson.M{"es_written": false})
	if err != nil {
		return nil, err
	}
	stats["failed_error_logs"] = failedError

	return stats, nil
}

// CleanupOldLogs manually removes logs older than the specified duration
func (mbs *MongoBackupService) CleanupOldLogs(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	filter := bson.M{"created_at": bson.M{"$lt": cutoff}}

	// Clean API logs
	apiCollection := mbs.database.Collection(APILogsCollection)
	apiResult, err := apiCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	// Clean error logs
	errorCollection := mbs.database.Collection(ErrorLogsCollection)
	errorResult, err := errorCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	mbs.logger.Info("Cleaned up old backup logs",
		zap.Int64("api_logs_deleted", apiResult.DeletedCount),
		zap.Int64("error_logs_deleted", errorResult.DeletedCount),
		zap.Duration("older_than", olderThan))

	return nil
}
