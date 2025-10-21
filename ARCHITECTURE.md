# AI-Powered Log Analysis System - Architecture

## 🎯 Project Overview

**Hackathon Challenge**: Elastic Challenge - Build AI-Powered Search  
**Repository**: https://github.com/ChrisChantszto/futurefrontier.git  
**Tech Stack**: Elasticsearch + Google Cloud Vertex AI + Go + React

### Problem Statement
Developers managing millions of API logs face challenges:
- Writing complex Elasticsearch queries
- Manually identifying failure patterns
- Detecting anomalies and security threats
- Finding root causes of recurring issues

### Solution
A conversational AI system that lets you ask questions like:
- "Show me the logs that are always failing in a pattern"
- "What APIs are failing due to high traffic?"
- "Analyze the root causes of 500 errors in the last hour"

---

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend Layer                          │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  React + TypeScript + TailwindCSS + shadcn/ui            │  │
│  │  - Conversational Chat Interface                          │  │
│  │  - Real-time Log Visualization                            │  │
│  │  - Pattern Detection Dashboard                            │  │
│  │  - Anomaly Alerts                                          │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │ REST API (JSON)
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend Layer (Go)                         │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                  Go Fiber Web Server                      │  │
│  │  ┌────────────────────────────────────────────────────┐  │  │
│  │  │  AI Handler (Enhanced)                             │  │  │
│  │  │  - /api/ai/chat-with-logs                          │  │  │
│  │  │  - /api/ai/detect-patterns                         │  │  │
│  │  │  - /api/ai/analyze-failures                        │  │  │
│  │  │  - /api/ai/root-cause-analysis                     │  │  │
│  │  │  - /api/ai/suggest-optimizations                   │  │  │
│  │  └────────────────────────────────────────────────────┘  │  │
│  │                                                            │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │  │
│  │  │ Vertex AI    │  │ Elasticsearch│  │ Logger       │   │  │
│  │  │ Service      │  │ Service      │  │ Service      │   │  │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘   │  │
│  └─────────┼──────────────────┼──────────────────┼──────────┘  │
└────────────┼──────────────────┼──────────────────┼─────────────┘
             │                  │                  │
             ▼                  ▼                  ▼
┌────────────────────┐  ┌──────────────────┐  ┌─────────────┐
│  Google Cloud      │  │  Elasticsearch   │  │  MongoDB    │
│  Vertex AI         │  │  (Local/Cloud)   │  │  (Metadata) │
│  - Gemini 2.0      │  │  - API Logs      │  │  - Users    │
│  - Flash Model     │  │  - Error Logs    │  │  - Settings │
│  - NLP Analysis    │  │  - Hybrid Search │  │             │
└────────────────────┘  └──────────────────┘  └─────────────┘
```

---

## 📊 Data Flow

### 1. Log Ingestion Flow
```
Application → Logger Middleware → Elasticsearch Service → Elasticsearch
                                                         ↓
                                                    Index: api-logs-YYYY-MM
                                                    Index: error-logs-YYYY-MM
```

### 2. AI Query Flow
```
User Question → Frontend → Backend API → Elasticsearch Query
                                              ↓
                                         Retrieve Logs
                                              ↓
                                         Vertex AI Analysis
                                              ↓
                                         AI Response → Frontend
