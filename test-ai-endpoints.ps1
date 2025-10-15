# PowerShell script to test AI endpoints
# Usage: .\test-ai-endpoints.ps1

$API_BASE = "http://localhost:8080/api"

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "AI Endpoints Test Script" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# Step 1: Login
Write-Host "Step 1: Logging in..." -ForegroundColor Yellow
$loginBody = @{
    email = "admin@example.com"
    password = "password"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$API_BASE/auth/login" `
        -Method Post `
        -ContentType "application/json" `
        -Body $loginBody
    
    $token = $loginResponse.data.access_token
    Write-Host "✓ Login successful!" -ForegroundColor Green
    Write-Host "Token: $($token.Substring(0, 20))..." -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Please update the email/password in this script" -ForegroundColor Yellow
    exit 1
}

# Headers for authenticated requests
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# Step 2: Test Generate Content
Write-Host "Step 2: Testing AI Content Generation..." -ForegroundColor Yellow
$generateBody = @{
    prompt = "Hello, Vertex AI! Please confirm you're working."
} | ConvertTo-Json

try {
    $generateResponse = Invoke-RestMethod -Uri "$API_BASE/ai/generate" `
        -Method Post `
        -Headers $headers `
        -Body $generateBody
    
    Write-Host "✓ AI Generation successful!" -ForegroundColor Green
    Write-Host "Response: $($generateResponse.data.text.Substring(0, 100))..." -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ AI Generation failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Error details: $($_.ErrorDetails.Message)" -ForegroundColor Red
    Write-Host ""
}

# Step 3: Test Analyze Logs
Write-Host "Step 3: Testing Log Analysis..." -ForegroundColor Yellow
$analyzeBody = @{
    query = "Provide a summary of API usage patterns"
    time_range = "1h"
    limit = 50
} | ConvertTo-Json

try {
    $analyzeResponse = Invoke-RestMethod -Uri "$API_BASE/ai/analyze-logs" `
        -Method Post `
        -Headers $headers `
        -Body $analyzeBody
    
    Write-Host "✓ Log Analysis successful!" -ForegroundColor Green
    Write-Host "Logs analyzed: $($analyzeResponse.data.logs_count)" -ForegroundColor Gray
    if ($analyzeResponse.data.analysis.text) {
        Write-Host "Analysis: $($analyzeResponse.data.analysis.text.Substring(0, 150))..." -ForegroundColor Gray
    }
    Write-Host ""
} catch {
    Write-Host "✗ Log Analysis failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message -match "No logs found") {
        Write-Host "Note: This is expected if you don't have logs in Elasticsearch yet" -ForegroundColor Yellow
    }
    Write-Host ""
}

# Step 4: Test Anomaly Detection
Write-Host "Step 4: Testing Anomaly Detection..." -ForegroundColor Yellow
$anomalyBody = @{
    time_range = "1h"
    limit = 100
} | ConvertTo-Json

try {
    $anomalyResponse = Invoke-RestMethod -Uri "$API_BASE/ai/detect-anomalies" `
        -Method Post `
        -Headers $headers `
        -Body $anomalyBody
    
    Write-Host "✓ Anomaly Detection successful!" -ForegroundColor Green
    Write-Host "Logs analyzed: $($anomalyResponse.data.logs_count)" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Anomaly Detection failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message -match "No logs found") {
        Write-Host "Note: This is expected if you don't have logs in Elasticsearch yet" -ForegroundColor Yellow
    }
    Write-Host ""
}

# Step 5: Test Optimization Suggestions
Write-Host "Step 5: Testing Optimization Suggestions..." -ForegroundColor Yellow
$optimizeBody = @{
    time_range = "24h"
    limit = 200
} | ConvertTo-Json

try {
    $optimizeResponse = Invoke-RestMethod -Uri "$API_BASE/ai/suggest-optimizations" `
        -Method Post `
        -Headers $headers `
        -Body $optimizeBody
    
    Write-Host "✓ Optimization Suggestions successful!" -ForegroundColor Green
    Write-Host "Logs analyzed: $($optimizeResponse.data.logs_count)" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Optimization Suggestions failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message -match "No logs found") {
        Write-Host "Note: This is expected if you don't have logs in Elasticsearch yet" -ForegroundColor Yellow
    }
    Write-Host ""
}

# Step 6: Test Chat Interface
Write-Host "Step 6: Testing Chat Interface..." -ForegroundColor Yellow
$chatBody = @{
    message = "How many API requests were made in the last hour?"
    time_range = "1h"
    limit = 100
} | ConvertTo-Json

try {
    $chatResponse = Invoke-RestMethod -Uri "$API_BASE/ai/chat" `
        -Method Post `
        -Headers $headers `
        -Body $chatBody
    
    Write-Host "✓ Chat Interface successful!" -ForegroundColor Green
    Write-Host "Response: $($chatResponse.data.response.Substring(0, 150))..." -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Chat Interface failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message -match "No logs found") {
        Write-Host "Note: This is expected if you don't have logs in Elasticsearch yet" -ForegroundColor Yellow
    }
    Write-Host ""
}

# Summary
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "All tests completed!" -ForegroundColor Green
Write-Host ""
Write-Host "Next Steps:" -ForegroundColor Yellow
Write-Host "1. Ensure Elasticsearch has some logs (make API calls to generate logs)" -ForegroundColor White
Write-Host "2. Try the endpoints with different queries" -ForegroundColor White
Write-Host "3. Check the server logs for detailed information" -ForegroundColor White
Write-Host ""
Write-Host "For more information, see:" -ForegroundColor Yellow
Write-Host "- VERTEX_AI_SETUP.md" -ForegroundColor White
Write-Host "- HACKATHON_QUICKSTART.md" -ForegroundColor White
Write-Host ""
