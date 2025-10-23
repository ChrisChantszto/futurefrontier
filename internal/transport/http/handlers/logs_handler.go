package handlers

import (
	"github.com/ChrisChantszto/futurefrontier/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LogsHandler struct {
	esService *service.ElasticsearchService
	log       *zap.Logger
}

func NewLogsHandler(esService *service.ElasticsearchService, log *zap.Logger) *LogsHandler {
	return &LogsHandler{
		esService: esService,
		log:       log,
	}
}

// GetLogs retrieves logs from Elasticsearch
func (h *LogsHandler) GetLogs(c *fiber.Ctx) error {
	// Get query parameters
	size := c.QueryInt("size", 10000) // Increased to support pagination
	searchQuery := c.Query("search", "")
	statusFilter := c.Query("status", "all")
	timeRange := c.Query("time_range", "24h")
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	h.log.Info("Fetching logs", 
		zap.Int("size", size),
		zap.String("search", searchQuery),
		zap.String("status", statusFilter),
		zap.String("time_range", timeRange),
	)

	// Fetch logs using SearchLogs or SearchLogsByDateRange
	var logs []map[string]interface{}
	var err error

	if timeRange == "custom" && startDate != "" && endDate != "" {
		logs, err = h.esService.SearchLogsByDateRange(c.Context(), startDate, endDate, size)
	} else {
		logs, err = h.esService.SearchLogs(c.Context(), timeRange, size)
	}
	if err != nil {
		h.log.Error("Failed to fetch logs", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch logs",
			"error":   err.Error(),
		})
	}

	// Apply filters in Go
	filteredLogs := make([]map[string]interface{}, 0)
	for _, log := range logs {
		// Search filter
		if searchQuery != "" {
			path, _ := log["path"].(string)
			errorMsg, _ := log["error_message"].(string)
			method, _ := log["method"].(string)
			
			if !contains(path, searchQuery) && !contains(errorMsg, searchQuery) && !contains(method, searchQuery) {
				continue
			}
		}

		// Status filter
		if statusFilter != "all" {
			statusCode, ok := log["status_code"].(float64)
			if !ok {
				continue
			}
			
			if statusFilter == "success" && statusCode >= 400 {
				continue
			}
			if statusFilter == "error" && statusCode < 400 {
				continue
			}
		}

		filteredLogs = append(filteredLogs, log)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"logs":  filteredLogs,
			"count": len(filteredLogs),
		},
	})
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		(s == substr || len(s) >= len(substr) && 
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
				containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetStats retrieves statistics from Elasticsearch
func (h *LogsHandler) GetStats(c *fiber.Ctx) error {
	h.log.Info("Fetching log statistics")

	// Get logs from last 7 days
	allLogs, err := h.esService.SearchLogs(c.Context(), "7d", 10000)
	
	if err != nil {
		h.log.Error("Failed to fetch logs for stats", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch statistics",
			"error":   err.Error(),
		})
	}

	// Calculate stats
	totalLogs := len(allLogs)
	errorCount := 0
	totalLatency := 0
	endpointsMap := make(map[string]bool)

	for _, log := range allLogs {
		// Extract status code
		if statusCode, ok := log["status_code"].(float64); ok && statusCode >= 400 {
			errorCount++
		}
		
		// Extract latency
		if latency, ok := log["latency_ms"].(float64); ok {
			totalLatency += int(latency)
		}
		
		// Extract path
		if path, ok := log["path"].(string); ok {
			endpointsMap[path] = true
		}
	}

	errorRate := 0.0
	avgLatency := 0
	if totalLogs > 0 {
		errorRate = float64(errorCount) / float64(totalLogs) * 100
		avgLatency = totalLatency / totalLogs
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total_logs":       totalLogs,
			"error_rate":       errorRate,
			"avg_latency":      avgLatency,
			"active_endpoints": len(endpointsMap),
		},
	})
}
