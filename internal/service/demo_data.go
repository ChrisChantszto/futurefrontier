package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ChrisChantszto/futurefrontier/internal/models"
	"go.uber.org/zap"
)

type DemoDataGenerator struct {
	esService *ElasticsearchService
	logger    *zap.Logger
}

func NewDemoDataGenerator(es *ElasticsearchService, logger *zap.Logger) *DemoDataGenerator {
	return &DemoDataGenerator{
		esService: es,
		logger:    logger,
	}
}

// GenerateRealisticLogs generates realistic API logs for demo purposes
func (d *DemoDataGenerator) GenerateRealisticLogs(ctx context.Context, count int) error {
	endpoints := []string{
		"/api/auth/login",
		"/api/auth/register",
		"/api/auth/logout",
		"/api/users/list",
		"/api/users/profile",
		"/api/users/update",
		"/api/pages/home",
		"/api/pages/about",
		"/api/pages/contact",
		"/api/settings",
		"/api/locale/en",
		"/api/locale/zh",
	}

	methods := map[string][]string{
		"/api/auth/login":    {"POST"},
		"/api/auth/register": {"POST"},
		"/api/auth/logout":   {"POST"},
		"/api/users/list":    {"GET"},
		"/api/users/profile": {"GET"},
		"/api/users/update":  {"PUT", "PATCH"},
		"/api/pages/home":    {"GET"},
		"/api/pages/about":   {"GET"},
		"/api/pages/contact": {"GET"},
		"/api/settings":      {"GET", "PATCH"},
		"/api/locale/en":     {"GET"},
		"/api/locale/zh":     {"GET"},
	}

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	ips := []string{
		"192.168.1.100",
		"192.168.1.101",
		"192.168.1.102",
		"10.0.0.50",
		"10.0.0.51",
		"172.16.0.10",
		"172.16.0.11",
	}

	d.logger.Info("Starting demo data generation", zap.Int("count", count))

	for i := 0; i < count; i++ {
		endpoint := endpoints[rand.Intn(len(endpoints))]
		methodList := methods[endpoint]
		method := methodList[rand.Intn(len(methodList))]

		// Generate realistic status codes with patterns
		statusCode := d.generateStatusCode(endpoint, i, count)
		latency := d.generateLatency(endpoint, statusCode)
		errorMsg := d.generateErrorMessage(statusCode, endpoint)

		// Calculate timestamp (spread over the last hour)
		timestamp := time.Now().Add(-time.Duration(count-i) * time.Minute)

		log := &models.APILogEntry{
			ProjectID:     "futurefrontier",
			Timestamp:     timestamp,
			RequestID:     fmt.Sprintf("req-%d-%d", time.Now().Unix(), i),
			Method:        method,
			Path:          endpoint,
			Route:         endpoint,
			StatusCode:    statusCode,
			Latency:       latency,
			ClientIP:      ips[rand.Intn(len(ips))],
			ClientPort:    fmt.Sprintf("%d", 50000+rand.Intn(10000)),
			UserAgent:     userAgents[rand.Intn(len(userAgents))],
			Referer:       "http://localhost:3000",
			Locale:        "en",
			ErrorMessage:  errorMsg,
			BytesSent:     int64(200 + rand.Intn(5000)),
			BytesReceived: int64(100 + rand.Intn(1000)),
		}

		if err := d.esService.WriteAPILog(ctx, log); err != nil {
			d.logger.Error("Failed to write demo log", zap.Error(err), zap.Int("index", i))
			// Continue generating other logs
		}

		// Log progress every 100 logs
		if (i+1)%100 == 0 {
			d.logger.Info("Demo data generation progress", zap.Int("generated", i+1), zap.Int("total", count))
		}
	}

	d.logger.Info("Demo data generation completed", zap.Int("count", count))
	return nil
}

func (d *DemoDataGenerator) generateStatusCode(endpoint string, index, total int) int {
	// Pattern 1: /api/users/list fails frequently (database timeout pattern)
	if endpoint == "/api/users/list" && rand.Float64() < 0.30 {
		return 500
	}

	// Pattern 2: /api/auth/login has authentication failures
	if endpoint == "/api/auth/login" && rand.Float64() < 0.15 {
		return 401
	}

	// Pattern 3: Traffic spike causes failures in the middle of the dataset
	// This simulates a peak hour scenario
	if index > total/3 && index < 2*total/3 {
		if rand.Float64() < 0.25 {
			// Service unavailable during peak
			return 503
		}
	}

	// Pattern 4: /api/users/update sometimes has validation errors
	if endpoint == "/api/users/update" && rand.Float64() < 0.10 {
		return 400
	}

	// Pattern 5: Occasional 404s for pages
	if (endpoint == "/api/pages/home" || endpoint == "/api/pages/about") && rand.Float64() < 0.05 {
		return 404
	}

	// Normal success (85% of remaining requests)
	if rand.Float64() < 0.85 {
		return 200
	}

	// Random errors (15% of remaining)
	errorCodes := []int{400, 404, 500, 502, 503}
	return errorCodes[rand.Intn(len(errorCodes))]
}

func (d *DemoDataGenerator) generateLatency(endpoint string, statusCode int) int64 {
	baseLatency := 100

	// Slow endpoints (realistic patterns)
	switch endpoint {
	case "/api/users/list":
		baseLatency = 450 // Database query is slow
	case "/api/users/profile":
		baseLatency = 200
	case "/api/pages/home":
		baseLatency = 150
	case "/api/auth/login":
		baseLatency = 180 // Authentication takes time
	case "/api/settings":
		baseLatency = 120
	default:
		baseLatency = 100
	}

	// Errors are often slower (timeouts, retries)
	if statusCode >= 500 {
		baseLatency += 300
	} else if statusCode >= 400 {
		baseLatency += 50
	}

	// Add random variance (±30%)
	variance := int(float64(baseLatency) * 0.3)
	adjustment := rand.Intn(2*variance) - variance
	
	result := baseLatency + adjustment
	if result < 10 {
		result = 10
	}

	return int64(result)
}

func (d *DemoDataGenerator) generateErrorMessage(statusCode int, endpoint string) string {
	if statusCode == 200 {
		return ""
	}

	errorMessages := map[int][]string{
		400: {
			"Invalid request body",
			"Missing required field: email",
			"Validation failed: password too short",
			"Invalid JSON format",
		},
		401: {
			"Invalid credentials",
			"Token expired",
			"Authentication failed",
			"Invalid email or password",
		},
		404: {
			"Resource not found",
			"Endpoint not found",
			"Page does not exist",
		},
		500: {
			"Database connection timeout",
			"Internal server error",
			"Null pointer exception",
			"Failed to execute query",
			"Connection pool exhausted",
		},
		502: {
			"Bad gateway",
			"Upstream server error",
			"Failed to connect to database",
		},
		503: {
			"Service unavailable",
			"Too many connections",
			"Database pool exhausted",
			"Server overloaded",
			"Rate limit exceeded",
		},
	}

	messages, ok := errorMessages[statusCode]
	if !ok {
		return "Unknown error"
	}

	// Add endpoint-specific context
	msg := messages[rand.Intn(len(messages))]
	
	if endpoint == "/api/users/list" && statusCode == 500 {
		// Make database timeout pattern more obvious
		return "Database connection timeout: failed to fetch user list"
	}

	return msg
}
