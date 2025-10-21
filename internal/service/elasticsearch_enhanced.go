package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// SearchLogsByPattern searches logs matching a specific pattern
func (es *ElasticsearchService) SearchLogsByPattern(ctx context.Context, pattern string, timeRange string, limit int) ([]map[string]interface{}, error) {
	if es.config.LocalMode || es.client == nil {
		es.logger.Warn("Cannot search logs in local mode")
		return []map[string]interface{}{}, nil
	}

	duration, err := parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	fromTime := time.Now().Add(-duration)
	
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"timestamp": map[string]interface{}{
								"gte": fromTime.Format(time.RFC3339),
								"lte": time.Now().Format(time.RFC3339),
							},
						},
					},
					{
						"multi_match": map[string]interface{}{
							"query":  pattern,
							"fields": []string{"path", "error_message", "request_body", "response_body"},
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{"timestamp": map[string]string{"order": "desc"}},
		},
		"size": limit,
	}

	return es.executeSearch(ctx, query)
}

// AggregateByEndpoint aggregates logs by API endpoint with statistics
func (es *ElasticsearchService) AggregateByEndpoint(ctx context.Context, timeRange string) (map[string]interface{}, error) {
	if es.config.LocalMode || es.client == nil {
		es.logger.Warn("Cannot aggregate logs in local mode")
		return map[string]interface{}{}, nil
	}

	duration, err := parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	fromTime := time.Now().Add(-duration)
	
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"gte": fromTime.Format(time.RFC3339),
					"lte": time.Now().Format(time.RFC3339),
				},
			},
		},
		"aggs": map[string]interface{}{
			"endpoints": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "path",
					"size":  50,
				},
				"aggs": map[string]interface{}{
					"status_codes": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "status_code",
						},
					},
					"avg_latency": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "latency_ms",
						},
					},
					"error_count": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"status_code": map[string]interface{}{
									"gte": 400,
								},
							},
						},
					},
				},
			},
		},
		"size": 0,
	}

	return es.executeAggregation(ctx, query)
}

// SearchFailingAPIs finds APIs with high failure rates
func (es *ElasticsearchService) SearchFailingAPIs(ctx context.Context, timeRange string, minFailureRate float64) ([]map[string]interface{}, error) {
	aggregation, err := es.AggregateByEndpoint(ctx, timeRange)
	if err != nil {
		return nil, err
	}

	// Extract endpoints with high failure rates
	failingAPIs := make([]map[string]interface{}, 0)
	
	if aggs, ok := aggregation["aggregations"].(map[string]interface{}); ok {
		if endpoints, ok := aggs["endpoints"].(map[string]interface{}); ok {
			if buckets, ok := endpoints["buckets"].([]interface{}); ok {
				for _, bucket := range buckets {
					b := bucket.(map[string]interface{})
					docCount := b["doc_count"].(float64)
					
					// Get error count
					errorCount := 0.0
					if errorCountAgg, ok := b["error_count"].(map[string]interface{}); ok {
						if ec, ok := errorCountAgg["doc_count"].(float64); ok {
							errorCount = ec
						}
					}
					
					// Calculate failure rate
					failureRate := errorCount / docCount
					
					if failureRate >= minFailureRate {
						failingAPIs = append(failingAPIs, map[string]interface{}{
							"endpoint":      b["key"],
							"total_requests": docCount,
							"error_count":    errorCount,
							"failure_rate":   failureRate,
							"avg_latency":    b["avg_latency"],
						})
					}
				}
			}
		}
	}

	return failingAPIs, nil
}

// SearchLogsByEndpoint searches logs for a specific endpoint
func (es *ElasticsearchService) SearchLogsByEndpoint(ctx context.Context, endpoint string, timeRange string, limit int) ([]map[string]interface{}, error) {
	if es.config.LocalMode || es.client == nil {
		es.logger.Warn("Cannot search logs in local mode")
		return []map[string]interface{}{}, nil
	}

	duration, err := parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	fromTime := time.Now().Add(-duration)
	
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"timestamp": map[string]interface{}{
								"gte": fromTime.Format(time.RFC3339),
								"lte": time.Now().Format(time.RFC3339),
							},
						},
					},
					{
						"term": map[string]interface{}{
							"path": endpoint,
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{"timestamp": map[string]string{"order": "desc"}},
		},
		"size": limit,
	}

	return es.executeSearch(ctx, query)
}

// Helper method to execute search
func (es *ElasticsearchService) executeSearch(ctx context.Context, query map[string]interface{}) ([]map[string]interface{}, error) {
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{"api-logs-*"},
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

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract hits
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}

	hitsList, ok := hits["hits"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}

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

	return logs, nil
}

// Helper method to execute aggregation
func (es *ElasticsearchService) executeAggregation(ctx context.Context, query map[string]interface{}) (map[string]interface{}, error) {
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{"api-logs-*"},
		Body:  bytes.NewReader(queryJSON),
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return nil, fmt.Errorf("aggregation request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch aggregation error: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
