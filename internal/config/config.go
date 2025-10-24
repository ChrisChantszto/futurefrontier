package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port             string
	MongoURI         string
	DBName           string
	JWTSecret        string
	JWTRefreshSecret string
	CookieDomain     string
	SecureCookies    bool
	CORSOrigins      string

	SuperAdminEmail    string
	SuperAdminPassword string

	PageLimit     int
	UserLimit     int
	LanguageLimit int

	AllowCustomSMTP bool
	CompanySMTP     SMTPConfig
	CustomSMTP      SMTPConfig

	// OTP Configuration
	AppHMACSecret         string
	OTPTTLMinutes         int
	OTPMaxAttempts        int
	OTPResendCooldownSecs int

	// i18n Configuration
	I18n I18nConfig

	// Logging Configuration
	Logging LoggingConfig

	// Google Cloud Configuration
	GoogleCloud GoogleCloudConfig
}

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromName  string
	FromEmail string
	ReplyTo   string
	BCC       string
	UseTLS    bool
}

type LoggingConfig struct {
	ProjectID           string
	ElasticsearchURL    string
	ElasticsearchUser   string
	ElasticsearchPass   string
	ElasticsearchAPIKey string // For Elastic Cloud API Key authentication
	EnableRequestBody   bool
	EnableResponseBody  bool
	EnableHeaders       bool
	MaxBodySize         int64
	QueueSize           int
	WorkerCount         int
	RetryAttempts       int
	RetryInterval       int  // seconds
	LocalMode           bool // for development - logs to file instead of ES
}

type I18nConfig struct {
	DefaultLocale    string   // Default locale if none is specified
	SupportedLocales []string // List of supported locales
	TemplatesDir     string   // Directory where email templates are stored
}

