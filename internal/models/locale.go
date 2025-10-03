package models

import "time"

// Locale represents a language/locale configuration
type Locale struct {
	ID          string    `bson:"_id" json:"id"`                   // e.g., "en", "zh-hans", "zh-hant", "de"
	Code        string    `bson:"code" json:"code"`                // Same as ID, canonical locale code
	Name        string    `bson:"name" json:"name"`                // Display name, e.g., "English", "简体中文"
	NativeName  string    `bson:"nativeName" json:"nativeName"`    // Native name, e.g., "English", "简体中文"
	IsDefault   bool      `bson:"isDefault" json:"isDefault"`      // Is this the default locale?
	IsEnabled   bool      `bson:"isEnabled" json:"isEnabled"`      // Is this locale active?
	SortOrder   int       `bson:"sortOrder" json:"sortOrder"`      // Display order
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}

// LocaleConfig represents the frontend routing configuration
type LocaleConfig struct {
	Locales       []string `json:"locales"`       // Array of enabled locale codes
	DefaultLocale string   `json:"defaultLocale"` // Default locale code
}

// CreateLocaleRequest for adding new locales
type CreateLocaleRequest struct {
	Code       string `json:"code" validate:"required"`       // e.g., "de", "fr", "ja"
	Name       string `json:"name" validate:"required"`       // e.g., "German"
	NativeName string `json:"nativeName" validate:"required"` // e.g., "Deutsch"
	IsEnabled  bool   `json:"isEnabled"`
	SortOrder  int    `json:"sortOrder"`
}

// UpdateLocaleRequest for updating existing locales
type UpdateLocaleRequest struct {
	Name       *string `json:"name,omitempty"`
	NativeName *string `json:"nativeName,omitempty"`
	IsEnabled  *bool   `json:"isEnabled,omitempty"`
	IsDefault  *bool   `json:"isDefault,omitempty"`
	SortOrder  *int    `json:"sortOrder,omitempty"`
}
