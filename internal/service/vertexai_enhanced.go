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
