# AI-Powered Log Analysis System - Complete Setup Guide

## 🎯 Project Summary

**Challenge**: Elastic Challenge - AI-Powered Search  
**Repository**: https://github.com/ChrisChantszto/futurefrontier.git  
**Tech Stack**: Elasticsearch + Google Cloud Vertex AI + Go Fiber + React

### What This System Does

Ask natural language questions about your API logs:
- **"Show me logs that are always failing in a pattern"** → AI identifies recurring failure patterns
- **"Which APIs fail due to high traffic?"** → AI analyzes traffic-related failures
- **"Why is /api/users/list failing?"** → AI performs root cause analysis

---

## 🚀 Quick Start (30 minutes)

### Step 1: Set Up Elasticsearch (5 minutes)

```powershell
# Run the setup script
.\setup-elasticsearch.ps1

# Or manually:
docker pull docker.elastic.co/elasticsearch/elasticsearch:8.11.0
docker run -d --name elasticsearch-local -p 9200:9200 -p 9300:9300 `
  -e "discovery.type=single-node" `
  -e "xpack.security.enabled=false" `
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" `
  docker.elastic.co/elasticsearch/elasticsearch:8.11.0

# Verify
curl http://localhost:9200
```

### Step 2: Set Up MongoDB (5 minutes)

```powershell
# Start MongoDB (if not already running)
docker run -d --name mongo-local -p 27018:27017 mongo:7.0

# Or if it exists:
docker start mongo-local
```

### Step 3: Verify Configuration (2 minutes)

Check `.env` file has these settings:
```env
LOG_PROJECT_ID=futurefrontier
ELASTICSEARCH_URL=http://localhost:9200
LOG_LOCAL_MODE=false
MONGODB_URI=mongodb://localhost:27018/futurefrontier?directConnection=true
DB_NAME=futurefrontier
GCP_PROJECT_ID=grand-principle-475206-b5
```

### Step 4: Start Backend (5 minutes)

```powershell
# Install dependencies
go mod download

# Run the server
go run .

# You should see:
# "listening on :8080"
# "Successfully connected to Elasticsearch"
```

### Step 5: Initialize User & Login (3 minutes)

```powershell
# Create first user
curl -X POST http://localhost:8080/api/auth/init-first-user `
  -H "Content-Type: application/json" `
  -d '{\"Email\":\"admin@example.com\",\"Password\":\"P@ssw0rd!\"}'

# Login (saves cookies)
curl -X POST http://localhost:8080/api/auth/login `
  -H "Content-Type: application/json" `
  -c cookies.txt `
  -d '{\"Email\":\"admin@example.com\",\"Password\":\"P@ssw0rd!\"}'
```

### Step 6: Generate Demo Data (5 minutes)

```powershell
# Generate 1000 sample logs with realistic patterns
curl -X POST "http://localhost:8080/api/demo/generate?count=1000" `
  -H "Content-Type: application/json" `
  -b cookies.txt

# This creates logs with:
# - /api/users/list failing 30% of the time (database timeout pattern)
# - /api/auth/login with 15% authentication failures
# - Traffic spike in the middle causing 503 errors
```

### Step 7: Test AI Endpoints (5 minutes)

```powershell
# 1. Chat with logs
curl -X POST http://localhost:8080/api/ai/chat `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{\"message\":\"Show me failing APIs\",\"time_range\":\"1h\",\"limit\":100}'

# 2. Detect patterns
curl -X POST http://localhost:8080/api/ai/detect-patterns `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{\"time_range\":\"1h\",\"limit\":500}'

# 3. Analyze failures for specific endpoint
curl -X POST http://localhost:8080/api/ai/analyze-failures `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{\"endpoint\":\"/api/users/list\",\"time_range\":\"1h\"}'

# 4. Root cause analysis
curl -X POST http://localhost:8080/api/ai/root-cause-analysis `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{\"error_pattern\":\"Database connection timeout\",\"time_range\":\"1h\"}'

