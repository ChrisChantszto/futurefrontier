package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// CreateOTPIndexes creates the necessary indexes for the otp_codes collection
func CreateOTPIndexes(ctx context.Context, db *mongo.Database, logger *zap.Logger) error {
	collection := db.Collection("otp_codes")
	
	indexes := []mongo.IndexModel{
		// Compound index for email + purpose + consumedAt + expiresAt (for finding active OTPs)
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
				{Key: "purpose", Value: 1},
				{Key: "consumedAt", Value: 1},
				{Key: "expiresAt", Value: 1},
			},
			Options: options.Index().SetName("email_purpose_consumed_expires"),
		},
		// Index for requestId (for exact lookups)
		{
			Keys: bson.D{{Key: "requestId", Value: 1}},
			Options: options.Index().SetName("requestId").SetUnique(true),
		},
		// TTL index for automatic cleanup of expired documents
		{
			Keys: bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetName("ttl_expires").SetExpireAfterSeconds(3600), // Clean up 1 hour after expiry
		},
		// Index for createdAt (for sorting and cleanup)
		{
			Keys: bson.D{{Key: "createdAt", Value: -1}},
			Options: options.Index().SetName("createdAt_desc"),
		},
		// Index for lastSentAt (for cooldown checks)
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
				{Key: "purpose", Value: 1},
				{Key: "lastSentAt", Value: -1},
			},
			Options: options.Index().SetName("email_purpose_lastSent"),
		},
	}
	
	// Create indexes
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		logger.Error("Failed to create OTP indexes", zap.Error(err))
		return err
	}
	
	logger.Info("OTP collection indexes created successfully")
	return nil
}

// CreateLocaleIndexes creates the necessary indexes for the locales collection
func CreateLocaleIndexes(ctx context.Context, db *mongo.Database, logger *zap.Logger) error {
	collection := db.Collection("locales")
	
	indexes := []mongo.IndexModel{
		// Index for sortOrder (for ordering locales)
		{
			Keys: bson.D{
				{Key: "sortOrder", Value: 1},
				{Key: "code", Value: 1},
			},
			Options: options.Index().SetName("sortOrder_code"),
		},
		// Index for isEnabled (for filtering enabled locales)
		{
			Keys: bson.D{{Key: "isEnabled", Value: 1}},
			Options: options.Index().SetName("isEnabled"),
		},
		// Index for isDefault (for finding default locale)
		{
			Keys: bson.D{{Key: "isDefault", Value: 1}},
			Options: options.Index().SetName("isDefault"),
		},
	}
	
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		logger.Error("Failed to create locale indexes", zap.Error(err))
		return err
	}
	
	logger.Info("Locale collection indexes created successfully")
	return nil
}

// InitializeDatabase sets up all necessary database indexes
func InitializeDatabase(ctx context.Context, db *mongo.Database, logger *zap.Logger) error {
	if err := CreateOTPIndexes(ctx, db, logger); err != nil {
		return err
	}
	
	if err := CreateLocaleIndexes(ctx, db, logger); err != nil {
		return err
	}
	
	logger.Info("Database initialization completed")
	return nil
}
