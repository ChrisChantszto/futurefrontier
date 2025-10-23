package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// DetectPatterns detects recurring failure patterns in logs
func (h *AIHandler) DetectPatterns(c *fiber.Ctx) error {
	var req struct {
		TimeRange      string `json:"time_range"`
		StartDate      string `json:"start_date"`
		EndDate        string `json:"end_date"`
		Limit          int    `json:"limit"`
		MinOccurrences int    `json:"min_occurrences"`
	}

	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body", 0, err.Error())
	}

	// Set defaults
	if req.TimeRange == "" {
		req.TimeRange = "24h"
	}
	if req.Limit == 0 {
		req.Limit = 1000
	}

	// Fetch logs from Elasticsearch
	var logs []map[string]interface{}
	var err error

	if req.TimeRange == "custom" && req.StartDate != "" && req.EndDate != "" {
		logs, err = h.esClient.SearchLogsByDateRange(context.Background(), req.StartDate, req.EndDate, req.Limit)
	} else {
		logs, err = h.esClient.SearchLogs(context.Background(), req.TimeRange, req.Limit)
	}

	if err != nil {
		h.logger.Error("Failed to fetch logs", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to fetch logs", 0, err.Error())
	}

	if len(logs) == 0 {
		return JSONError(c, fiber.StatusNotFound, "No logs found for the specified time range", 0, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Detect patterns with AI
	resp, err := h.vertexAI.DetectPatterns(ctx, logs)
	if err != nil {
		h.logger.Error("Failed to detect patterns", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to detect patterns", 0, err.Error())
	}

	return JSONSuccess(c, "Patterns detected successfully", fiber.Map{
		"analysis": fiber.Map{
			"text": resp.Text,
		},
		"logs_count": len(logs),
		"time_range": req.TimeRange,
	})
}

// AnalyzeFailures analyzes specific failure patterns
func (h *AIHandler) AnalyzeFailures(c *fiber.Ctx) error {
	var req struct {
		Endpoint    string `json:"endpoint"`
		TimeRange   string `json:"time_range"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		StatusCodes []int  `json:"status_codes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body", 0, err.Error())
	}

	// Set defaults
	if req.TimeRange == "" {
		req.TimeRange = "24h"
	}

	// Fetch logs
	var logs []map[string]interface{}
	var err error

	if req.Endpoint != "" {
		// Fetch logs for specific endpoint
		logs, err = h.esClient.SearchLogsByEndpoint(context.Background(), req.Endpoint, req.TimeRange, 500)
	} else {
		// Fetch all logs
		logs, err = h.esClient.SearchLogs(context.Background(), req.TimeRange, 500)
	}

	if err != nil {
		h.logger.Error("Failed to fetch logs", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to fetch logs", 0, err.Error())
	}

	// Filter by status codes if specified
	if len(req.StatusCodes) > 0 {
		logs = filterLogsByStatusCodes(logs, req.StatusCodes)
	}

	if len(logs) == 0 {
		return JSONError(c, fiber.StatusNotFound, "No matching logs found", 0, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Analyze with AI
	var resp interface{}
	if req.Endpoint != "" {
		resp, err = h.vertexAI.AnalyzeFailuresByEndpoint(ctx, logs, req.Endpoint)
	} else {
		resp, err = h.vertexAI.AnalyzeAPILogs(ctx, logs, "Analyze these failures and provide insights")
	}

	if err != nil {
		h.logger.Error("Failed to analyze failures", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to analyze failures", 0, err.Error())
	}

	return JSONSuccess(c, "Failures analyzed successfully", fiber.Map{
		"analysis":   resp,
		"logs_count": len(logs),
		"endpoint":   req.Endpoint,
	})
}

// RootCauseAnalysis performs root cause analysis on error patterns
func (h *AIHandler) RootCauseAnalysis(c *fiber.Ctx) error {
	var req struct {
		ErrorPattern string `json:"error_pattern"`
		TimeRange    string `json:"time_range"`
		StartDate    string `json:"start_date"`
		EndDate      string `json:"end_date"`
		Limit        int    `json:"limit"`
	}

	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body", 0, err.Error())
	}

	if req.ErrorPattern == "" {
		return JSONError(c, fiber.StatusBadRequest, "Error pattern is required", 0, nil)
	}

	// Set defaults
	if req.TimeRange == "" {
		req.TimeRange = "1h"
	}
	if req.Limit == 0 {
		req.Limit = 200
	}

	// Search logs matching the error pattern
	logs, err := h.esClient.SearchLogsByPattern(context.Background(), req.ErrorPattern, req.TimeRange, req.Limit)
	if err != nil {
		h.logger.Error("Failed to search logs by pattern", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to search logs", 0, err.Error())
	}

	if len(logs) == 0 {
		return JSONError(c, fiber.StatusNotFound, "No logs found matching the error pattern", 0, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Perform root cause analysis
	resp, err := h.vertexAI.RootCauseAnalysis(ctx, logs, req.ErrorPattern)
	if err != nil {
		h.logger.Error("Failed to perform root cause analysis", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to perform root cause analysis", 0, err.Error())
	}

	return JSONSuccess(c, "Root cause analysis completed", fiber.Map{
		"analysis":      resp,
		"logs_analyzed": len(logs),
		"error_pattern": req.ErrorPattern,
	})
}

// AnalyzeTrafficFailures analyzes failures related to high traffic
func (h *AIHandler) AnalyzeTrafficFailures(c *fiber.Ctx) error {
	var req struct {
		TimeRange string  `json:"time_range"`
		StartDate string  `json:"start_date"`
		EndDate   string  `json:"end_date"`
		Limit     int     `json:"limit"`
		MinRate   float64 `json:"min_failure_rate"`
	}

	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body", 0, err.Error())
	}

	// Set defaults
	if req.TimeRange == "" {
		req.TimeRange = "1h"
	}
	if req.Limit == 0 {
		req.Limit = 500
	}
	if req.MinRate == 0 {
		req.MinRate = 0.1 // 10% failure rate
	}

	// Get failing APIs
	failingAPIs, err := h.esClient.SearchFailingAPIs(context.Background(), req.TimeRange, req.MinRate)
	if err != nil {
		h.logger.Error("Failed to search failing APIs", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to search failing APIs", 0, err.Error())
	}

	if len(failingAPIs) == 0 {
		return JSONSuccess(c, "No traffic-related failures detected", fiber.Map{
			"message": "All APIs are performing within normal parameters",
		})
	}

	// Fetch detailed logs for failing APIs
	var logs []map[string]interface{}

	if req.TimeRange == "custom" && req.StartDate != "" && req.EndDate != "" {
		logs, err = h.esClient.SearchLogsByDateRange(context.Background(), req.StartDate, req.EndDate, req.Limit)
	} else {
		logs, err = h.esClient.SearchLogs(context.Background(), req.TimeRange, req.Limit)
	}

	if err != nil {
		h.logger.Error("Failed to fetch logs", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to fetch logs", 0, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Analyze traffic failures with AI
	resp, err := h.vertexAI.AnalyzeTrafficFailures(ctx, logs)
	if err != nil {
		h.logger.Error("Failed to analyze traffic failures", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to analyze traffic failures", 0, err.Error())
	}

	return JSONSuccess(c, "Traffic failures analyzed successfully", fiber.Map{
		"analysis":     resp,
		"failing_apis": failingAPIs,
		"logs_count":   len(logs),
	})
}

// Helper function to filter logs by status codes
func filterLogsByStatusCodes(logs []map[string]interface{}, statusCodes []int) []map[string]interface{} {
	filtered := make([]map[string]interface{}, 0)

	for _, log := range logs {
		statusCode, ok := log["status_code"].(float64)
		if !ok {
			continue
		}

		for _, code := range statusCodes {
			if int(statusCode) == code {
				filtered = append(filtered, log)
				break
			}
		}
	}

	return filtered
}
