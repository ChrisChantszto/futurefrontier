package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port              string
	MongoURI          string
	DBName            string
	JWTSecret         string
	JWTRefreshSecret  string
	CookieDomain      string
	SecureCookies     bool

	SuperAdminEmail    string
	SuperAdminPassword string

	PageLimit    int
	UserLimit    int
	LanguageLimit int

	AllowCustomSMTP bool
	CompanySMTP     SMTPConfig
}

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromName  string
	FromEmail string
	UseTLS    bool
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getbool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	if v == "1" || v == "true" || v == "TRUE" {
		return true
	}
	return false
}

func atoi(s string, def int) int {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return def
	}
	return n
}

func Load() Config {
	return Config{
		Port:              getenv("PORT", "8080"),
		MongoURI:          getenv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:            getenv("DB_NAME", "onetake"),
		JWTSecret:         getenv("JWT_SECRET", "dev-secret"),
		JWTRefreshSecret:  getenv("JWT_REFRESH_SECRET", "dev-refresh-secret"),
		CookieDomain:      getenv("COOKIE_DOMAIN", ""),
		SecureCookies:     getbool("SECURE_COOKIES", false),
		SuperAdminEmail:   getenv("SUPERADMIN_EMAIL", ""),
		SuperAdminPassword: getenv("SUPERADMIN_PASSWORD", ""),
		PageLimit:          atoi(getenv("PAGE_LIMIT", "100"), 100),
		UserLimit:          atoi(getenv("USER_LIMIT", "100"), 100),
		LanguageLimit:      atoi(getenv("LANGUAGE_LIMIT", "3"), 3),
		AllowCustomSMTP:    getbool("ALLOW_CUSTOM_SMTP", false),
		CompanySMTP: SMTPConfig{
			Host:      getenv("ONETAKE_SMTP_HOST", ""),
			Port:      atoi(getenv("ONETAKE_SMTP_PORT", "587"), 587),
			Username:  getenv("ONETAKE_SMTP_USER", ""),
			Password:  getenv("ONETAKE_SMTP_PASS", ""),
			FromName:  getenv("ONETAKE_SMTP_FROM_NAME", "Onetake CMS"),
			FromEmail: getenv("ONETAKE_SMTP_FROM_EMAIL", ""),
			UseTLS:    getbool("ONETAKE_SMTP_TLS", true),
		},
	}
}