# Vertex AI Integration Guide

This guide will help you set up and use Google Cloud Vertex AI with your Go Fiber backend for the AI Accelerate Hackathon.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Google Cloud Setup](#google-cloud-setup)
- [Backend Configuration](#backend-configuration)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [Hackathon Strategy](#hackathon-strategy)

## Prerequisites

- Google Cloud Platform account
- Vertex AI API enabled in your GCP project
- Service account with Vertex AI permissions
- Elasticsearch running with API logs

## Google Cloud Setup

### 1. Create a GCP Project

```bash
# Install gcloud CLI if not already installed
# Visit: https://cloud.google.com/sdk/docs/install

# Login to Google Cloud
gcloud auth login

# Create a new project (or use existing)
gcloud projects create your-project-id --name="AI Accelerate Hackathon"

# Set the project as default
gcloud config set project your-project-id
```

### 2. Enable Vertex AI API

```bash
# Enable Vertex AI API
gcloud services enable aiplatform.googleapis.com

# Verify the API is enabled
gcloud services list --enabled | grep aiplatform
```

### 3. Create Service Account

```bash
# Create service account
gcloud iam service-accounts create vertex-ai-service \
    --display-name="Vertex AI Service Account"

# Grant Vertex AI User role
gcloud projects add-iam-policy-binding your-project-id \
    --member="serviceAccount:vertex-ai-service@your-project-id.iam.gserviceaccount.com" \
    --role="roles/aiplatform.user"

# Create and download key file
gcloud iam service-accounts keys create gcp-service-account-key.json \
    --iam-account=vertex-ai-service@your-project-id.iam.gserviceaccount.com
```

### 4. Move Key File to Project

```bash
# Move the key file to your backend project root
mv gcp-service-account-key.json /path/to/onetake-corpsite-backend/
```

## Backend Configuration

### 1. Update .env File

Edit your `.env` file with your Google Cloud credentials:

```env
# Google Cloud Configuration
GCP_PROJECT_ID=your-project-id
GCP_LOCATION=us-central1
GOOGLE_APPLICATION_CREDENTIALS=./gcp-service-account-key.json
VERTEX_AI_MODEL=gemini-1.5-flash
VERTEX_AI_ENDPOINT=
VERTEX_AI_MAX_TOKENS=8192
VERTEX_AI_TEMPERATURE=0.7
VERTEX_AI_TOP_P=0.95
VERTEX_AI_TOP_K=40
```

**Important Configuration Notes:**
- `GCP_PROJECT_ID`: Your actual GCP project ID
- `GCP_LOCATION`: Choose a location close to you:
  - `us-central1` (Iowa, USA)
  - `us-east4` (Virginia, USA)
  - `asia-southeast1` (Singapore)
  - `europe-west4` (Netherlands)
- `VERTEX_AI_MODEL`: Available models:
  - `gemini-1.5-pro` (More powerful, slower)
  - `gemini-1.5-flash` (Faster, cost-effective - recommended)
  - `gemini-1.0-pro`

### 2. Verify Installation

```bash
# Build the project
go build

# Run the server
./onetake-backend.exe
```

You should see in the logs:
```
Vertex AI service initialized
AI routes registered successfully
```

## API Endpoints

All AI endpoints are protected and require authentication. Base path: `/api/ai`

### 1. Generate Content (General AI)

**Endpoint:** `POST /api/ai/generate`

Generate AI content with custom prompts.

**Request:**
```json
{
  "prompt": "Explain how API rate limiting works",
  "max_tokens": 1000,
  "temperature": 0.7
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "text": "API rate limiting is a technique...",
    "finish_reason": "STOP",
    "token_count": 245
  }
}
```

### 2. Analyze API Logs

**Endpoint:** `POST /api/ai/analyze-logs`

Analyze API logs with AI to get insights.

**Request:**
```json
{
  "query": "What are the most common errors in the last hour?",
  "time_range": "1h",
  "limit": 100
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "analysis": {
      "text": "Based on the logs, the most common errors are...",
      "finish_reason": "STOP"
    },
    "logs_count": 87,
    "time_range": "1h"
  }
}
```

### 3. Detect Anomalies

**Endpoint:** `POST /api/ai/detect-anomalies`

Use AI to detect anomalies and security threats in API logs.

**Request:**
```json
{
  "time_range": "24h",
  "limit": 200
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "anomalies": {
      "text": "CRITICAL: Detected unusual spike in 401 errors from IP 192.168.1.100...",
      "safety_ratings": []
    },
    "logs_count": 156,
    "time_range": "24h"
  }
}
```

### 4. Suggest Optimizations

**Endpoint:** `POST /api/ai/suggest-optimizations`

Get AI-powered optimization recommendations.

**Request:**
```json
{
  "time_range": "7d",
  "limit": 500
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "suggestions": {
      "text": "1. Implement caching for /api/pages endpoint (avg 450ms response time)...",
      "finish_reason": "STOP"
    },
    "logs_count": 432,
    "time_range": "7d"
  }
}
```

### 5. Chat with Logs

**Endpoint:** `POST /api/ai/chat`

Conversational interface to query your API logs.

**Request:**
```json
{
  "message": "Show me all failed login attempts from the last 2 hours",
  "time_range": "2h",
  "limit": 100
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "response": "I found 12 failed login attempts in the last 2 hours...",
    "logs_count": 98,
    "time_range": "2h"
  }
}
```

## Usage Examples

### Using cURL

```bash
# First, login to get auth token
TOKEN=$(curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' \
  | jq -r '.data.access_token')

# Analyze logs
curl -X POST http://localhost:8080/api/ai/analyze-logs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "What are the slowest API endpoints?",
    "time_range": "1h",
    "limit": 100
  }'

# Detect anomalies
curl -X POST http://localhost:8080/api/ai/detect-anomalies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "time_range": "24h",
    "limit": 200
  }'

# Chat with logs
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "How many requests did we receive in the last hour?",
    "time_range": "1h"
  }'
```

### Using JavaScript/Fetch

```javascript
const API_BASE = 'http://localhost:8080/api';
let authToken = '';

// Login
async function login() {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: 'admin@example.com',
      password: 'password'
    })
  });
  const data = await response.json();
  authToken = data.data.access_token;
}

// Analyze logs
async function analyzeLogs(query) {
  const response = await fetch(`${API_BASE}/ai/analyze-logs`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${authToken}`
    },
    body: JSON.stringify({
      query: query,
      time_range: '1h',
      limit: 100
    })
  });
  return await response.json();
}

