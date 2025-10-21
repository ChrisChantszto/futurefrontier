package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ChrisChantszto/futurefrontier/internal/config"
	"go.uber.org/zap"
	"google.golang.org/genai"
)

// VertexAIService handles interactions with Google Cloud Vertex AI (Generative API)
type VertexAIService struct {
	config config.GoogleCloudConfig
	logger *zap.Logger
	client *genai.Client
	model  string
}

// NewVertexAIService creates a new Vertex AI service instance
func NewVertexAIService(cfg config.GoogleCloudConfig, logger *zap.Logger) (*VertexAIService, error) {
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT_ID is required")
	}

	ctx := context.Background()

	// Create Generative AI client using Vertex AI backend
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  cfg.ProjectID,
		Location: cfg.Location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Vertex AI genai client: %w", err)
	}

	logger.Info("Vertex AI service initialized",
		zap.String("project_id", cfg.ProjectID),
		zap.String("location", cfg.Location),
		zap.String("model", cfg.VertexAIModel))

	return &VertexAIService{
		config: cfg,
		logger: logger,
		client: client,
		model:  cfg.VertexAIModel,
	}, nil
}

// GenerateContentRequest represents a request to generate content
type GenerateContentRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   *int    `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	TopP        *float64 `json:"top_p,omitempty"`
	TopK        *int    `json:"top_k,omitempty"`
}

// GenerateContentResponse represents the response from content generation
type GenerateContentResponse struct {
	Text            string                 `json:"text"`
	FinishReason    string                 `json:"finish_reason,omitempty"`
	SafetyRatings   []map[string]interface{} `json:"safety_ratings,omitempty"`
	TokenCount      int                    `json:"token_count,omitempty"`
}

// GenerateContent generates text content using Vertex AI
func (s *VertexAIService) GenerateContent(ctx context.Context, req GenerateContentRequest) (*GenerateContentResponse, error) {
	if req.Prompt == "" {
		return nil, fmt.Errorf("prompt is required")
	}
	// Call Generative API
	parts := []*genai.Part{{Text: req.Prompt}}
	contents := []*genai.Content{{Role: "user", Parts: parts}}
	gcResp, err := s.client.Models.GenerateContent(ctx, s.model, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("prediction failed: %w", err)
	}

	// Extract text from candidates
	var b strings.Builder
	if gcResp != nil && len(gcResp.Candidates) > 0 {
		for _, cand := range gcResp.Candidates {
			if cand == nil || cand.Content == nil {
				continue
			}
			for _, p := range cand.Content.Parts {
				if p != nil && p.Text != "" {
					b.WriteString(p.Text)
				}
			}
		}
	}

	text := b.String()
	if text == "" {
		return nil, fmt.Errorf("failed to extract text from response")
	}

	response := &GenerateContentResponse{Text: text}
	s.logger.Debug("Received response from Vertex AI",
		zap.Int("text_length", len(text)))
	return response, nil
}

// AnalyzeAPILogs analyzes API logs using Vertex AI to provide insights
func (s *VertexAIService) AnalyzeAPILogs(ctx context.Context, logs []map[string]interface{}, query string) (*GenerateContentResponse, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs provided")
	}

	// Convert logs to JSON for analysis
	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}

	// Build analysis prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString("You are an expert API monitoring and analytics assistant. ")
	promptBuilder.WriteString("Analyze the following API logs and provide insights.\n\n")
	
	if query != "" {
		promptBuilder.WriteString(fmt.Sprintf("User Question: %s\n\n", query))
	} else {
		promptBuilder.WriteString("Provide a comprehensive analysis including:\n")
		promptBuilder.WriteString("1. Summary of API usage patterns\n")
		promptBuilder.WriteString("2. Any anomalies or errors detected\n")
		promptBuilder.WriteString("3. Performance insights\n")
		promptBuilder.WriteString("4. Recommendations for optimization\n\n")
	}

	promptBuilder.WriteString("API Logs:\n")
	promptBuilder.WriteString(string(logsJSON))

	return s.GenerateContent(ctx, GenerateContentRequest{
		Prompt: promptBuilder.String(),
	})
}

// DetectAnomalies uses AI to detect anomalies in API behavior
func (s *VertexAIService) DetectAnomalies(ctx context.Context, logs []map[string]interface{}) (*GenerateContentResponse, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs provided")
	}

	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}

	prompt := fmt.Sprintf(`You are an expert in API security and anomaly detection. 
Analyze the following API logs and identify any anomalies, suspicious patterns, or security concerns.

Focus on:
1. Unusual request patterns or frequencies
2. Error rate spikes
3. Suspicious IP addresses or user agents
4. Potential security threats
5. Performance degradation patterns

Provide a structured analysis with severity levels (LOW, MEDIUM, HIGH, CRITICAL) for each finding.

API Logs:
%s`, string(logsJSON))

	return s.GenerateContent(ctx, GenerateContentRequest{
		Prompt: prompt,
	})
}

// SuggestOptimizations provides AI-powered optimization suggestions
func (s *VertexAIService) SuggestOptimizations(ctx context.Context, logs []map[string]interface{}) (*GenerateContentResponse, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs provided")
	}

	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}

	prompt := fmt.Sprintf(`You are an expert backend performance engineer. 
Analyze the following API logs and provide actionable optimization recommendations.

Consider:
1. Response time improvements
2. Caching opportunities
3. Database query optimization
4. API endpoint consolidation
5. Rate limiting strategies
6. Error handling improvements

Provide specific, actionable recommendations with expected impact.

API Logs:
%s`, string(logsJSON))

	return s.GenerateContent(ctx, GenerateContentRequest{
		Prompt: prompt,
	})
}

// Close closes the Vertex AI client
func (s *VertexAIService) Close() error {
    return nil
}
