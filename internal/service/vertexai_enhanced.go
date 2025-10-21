package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// DetectPatterns analyzes logs to find recurring failure patterns
func (s *VertexAIService) DetectPatterns(ctx context.Context, logs []map[string]interface{}) (*GenerateContentResponse, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs provided")
	}

	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}

	prompt := fmt.Sprintf(`You are an expert in API monitoring and pattern recognition.

Analyze these API logs and identify RECURRING FAILURE PATTERNS:

Instructions:
1. Group similar errors together
2. Identify which endpoints fail consistently
3. Look for time-based patterns (e.g., failures during peak hours)
4. Find correlations between failures (e.g., one failure causing others)
5. Determine if failures are traffic-related

For each pattern found, provide:
- Pattern description
- Affected endpoints
- Frequency and timing
- Likely root cause
- Severity level (LOW/MEDIUM/HIGH/CRITICAL)
- Recommended actions

API Logs:
%s

Format your response as a structured analysis with clear sections.`, string(logsJSON))

	return s.GenerateContent(ctx, GenerateContentRequest{
		Prompt: prompt,
	})
}

// AnalyzeTrafficFailures analyzes failures related to high traffic
func (s *VertexAIService) AnalyzeTrafficFailures(ctx context.Context, logs []map[string]interface{}) (*GenerateContentResponse, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs provided")
	}

	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}

	prompt := fmt.Sprintf(`You are an expert in API scalability and performance under load.

Analyze these API logs to identify TRAFFIC-RELATED FAILURES:

Focus on:
1. Endpoints with high request volume AND high error rates
2. Timeouts and connection errors during peak traffic
3. Database connection pool exhaustion
4. Memory or CPU constraints
5. Rate limiting issues

For each traffic-related issue:
- Endpoint and error type
- Request volume when failures occur
- Error rate percentage
- Performance degradation pattern
- Root cause (e.g., insufficient resources, poor scaling)
- Scaling recommendations (horizontal/vertical)
- Immediate mitigation steps

API Logs:
%s

Provide actionable recommendations for handling high traffic.`, string(logsJSON))

	return s.GenerateContent(ctx, GenerateContentRequest{
		Prompt: prompt,
	})
}

// RootCauseAnalysis performs deep analysis to identify root causes
func (s *VertexAIService) RootCauseAnalysis(ctx context.Context, logs []map[string]interface{}, errorPattern string) (*GenerateContentResponse, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs provided")
	}

	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}

	prompt := fmt.Sprintf(`You are an expert software engineer performing ROOT CAUSE ANALYSIS.

Error Pattern: %s

Analyze these logs to determine the ROOT CAUSE:

Investigation steps:
1. Identify the primary error
2. Trace back to the originating issue
3. Distinguish between symptoms and root causes
4. Check for cascading failures
5. Consider external dependencies

Provide:
- Primary root cause
- Contributing factors
- Evidence from logs
- Why this is the root cause (not just a symptom)
- Specific code/infrastructure changes needed
- Prevention strategies

API Logs:
%s

Be specific and technical in your analysis.`, errorPattern, string(logsJSON))

	return s.GenerateContent(ctx, GenerateContentRequest{
		Prompt: prompt,
	})
}

// AnalyzeFailuresByEndpoint analyzes failures for a specific endpoint
func (s *VertexAIService) AnalyzeFailuresByEndpoint(ctx context.Context, logs []map[string]interface{}, endpoint string) (*GenerateContentResponse, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs provided")
	}

	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}

	var promptBuilder strings.Builder
	promptBuilder.WriteString("You are an expert API debugging specialist.\n\n")
	promptBuilder.WriteString(fmt.Sprintf("Analyze failures for endpoint: %s\n\n", endpoint))
	promptBuilder.WriteString("Provide:\n")
	promptBuilder.WriteString("1. Summary of error types and frequencies\n")
	promptBuilder.WriteString("2. Common error messages and their meanings\n")
	promptBuilder.WriteString("3. Performance characteristics (latency patterns)\n")
	promptBuilder.WriteString("4. Potential root causes\n")
	promptBuilder.WriteString("5. Specific fixes and optimizations\n\n")
	promptBuilder.WriteString("API Logs:\n")
	promptBuilder.WriteString(string(logsJSON))

	return s.GenerateContent(ctx, GenerateContentRequest{
		Prompt: promptBuilder.String(),
	})
}

