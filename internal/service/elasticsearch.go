package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.uber.org/zap"

	"github.com/ChrisChantszto/futurefrontier/internal/config"
	"github.com/ChrisChantszto/futurefrontier/internal/models"
)

type ElasticsearchService struct {
	client *elasticsearch.Client
	logger *zap.Logger
	config config.LoggingConfig
}

// NewElasticsearchService creates a new Elasticsearch service
func NewElasticsearchService(cfg config.LoggingConfig, logger *zap.Logger) (*ElasticsearchService, error) {
	if cfg.LocalMode {
		logger.Info("Elasticsearch service running in local mode - logs will be written to files")
		return &ElasticsearchService{
			client: nil,
			logger: logger,
			config: cfg,
		}, nil
	}

	esConfig := elasticsearch.Config{
		Addresses: []string{cfg.ElasticsearchURL},
	}

	if cfg.ElasticsearchUser != "" && cfg.ElasticsearchPass != "" {
		esConfig.Username = cfg.ElasticsearchUser
		esConfig.Password = cfg.ElasticsearchPass
	}

	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := esapi.InfoRequest{}
	res, err := req.Do(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch connection error: %s", res.Status())
	}

	logger.Info("Successfully connected to Elasticsearch", zap.String("url", cfg.ElasticsearchURL))

	service := &ElasticsearchService{
		client: client,
		logger: logger,
		config: cfg,
	}

	// Create indexes if they don't exist
	if err := service.createIndexes(ctx); err != nil {
		logger.Warn("Failed to create elasticsearch indexes", zap.Error(err))
	}

	return service, nil
}

// createIndexes creates the required indexes for API and error logs
func (es *ElasticsearchService) createIndexes(ctx context.Context) error {
	apiIndexName := fmt.Sprintf("api-logs-%s", time.Now().Format("2006-01"))
	errorIndexName := fmt.Sprintf("error-logs-%s", time.Now().Format("2006-01"))

	// API logs index mapping
	apiMapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"project_id":      map[string]interface{}{"type": "keyword"},
				"timestamp":       map[string]interface{}{"type": "date"},
				"request_id":      map[string]interface{}{"type": "keyword"},
				"method":          map[string]interface{}{"type": "keyword"},
				"path":            map[string]interface{}{"type": "keyword"},
				"route":           map[string]interface{}{"type": "keyword"},
				"status_code":     map[string]interface{}{"type": "integer"},
				"latency_ms":      map[string]interface{}{"type": "long"},
				"client_ip":       map[string]interface{}{"type": "ip"},
				"client_port":     map[string]interface{}{"type": "keyword"},
				"user_agent":      map[string]interface{}{"type": "text"},
				"referer":         map[string]interface{}{"type": "keyword"},
				"locale":          map[string]interface{}{"type": "keyword"},
				"query_params":    map[string]interface{}{"type": "object"},
				"headers":         map[string]interface{}{"type": "object"},
				"request_body":    map[string]interface{}{"type": "text"},
				"response_body":   map[string]interface{}{"type": "text"},
				"bytes_sent":      map[string]interface{}{"type": "long"},
				"bytes_received":  map[string]interface{}{"type": "long"},
				"error_message":   map[string]interface{}{"type": "text"},
			},
		},
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
	}

	// Error logs index mapping
	errorMapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"project_id":     map[string]interface{}{"type": "keyword"},
				"timestamp":      map[string]interface{}{"type": "date"},
				"request_id":     map[string]interface{}{"type": "keyword"},
				"method":         map[string]interface{}{"type": "keyword"},
				"path":           map[string]interface{}{"type": "keyword"},
				"status_code":    map[string]interface{}{"type": "integer"},
				"error_message":  map[string]interface{}{"type": "text"},
				"error_type":     map[string]interface{}{"type": "keyword"},
				"stack_trace":    map[string]interface{}{"type": "text"},
				"client_ip":      map[string]interface{}{"type": "ip"},
				"user_agent":     map[string]interface{}{"type": "text"},
				"locale":         map[string]interface{}{"type": "keyword"},
				"request_body":   map[string]interface{}{"type": "text"},
			},
		},
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
	}

	// Create API logs index
	if err := es.createIndex(ctx, apiIndexName, apiMapping); err != nil {
		return fmt.Errorf("failed to create API logs index: %w", err)
	}

	// Create error logs index
	if err := es.createIndex(ctx, errorIndexName, errorMapping); err != nil {
		return fmt.Errorf("failed to create error logs index: %w", err)
	}

	return nil
}

