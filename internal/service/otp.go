package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/config"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/models"
)

type OTPService struct {
	db          *mongo.Database
	config      config.Config
	smtpService *SMTPService
}

func NewOTPService(db *mongo.Database, cfg config.Config, smtpService *SMTPService) *OTPService {
	return &OTPService{
		db:          db,
		config:      cfg,
		smtpService: smtpService,
	}
}

// generateOTPCode generates a 6-digit OTP code using crypto/rand
func (s *OTPService) generateOTPCode() (string, error) {
	bytes := make([]byte, 3) // 3 bytes = 24 bits, enough for 6 digits
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	// Convert to 6-digit number
	num := int(bytes[0])<<16 | int(bytes[1])<<8 | int(bytes[2])
	code := fmt.Sprintf("%06d", num%1000000)
	return code, nil
}

// computeHMAC computes HMAC-SHA256 of the code
func (s *OTPService) computeHMAC(code string) (string, error) {
	if s.config.AppHMACSecret == "" {
		return "", fmt.Errorf("APP_HMAC_SECRET not configured")
	}
	
	// Decode the base64 secret
	secret, err := base64.StdEncoding.DecodeString(s.config.AppHMACSecret)
	if err != nil {
		// Try hex decoding if base64 fails
		secret, err = hex.DecodeString(s.config.AppHMACSecret)
		if err != nil {
			// Use as plain string if both fail
			secret = []byte(s.config.AppHMACSecret)
		}
	}
	
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(code))
	return hex.EncodeToString(h.Sum(nil)), nil
}

// verifyHMAC verifies the HMAC using constant-time comparison
func (s *OTPService) verifyHMAC(code, expectedHMAC string) (bool, error) {
	computedHMAC, err := s.computeHMAC(code)
	if err != nil {
		return false, err
	}
	
	return subtle.ConstantTimeCompare([]byte(computedHMAC), []byte(expectedHMAC)) == 1, nil
}

// normalizeEmail normalizes email address (trim and lowercase)
func (s *OTPService) normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// RequestOTP handles OTP request logic
func (s *OTPService) RequestOTP(ctx context.Context, email string, purpose models.OTPPurpose) (string, error) {
	email = s.normalizeEmail(email)
	
	if !purpose.IsValid() {
		return "", fmt.Errorf("invalid OTP purpose")
	}
	
	// Check resend cooldown
	if err := s.checkResendCooldown(ctx, email, purpose); err != nil {
		return "", err
	}
	
	// Generate request ID and OTP code
	requestID := uuid.New().String()
	code, err := s.generateOTPCode()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP code: %w", err)
	}
	
	// Compute HMAC
	codeHMAC, err := s.computeHMAC(code)
	if err != nil {
		return "", fmt.Errorf("failed to compute HMAC: %w", err)
	}
	
	// Create OTP document
	now := time.Now().UTC()
	otpDoc := &models.OTPCode{
		Email:       email,
		RequestID:   requestID,
		Purpose:     purpose,
		CodeHMAC:    codeHMAC,
		ExpiresAt:   now.Add(time.Duration(s.config.OTPTTLMinutes) * time.Minute),
		Attempts:    0,
		MaxAttempts: s.config.OTPMaxAttempts,
		CreatedAt:   now,
		ConsumedAt:  nil,
		LastSentAt:  now,
	}
	
	// Insert into database
	collection := s.db.Collection("otp_codes")
	_, err = collection.InsertOne(ctx, otpDoc)
	if err != nil {
		return "", fmt.Errorf("failed to save OTP: %w", err)
	}
	
	// Optionally mark previous active OTPs as consumed
	s.consumePreviousOTPs(ctx, email, purpose, requestID)
	
	// Send email
	if err := s.smtpService.SendOTPEmail(email, code); err != nil {
		return "", fmt.Errorf("failed to send OTP email: %w", err)
	}
	
	return requestID, nil
}