// ChatWithContext implements RAG (Retrieval-Augmented Generation) for context-aware chat
func (s *VertexAIService) ChatWithContext(ctx context.Context, userMessage string, logs []map[string]interface{}) (*GenerateContentResponse, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs provided for context")
	}

	// Build context from logs
	context := s.buildLogContext(logs)

	prompt := fmt.Sprintf(`You are an AI assistant specializing in API log analysis and troubleshooting.

IMPORTANT: You have access to ACTUAL log data below. Use this specific data in your response. Reference actual endpoints, error messages, status codes, and timestamps from the data.

=== CONTEXT: RECENT LOG DATA ===
%s
================================

USER QUESTION: %s

Instructions:
1. Analyze the actual log data provided above
2. Reference specific endpoints, errors, and patterns you see in the data
3. Provide concrete, data-driven insights
4. Include specific examples from the logs
5. Give actionable recommendations based on the real data

Be specific and reference actual data points. Don't give generic responses.`, context, userMessage)

	return s.GenerateContent(ctx, GenerateContentRequest{
		Prompt: prompt,
	})
}

// buildLogContext creates a structured summary of logs for RAG
func (s *VertexAIService) buildLogContext(logs []map[string]interface{}) string {
	var builder strings.Builder

	// Statistics
	totalLogs := len(logs)
	errorCount := 0
	endpoints := make(map[string]int)
	statusCodes := make(map[float64]int)
	errors := make(map[string]int)
	latencies := []float64{}

	// Analyze logs
	for _, log := range logs {
		// Count errors
		if statusCode, ok := log["status_code"].(float64); ok {
			statusCodes[statusCode]++
			if statusCode >= 400 {
				errorCount++
			}
		}

		// Track endpoints
		if path, ok := log["path"].(string); ok {
			endpoints[path]++
		}

		// Track error messages
		if errMsg, ok := log["error_message"].(string); ok && errMsg != "" {
			errors[errMsg]++
		}

		// Track latencies
		if latency, ok := log["latency_ms"].(float64); ok {
			latencies = append(latencies, latency)
		}
	}

	// Calculate average latency
	avgLatency := 0.0
	if len(latencies) > 0 {
		sum := 0.0
		for _, l := range latencies {
			sum += l
		}
		avgLatency = sum / float64(len(latencies))
	}

	// Build context string
	builder.WriteString(fmt.Sprintf("Total Logs Analyzed: %d\n", totalLogs))
	builder.WriteString(fmt.Sprintf("Error Count: %d (%.1f%% error rate)\n", errorCount, float64(errorCount)/float64(totalLogs)*100))
	builder.WriteString(fmt.Sprintf("Average Latency: %.0fms\n\n", avgLatency))

	// Top failing endpoints
	builder.WriteString("Top Endpoints by Request Count:\n")
	count := 0
	for endpoint, reqCount := range endpoints {
		if count >= 10 {
			break
		}
		builder.WriteString(fmt.Sprintf("  - %s: %d requests\n", endpoint, reqCount))
		count++
	}

	// Status code breakdown
	builder.WriteString("\nStatus Code Distribution:\n")
	for code, count := range statusCodes {
		builder.WriteString(fmt.Sprintf("  - %d: %d requests (%.1f%%)\n", int(code), count, float64(count)/float64(totalLogs)*100))
	}

	// Common errors
	if len(errors) > 0 {
		builder.WriteString("\nCommon Error Messages:\n")
		count = 0
		for errMsg, errCount := range errors {
			if count >= 5 {
				break
			}
			builder.WriteString(fmt.Sprintf("  - \"%s\": %d occurrences\n", errMsg, errCount))
			count++
		}
	}

	// Sample recent logs
	builder.WriteString("\nRecent Log Samples:\n")
	sampleCount := 3
	if len(logs) < sampleCount {
		sampleCount = len(logs)
	}
	for i := 0; i < sampleCount; i++ {
		log := logs[i]
		builder.WriteString(fmt.Sprintf("\nLog %d:\n", i+1))
		if timestamp, ok := log["timestamp"].(string); ok {
			builder.WriteString(fmt.Sprintf("  Time: %s\n", timestamp))
		}
		if method, ok := log["method"].(string); ok {
			builder.WriteString(fmt.Sprintf("  Method: %s\n", method))
		}
		if path, ok := log["path"].(string); ok {
			builder.WriteString(fmt.Sprintf("  Path: %s\n", path))
		}
		if statusCode, ok := log["status_code"].(float64); ok {
			builder.WriteString(fmt.Sprintf("  Status: %d\n", int(statusCode)))
		}
		if latency, ok := log["latency_ms"].(float64); ok {
			builder.WriteString(fmt.Sprintf("  Latency: %.0fms\n", latency))
		}
		if errMsg, ok := log["error_message"].(string); ok && errMsg != "" {
			builder.WriteString(fmt.Sprintf("  Error: %s\n", errMsg))
		}
	}

	return builder.String()
}