```

### 3. Pattern Detection Flow
```
Scheduled Job → Fetch Recent Logs → AI Pattern Analysis → Alert System
```

---

## 🔧 Component Details

### Backend Components

#### 1. **Elasticsearch Service** (`internal/service/elasticsearch.go`)
- **Purpose**: Log storage and retrieval
- **Features**:
  - Time-based indexing (monthly rotation)
  - Hybrid search capabilities
  - Aggregation queries for pattern detection
  - Full-text search on error messages

**Key Methods**:
```go
- SearchLogs(timeRange, limit) - Basic time-range search
- SearchLogsByPattern(pattern) - Pattern-based search
- AggregateByEndpoint() - Group logs by API endpoint
- AggregateByStatusCode() - Group by HTTP status
- SearchFailingAPIs(threshold) - Find APIs with high failure rates
```

#### 2. **Vertex AI Service** (`internal/service/vertexai.go`)
- **Purpose**: AI-powered log analysis
- **Model**: Gemini 2.0 Flash (fast, cost-effective)
- **Features**:
  - Natural language understanding
  - Pattern recognition
  - Root cause analysis
  - Optimization suggestions

**Key Methods**:
```go
- AnalyzeAPILogs(logs, query) - General analysis
- DetectPatterns(logs) - Find recurring patterns
- AnalyzeFailures(logs) - Failure analysis
- RootCauseAnalysis(logs) - Identify root causes
- SuggestOptimizations(logs) - Performance recommendations
```

#### 3. **AI Handler** (`internal/transport/http/handlers/ai_handler.go`)
- **Purpose**: HTTP endpoints for AI features
- **Endpoints**:
  - `POST /api/ai/chat-with-logs` - Conversational interface
  - `POST /api/ai/detect-patterns` - Pattern detection
  - `POST /api/ai/analyze-failures` - Failure analysis
  - `POST /api/ai/root-cause-analysis` - Root cause identification
  - `POST /api/ai/suggest-optimizations` - Performance tips

### Frontend Components

#### 1. **Chat Interface** (`frontend/src/components/ChatInterface.tsx`)
- Natural language input
- Streaming responses
- Context-aware suggestions
- History management

#### 2. **Log Viewer** (`frontend/src/components/LogViewer.tsx`)
- Real-time log display
- Filtering and sorting
- Syntax highlighting
- Export functionality

#### 3. **Pattern Dashboard** (`frontend/src/components/PatternDashboard.tsx`)
- Visual pattern detection
- Failure rate charts
- Anomaly alerts
- Time-series graphs

#### 4. **Analytics Panel** (`frontend/src/components/AnalyticsPanel.tsx`)
- Performance metrics
- API usage statistics
- Error rate trends
- Optimization recommendations

---

## 🗄️ Data Models

### API Log Entry
```json
{
  "project_id": "futurefrontier",
  "timestamp": "2025-10-21T15:30:00Z",
  "request_id": "uuid",
  "method": "POST",
  "path": "/api/auth/login",
  "route": "/api/auth/login",
  "status_code": 401,
  "latency_ms": 145,
  "client_ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "error_message": "Invalid credentials",
  "request_body": "{...}",
  "response_body": "{...}"
}
```

### Error Log Entry
```json
{
  "project_id": "futurefrontier",
  "timestamp": "2025-10-21T15:30:00Z",
  "request_id": "uuid",
  "method": "GET",
  "path": "/api/users/list",
  "status_code": 500,
  "error_message": "Database connection timeout",
  "error_type": "DatabaseError",
  "stack_trace": "...",
  "client_ip": "192.168.1.100"
}
```

---

## 🔍 Key Features Implementation

### 1. Pattern Detection
**Use Case**: "Show me logs that are always failing in a pattern"

**Implementation**:
```go
// 1. Aggregate logs by endpoint and status code
aggregation := es.AggregateByEndpoint(timeRange)

// 2. Filter endpoints with high failure rates
failingAPIs := filterByFailureRate(aggregation, threshold)

// 3. Send to AI for pattern analysis
patterns := vertexAI.DetectPatterns(failingAPIs)

// 4. Return structured insights
return patterns
```

**AI Prompt**:
```
Analyze these API logs and identify recurring failure patterns:
- Which endpoints fail consistently?
- What are the common error messages?
- Are there time-based patterns (e.g., peak hours)?
- What are the likely root causes?
```

### 2. Root Cause Analysis
**Use Case**: "Analyze why /api/users/list is failing"

**Implementation**:
```go
// 1. Fetch logs for specific endpoint
logs := es.SearchLogsByEndpoint("/api/users/list", timeRange)

// 2. Group by error type
errorGroups := groupByErrorType(logs)

// 3. AI analysis
rootCauses := vertexAI.RootCauseAnalysis(errorGroups)

// 4. Return with recommendations
return rootCauses
```

**AI Prompt**:
```
Perform root cause analysis on these error logs:
- What is the primary cause of failures?
- Are there secondary contributing factors?
- Is this a code issue, infrastructure issue, or external dependency?
- What specific actions should be taken to fix this?
```

### 3. Traffic-Based Failure Detection
**Use Case**: "Show APIs failing due to high traffic"

**Implementation**:
```go
// 1. Aggregate by endpoint with request count and error rate
stats := es.AggregateWithTrafficStats(timeRange)

