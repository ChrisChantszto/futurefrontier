# Railway Deployment Guide

## Issues Fixed

### 1. Vertex AI Model Error (404)
**Problem:** `gemini-1.5-flash` model not found
**Solution:** Updated to `gemini-2.0-flash-exp`

### 2. Elasticsearch 404 Error
**Problem:** `/api-logs-*/_search` returns 404 because indexes don't exist
**Solution:** Indexes are created automatically on first API log, but you need to:
- Either wait for natural traffic to create indexes
- Or manually generate demo data after deployment

## Required Railway Environment Variables

### Core Application
```bash
PORT=8080
MONGODB_URI=<your-mongodb-connection-string>
DB_NAME=futurefrontier
JWT_SECRET=<your-jwt-secret>
JWT_REFRESH_SECRET=<your-jwt-refresh-secret>
CORS_ORIGINS=https://your-frontend-domain.com
```

### Elasticsearch (Elastic Cloud)
```bash
ELASTICSEARCH_URL=https://your-cluster.es.us-central1.gcp.cloud.es.io:443
ELASTICSEARCH_API_KEY=<your-elastic-cloud-api-key>
LOG_PROJECT_ID=futurefrontier
LOG_LOCAL_MODE=false
```

### Google Cloud Vertex AI
```bash
GCP_PROJECT_ID=grand-principle-475206-b5
GCP_LOCATION=us-central1
VERTEX_AI_MODEL=gemini-2.0-flash-exp
GOOGLE_APPLICATION_CREDENTIALS=<base64-encoded-service-account-json>
```

**Important:** For Railway, you need to base64 encode your service account JSON:
```bash
# On your local machine
base64 -i service-account-key.json
# Copy the output and set it as GOOGLE_APPLICATION_CREDENTIALS
```

Then in your code, decode it at runtime (already handled in the codebase).

## Post-Deployment Setup

### Step 1: Verify Deployment
```bash
curl https://futurefrontier-production.up.railway.app/api/healthz
```

Expected response:
```json
{"ok": true}
```

### Step 2: Authenticate
Login to get an auth token:
```bash
curl -X POST https://futurefrontier-production.up.railway.app/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-email@example.com",
    "password": "your-password"
  }'
```

Save the `access_token` from the response.

### Step 3: Generate Demo Data (Creates Elasticsearch Indexes)
```bash
curl -X POST "https://futurefrontier-production.up.railway.app/api/demo/generate?count=1000" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

This will:
- Create the `api-logs-*` indexes in Elasticsearch
- Generate 1000 sample log entries with realistic failure patterns
- Enable all AI endpoints to work

### Step 4: Test AI Chat Endpoint
```bash
curl -X POST https://futurefrontier-production.up.railway.app/api/ai/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "message": "Show me failing APIs",
    "time_range": "1h",
    "limit": 100
  }'
```

## Available AI Endpoints

All endpoints require authentication (`Authorization: Bearer <token>`):

### 1. Chat with Logs (RAG)
```bash
POST /api/ai/chat
{
  "message": "Show me failing APIs",
  "time_range": "1h",
  "limit": 100
}
```

### 2. Detect Patterns
```bash
POST /api/ai/detect-patterns
{
  "time_range": "24h",
  "limit": 1000
}
```

### 3. Analyze Failures
```bash
POST /api/ai/analyze-failures
{
  "endpoint": "/api/users/list",
  "time_range": "24h"
}
```

### 4. Root Cause Analysis
```bash
POST /api/ai/root-cause-analysis
{
  "error_pattern": "database timeout",
  "time_range": "1h",
  "limit": 200
}
```

### 5. Analyze Traffic Failures
```bash
POST /api/ai/analyze-traffic-failures
{
  "time_range": "1h",
  "min_failure_rate": 0.1
}
```

## Troubleshooting

### Issue: "No logs found"
**Cause:** Elasticsearch indexes don't exist yet
**Solution:** Run the demo data generation endpoint (Step 3 above)

### Issue: "Publisher Model not found"
**Cause:** Wrong Vertex AI model name or missing GCP credentials
**Solution:** 
- Verify `VERTEX_AI_MODEL=gemini-2.0-flash-exp`
- Verify `GCP_PROJECT_ID` is correct
- Verify `GOOGLE_APPLICATION_CREDENTIALS` is properly base64 encoded

### Issue: "Elasticsearch connection error"
**Cause:** Invalid Elasticsearch credentials or URL
**Solution:**
- Verify `ELASTICSEARCH_URL` includes `https://` and port `:443`
- Verify `ELASTICSEARCH_API_KEY` is valid
- Test connection from Railway logs

### Issue: "Authentication required"
**Cause:** Missing or invalid JWT token
**Solution:** Login first and include `Authorization: Bearer <token>` header

## Monitoring

Check Railway logs for:
- `AI routes registered successfully` - Confirms Vertex AI is initialized
- `Successfully connected to Elasticsearch` - Confirms ES connection
- `Created elasticsearch index` - Confirms indexes are being created
- `Retrieved logs from Elasticsearch` - Confirms log queries work

## Demo Data Patterns

The demo data generator creates realistic patterns:
- `/api/users/list`: 30% failure rate (database timeout)
- `/api/auth/login`: 15% authentication failures
- Traffic spike in the middle causing 503 errors
- Various latency patterns

This allows you to test all AI analysis features immediately.
