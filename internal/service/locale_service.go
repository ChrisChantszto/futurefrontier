package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/models"
)

type LocaleService struct {
	Coll *mongo.Collection
}

func NewLocaleService(db *mongo.Database) *LocaleService {
	return &LocaleService{Coll: db.Collection("locales")}
}

// InitializeDefaultLocales creates default locales if none exist
func (s *LocaleService) InitializeDefaultLocales(ctx context.Context) error {
	count, err := s.Coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	
	if count > 0 {
		return nil // Already initialized
	}

	defaultLocales := []models.Locale{
		{
			ID:         "en",
			Code:       "en",
			Name:       "English",
			NativeName: "English",
			IsDefault:  true,
			IsEnabled:  true,
			SortOrder:  1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "zh-hans",
			Code:       "zh-hans",
			Name:       "Simplified Chinese",
			NativeName: "简体中文",
			IsDefault:  false,
			IsEnabled:  true,
			SortOrder:  2,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "zh-hant",
			Code:       "zh-hant",
			Name:       "Traditional Chinese",
			NativeName: "繁體中文",
			IsDefault:  false,
			IsEnabled:  true,
			SortOrder:  3,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	var docs []interface{}
	for _, loc := range defaultLocales {
		docs = append(docs, loc)
	}

	_, err = s.Coll.InsertMany(ctx, docs)
	return err
}

// GetAll returns all locales, optionally filtered by enabled status
func (s *LocaleService) GetAll(ctx context.Context, enabledOnly bool) ([]models.Locale, error) {
	filter := bson.M{}
	if enabledOnly {
		filter["isEnabled"] = true
	}

	opts := options.Find().SetSort(bson.D{{Key: "sortOrder", Value: 1}, {Key: "code", Value: 1}})
	cursor, err := s.Coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var locales []models.Locale
	if err := cursor.All(ctx, &locales); err != nil {
		return nil, err
	}

	return locales, nil
}

// GetByID returns a single locale by ID/code
func (s *LocaleService) GetByID(ctx context.Context, id string) (*models.Locale, error) {
	var locale models.Locale
	err := s.Coll.FindOne(ctx, bson.M{"_id": id}).Decode(&locale)
	if err != nil {
		return nil, err
	}
	return &locale, nil
}

// GetConfig returns the frontend-compatible locale configuration
func (s *LocaleService) GetConfig(ctx context.Context) (*models.LocaleConfig, error) {
	locales, err := s.GetAll(ctx, true) // Only enabled
	if err != nil {
		return nil, err
	}

	var codes []string
	var defaultCode string

	for _, loc := range locales {
		codes = append(codes, loc.Code)
		if loc.IsDefault {
			defaultCode = loc.Code
		}
	}

	// Fallback to first locale if no default is set
	if defaultCode == "" && len(codes) > 0 {
		defaultCode = codes[0]
	}

	return &models.LocaleConfig{
		Locales:       codes,
		DefaultLocale: defaultCode,
	}, nil
}

// Create adds a new locale
func (s *LocaleService) Create(ctx context.Context, req models.CreateLocaleRequest) (*models.Locale, error) {
	// Normalize code to lowercase
	code := strings.ToLower(strings.TrimSpace(req.Code))
	if code == "" {
		return nil, fmt.Errorf("locale code cannot be empty")
	}

	// Check if locale already exists
	existing, _ := s.GetByID(ctx, code)
	if existing != nil {
		return nil, fmt.Errorf("locale '%s' already exists", code)
	}

	locale := models.Locale{
		ID:         code,
		Code:       code,
		Name:       strings.TrimSpace(req.Name),
		NativeName: strings.TrimSpace(req.NativeName),
		IsDefault:  false, // New locales are never default by default
		IsEnabled:  req.IsEnabled,
		SortOrder:  req.SortOrder,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := s.Coll.InsertOne(ctx, locale)
	if err != nil {
		return nil, err
	}

	return &locale, nil
}

// Update modifies an existing locale
func (s *LocaleService) Update(ctx context.Context, id string, req models.UpdateLocaleRequest) (*models.Locale, error) {
	update := bson.M{
		"updatedAt": time.Now(),
	}

	if req.Name != nil {
		update["name"] = strings.TrimSpace(*req.Name)
	}
	if req.NativeName != nil {
		update["nativeName"] = strings.TrimSpace(*req.NativeName)
	}
	if req.IsEnabled != nil {
		update["isEnabled"] = *req.IsEnabled
	}
	if req.SortOrder != nil {
		update["sortOrder"] = *req.SortOrder
	}

	// Handle setting default locale
	if req.IsDefault != nil && *req.IsDefault {
		// First, unset all other defaults
		_, err := s.Coll.UpdateMany(ctx, bson.M{}, bson.M{"$set": bson.M{"isDefault": false}})
		if err != nil {
			return nil, err
		}
		update["isDefault"] = true
	}

	_, err := s.Coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

// Delete removes a locale (with safety checks)
func (s *LocaleService) Delete(ctx context.Context, id string) error {
	// Check if it's the default locale
	locale, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if locale.IsDefault {
		return fmt.Errorf("cannot delete the default locale")
	}

	// Check if it's the last enabled locale
	enabledLocales, err := s.GetAll(ctx, true)
	if err != nil {
		return err
	}

	if len(enabledLocales) == 1 && enabledLocales[0].ID == id {
		return fmt.Errorf("cannot delete the last enabled locale")
	}

	_, err = s.Coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// SetDefault sets a locale as the default (and unsets others)
func (s *LocaleService) SetDefault(ctx context.Context, id string) error {
	// Check if locale exists and is enabled
	locale, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !locale.IsEnabled {
		return fmt.Errorf("cannot set disabled locale as default")
	}

	// Unset all defaults
	_, err = s.Coll.UpdateMany(ctx, bson.M{}, bson.M{"$set": bson.M{"isDefault": false}})
	if err != nil {
		return err
	}

	// Set new default
	_, err = s.Coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isDefault": true, "updatedAt": time.Now()}})
	return err
}