# 5. Traffic-related failures
curl -X POST http://localhost:8080/api/ai/analyze-traffic-failures `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{\"time_range\":\"1h\",\"min_failure_rate\":0.1}'
```

---

## 📁 Project Structure

```
futurefrontier/
├── internal/
│   ├── service/
│   │   ├── elasticsearch.go              # Base Elasticsearch service
│   │   ├── elasticsearch_enhanced.go     # Pattern detection queries
│   │   ├── vertexai.go                   # Base Vertex AI service
│   │   ├── vertexai_enhanced.go          # AI pattern analysis
│   │   └── demo_data.go                  # Demo data generator
│   ├── transport/http/
│   │   ├── handlers/
│   │   │   ├── ai_handler.go             # Base AI endpoints
│   │   │   ├── ai_handler_enhanced.go    # Enhanced AI endpoints
│   │   │   └── demo_handler.go           # Demo data endpoint
│   │   └── routes.go                     # Route registration
│   ├── models/
│   │   └── log.go                        # Log data models
│   └── config/
│       └── config.go                     # Configuration
├── .env                                  # Environment configuration
├── ARCHITECTURE.md                       # System architecture
├── IMPLEMENTATION_GUIDE.md               # Implementation details
├── HACKATHON_SETUP.md                    # This file
└── setup-elasticsearch.ps1               # Elasticsearch setup script
```

---

## 🔧 API Endpoints

### Demo Data
- `POST /api/demo/generate?count=1000` - Generate demo logs

### AI Endpoints (All require authentication)

#### Basic AI
- `POST /api/ai/generate` - General AI content generation
- `POST /api/ai/chat` - Chat with logs
- `POST /api/ai/analyze-logs` - Analyze logs with query
- `POST /api/ai/detect-anomalies` - Detect anomalies
- `POST /api/ai/suggest-optimizations` - Get optimization suggestions

#### Enhanced AI (Pattern Detection)
- `POST /api/ai/detect-patterns` - Find recurring failure patterns
- `POST /api/ai/analyze-failures` - Analyze specific failures
- `POST /api/ai/root-cause-analysis` - Perform root cause analysis
- `POST /api/ai/analyze-traffic-failures` - Analyze traffic-related failures

---

## 🎬 Demo Scenarios

### Scenario 1: Pattern Detection
**Question**: "Show me logs that are always failing in a pattern"

**API Call**:
```powershell
curl -X POST http://localhost:8080/api/ai/detect-patterns `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{\"time_range\":\"1h\",\"limit\":500}'
```

**Expected AI Response**:
```
PATTERN DETECTED:

1. HIGH SEVERITY: /api/users/list - Consistent Database Timeouts
   - Frequency: 30% failure rate (150 out of 500 requests)
   - Error: "Database connection timeout: failed to fetch user list"
   - Timing: Consistent across all time periods
   - Root Cause: Database connection pool exhaustion
   - Recommendation: Increase connection pool size or add caching

2. MEDIUM SEVERITY: /api/auth/login - Authentication Failures
   - Frequency: 15% failure rate
   - Error: "Invalid credentials"
   - Pattern: Multiple failed attempts from same IPs
   - Recommendation: Implement rate limiting

3. CRITICAL: Traffic Spike Pattern
   - Time: Peak hours (middle of dataset)
   - Errors: 503 Service Unavailable
   - Affected: Multiple endpoints
   - Root Cause: Insufficient server capacity
   - Recommendation: Implement auto-scaling
```

### Scenario 2: Traffic Analysis
**Question**: "Which APIs fail due to high traffic?"

**API Call**:
```powershell
curl -X POST http://localhost:8080/api/ai/analyze-traffic-failures `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{\"time_range\":\"1h\",\"min_failure_rate\":0.1}'
```

**Expected AI Response**:
```
TRAFFIC-RELATED FAILURES:

1. /api/users/list
   - Request Volume: 500 requests/hour
   - Error Rate: 30%
   - Performance: Avg latency 750ms (normal: 450ms)
   - Issue: Database connection pool exhausted under load
   - Scaling Recommendation:
     * Horizontal: Add read replicas
     * Vertical: Increase connection pool from 10 to 50
     * Immediate: Implement Redis caching (TTL: 5 minutes)

2. Multiple endpoints during peak (12:00-12:30)
   - Error: 503 Service Unavailable
   - Pattern: Server overload
   - Recommendation: Implement load balancing and auto-scaling
```

### Scenario 3: Root Cause Analysis
**Question**: "Why is /api/users/list failing?"

**API Call**:
```powershell
curl -X POST http://localhost:8080/api/ai/root-cause-analysis `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{\"error_pattern\":\"Database connection timeout\",\"time_range\":\"1h\"}'
```

**Expected AI Response**:
```
ROOT CAUSE ANALYSIS:

Primary Root Cause:
Database connection pool exhaustion

Evidence:
- Error message: "Database connection timeout: failed to fetch user list"
- Frequency: 150 occurrences in 1 hour
- Latency: 750ms average (vs 450ms normal)
- Pattern: Consistent across all time periods (not traffic-related)

Contributing Factors:
1. Inefficient query (likely missing indexes)
2. Connection pool size too small (default: 10)
3. No query timeout configured
4. Missing caching layer

Why This Is Root Cause (Not Symptom):
- Failures occur even during low traffic
- Latency is consistently high
- Error message explicitly mentions connection timeout
- Other endpoints don't show this pattern

Specific Fixes:
1. CODE: Add database indexes on user table
2. CONFIG: Increase connection pool to 50
3. CODE: Implement Redis caching with 5-minute TTL
4. CONFIG: Set query timeout to 5 seconds
5. MONITORING: Add connection pool metrics

Prevention:
- Load testing before deployment
- Connection pool monitoring
- Query performance profiling
```

---

## 🎨 Frontend (Optional - For Full Demo)

### Quick Frontend Setup

```powershell
cd frontend