// 2. Identify high-traffic endpoints with elevated error rates
candidates := findTrafficRelatedFailures(stats)

// 3. AI analysis
analysis := vertexAI.AnalyzeTrafficFailures(candidates)

// 4. Return insights with scaling recommendations
return analysis
```

---

## 🚀 Setup Instructions

### Prerequisites
- Go 1.24+
- Node.js 18+
- Docker Desktop
- Google Cloud account with Vertex AI enabled

### 1. Set Up Local Elasticsearch

```powershell
# Pull Elasticsearch 8.x
docker pull docker.elastic.co/elasticsearch/elasticsearch:8.11.0

# Run Elasticsearch (single-node, no security for local dev)
docker run -d `
  --name elasticsearch-local `
  -p 9200:9200 `
  -p 9300:9300 `
  -e "discovery.type=single-node" `
  -e "xpack.security.enabled=false" `
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" `
  docker.elastic.co/elasticsearch/elasticsearch:8.11.0

# Verify
curl http://localhost:9200
```

### 2. Configure Environment

Update `.env`:
```env
# Elasticsearch (Local)
ELASTICSEARCH_URL=http://localhost:9200
ELASTICSEARCH_USER=
ELASTICSEARCH_PASS=
LOG_LOCAL_MODE=false

# Google Cloud
GCP_PROJECT_ID=grand-principle-475206-b5
GCP_LOCATION=us-central1
GOOGLE_APPLICATION_CREDENTIALS=./grand-principle-475206-b5-da3a41415c39.json
VERTEX_AI_MODEL=gemini-2.0-flash-001

# Project ID
LOG_PROJECT_ID=futurefrontier
```

### 3. Run Backend

```powershell
# Install dependencies
go mod download

# Run server
go run .
```

### 4. Run Frontend

```powershell
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev
```

---

## 📝 API Endpoints

### AI Endpoints

#### 1. Chat with Logs
```http
POST /api/ai/chat-with-logs
Content-Type: application/json

{
  "message": "Show me all failed login attempts in the last hour",
  "time_range": "1h",
  "limit": 100
}
```

**Response**:
```json
{
  "success": true,
  "message": "Message processed successfully",
  "data": {
    "response": "I found 12 failed login attempts...",
    "logs_count": 98,
    "time_range": "1h"
  }
}
```

#### 2. Detect Patterns
```http
POST /api/ai/detect-patterns
Content-Type: application/json

{
  "time_range": "24h",
  "limit": 500,
  "min_occurrences": 5
}
```

#### 3. Analyze Failures
```http
POST /api/ai/analyze-failures
Content-Type: application/json

{
  "endpoint": "/api/users/list",
  "time_range": "24h",
  "status_codes": [500, 502, 503]
}
```

#### 4. Root Cause Analysis
```http
POST /api/ai/root-cause-analysis
Content-Type: application/json