// checkResendCooldown checks if the user is within the resend cooldown period
func (s *OTPService) checkResendCooldown(ctx context.Context, email string, purpose models.OTPPurpose) error {
	collection := s.db.Collection("otp_codes")
	
	// Find latest active OTP
	filter := bson.M{
		"email":      email,
		"purpose":    purpose,
		"consumedAt": nil,
		"expiresAt":  bson.M{"$gt": time.Now().UTC()},
	}
	
	opts := options.FindOne().SetSort(bson.M{"createdAt": -1})
	var lastOTP models.OTPCode
	err := collection.FindOne(ctx, filter, opts).Decode(&lastOTP)
	
	if err == mongo.ErrNoDocuments {
		return nil // No active OTP, cooldown doesn't apply
	}
	if err != nil {
		return fmt.Errorf("failed to check existing OTP: %w", err)
	}
	
	// Check cooldown
	cooldownEnd := lastOTP.LastSentAt.Add(time.Duration(s.config.OTPResendCooldownSecs) * time.Second)
	if time.Now().UTC().Before(cooldownEnd) {
		return fmt.Errorf("please wait before requesting another code")
	}
	
	return nil
}

// consumePreviousOTPs marks previous active OTPs as consumed
func (s *OTPService) consumePreviousOTPs(ctx context.Context, email string, purpose models.OTPPurpose, excludeRequestID string) {
	collection := s.db.Collection("otp_codes")
	
	filter := bson.M{
		"email":     email,
		"purpose":   purpose,
		"requestId": bson.M{"$ne": excludeRequestID},
		"consumedAt": nil,
	}
	
	update := bson.M{
		"$set": bson.M{
			"consumedAt": time.Now().UTC(),
		},
	}
	
	collection.UpdateMany(ctx, filter, update)
}

// VerifyOTP verifies an OTP code
func (s *OTPService) VerifyOTP(ctx context.Context, email, code, requestID string, purpose models.OTPPurpose) (*models.OTPCode, error) {
	email = s.normalizeEmail(email)
	
	if !purpose.IsValid() {
		return nil, fmt.Errorf("invalid OTP purpose")
	}
	
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid OTP code format")
	}
	
	// Find OTP document
	otpDoc, err := s.findOTPForVerification(ctx, email, purpose, requestID)
	if err != nil {
		return nil, err
	}
	
	// Verify HMAC
	valid, err := s.verifyHMAC(code, otpDoc.CodeHMAC)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OTP: %w", err)
	}
	
	if !valid {
		// Increment attempts
		if err := s.incrementOTPAttempts(ctx, otpDoc.ID); err != nil {
			return nil, fmt.Errorf("failed to update OTP attempts: %w", err)
		}
		return nil, fmt.Errorf("invalid OTP code")
	}
	
	// Mark as consumed
	if err := s.consumeOTP(ctx, otpDoc.ID); err != nil {
		return nil, fmt.Errorf("failed to consume OTP: %w", err)
	}
	
	return otpDoc, nil
}

// findOTPForVerification finds the appropriate OTP for verification
func (s *OTPService) findOTPForVerification(ctx context.Context, email string, purpose models.OTPPurpose, requestID string) (*models.OTPCode, error) {
	collection := s.db.Collection("otp_codes")
	
	filter := bson.M{
		"email":      email,
		"purpose":    purpose,
		"consumedAt": nil,
		"expiresAt":  bson.M{"$gt": time.Now().UTC()},
		"attempts":   bson.M{"$lt": s.config.OTPMaxAttempts},
	}
	
	// If requestID is provided, use it for exact match
	if requestID != "" {
		filter["requestId"] = requestID
	}
	
	opts := options.FindOne().SetSort(bson.M{"createdAt": -1})
	var otpDoc models.OTPCode
	err := collection.FindOne(ctx, filter, opts).Decode(&otpDoc)
	
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no valid OTP found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find OTP: %w", err)
	}
	
	return &otpDoc, nil
}

// incrementOTPAttempts increments the attempt count for an OTP
func (s *OTPService) incrementOTPAttempts(ctx context.Context, otpID string) error {
	collection := s.db.Collection("otp_codes")
	
	filter := bson.M{"_id": otpID}
	update := bson.M{"$inc": bson.M{"attempts": 1}}
	
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// consumeOTP marks an OTP as consumed
func (s *OTPService) consumeOTP(ctx context.Context, otpID string) error {
	collection := s.db.Collection("otp_codes")
	
	filter := bson.M{"_id": otpID}
	update := bson.M{"$set": bson.M{"consumedAt": time.Now().UTC()}}
	
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// ResendOTP handles OTP resend by creating a new OTP and consuming the old one
func (s *OTPService) ResendOTP(ctx context.Context, email string, purpose models.OTPPurpose) (string, error) {
	// This is essentially the same as RequestOTP, but we explicitly consume previous OTPs first
	return s.RequestOTP(ctx, email, purpose)
}
