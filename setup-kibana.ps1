# Optional: Setup Kibana for log visualization
# This is NOT required for the hackathon, but useful for debugging

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Kibana Setup (Optional)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "NOTE: Kibana is optional. Your backend API already provides" -ForegroundColor Yellow
Write-Host "all the functionality needed for the hackathon." -ForegroundColor Yellow
Write-Host ""

# Pull Kibana image
Write-Host "Pulling Kibana 8.11.0 image..." -ForegroundColor Yellow
docker pull docker.elastic.co/kibana/kibana:8.11.0

# Run Kibana container
Write-Host ""
Write-Host "Starting Kibana container..." -ForegroundColor Yellow
docker run -d `
    --name kibana-local `
    -p 5601:5601 `
    -e "ELASTICSEARCH_HOSTS=http://host.docker.internal:9200" `
    docker.elastic.co/kibana/kibana:8.11.0

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Kibana container started" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to start Kibana container" -ForegroundColor Red
    exit 1
}

# Wait for Kibana to be ready
Write-Host ""
Write-Host "Waiting for Kibana to be ready (this takes 1-2 minutes)..." -ForegroundColor Yellow
$maxAttempts = 60
$attempt = 0
$ready = $false

while ($attempt -lt $maxAttempts -and -not $ready) {
    Start-Sleep -Seconds 2
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:5601/api/status" -UseBasicParsing -ErrorAction SilentlyContinue
        if ($response.StatusCode -eq 200) {
            $ready = $true
        }
    } catch {
        # Continue waiting
    }
    $attempt++
    Write-Host "." -NoNewline -ForegroundColor Yellow
}

Write-Host ""
if ($ready) {
    Write-Host "✓ Kibana is ready!" -ForegroundColor Green
} else {
    Write-Host "⚠ Kibana is still starting. Check logs with: docker logs kibana-local" -ForegroundColor Yellow
}

# Summary
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Kibana Setup Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Kibana is running at: http://localhost:5601" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Open http://localhost:5601 in your browser" -ForegroundColor White
Write-Host "2. Go to 'Discover' to view logs" -ForegroundColor White
Write-Host "3. Create index pattern: api-logs-*" -ForegroundColor White
Write-Host ""
Write-Host "Useful commands:" -ForegroundColor Yellow
Write-Host "  Stop:    docker stop kibana-local" -ForegroundColor White
Write-Host "  Start:   docker start kibana-local" -ForegroundColor White
Write-Host "  Logs:    docker logs kibana-local" -ForegroundColor White
Write-Host "  Remove:  docker rm -f kibana-local" -ForegroundColor White
Write-Host ""
Write-Host "Remember: Kibana is optional for the hackathon!" -ForegroundColor Yellow
Write-Host "Your backend API already provides all needed functionality." -ForegroundColor Yellow
Write-Host ""