// createIndex creates an index with the given mapping
func (es *ElasticsearchService) createIndex(ctx context.Context, indexName string, mapping map[string]interface{}) error {
	if es.client == nil {
		return nil // Skip in local mode
	}

	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("failed to marshal mapping: %w", err)
	}

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader(mappingJSON),
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && !strings.Contains(res.Status(), "400") { // 400 means index already exists
		return fmt.Errorf("elasticsearch index creation error: %s", res.Status())
	}

	es.logger.Info("Created elasticsearch index", zap.String("index", indexName))
	return nil
}

// WriteAPILog writes an API log entry to Elasticsearch
func (es *ElasticsearchService) WriteAPILog(ctx context.Context, logEntry *models.APILogEntry) error {
	if es.config.LocalMode {
		return es.writeToLocalFile("api", logEntry)
	}

	if es.client == nil {
		return fmt.Errorf("elasticsearch client not initialized")
	}

	indexName := fmt.Sprintf("api-logs-%s", time.Now().Format("2006-01"))
	
	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	req := esapi.IndexRequest{
		Index: indexName,
		Body:  bytes.NewReader(logJSON),
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch indexing error: %s", res.Status())
	}

	return nil
}

// WriteErrorLog writes an error log entry to Elasticsearch
func (es *ElasticsearchService) WriteErrorLog(ctx context.Context, logEntry *models.ErrorLogEntry) error {
	if es.config.LocalMode {
		return es.writeToLocalFile("error", logEntry)
	}

	if es.client == nil {
		return fmt.Errorf("elasticsearch client not initialized")
	}

	indexName := fmt.Sprintf("error-logs-%s", time.Now().Format("2006-01"))
	
	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	req := esapi.IndexRequest{
		Index: indexName,
		Body:  bytes.NewReader(logJSON),
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch indexing error: %s", res.Status())
	}

	return nil
}

// writeToLocalFile writes logs to local files for development
func (es *ElasticsearchService) writeToLocalFile(logType string, logEntry interface{}) error {
	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// In local mode, just log to zap logger with structured data
	switch logType {
	case "api":
		es.logger.Info("API Log", zap.String("type", "api_log"), zap.ByteString("data", logJSON))
	case "error":
		es.logger.Error("Error Log", zap.String("type", "error_log"), zap.ByteString("data", logJSON))
	}

	return nil
}

// SearchLogs searches for logs in Elasticsearch based on time range
func (es *ElasticsearchService) SearchLogs(ctx context.Context, timeRange string, limit int) ([]map[string]interface{}, error) {
	if es.config.LocalMode || es.client == nil {
		es.logger.Warn("Cannot search logs in local mode")
		return []map[string]interface{}{}, nil
	}

	// Parse time range (e.g., "1h", "24h", "7d")
	duration, err := parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	fromTime := time.Now().Add(-duration)
	
	// Build search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"gte": fromTime.Format(time.RFC3339),
					"lte": time.Now().Format(time.RFC3339),
				},
			},
		},
		"sort": []map[string]interface{}{
			{"timestamp": map[string]string{"order": "desc"}},
		},
		"size": limit,
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Search across all API log indexes
	indexPattern := "api-logs-*"
	
	req := esapi.SearchRequest{
		Index: []string{indexPattern},
		Body:  bytes.NewReader(queryJSON),
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search error: %s", res.Status())
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract hits
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	hitsList, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hits format")
	}

	// Extract source documents
	logs := make([]map[string]interface{}, 0, len(hitsList))
	for _, hit := range hitsList {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}
		
		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}
		
		logs = append(logs, source)
	}

	es.logger.Info("Retrieved logs from Elasticsearch",
		zap.Int("count", len(logs)),
		zap.String("time_range", timeRange))

	return logs, nil
}

// parseTimeRange parses time range strings like "1h", "24h", "7d" into duration
func parseTimeRange(timeRange string) (time.Duration, error) {
	if len(timeRange) < 2 {
		return 0, fmt.Errorf("invalid time range format")
	}

	unit := timeRange[len(timeRange)-1:]
	valueStr := timeRange[:len(timeRange)-1]
	
	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return 0, fmt.Errorf("invalid time range value: %w", err)
	}

	switch unit {
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid time range unit: %s (use m, h, or d)", unit)
	}
}

// IsHealthy checks if Elasticsearch is healthy
func (es *ElasticsearchService) IsHealthy(ctx context.Context) bool {
	if es.config.LocalMode || es.client == nil {
		return true // Always healthy in local mode
	}

	res, err := es.client.Cluster.Health(
		es.client.Cluster.Health.WithContext(ctx),
		es.client.Cluster.Health.WithTimeout(5*time.Second),
	)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	return !res.IsError()
}
