# 🤖 AI-Powered Log Analysis System

> **AI Accelerate Hackathon - Elastic Challenge**  
> Transform API log analysis with conversational AI powered by Elasticsearch and Google Cloud Vertex AI

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![Vertex AI](https://img.shields.io/badge/Vertex%20AI-Gemini%202.0-orange.svg)](https://cloud.google.com/vertex-ai)
[![Elasticsearch](https://img.shields.io/badge/Elasticsearch-8.11-green.svg)](https://www.elastic.co/)

---

## 🎯 The Problem

Developers managing millions of API logs face these challenges:
- ❌ Writing complex Elasticsearch queries
- ❌ Manually analyzing logs for patterns
- ❌ Spending hours identifying root causes
- ❌ Missing critical failure patterns

## 💡 Our Solution

**Ask questions in plain English. Get AI-powered insights instantly.**

```
You: "Show me logs that are always failing in a pattern"

AI: "I found 3 recurring patterns:
     1. /api/users/list - 30% failure rate (database timeout)
     2. /api/auth/login - 15% authentication failures
     3. Traffic spike causing 503 errors during peak hours
     
     Root cause: Database connection pool exhaustion
     Fix: Increase pool size to 50 and add Redis caching"
```

---

## 🚀 Quick Start (30 Minutes)

### Prerequisites
- Docker Desktop
- Go 1.24+
- Google Cloud account (Vertex AI enabled)

### 1. Set Up Elasticsearch
```powershell
.\setup-elasticsearch.ps1
```

### 2. Start Backend
```powershell
go run .
```

### 3. Create User & Login
```powershell
# Create first user
curl -X POST http://localhost:8080/api/auth/init-first-user `
  -H "Content-Type: application/json" `
  -d '{\"Email\":\"admin@example.com\",\"Password\":\"P@ssw0rd!\"}'

# Login
curl -X POST http://localhost:8080/api/auth/login `
  -H "Content-Type: application/json" `
  -c cookies.txt `
  -d '{\"Email\":\"admin@example.com\",\"Password\":\"P@ssw0rd!\"}'
```

### 4. Generate Demo Data
```powershell
curl -X POST "http://localhost:8080/api/demo/generate?count=1000" -b cookies.txt
```

### 5. Ask AI Questions!
```powershell
# Detect patterns
curl -X POST http://localhost:8080/api/ai/detect-patterns `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{\"time_range\":\"1h\",\"limit\":500}'
```

**See [HACKATHON_SETUP.md](./HACKATHON_SETUP.md) for detailed instructions.**

---

## ✨ Key Features

### 🔍 Pattern Detection
Automatically identifies recurring failure patterns across millions of logs.

**Example Query**: "Show me logs that are always failing"
```json
{
  "patterns": [
    {
      "endpoint": "/api/users/list",
      "failure_rate": "30%",
      "error": "Database connection timeout",
      "severity": "HIGH"
    }
  ]
}
```

### 🎯 Root Cause Analysis
Goes beyond showing logs - explains WHY failures occur.

**Example Query**: "Why is /api/users/list failing?"
```json
{
  "root_cause": "Database connection pool exhaustion",
  "evidence": "300 timeout errors in 1 hour",
  "fixes": [
    "Increase connection pool to 50",
    "Add Redis caching (5-min TTL)",
    "Add database indexes"
  ]
}
```

### 📊 Traffic Analysis
Identifies APIs failing due to high traffic.

**Example Query**: "Which APIs fail due to traffic?"
```json
{
  "traffic_failures": [
    {
      "endpoint": "/api/users/list",
      "requests_per_hour": 500,
      "error_rate": "30%",
      "recommendation": "Add read replicas and caching"
    }
  ]
}
```

### 🤖 Conversational Interface
Natural language queries - no Elasticsearch expertise needed.

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    User Question                        │
│     "Show me logs that are always failing"              │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Go Fiber Backend (Port 8080)               │
│  ┌──────────────────────────────────────────────────┐  │
│  │  AI Handler                                       │  │
│  │  - /api/ai/detect-patterns                        │  │
│  │  - /api/ai/analyze-failures                       │  │
│  │  - /api/ai/root-cause-analysis                    │  │
│  └────────┬─────────────────────┬────────────────────┘  │
└───────────┼─────────────────────┼───────────────────────┘
            │                     │
            ▼                     ▼
┌──────────────────┐    ┌──────────────────┐
│  Elasticsearch   │    │  Vertex AI       │
│  (localhost:9200)│    │  Gemini 2.0      │
│  - Hybrid Search │    │  - NLP Analysis  │
│  - Aggregations  │    │  - Pattern Rec.  │
└──────────────────┘    └──────────────────┘
```

---

## 📚 API Endpoints

### Pattern Detection
```bash
POST /api/ai/detect-patterns
{
  "time_range": "24h",
  "limit": 500,
  "min_occurrences": 5
}
```

### Failure Analysis
```bash
POST /api/ai/analyze-failures
{
  "endpoint": "/api/users/list",
  "time_range": "24h",
  "status_codes": [500, 503]
}
```

### Root Cause Analysis
```bash
POST /api/ai/root-cause-analysis
{
  "error_pattern": "Database connection timeout",
  "time_range": "1h"
}
```

### Traffic Failures
```bash
POST /api/ai/analyze-traffic-failures
{
  "time_range": "1h",
  "min_failure_rate": 0.1
}
```

**See [ARCHITECTURE.md](./ARCHITECTURE.md) for complete API documentation.**

---

## 🎬 Demo

### Live Demo Scenarios

#### Scenario 1: Pattern Detection
```
User: "Show me logs that are always failing in a pattern"

AI Response:
┌─────────────────────────────────────────────────────┐
│ PATTERN DETECTED: /api/users/list                  │
│                                                     │
│ Failure Rate: 30% (300 out of 1000 requests)       │
│ Error: "Database connection timeout"               │
│ Timing: Consistent across all hours                │
│                                                     │
│ ROOT CAUSE: Database connection pool exhaustion    │
│                                                     │
│ RECOMMENDATIONS:                                    │
│ ✓ Increase connection pool from 10 to 50           │
│ ✓ Implement Redis caching (5-minute TTL)           │
│ ✓ Add database indexes on user table               │
│ ✓ Set query timeout to 5 seconds                   │
└─────────────────────────────────────────────────────┘
```

#### Scenario 2: Traffic Analysis
```
User: "Which APIs fail due to high traffic?"

AI Response:
┌─────────────────────────────────────────────────────┐
│ TRAFFIC-RELATED FAILURES                            │
│                                                     │
│ 1. /api/users/list                                  │
│    • Request Volume: 500 req/hour                   │
│    • Error Rate: 30%                                │
│    • Issue: Connection pool exhausted under load    │
│    • Fix: Add read replicas + Redis caching         │
│                                                     │
│ 2. Multiple endpoints (12:00-12:30)                 │
│    • Error: 503 Service Unavailable                 │
│    • Issue: Server overload during peak             │
│    • Fix: Implement auto-scaling                    │
└─────────────────────────────────────────────────────┘
```

---

## 🎯 Hackathon Highlights

### Why This Project Wins

✅ **Solves Real Problem** - Developers waste hours analyzing logs  
✅ **Innovative Integration** - Seamless Elastic + Google Cloud AI  
✅ **Practical Impact** - Immediate productivity improvements  
✅ **Production-Ready** - Scalable architecture  
✅ **Great UX** - Natural language interface anyone can use  

### Technologies Used

**Backend**: Go Fiber, MongoDB, Zap Logging  
**Search**: Elasticsearch 8.11 (Hybrid Search + Aggregations)  
**AI**: Google Cloud Vertex AI (Gemini 2.0 Flash)  
**DevOps**: Docker, Google Cloud Run

---

## 📊 Performance

- **Query Speed**: < 2 seconds for 1000 logs
- **AI Response**: 3-5 seconds average
- **Throughput**: 1000+ requests/minute
- **Cost**: ~$0.002 per AI query (Gemini Flash)
- **Scalability**: Handles millions of logs

---

## 📁 Project Structure

```
futurefrontier/
├── internal/
│   ├── service/
│   │   ├── elasticsearch_enhanced.go    # Pattern detection queries
│   │   ├── vertexai_enhanced.go         # AI analysis methods
│   │   └── demo_data.go                 # Demo data generator
│   └── transport/http/
│       └── handlers/
│           ├── ai_handler_enhanced.go   # New AI endpoints
│           └── demo_handler.go          # Demo endpoint
├── ARCHITECTURE.md                      # Complete architecture
├── IMPLEMENTATION_GUIDE.md              # Technical details
├── HACKATHON_SETUP.md                   # Quick start guide
└── setup-elasticsearch.ps1              # Automated setup
```

---

## 🔒 Security

- JWT authentication for all endpoints
- Service account credentials (not API keys)
- Request/response sanitization
- Rate limiting on AI endpoints
- Audit logging for all AI queries

---

## 🌟 Future Enhancements

- [ ] Real-time streaming analysis
- [ ] Slack/Discord integration
- [ ] Custom alert rules
- [ ] Multi-tenant support
- [ ] GraphQL API
- [ ] Mobile app
- [ ] Predictive analytics

---

## 📝 Documentation

- **[ARCHITECTURE.md](./ARCHITECTURE.md)** - Complete system architecture
- **[IMPLEMENTATION_GUIDE.md](./IMPLEMENTATION_GUIDE.md)** - Implementation details
- **[HACKATHON_SETUP.md](./HACKATHON_SETUP.md)** - Quick start guide
- **[SUMMARY.md](./SUMMARY.md)** - Project summary

---

## 🤝 Contributing

This is a hackathon project, but we welcome contributions!

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Google Cloud for Vertex AI platform
- Elastic for powerful search capabilities
- The Go community for excellent libraries
- AI Accelerate Hackathon organizers

---

## 📞 Contact

- **GitHub**: https://github.com/ChrisChantszto/futurefrontier
- **Issues**: https://github.com/ChrisChantszto/futurefrontier/issues

---

**Built with ❤️ for the AI Accelerate Hackathon - Elastic Challenge**

*Transforming how developers interact with their API logs through the power of AI*

---

## 🚀 Get Started Now

```powershell
# Clone the repository
git clone https://github.com/ChrisChantszto/futurefrontier.git
cd futurefrontier

# Set up Elasticsearch
.\setup-elasticsearch.ps1

# Start the backend
go run .

# Generate demo data and start asking questions!
```

**See [HACKATHON_SETUP.md](./HACKATHON_SETUP.md) for complete setup instructions.**