{
  "error_pattern": "Database connection timeout",
  "time_range": "1h"
}
```

---

## 🎨 Frontend Architecture

### Tech Stack
- **Framework**: React 18 + TypeScript
- **Styling**: TailwindCSS + shadcn/ui
- **State Management**: React Query + Zustand
- **Charts**: Recharts
- **Icons**: Lucide React
- **HTTP Client**: Axios

### Component Structure
```
frontend/
├── src/
│   ├── components/
│   │   ├── ui/              # shadcn/ui components
│   │   ├── ChatInterface.tsx
│   │   ├── LogViewer.tsx
│   │   ├── PatternDashboard.tsx
│   │   ├── AnalyticsPanel.tsx
│   │   └── AnomalyAlerts.tsx
│   ├── hooks/
│   │   ├── useAIChat.ts
│   │   ├── useLogs.ts
│   │   └── usePatterns.ts
│   ├── services/
│   │   └── api.ts
│   ├── types/
│   │   └── index.ts
│   └── App.tsx
```

---

## 🧪 Demo Data Generation

### Generate Sample Logs
```go
// internal/service/demo_data.go
func GenerateDemoLogs(count int) {
    // Generate realistic API logs with:
    // - Various endpoints
    // - Different status codes
    // - Realistic latencies
    // - Error patterns
    // - Traffic spikes
}
```

### Demo Scenarios
1. **High Traffic Failure**: Simulate 500 errors during peak hours
2. **Authentication Issues**: Generate failed login patterns
3. **Database Timeouts**: Create connection timeout errors
4. **Rate Limiting**: Show 429 errors for specific IPs
5. **Slow Endpoints**: Generate high-latency requests

---

## 📊 Hackathon Deliverables

### 1. Working Application
- ✅ Backend API with AI endpoints
- ✅ Frontend UI with conversational interface
- ✅ Local Elasticsearch integration
- ✅ Google Cloud Vertex AI integration

### 2. Open Source Repository
- ✅ Public GitHub repo: https://github.com/ChrisChantszto/futurefrontier.git
- ✅ MIT License
- ✅ Complete source code
- ✅ Setup instructions

### 3. Demo Video (3 minutes)
**Script**:
- 0:00-0:30: Problem statement
- 0:30-1:00: Solution overview
- 1:00-2:00: Live demo (conversational queries)
- 2:00-2:30: Pattern detection showcase
- 2:30-3:00: Impact and future plans

### 4. Hosted Project
**Options**:
- Google Cloud Run (Backend)
- Vercel/Netlify (Frontend)
- Elastic Cloud (Elasticsearch)

---

## 🎯 Unique Value Propositions

### 1. Conversational Interface
Unlike traditional log analysis tools, users can ask questions naturally:
- "What's causing the spike in 500 errors?"
- "Show me slow endpoints in the last hour"
- "Are there any security threats?"

### 2. Pattern Recognition
AI automatically detects:
- Recurring failure patterns
- Time-based anomalies
- Traffic-related issues
- Security threats

### 3. Root Cause Analysis
Goes beyond showing logs - explains WHY things fail:
- Database connection issues
- Memory leaks
- External API failures
- Rate limiting

### 4. Actionable Recommendations
Provides specific fixes:
- "Add Redis caching for /api/users/list"
- "Increase database connection pool"
- "Implement exponential backoff for external APIs"

---

## 🔮 Future Enhancements

1. **Real-time Streaming**: WebSocket-based live log analysis
2. **Alert System**: Proactive notifications for anomalies
3. **Multi-tenant Support**: Separate log spaces for different projects
4. **Custom Dashboards**: User-configurable analytics views
5. **Integration**: Slack, Discord, PagerDuty notifications
6. **ML Models**: Train custom models for specific use cases
7. **Predictive Analytics**: Forecast potential issues before they occur

---

## 📈 Performance Metrics

### Target Metrics
- **Query Response Time**: < 2 seconds for 1000 logs
- **AI Analysis Time**: 3-5 seconds average
- **Throughput**: 1000+ requests/minute
- **Cost per Query**: ~$0.002 (Gemini Flash)
- **Elasticsearch Index Size**: ~1GB per million logs

### Optimization Strategies
1. **Caching**: Redis cache for frequent queries
2. **Pagination**: Limit log retrieval to relevant subset
3. **Aggregation**: Use Elasticsearch aggregations for statistics
4. **Async Processing**: Background jobs for heavy analysis
5. **Rate Limiting**: Prevent abuse of AI endpoints

---

## 🔒 Security Considerations

1. **Authentication**: JWT-based auth for all AI endpoints
2. **Rate Limiting**: Prevent excessive AI usage
3. **Data Sanitization**: Remove sensitive data from logs
4. **Access Control**: Role-based access to logs
5. **Audit Logging**: Track all AI queries
6. **Encryption**: TLS for data in transit

---

## 📚 Documentation Structure

```
docs/
├── SETUP.md              # Detailed setup guide
├── API.md                # API documentation
├── DEPLOYMENT.md         # Deployment instructions
├── DEMO_SCRIPT.md        # Demo video script
├── TROUBLESHOOTING.md    # Common issues
└── CONTRIBUTING.md       # Contribution guidelines
```

---

## ✅ Hackathon Checklist

- [ ] Local Elasticsearch running
- [ ] Backend API functional
- [ ] Frontend UI built
- [ ] AI endpoints working
- [ ] Demo data generated
- [ ] Video recorded (3 min)
- [ ] GitHub repo public
- [ ] MIT License added
- [ ] README.md complete
- [ ] Application deployed
- [ ] Devpost submission ready

---

**Built for AI Accelerate Hackathon - Elastic Challenge**  
*Transforming log analysis through conversational AI*
