# Implementation Guide - AI-Powered Log Analysis System

## 🎯 Quick Start Steps

### 1. Set Up Local Elasticsearch (5 minutes)

```powershell
docker pull docker.elastic.co/elasticsearch/elasticsearch:8.11.0

docker run -d --name elasticsearch-local -p 9200:9200 -p 9300:9300 `
  -e "discovery.type=single-node" `
  -e "xpack.security.enabled=false" `
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" `
  docker.elastic.co/elasticsearch/elasticsearch:8.11.0

# Verify
curl http://localhost:9200
```

### 2. Update .env Configuration

```env
# Change project references
LOG_PROJECT_ID=futurefrontier
MONGODB_URI=mongodb://localhost:27018/futurefrontier?directConnection=true
DB_NAME=futurefrontier

# Local Elasticsearch
ELASTICSEARCH_URL=http://localhost:9200
LOG_LOCAL_MODE=false
```

### 3. Test Backend

```powershell
go mod download
go run .
```

### 4. Generate Demo Data

```powershell
# Generate 1000 sample logs
curl -X POST "http://localhost:8080/api/demo/generate?count=1000"
```

### 5. Test AI Endpoints

```powershell
# Login first
curl -X POST http://localhost:8080/api/auth/login `
  -H "Content-Type: application/json" `
  -d '{"email":"admin@example.com","password":"P@ssw0rd!"}'

# Chat with logs
curl -X POST http://localhost:8080/api/ai/chat `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -d '{"message":"Show me failing APIs","time_range":"1h","limit":100}'
```

---

## 📁 Files to Create/Modify

### Backend Files

1. **internal/service/elasticsearch_enhanced.go** - Enhanced search methods
2. **internal/service/vertexai_enhanced.go** - Pattern detection methods  
3. **internal/service/demo_data.go** - Demo data generator
4. **internal/transport/http/handlers/demo_handler.go** - Demo endpoint

### Frontend Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── ChatInterface.tsx
│   │   ├── LogViewer.tsx
│   │   └── PatternDashboard.tsx
│   ├── services/
│   │   └── api.ts
│   └── App.tsx
├── package.json
└── vite.config.ts
```

---

## 🎬 Demo Scenarios

### Scenario 1: Pattern Detection
**Query**: "Show me logs that are always failing in a pattern"
**Expected**: AI identifies /api/users/list with 30% failure rate

### Scenario 2: Traffic Analysis  
**Query**: "Which APIs fail due to high traffic?"
**Expected**: AI shows 503 errors during peak hours

### Scenario 3: Root Cause
**Query**: "Why is /api/users/list failing?"
**Expected**: AI identifies "Database connection timeout" pattern

---

## ✅ Testing Checklist

- [ ] Elasticsearch running on localhost:9200
- [ ] Backend starts without errors
- [ ] Demo data generates successfully
- [ ] AI endpoints return responses
- [ ] Frontend connects to backend
- [ ] Chat interface works
- [ ] Pattern detection works

---

## 🚀 Deployment Options

### Option 1: Google Cloud Run
- Deploy backend as container
- Use Elastic Cloud for Elasticsearch
- Frontend on Vercel

### Option 2: All-in-One
- Single VM with Docker Compose
- Nginx reverse proxy
- Let's Encrypt SSL

---

## 📝 Next Steps

1. Review ARCHITECTURE.md for complete system design
2. Follow Phase 1-4 implementation steps
3. Generate demo data
4. Record 3-minute demo video
5. Deploy to cloud
6. Submit to Devpost

---

**Need Help?** Check TROUBLESHOOTING.md or open an issue on GitHub.
