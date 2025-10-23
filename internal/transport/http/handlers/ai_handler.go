package handlers

import (
	"context"
	"strings"
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

// ChatWithLogs provides conversational interface to query logs with RAG
func (h *AIHandler) ChatWithLogs(c *fiber.Ctx) error {
	var req struct {
		Message   string `json:"message" validate:"required"`
		TimeRange string `json:"time_range"`
		StartDate string `json:"start_date"` // For custom range
		EndDate   string `json:"end_date"`   // For custom range
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
		req.Limit = 200
	}

	// Check if this is a general question (not log-related)
	isGeneral := isGeneralQuestion(req.Message)

	if isGeneral {
		// Handle general questions without log context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		resp, err := h.vertexAI.GenerateContent(ctx, service.GenerateContentRequest{
			Prompt: req.Message,
		})
		if err != nil {
			h.logger.Error("Failed to generate response", zap.Error(err))
			return JSONError(c, fiber.StatusInternalServerError, "Failed to generate response", 0, err.Error())
		}

		return JSONSuccess(c, "Response generated successfully", fiber.Map{
			"response":     resp.Text,
			"context_used": false,
		})
	}

	// Fetch logs from Elasticsearch for context
	var logs []map[string]interface{}
	var err error

	if req.TimeRange == "custom" && req.StartDate != "" && req.EndDate != "" {
		// Use custom date range
		logs, err = h.esClient.SearchLogsByDateRange(context.Background(), req.StartDate, req.EndDate, req.Limit)
	} else {
		// Use predefined time range
		logs, err = h.esClient.SearchLogs(context.Background(), req.TimeRange, req.Limit)
	}

	if err != nil {
		h.logger.Error("Failed to fetch logs from Elasticsearch", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to fetch logs", 0, err.Error())
	}

	if len(logs) == 0 {
		return JSONSuccess(c, "No logs found", fiber.Map{
			"response":     "I couldn't find any logs in the specified time range. Please try a different time range or check if logs are being collected.",
			"logs_count":   0,
			"context_used": false,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Use RAG to generate context-aware response
	resp, err := h.vertexAI.ChatWithContext(ctx, req.Message, logs)
	if err != nil {
		h.logger.Error("Failed to process chat message with context", zap.Error(err))
		return JSONError(c, fiber.StatusInternalServerError, "Failed to analyze logs", 0, err.Error())
	}

	return JSONSuccess(c, "Analysis completed successfully", fiber.Map{
		"response":     resp.Text,
		"logs_count":   len(logs),
		"context_used": true,
		"time_range":   req.TimeRange,
	})
}

// isGeneralQuestion checks if the user's message is a general question
func isGeneralQuestion(message string) bool {
	generalPhrases := []string{
		"hello", "hi", "hey", "how are you", "what can you do",
		"help", "who are you", "what are you", "introduce yourself",
		"good morning", "good afternoon", "good evening",
		"thanks", "thank you", "bye", "goodbye",
	}

	lowerMsg := strings.ToLower(message)
	for _, phrase := range generalPhrases {
		if strings.Contains(lowerMsg, phrase) {
			return true
		}
	}

	// If message is very short and doesn't mention logs/errors/api, it's likely general
	if len(message) < 20 && !strings.Contains(lowerMsg, "log") &&
		!strings.Contains(lowerMsg, "error") && !strings.Contains(lowerMsg, "api") &&
		!strings.Contains(lowerMsg, "fail") && !strings.Contains(lowerMsg, "issue") {
		return true
	}

	return false
}
