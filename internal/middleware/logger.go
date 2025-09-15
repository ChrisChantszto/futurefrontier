package middleware

import (
	"net"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/models"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/service"
)

// LoggerMiddleware creates a new logging middleware
func LoggerMiddleware(loggerService *service.LoggerService, config models.LogConfig, logger *zap.Logger) fiber.Handler {
	// Sensitive headers to exclude from logging
	sensitiveHeaders := map[string]bool{
		"authorization": true,
		"cookie":        true,
		"set-cookie":    true,
		"x-api-key":     true,
		"x-auth-token":  true,
	}

	// Add custom sensitive headers from config
	for _, header := range config.SensitiveHeaders {
		sensitiveHeaders[strings.ToLower(header)] = true
	}

	// Sensitive paths where we shouldn't log request/response bodies
	sensitivePaths := make(map[string]bool)
	for _, path := range config.SensitivePaths {
		sensitivePaths[path] = true
	}

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Generate or get request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("X-Request-ID", requestID)
		}

		// Store original body for logging
		var requestBody string
		if config.EnableRequestBody && !sensitivePaths[c.Path()] {
			if c.Body() != nil && len(c.Body()) > 0 && int64(len(c.Body())) <= config.MaxBodySize {
				requestBody = string(c.Body())
			}
		}

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get response body if enabled (simplified approach)
		var responseBody string
		if config.EnableResponseBody && !sensitivePaths[c.Path()] {
			if c.Response().Body() != nil && int64(len(c.Response().Body())) <= config.MaxBodySize {
				responseBody = string(c.Response().Body())
			}
		}

		// Extract client IP and port
		clientIP, clientPort := extractIPAndPort(c)

		// Build query parameters map
		queryParams := make(map[string]string)
		c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
			queryParams[string(key)] = string(value)
		})

		// Build headers map (excluding sensitive ones)
		var headers map[string]string
		if config.EnableHeaders {
			headers = make(map[string]string)
			c.Request().Header.VisitAll(func(key, value []byte) {
				headerKey := strings.ToLower(string(key))
				if !sensitiveHeaders[headerKey] {
					headers[string(key)] = string(value)
				}
			})
		}

		// Get route name/pattern
		route := c.Route().Path
		if route == "" {
			route = c.Path()
		}

		// Create API log entry
		logEntry := &models.APILogEntry{
			ProjectID:     config.ProjectID,
			Timestamp:     start,
			RequestID:     requestID,
			Method:        c.Method(),
			Path:          c.Path(),
			Route:         route,
			StatusCode:    c.Response().StatusCode(),
			Latency:       latency.Milliseconds(),
			ClientIP:      clientIP,
			ClientPort:    clientPort,
			UserAgent:     c.Get("User-Agent"),
			Referer:       c.Get("Referer"),
			QueryParams:   queryParams,
			Headers:       headers,
			RequestBody:   requestBody,
			ResponseBody:  responseBody,
			BytesSent:     int64(len(c.Response().Body())),
			BytesReceived: int64(len(c.Body())),
		}

		// Add error information if request failed
		if err != nil {
			logEntry.ErrorMessage = err.Error()
			
			// Also create error log entry for 5xx errors
			if c.Response().StatusCode() >= 500 {
				errorLogEntry := &models.ErrorLogEntry{
					ProjectID:    config.ProjectID,
					Timestamp:    start,
					RequestID:    requestID,
					Method:       c.Method(),
					Path:         c.Path(),
					StatusCode:   c.Response().StatusCode(),
					ErrorMessage: err.Error(),
					ErrorType:    "server_error",
					ClientIP:     clientIP,
					UserAgent:    c.Get("User-Agent"),
					RequestBody:  requestBody,
				}

				// Log error asynchronously
				loggerService.LogError(errorLogEntry)
			}
		}

		// Log API request asynchronously
		loggerService.LogAPI(logEntry)

		// Log to structured logger for immediate visibility
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", latency),
			zap.String("ip", clientIP),
			zap.String("user_agent", c.Get("User-Agent")),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Error("Request completed with error", fields...)
		} else {
			logger.Info("Request completed", fields...)
		}

		return err
	}
}

// extractIPAndPort extracts client IP and port from the request
func extractIPAndPort(c *fiber.Ctx) (string, string) {
	// Try to get real IP from headers (for proxy/load balancer scenarios)
	ip := c.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if idx := strings.Index(ip, ","); idx != -1 {
			ip = strings.TrimSpace(ip[:idx])
		}
	}

	if ip == "" {
		ip = c.Get("X-Real-IP")
	}

	if ip == "" {
		ip = c.Get("CF-Connecting-IP") // Cloudflare
	}

	if ip == "" {
		// Fallback to remote address
		remoteAddr := c.Context().RemoteAddr().String()
		host, port, err := net.SplitHostPort(remoteAddr)
		if err == nil {
			return host, port
		}
		return remoteAddr, ""
	}

	// For forwarded IPs, we don't have port information
	return ip, ""
}

// RequestIDMiddleware adds a request ID to each request if not present
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("X-Request-ID", requestID)
		}
		
		// Also set it in response headers for tracing
		c.Set("X-Request-ID", requestID)
		
		return c.Next()
	}
}

// ErrorHandlerMiddleware creates structured error logs for panics and errors
func ErrorHandlerMiddleware(loggerService *service.LoggerService, config models.LogConfig) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		// Check if it's a Fiber error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		// Get request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Extract client info
		clientIP, _ := extractIPAndPort(c)

		// Get request body if enabled
		var requestBody string
		if config.EnableRequestBody && c.Body() != nil && int64(len(c.Body())) <= config.MaxBodySize {
			requestBody = string(c.Body())
		}

		// Create error log entry
		errorLogEntry := &models.ErrorLogEntry{
			ProjectID:    config.ProjectID,
			Timestamp:    time.Now(),
			RequestID:    requestID,
			Method:       c.Method(),
			Path:         c.Path(),
			StatusCode:   code,
			ErrorMessage: err.Error(),
			ErrorType:    getErrorType(code),
			ClientIP:     clientIP,
			UserAgent:    c.Get("User-Agent"),
			RequestBody:  requestBody,
		}

		// Log error asynchronously
		loggerService.LogError(errorLogEntry)

		// Set standardized response envelope while preserving HTTP status code
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(code).JSON(fiber.Map{
			"success":    false,
			"message":    message,
			"data":       fiber.Map{},
			"request_id": requestID,
		})
	}
}

// getErrorType categorizes errors based on status code
func getErrorType(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "server_error"
	case statusCode >= 400:
		return "client_error"
	case statusCode >= 300:
		return "redirect"
	default:
		return "unknown"
	}
}

// LoggingStatsHandler returns logging statistics
func LoggingStatsHandler(loggerService *service.LoggerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		processed, failed, queued := loggerService.GetStats()
		
		return c.JSON(fiber.Map{
			"logging_stats": fiber.Map{
				"processed_logs": processed,
				"failed_logs":    failed,
				"queued_logs":    queued,
			},
		})
	}
}