type GoogleCloudConfig struct {
	ProjectID        string  // GCP Project ID
	Location         string  // Vertex AI location (e.g., us-central1, asia-southeast1)
	CredentialsPath  string  // Path to service account JSON key file
	VertexAIModel    string  // Model name (e.g., gemini-1.5-pro, gemini-1.5-flash)
	VertexAIEndpoint string  // Optional custom endpoint
	MaxTokens        int     // Max tokens for generation
	Temperature      float64 // Temperature for generation (0.0-1.0)
	TopP             float64 // Top-p for generation
	TopK             int     // Top-k for generation
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

func atof(s string, def float64) float64 {
	var n float64
	_, err := fmt.Sscanf(s, "%f", &n)
	if err != nil {
		return def
	}
	return n
}

func Load() Config {
	return Config{
		Port:               getenv("PORT", "8080"),
		MongoURI:           getenv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:             getenv("DB_NAME", "onetake"),
		JWTSecret:          getenv("JWT_SECRET", "dev-secret"),
		JWTRefreshSecret:   getenv("JWT_REFRESH_SECRET", "dev-refresh-secret"),
		CookieDomain:       getenv("COOKIE_DOMAIN", ""),
		SecureCookies:      getbool("SECURE_COOKIES", false),
		CORSOrigins:        getenv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		SuperAdminEmail:    getenv("SUPERADMIN_EMAIL", ""),
		SuperAdminPassword: getenv("SUPERADMIN_PASSWORD", ""),
		PageLimit:          atoi(getenv("PAGE_LIMIT", "100"), 100),
		UserLimit:          atoi(getenv("USER_LIMIT", "100"), 100),
		LanguageLimit:      atoi(getenv("LANGUAGE_LIMIT", "3"), 3),
		AllowCustomSMTP:    getbool("ALLOW_CUSTOM_SMTP", false),

		// Custom SMTP (user-configured)
		CustomSMTP: SMTPConfig{
			Host:      getenv("SMTP_HOST", ""),
			Port:      atoi(getenv("SMTP_PORT", "465"), 465),
			Username:  getenv("SMTP_USER", ""),
			Password:  getenv("SMTP_PASS", ""),
			FromName:  getenv("SMTP_FROM_NAME", ""),
			FromEmail: getenv("SMTP_FROM_EMAIL", ""),
			ReplyTo:   getenv("SMTP_REPLY_TO", ""),
			BCC:       getenv("SMTP_BCC", ""),
			UseTLS:    true,
		},

		// Company SMTP (fallback)
		CompanySMTP: SMTPConfig{
			Host:      getenv("ONETAKE_SMTP_HOST", "smtpdm-ap-southeast-1.aliyun.com"),
			Port:      atoi(getenv("ONETAKE_SMTP_PORT", "465"), 465),
			Username:  getenv("ONETAKE_SMTP_USER", "system@onetakesolutions.com.hk"),
			Password:  getenv("ONETAKE_SMTP_PASS", "OneTake07102022"),
			FromName:  getenv("ONETAKE_SMTP_FROM_NAME", "OTS system"),
			FromEmail: getenv("ONETAKE_SMTP_FROM_EMAIL", "system@onetakesolutions.com.hk"),
			ReplyTo:   getenv("ONETAKE_SMTP_REPLY_TO", "info@onetakesolutions.com.hk"),
			BCC:       getenv("ONETAKE_SMTP_BCC", "tim.ho@onetakesolutions.com.hk"),
			UseTLS:    getbool("ONETAKE_SMTP_TLS", true),
		},

		// OTP Configuration
		AppHMACSecret:         getenv("APP_HMAC_SECRET", ""),
		OTPTTLMinutes:         atoi(getenv("OTP_TTL_MINUTES", "20"), 20),
		OTPMaxAttempts:        atoi(getenv("OTP_MAX_ATTEMPTS", "5"), 5),
		OTPResendCooldownSecs: atoi(getenv("OTP_RESEND_COOLDOWN_SECONDS", "60"), 60),

		// i18n Configuration
		I18n: I18nConfig{
			DefaultLocale:    getenv("I18N_DEFAULT_LOCALE", "en"),
			SupportedLocales: strings.Split(getenv("I18N_SUPPORTED_LOCALES", "en,zh-hans,zh-hant"), ","),
			TemplatesDir:     getenv("I18N_TEMPLATES_DIR", "./templates"),
		},

		// Logging Configuration
		Logging: LoggingConfig{
			ProjectID:           getenv("LOG_PROJECT_ID", "futurefrontier"),
			ElasticsearchURL:    getenv("ELASTICSEARCH_URL", "http://localhost:9200"),
			ElasticsearchUser:   getenv("ELASTICSEARCH_USER", ""),
			ElasticsearchPass:   getenv("ELASTICSEARCH_PASS", ""),
			ElasticsearchAPIKey: getenv("ELASTICSEARCH_API_KEY", ""),
			EnableRequestBody:   getbool("LOG_ENABLE_REQUEST_BODY", false),
			EnableResponseBody:  getbool("LOG_ENABLE_RESPONSE_BODY", false),
			EnableHeaders:       getbool("LOG_ENABLE_HEADERS", false),
			MaxBodySize:         int64(atoi(getenv("LOG_MAX_BODY_SIZE", "1024"), 1024)),
			QueueSize:           atoi(getenv("LOG_QUEUE_SIZE", "1000"), 1000),
			WorkerCount:         atoi(getenv("LOG_WORKER_COUNT", "3"), 3),
			RetryAttempts:       atoi(getenv("LOG_RETRY_ATTEMPTS", "3"), 3),
			RetryInterval:       atoi(getenv("LOG_RETRY_INTERVAL", "30"), 30),
			LocalMode:           getbool("LOG_LOCAL_MODE", false), // Changed to false for production
		},

		// Google Cloud Configuration
		GoogleCloud: GoogleCloudConfig{
			ProjectID:        getenv("GCP_PROJECT_ID", ""),
			Location:         getenv("GCP_LOCATION", "us-central1"),
			CredentialsPath:  getenv("GOOGLE_APPLICATION_CREDENTIALS", ""),
			VertexAIModel:    getenv("VERTEX_AI_MODEL", "gemini-2.0-flash"),
			VertexAIEndpoint: getenv("VERTEX_AI_ENDPOINT", ""),
			MaxTokens:        atoi(getenv("VERTEX_AI_MAX_TOKENS", "8192"), 8192),
			Temperature:      atof(getenv("VERTEX_AI_TEMPERATURE", "0.7"), 0.7),
			TopP:             atof(getenv("VERTEX_AI_TOP_P", "0.95"), 0.95),
			TopK:             atoi(getenv("VERTEX_AI_TOP_K", "40"), 40),
		},
	}
}
