# Setup script for local Elasticsearch
# Run this script to set up Elasticsearch for the AI-powered log analysis system

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Elasticsearch Setup for FutureFrontier" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Docker is running
Write-Host "Checking Docker..." -ForegroundColor Yellow
try {
    docker info | Out-Null
    Write-Host "✓ Docker is running" -ForegroundColor Green
} catch {
    Write-Host "✗ Docker is not running. Please start Docker Desktop first." -ForegroundColor Red
    exit 1
}

# Check if Elasticsearch container already exists
Write-Host ""
Write-Host "Checking for existing Elasticsearch container..." -ForegroundColor Yellow
$existingContainer = docker ps -a --filter "name=elasticsearch-local" --format "{{.Names}}"

if ($existingContainer -eq "elasticsearch-local") {
    Write-Host "Found existing container. Removing it..." -ForegroundColor Yellow
    docker rm -f elasticsearch-local | Out-Null
    Write-Host "✓ Removed existing container" -ForegroundColor Green
}

# Pull Elasticsearch image
Write-Host ""
Write-Host "Pulling Elasticsearch 8.11.0 image..." -ForegroundColor Yellow
docker pull docker.elastic.co/elasticsearch/elasticsearch:8.11.0

# Run Elasticsearch container
Write-Host ""
Write-Host "Starting Elasticsearch container..." -ForegroundColor Yellow
docker run -d `
    --name elasticsearch-local `
    -p 9200:9200 `
    -p 9300:9300 `
    -e "discovery.type=single-node" `
    -e "xpack.security.enabled=false" `
    -e "xpack.security.http.ssl.enabled=false" `
    -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" `
    docker.elastic.co/elasticsearch/elasticsearch:8.11.0

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Elasticsearch container started" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to start Elasticsearch container" -ForegroundColor Red
    exit 1
}

# Wait for Elasticsearch to be ready
Write-Host ""
Write-Host "Waiting for Elasticsearch to be ready..." -ForegroundColor Yellow
$maxAttempts = 30
$attempt = 0
$ready = $false

while ($attempt -lt $maxAttempts -and -not $ready) {
    Start-Sleep -Seconds 2
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:9200" -UseBasicParsing -ErrorAction SilentlyContinue
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
    Write-Host "✓ Elasticsearch is ready!" -ForegroundColor Green
} else {
    Write-Host "✗ Elasticsearch did not start in time" -ForegroundColor Red
    Write-Host "Check logs with: docker logs elasticsearch-local" -ForegroundColor Yellow
    exit 1
}

# Verify Elasticsearch
Write-Host ""
Write-Host "Verifying Elasticsearch..." -ForegroundColor Yellow
try {
    $info = Invoke-RestMethod -Uri "http://localhost:9200" -Method Get
    Write-Host "✓ Elasticsearch is running" -ForegroundColor Green
    Write-Host "  Cluster: $($info.cluster_name)" -ForegroundColor Cyan
    Write-Host "  Version: $($info.version.number)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed to verify Elasticsearch" -ForegroundColor Red
    exit 1
}

# Summary
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Setup Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Elasticsearch is running at: http://localhost:9200" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Start MongoDB: docker start mongo-local" -ForegroundColor White
Write-Host "2. Run backend: go run ." -ForegroundColor White
Write-Host "3. Generate demo data: Invoke-WebRequest -Uri 'http://localhost:8080/api/demo/generate?count=1000' -Method POST" -ForegroundColor White
Write-Host ""
Write-Host "Useful commands:" -ForegroundColor Yellow
Write-Host "  Stop:    docker stop elasticsearch-local" -ForegroundColor White
Write-Host "  Start:   docker start elasticsearch-local" -ForegroundColor White
Write-Host "  Logs:    docker logs elasticsearch-local" -ForegroundColor White
Write-Host "  Remove:  docker rm -f elasticsearch-local" -ForegroundColor White
Write-Host ""
