package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/models"
)

type SettingsService struct {
	Coll *mongo.Collection
}

func NewSettingsService(db *mongo.Database) *SettingsService {
	return &SettingsService{Coll: db.Collection("settings")}
}

func (s *SettingsService) Get(ctx context.Context) (*models.Settings, error) {
	var set models.Settings
	err := s.Coll.FindOne(ctx, bson.M{"_id": "global"}).Decode(&set)
	if err == mongo.ErrNoDocuments {
		// default settings
		set = models.Settings{
			ID:           "global",
			Theme:        "light",
			Languages:    []string{"en", "zh-Hant", "zh-Hans"},
			PageLimit:    100,
			UserLimit:    100,
			LanguageLimit: 3,
			SMTP: models.SMTPSettings{
				AllowCustom: false,
				Selected:    "company",
			},
		}
		_, _ = s.Coll.InsertOne(ctx, set)
		return &set, nil
	}
	return &set, err
}

func (s *SettingsService) Update(ctx context.Context, update bson.M) (*models.Settings, error) {
	_, err := s.Coll.UpdateByID(ctx, "global", bson.M{"$set": update})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx)
}