package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APILogEntry represents a complete API request/response log entry
type APILogEntry struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProjectID      string             `json:"project_id" bson:"project_id"`
	Timestamp      time.Time          `json:"timestamp" bson:"timestamp"`
	RequestID      string             `json:"request_id" bson:"request_id"`
	Method         string             `json:"method" bson:"method"`
	Path           string             `json:"path" bson:"path"`
	Route          string             `json:"route" bson:"route"`
	StatusCode     int                `json:"status_code" bson:"status_code"`
	Latency        int64              `json:"latency_ms" bson:"latency_ms"` // in milliseconds
	ClientIP       string             `json:"client_ip" bson:"client_ip"`
	ClientPort     string             `json:"client_port" bson:"client_port"`
	UserAgent      string             `json:"user_agent" bson:"user_agent"`
	Referer        string             `json:"referer" bson:"referer"`
	Locale         string             `json:"locale,omitempty" bson:"locale,omitempty"`
	QueryParams    map[string]string  `json:"query_params" bson:"query_params"`
	Headers        map[string]string  `json:"headers,omitempty" bson:"headers,omitempty"`
	RequestBody    string             `json:"request_body,omitempty" bson:"request_body,omitempty"`
	ResponseBody   string             `json:"response_body,omitempty" bson:"response_body,omitempty"`
	BytesSent      int64              `json:"bytes_sent" bson:"bytes_sent"`
	BytesReceived  int64              `json:"bytes_received" bson:"bytes_received"`
	ErrorMessage   string             `json:"error_message,omitempty" bson:"error_message,omitempty"`
	
	// Metadata for backup system
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	ESWritten      bool               `json:"es_written" bson:"es_written"`
	RetryCount     int                `json:"retry_count" bson:"retry_count"`
	LastRetryAt    *time.Time         `json:"last_retry_at,omitempty" bson:"last_retry_at,omitempty"`
}

// ErrorLogEntry represents error-specific log entries
type ErrorLogEntry struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProjectID      string             `json:"project_id" bson:"project_id"`
	Timestamp      time.Time          `json:"timestamp" bson:"timestamp"`
	RequestID      string             `json:"request_id" bson:"request_id"`
	Method         string             `json:"method" bson:"method"`
	Path           string             `json:"path" bson:"path"`
	StatusCode     int                `json:"status_code" bson:"status_code"`
	ErrorMessage   string             `json:"error_message" bson:"error_message"`
	ErrorType      string             `json:"error_type" bson:"error_type"`
	StackTrace     string             `json:"stack_trace,omitempty" bson:"stack_trace,omitempty"`
	ClientIP       string             `json:"client_ip" bson:"client_ip"`
	UserAgent      string             `json:"user_agent" bson:"user_agent"`
	Locale         string             `json:"locale,omitempty" bson:"locale,omitempty"`
	RequestBody    string             `json:"request_body,omitempty" bson:"request_body,omitempty"`
	
	// Metadata for backup system
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	ESWritten      bool               `json:"es_written" bson:"es_written"`
	RetryCount     int                `json:"retry_count" bson:"retry_count"`
	LastRetryAt    *time.Time         `json:"last_retry_at,omitempty" bson:"last_retry_at,omitempty"`
}

// LogConfig represents configuration for the logging system
type LogConfig struct {
	ProjectID           string
	EnableRequestBody   bool
	EnableResponseBody  bool
	EnableHeaders       bool
	MaxBodySize         int64 // Maximum size of request/response body to log
	SensitiveHeaders    []string // Headers to exclude from logging
	SensitivePaths      []string // Paths to exclude body logging
}