// Chat with logs
async function chatWithLogs(message) {
  const response = await fetch(`${API_BASE}/ai/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${authToken}`
    },
    body: JSON.stringify({
      message: message,
      time_range: '1h'
    })
  });
  return await response.json();
}

// Usage
await login();
const analysis = await analyzeLogs('What are the most common errors?');
console.log(analysis.data.analysis.text);

const chat = await chatWithLogs('Show me slow endpoints');
console.log(chat.data.response);
```

## Hackathon Strategy

### Winning the Elastic Challenge

Your project integrates:
1. **Elastic Search**: All API calls logged to Elasticsearch
2. **Google Cloud Vertex AI**: Gemini models for intelligent analysis
3. **Conversational AI**: Natural language interface to query logs

### Key Features to Highlight

1. **AI-Powered API Monitoring**
   - Real-time anomaly detection
   - Predictive analytics on API performance
   - Security threat identification

2. **Conversational Log Analysis**
   - Ask questions in natural language
   - Get instant insights from millions of logs
   - No need to write complex Elasticsearch queries

3. **Automated Optimization**
   - AI suggests performance improvements
   - Identifies caching opportunities
   - Recommends rate limiting strategies

4. **Hybrid Search + AI**
   - Elastic's powerful search capabilities
   - Enhanced with Google Cloud's generative AI
   - Best of both worlds

### Demo Video Script (3 minutes)

**Minute 1: Problem Statement**
- Show the challenge: Managing and analyzing millions of API logs
- Traditional methods are time-consuming and require expertise

**Minute 2: Solution Demo**
- Show the conversational interface
- Ask: "What are the slowest endpoints in the last hour?"
- Show AI analyzing logs and providing insights
- Demonstrate anomaly detection catching a security threat

**Minute 3: Impact & Innovation**
- Highlight the seamless integration of Elastic + Google Cloud
- Show real-world use case: Preventing API downtime
- Emphasize developer experience improvement

### Deployment Checklist

- [ ] Deploy backend to Google Cloud Run or App Engine
- [ ] Ensure Elasticsearch is accessible
- [ ] Create demo frontend (optional but recommended)
- [ ] Record 3-minute demo video
- [ ] Prepare GitHub repository with:
  - [ ] Clear README
  - [ ] Open source license
  - [ ] Setup instructions
  - [ ] Architecture diagram
- [ ] Test all API endpoints
- [ ] Prepare presentation slides

### Cost Optimization

**Free Tier Usage:**
- Vertex AI: First 1000 requests/month free
- Use `gemini-1.5-flash` for cost efficiency
- Elasticsearch: Use existing server

**Estimated Costs for Hackathon:**
- Vertex AI: ~$5-10 (with flash model)
- Google Cloud Run: Free tier sufficient
- Total: Under $15 for entire hackathon

## Troubleshooting

### Error: "failed to create Vertex AI client"

**Solution:**
1. Verify `GCP_PROJECT_ID` is correct
2. Check service account key file exists
3. Ensure Vertex AI API is enabled

```bash
gcloud services enable aiplatform.googleapis.com
```

### Error: "permission denied"

**Solution:**
Grant proper IAM roles:

```bash
gcloud projects add-iam-policy-binding your-project-id \
    --member="serviceAccount:vertex-ai-service@your-project-id.iam.gserviceaccount.com" \
    --role="roles/aiplatform.user"
```

### Error: "Cannot search logs in local mode"

**Solution:**
Set `LOG_LOCAL_MODE=false` in `.env` and ensure Elasticsearch is running.

### Slow Response Times

**Solution:**
1. Use `gemini-1.5-flash` instead of `gemini-1.5-pro`
2. Reduce `VERTEX_AI_MAX_TOKENS`
3. Limit the number of logs analyzed (reduce `limit` parameter)

## Additional Resources

- [Vertex AI Documentation](https://cloud.google.com/vertex-ai/docs)
- [Gemini API Guide](https://cloud.google.com/vertex-ai/docs/generative-ai/model-reference/gemini)
- [Elastic Search Documentation](https://www.elastic.co/guide/index.html)
- [AI Accelerate Hackathon Rules](https://devpost.com/hackathons)

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review Google Cloud logs: `gcloud logging read`
3. Check application logs for detailed error messages

---

**Good luck with the AI Accelerate Hackathon! 🚀**