# Initialize React + TypeScript
npm create vite@latest . -- --template react-ts

# Install dependencies
npm install axios react-query zustand
npm install recharts lucide-react
npm install -D tailwindcss postcss autoprefixer

# Initialize Tailwind
npx tailwindcss init -p

# Install shadcn/ui
npx shadcn-ui@latest init
npx shadcn-ui@latest add button card input textarea badge alert tabs

# Start dev server
npm run dev
```

### Simple Chat Interface (src/App.tsx)

```typescript
import { useState } from 'react';
import axios from 'axios';

function App() {
  const [message, setMessage] = useState('');
  const [response, setResponse] = useState('');
  const [loading, setLoading] = useState(false);

  const sendMessage = async () => {
    setLoading(true);
    try {
      const res = await axios.post('http://localhost:8080/api/ai/chat', {
        message,
        time_range: '1h',
        limit: 100
      }, { withCredentials: true });
      
      setResponse(res.data.data.response);
    } catch (error) {
      console.error(error);
      setResponse('Error: ' + error.message);
    }
    setLoading(false);
  };

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-4">AI Log Analysis</h1>
      <textarea
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        placeholder="Ask a question about your logs..."
        className="w-full p-4 border rounded mb-4"
        rows={4}
      />
      <button
        onClick={sendMessage}
        disabled={loading}
        className="px-6 py-2 bg-blue-500 text-white rounded"
      >
        {loading ? 'Analyzing...' : 'Ask AI'}
      </button>
      {response && (
        <div className="mt-8 p-4 bg-gray-100 rounded">
          <pre className="whitespace-pre-wrap">{response}</pre>
        </div>
      )}
    </div>
  );
}

export default App;
```

---

## 🎥 Demo Video Script (3 minutes)

### 0:00-0:30 - Problem Statement
- Show traditional log analysis (complex Elasticsearch queries)
- Highlight pain points: time-consuming, requires expertise
- "What if you could just ask questions?"

### 0:30-1:00 - Solution Overview
- Show architecture diagram
- Explain: Elasticsearch + Vertex AI + Conversational Interface
- Key features: Pattern detection, root cause analysis, natural language

### 1:00-2:00 - Live Demo
1. Generate demo data (show command)
2. Ask: "Show me failing APIs" (show AI response)
3. Ask: "Why is /api/users/list failing?" (show root cause analysis)
4. Ask: "Which APIs fail due to traffic?" (show traffic analysis)

### 2:00-2:30 - Impact & Results
- Show metrics: 30% failure pattern detected
- Show AI recommendations: specific, actionable
- Highlight time saved: minutes vs hours

### 2:30-3:00 - Future & Call to Action
- Future enhancements: real-time alerts, Slack integration
- Open source: GitHub link
- Call to action: Try it yourself!

---

## ✅ Hackathon Submission Checklist

- [ ] **Working Application**
  - [ ] Backend running on localhost:8080
  - [ ] Elasticsearch running on localhost:9200
  - [ ] All AI endpoints functional
  - [ ] Demo data generated successfully

- [ ] **Open Source Repository**
  - [ ] Code pushed to GitHub: https://github.com/ChrisChantszto/futurefrontier.git
  - [ ] README.md updated
  - [ ] MIT License added
  - [ ] Architecture documentation complete

- [ ] **Demo Video (3 minutes)**
  - [ ] Recorded and edited
  - [ ] Uploaded to YouTube
  - [ ] Link added to README

- [ ] **Hosted Application** (Optional but recommended)
  - [ ] Backend deployed (Google Cloud Run / Heroku)
  - [ ] Frontend deployed (Vercel / Netlify)
  - [ ] Elasticsearch (Elastic Cloud / self-hosted)

- [ ] **Devpost Submission**
  - [ ] Project description
  - [ ] Screenshots
  - [ ] Video link
  - [ ] GitHub link
  - [ ] Challenge selected: Elastic Challenge

---

## 🐛 Troubleshooting

### Elasticsearch won't start
```powershell
# Check Docker
docker ps

# Check logs
docker logs elasticsearch-local

# Restart
docker restart elasticsearch-local
```

### Backend can't connect to Elasticsearch
- Verify `.env` has `ELASTICSEARCH_URL=http://localhost:9200`
- Check Elasticsearch is running: `curl http://localhost:9200`
- Check `LOG_LOCAL_MODE=false`

### AI endpoints return errors
- Verify Google Cloud credentials exist
- Check `GCP_PROJECT_ID` in `.env`
- Ensure Vertex AI API is enabled in Google Cloud Console

### No demo data appears
- Check Elasticsearch is running
- Verify logs with: `curl http://localhost:9200/api-logs-*/_search`
- Check backend logs for errors

---

## 📞 Support

- **GitHub Issues**: https://github.com/ChrisChantszto/futurefrontier/issues
- **Documentation**: See ARCHITECTURE.md and IMPLEMENTATION_GUIDE.md

---

**Good luck with your hackathon submission! 🚀**
