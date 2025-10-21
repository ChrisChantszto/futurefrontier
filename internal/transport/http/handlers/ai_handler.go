package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ChrisChantszto/futurefrontier/internal/service"
	"go.uber.org/zap"
)

// AIHandler handles AI-powered endpoints
type AIHandler struct {
	vertexAI *service.VertexAIService
	esClient *service.ElasticsearchService
	logger   *zap.Logger
}

// NewAIHandler creates a new AI handler
func NewAIHandler(vertexAI *service.VertexAIService, esClient *service.ElasticsearchService, logger *zap.Logger) *AIHandler {
	return &AIHandler{
		vertexAI: vertexAI,
		esClient: esClient,
		logger:   logger,
	}
}

// GenerateContentRequest represents the request body for content generation
type GenerateContentRequest struct {
	Prompt      string   `json:"prompt" validate:"required"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	TopP        *float64 `json:"top_p,omitempty"`
	TopK        *int     `json:"top_k,omitempty"`
}

// AnalyzeLogsRequest represents the request body for log analysis
type AnalyzeLogsRequest struct {
	Query     string `json:"query"`
	TimeRange string `json:"time_range"` // e.g., "1h", "24h", "7d"
	Limit     int    `json:"limit"`
}

// GenerateContent handles AI content generation
func (h *AIHandler) GenerateContent(c *fiber.Ctx) error {
	var req GenerateContentRequest
	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body", 0, err.Error())
	}

	if req.Prompt == "" {
		return JSONError(c, fiber.StatusBadRequest, "Prompt is required", 0, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := h.vertexAI.GenerateContent(ctx, service.GenerateContentRequest{
		Prompt:      req.Prompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		TopK:        req.TopK,
	})
	if err != nil {
		h.logger.Error("Failed to generate content", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to generate content", 0, err.Error())
	}

	return JSONSuccess(c, "Content generated successfully", resp)
}

// AnalyzeLogs analyzes API logs using AI
func (h *AIHandler) AnalyzeLogs(c *fiber.Ctx) error {
	var req AnalyzeLogsRequest
	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body", 0, err.Error())
	}

	// Set defaults
	if req.TimeRange == "" {
		req.TimeRange = "1h"
	}
	if req.Limit == 0 {
		req.Limit = 100
	}

	// Fetch logs from Elasticsearch
	logs, err := h.esClient.SearchLogs(context.Background(), req.TimeRange, req.Limit)
	if err != nil {
		h.logger.Error("Failed to fetch logs from Elasticsearch", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to fetch logs", 0, err.Error())
	}

	if len(logs) == 0 {
		return JSONError(c, fiber.StatusNotFound, "No logs found for the specified time range", 0, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Analyze logs with AI
	resp, err := h.vertexAI.AnalyzeAPILogs(ctx, logs, req.Query)
	if err != nil {
		h.logger.Error("Failed to analyze logs", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to analyze logs", 0, err.Error())
	}

	return JSONSuccess(c, "Logs analyzed successfully", fiber.Map{
		"analysis":   resp,
		"logs_count": len(logs),
		"time_range": req.TimeRange,
	})
}

// DetectAnomalies detects anomalies in API logs
func (h *AIHandler) DetectAnomalies(c *fiber.Ctx) error {
	var req AnalyzeLogsRequest
	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body", 0, err.Error())
	}

	// Set defaults
	if req.TimeRange == "" {
		req.TimeRange = "1h"
	}
	if req.Limit == 0 {
		req.Limit = 200
	}

	// Fetch logs from Elasticsearch
	logs, err := h.esClient.SearchLogs(context.Background(), req.TimeRange, req.Limit)
	if err != nil {
		h.logger.Error("Failed to fetch logs from Elasticsearch", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to fetch logs", 0, err.Error())
	}

	if len(logs) == 0 {
		return JSONError(c, fiber.StatusNotFound, "No logs found for the specified time range", 0, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Detect anomalies with AI
	resp, err := h.vertexAI.DetectAnomalies(ctx, logs)
	if err != nil {
		h.logger.Error("Failed to detect anomalies", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to detect anomalies", 0, err.Error())
	}

	return JSONSuccess(c, "Anomalies detected successfully", fiber.Map{
		"anomalies":  resp,
		"logs_count": len(logs),
		"time_range": req.TimeRange,
	})
}

// SuggestOptimizations provides AI-powered optimization suggestions
func (h *AIHandler) SuggestOptimizations(c *fiber.Ctx) error {
	var req AnalyzeLogsRequest
	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body", 0, err.Error())
	}

	// Set defaults
	if req.TimeRange == "" {
		req.TimeRange = "24h"
	}
	if req.Limit == 0 {
		req.Limit = 500
	}

	// Fetch logs from Elasticsearch
	logs, err := h.esClient.SearchLogs(context.Background(), req.TimeRange, req.Limit)
	if err != nil {
		h.logger.Error("Failed to fetch logs from Elasticsearch", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to fetch logs", 0, err.Error())
	}

	if len(logs) == 0 {
		return JSONError(c, fiber.StatusNotFound, "No logs found for the specified time range", 0, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Get optimization suggestions with AI
	resp, err := h.vertexAI.SuggestOptimizations(ctx, logs)
	if err != nil {
		h.logger.Error("Failed to suggest optimizations", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to suggest optimizations", 0, err.Error())
	}

	return JSONSuccess(c, "Optimizations suggested successfully", fiber.Map{
		"suggestions": resp,
		"logs_count":  len(logs),
		"time_range":  req.TimeRange,
	})
}

// ChatWithLogs provides conversational interface to query logs
func (h *AIHandler) ChatWithLogs(c *fiber.Ctx) error {
	var req struct {
		Message   string `json:"message" validate:"required"`
		TimeRange string `json:"time_range"`
		Limit     int    `json:"limit"`
	}

	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body", 0, err.Error())
	}

	if req.Message == "" {
		return JSONError(c, fiber.StatusBadRequest, "Message is required", 0, nil)
	}

	// Set defaults
	if req.TimeRange == "" {
		req.TimeRange = "1h"
	}
	if req.Limit == 0 {
		req.Limit = 100
	}

	// Fetch logs from Elasticsearch
	logs, err := h.esClient.SearchLogs(context.Background(), req.TimeRange, req.Limit)
	if err != nil {
		h.logger.Error("Failed to fetch logs from Elasticsearch", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to fetch logs", 0, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Analyze logs with user's question
	resp, err := h.vertexAI.AnalyzeAPILogs(ctx, logs, req.Message)
	if err != nil {
		h.logger.Error("Failed to process chat message", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to process message", 0, err.Error())
	}

	return JSONSuccess(c, "Message processed successfully", fiber.Map{
		"response":   resp.Text,
		"logs_count": len(logs),
		"time_range": req.TimeRange,
	})
}
