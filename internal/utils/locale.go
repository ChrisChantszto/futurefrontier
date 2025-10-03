package utils

import "strings"

// NormalizeLocale maps various locale inputs to canonical values used by the app.
// Supported outputs: "en", "zh-hans", "zh-hant", "de", "fr", "ja", etc.
func NormalizeLocale(input string) string {
	if input == "" {
		return "en"
	}
	lc := strings.ToLower(strings.TrimSpace(input))

	// Strip quality and take first token if a comma-separated list is passed (e.g., Accept-Language)
	if idx := strings.Index(lc, ","); idx != -1 {
		lc = strings.TrimSpace(lc[:idx])
	}

	// Common aliases and script detection for Chinese
	if lc == "zh" || strings.HasPrefix(lc, "zh-cn") || strings.HasPrefix(lc, "zh-sg") || strings.Contains(lc, "hans") {
		return "zh-hans"
	}
	if strings.HasPrefix(lc, "zh-tw") || strings.HasPrefix(lc, "zh-hk") || strings.HasPrefix(lc, "zh-mo") || strings.Contains(lc, "hant") {
		return "zh-hant"
	}

	// Extract base language code (e.g., "en-US" -> "en", "de-DE" -> "de")
	if idx := strings.Index(lc, "-"); idx != -1 {
		base := lc[:idx]
		// Keep full code for Chinese variants
		if base == "zh" {
			return lc
		}
		return base
	}

	return lc
}

// IsValidLocaleCode checks if a locale code is in a valid format
func IsValidLocaleCode(code string) bool {
	if code == "" {
		return false
	}
	// Basic validation: lowercase letters and hyphens only
	lc := strings.ToLower(code)
	for _, r := range lc {
		if !((r >= 'a' && r <= 'z') || r == '-') {
			return false
		}
	}
	return true
}

// GetRequestLocaleFromRawHeaders prefers X-Locale and falls back to Accept-Language
// Accept-Language may contain multiple tags; we only consider the first tag
func GetRequestLocaleFromRawHeaders(xLocale string, acceptLanguage string) string {
	x := strings.TrimSpace(xLocale)
	if x != "" {
		return NormalizeLocale(x)
	}
	al := strings.TrimSpace(acceptLanguage)
	if al != "" {
		return NormalizeLocale(al)
	}
	return "en"
}
