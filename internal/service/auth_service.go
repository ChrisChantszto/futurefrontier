package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/models"
)

type AuthService struct {
	DB             *mongo.Database
	UsersColl      *mongo.Collection
	TokensColl     *mongo.Collection
}

func NewAuthService(db *mongo.Database) *AuthService {
	return &AuthService{
		DB:         db,
		UsersColl:  db.Collection("users"),
		TokensColl: db.Collection("reset_tokens"),
	}
}

func (s *AuthService) CountUsers(ctx context.Context) (int64, error) {
	return s.UsersColl.CountDocuments(ctx, bson.M{})
}

func (s *AuthService) CreateUser(ctx context.Context, email, password string, roles []models.Role) (*models.User, error) {
	now := time.Now().UTC()
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &models.User{
		Email:         email,
		PasswordHash:  string(hash),
		Roles:         roles,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	_, err := s.UsersColl.InsertOne(ctx, u)
	return u, err
}

func (s *AuthService) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := s.UsersColl.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	return &u, err
}

func (s *AuthService) VerifyPassword(u *models.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func (s *AuthService) GenerateResetToken(ctx context.Context, userID string, ttl time.Duration) (plain string, err error) {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	plain = hex.EncodeToString(buf)
	hashb, _ := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)

	doc := &models.ResetToken{
		UserID:    userID,
		Hash:      string(hashb),
		ExpiresAt: time.Now().UTC().Add(ttl),
		Used:      false,
	}
	_, err = s.TokensColl.InsertOne(ctx, doc)
	return
}

func (s *AuthService) ConsumeResetToken(ctx context.Context, plain string) (*models.ResetToken, error) {
	cur, err := s.TokensColl.Find(ctx, bson.M{"used": false, "expiresAt": bson.M{"$gt": time.Now().UTC()}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var t models.ResetToken
		if err := cur.Decode(&t); err != nil {
			continue
		}
		if bcrypt.CompareHashAndPassword([]byte(t.Hash), []byte(plain)) == nil {
			_, _ = s.TokensColl.UpdateByID(ctx, t.ID, bson.M{"$set": bson.M{"used": true}})
			return &t, nil
		}
	}
	return nil, mongo.ErrNoDocuments
}