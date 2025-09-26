package utils

import "strings"

// NormalizeLocale maps various locale inputs to canonical values used by the app.
// Supported outputs: "en", "zh-hans", "zh-hant" (extend as needed)
func NormalizeLocale(input string) string {
	if input == "" {
		return "en"
	}
	lc := strings.ToLower(strings.TrimSpace(input))

	// Strip quality and take first token if a comma-separated list is passed (e.g., Accept-Language)
	if idx := strings.Index(lc, ","); idx != -1 {
		lc = strings.TrimSpace(lc[:idx])
	}

	// Common aliases and script detection
	if lc == "zh" || strings.HasPrefix(lc, "zh-cn") || strings.HasPrefix(lc, "zh-sg") || strings.Contains(lc, "hans") {
		return "zh-hans"
	}
	if strings.HasPrefix(lc, "zh-tw") || strings.HasPrefix(lc, "zh-hk") || strings.HasPrefix(lc, "zh-mo") || strings.Contains(lc, "hant") {
		return "zh-hant"
	}

	// Already canonical (en, zh-hans, zh-hant)
	switch lc {
	case "en", "zh-hans", "zh-hant":
		return lc
	}

	return lc
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
