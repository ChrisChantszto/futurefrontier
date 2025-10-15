# AI Accelerate Hackathon - Quick Start Guide

Get your AI-powered API monitoring system running in **15 minutes**!

## 🚀 Quick Setup (5 Steps)

### Step 1: Google Cloud Setup (5 min)

```bash
# 1. Login to Google Cloud
gcloud auth login

# 2. Set your project (replace with your actual project ID)
gcloud config set project YOUR_PROJECT_ID

# 3. Enable Vertex AI API
gcloud services enable aiplatform.googleapis.com

# 4. Create service account
gcloud iam service-accounts create vertex-ai-service \
    --display-name="Vertex AI Service Account"

# 5. Grant permissions
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
    --member="serviceAccount:vertex-ai-service@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/aiplatform.user"

# 6. Create and download key
gcloud iam service-accounts keys create gcp-service-account-key.json \
    --iam-account=vertex-ai-service@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

### Step 2: Configure Backend (2 min)

1. Move the key file to your project:
```bash
mv gcp-service-account-key.json c:/Users/ttchan67/onetake-corpsite-backend/
```

2. Update `.env` file:
```env
# Replace these values
GCP_PROJECT_ID=YOUR_PROJECT_ID
GCP_LOCATION=us-central1
GOOGLE_APPLICATION_CREDENTIALS=./gcp-service-account-key.json
VERTEX_AI_MODEL=gemini-1.5-flash
```

### Step 3: Build & Run (2 min)

```bash
# Navigate to project
cd c:/Users/ttchan67/onetake-corpsite-backend

# Build
go build

# Run
./onetake-backend.exe
```

Look for these log messages:
```
✓ Vertex AI service initialized
✓ AI routes registered successfully
```

### Step 4: Test the API (3 min)

```bash
# 1. Login (replace with your credentials)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your-email","password":"your-password"}'

# Save the token from response

# 2. Test AI Chat
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Analyze my API performance",
    "time_range": "1h"
  }'
```

### Step 5: Verify Everything Works (3 min)

Test all endpoints:

```bash
# Set your token
TOKEN="your-access-token-here"

# 1. Generate AI content
curl -X POST http://localhost:8080/api/ai/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello, Vertex AI!"}'

# 2. Analyze logs
curl -X POST http://localhost:8080/api/ai/analyze-logs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"time_range": "1h", "limit": 50}'

# 3. Detect anomalies
curl -X POST http://localhost:8080/api/ai/detect-anomalies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"time_range": "1h"}'
```

## 🎯 What You've Built

### AI-Powered Features

1. **Conversational Log Analysis**
   - Ask questions in natural language
   - Get instant insights from API logs
   - Example: "Show me all errors in the last hour"

2. **Anomaly Detection**
   - AI automatically detects unusual patterns
   - Security threat identification
   - Performance degradation alerts

3. **Optimization Suggestions**
   - AI analyzes your API usage
   - Suggests caching strategies
   - Identifies slow endpoints

4. **Hybrid Search**
   - Elasticsearch for fast log retrieval
   - Vertex AI Gemini for intelligent analysis
   - Best of both worlds

### Technology Stack

- **Backend**: Go Fiber
- **AI**: Google Cloud Vertex AI (Gemini 1.5 Flash)
- **Search**: Elasticsearch
- **Database**: MongoDB
- **Logging**: Structured logging with Zap

## 📊 Available API Endpoints

All endpoints require authentication (Bearer token).

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/ai/generate` | POST | Generate AI content |
| `/api/ai/analyze-logs` | POST | Analyze API logs with AI |
| `/api/ai/detect-anomalies` | POST | Detect anomalies in logs |
| `/api/ai/suggest-optimizations` | POST | Get optimization suggestions |
| `/api/ai/chat` | POST | Chat with your logs |

## 🏆 Hackathon Submission Checklist

### Required Components

- [x] **Backend Integration**: Go Fiber + Vertex AI ✓
- [x] **Elastic Integration**: Logs stored in Elasticsearch ✓
- [x] **Google Cloud AI**: Vertex AI Gemini models ✓
- [x] **Conversational Interface**: Natural language queries ✓
- [ ] **Demo Video**: Record 3-minute demo
- [ ] **GitHub Repository**: Public with open source license
- [ ] **Hosted Project**: Deploy to Google Cloud Run
- [ ] **README**: Clear setup instructions

### Demo Video Outline (3 minutes)

**0:00-0:45 - Problem & Solution**
- Show the challenge: Analyzing millions of API logs
- Introduce your AI-powered solution

**0:45-2:00 - Live Demo**
- Show conversational interface
- Ask: "What are my slowest endpoints?"
- Demonstrate anomaly detection
- Show optimization suggestions

**2:00-2:45 - Technical Architecture**
- Explain Elastic + Google Cloud integration
- Highlight hybrid search capabilities
- Show code snippets

**2:45-3:00 - Impact & Future**
- Real-world applications
- Developer experience improvements
- Call to action

### GitHub Repository Structure

```
onetake-corpsite-backend/
├── README.md                    # Project overview
├── VERTEX_AI_SETUP.md          # Detailed setup guide
├── HACKATHON_QUICKSTART.md     # This file
├── LICENSE                      # Open source license
├── .env.example                 # Environment template
├── main.go                      # Entry point
├── go.mod                       # Dependencies
├── internal/
│   ├── config/                  # Configuration
│   ├── service/
│   │   └── vertexai.go         # Vertex AI integration
│   └── transport/
│       └── http/
│           └── handlers/
│               └── ai_handler.go # AI endpoints
└── docs/
    └── architecture.png         # System diagram
```

## 🎨 Create a Simple Frontend (Optional)

Create a basic HTML interface for the demo:

```html
<!DOCTYPE html>
<html>
<head>
    <title>AI Log Analyzer</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 p-8">
    <div class="max-w-4xl mx-auto bg-white rounded-lg shadow-lg p-6">
        <h1 class="text-3xl font-bold mb-6">AI-Powered Log Analyzer</h1>
        
        <div class="mb-4">
            <label class="block text-sm font-medium mb-2">Ask a question about your API logs:</label>
            <input type="text" id="question" 
                   class="w-full border rounded px-4 py-2"
                   placeholder="e.g., What are the slowest endpoints in the last hour?">
        </div>
        
        <button onclick="askAI()" 
                class="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600">
            Analyze
        </button>
        
        <div id="response" class="mt-6 p-4 bg-gray-50 rounded hidden">
            <h3 class="font-bold mb-2">AI Response:</h3>
            <p id="responseText" class="text-gray-700"></p>
        </div>
    </div>
    
    <script>
        const API_BASE = 'http://localhost:8080/api';
        const TOKEN = 'YOUR_TOKEN_HERE'; // Replace after login
        
        async function askAI() {
            const question = document.getElementById('question').value;
            const responseDiv = document.getElementById('response');
            const responseText = document.getElementById('responseText');
            
            responseText.textContent = 'Analyzing...';
            responseDiv.classList.remove('hidden');
            
            try {
                const response = await fetch(`${API_BASE}/ai/chat`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${TOKEN}`
                    },
                    body: JSON.stringify({
                        message: question,
                        time_range: '1h'
                    })
                });
                
                const data = await response.json();
                responseText.textContent = data.data.response;
            } catch (error) {
                responseText.textContent = 'Error: ' + error.message;
            }
        }
    </script>
</body>
</html>
```

Save as `demo.html` and open in browser.

## 🐛 Common Issues & Fixes

### Issue: "GCP_PROJECT_ID is required"

**Fix:** Update `.env` with your actual project ID:
```env
GCP_PROJECT_ID=your-actual-project-id
```

### Issue: "failed to create Vertex AI client"

**Fix:** Verify service account key exists:
```bash
ls gcp-service-account-key.json
```

### Issue: "Cannot search logs in local mode"

**Fix:** Ensure Elasticsearch is running and set:
```env
LOG_LOCAL_MODE=false
ELASTICSEARCH_URL=http://server.onetakesolutions.com.hk:9201
```

### Issue: "permission denied"

**Fix:** Re-run IAM binding:
```bash
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
    --member="serviceAccount:vertex-ai-service@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/aiplatform.user"
```

## 💡 Pro Tips

1. **Use Flash Model**: `gemini-1.5-flash` is faster and cheaper than `gemini-1.5-pro`
2. **Limit Log Count**: Start with `limit: 50-100` for faster responses
3. **Time Ranges**: Use shorter ranges (`1h`, `2h`) for quicker analysis
4. **Cache Responses**: Consider caching common queries
5. **Error Handling**: Always check response status codes

## 📈 Next Steps

1. **Deploy to Production**
   ```bash
   # Deploy to Google Cloud Run
   gcloud run deploy onetake-backend \
       --source . \
       --region us-central1 \
       --allow-unauthenticated
   ```

2. **Add More Features**
   - Real-time streaming analysis
   - Scheduled anomaly reports
   - Slack/email notifications
   - Custom dashboards

3. **Optimize Performance**
   - Implement response caching
   - Add rate limiting
   - Use batch processing for large datasets

4. **Enhance Security**
   - Add API key authentication
   - Implement request signing
   - Enable audit logging

## 🎓 Learning Resources

- [Vertex AI Documentation](https://cloud.google.com/vertex-ai/docs)
- [Gemini API Reference](https://cloud.google.com/vertex-ai/docs/generative-ai/model-reference/gemini)
- [Elasticsearch Guide](https://www.elastic.co/guide/index.html)
- [Go Fiber Documentation](https://docs.gofiber.io/)

## 🏅 Winning Strategy

### What Makes Your Project Stand Out

1. **Real-World Application**: Solves actual developer pain points
2. **Seamless Integration**: Elastic + Google Cloud working together
3. **User Experience**: Natural language interface
4. **Innovation**: Hybrid search + generative AI
5. **Scalability**: Production-ready architecture

### Key Points for Judges

- **Technical Excellence**: Clean code, proper error handling
- **Innovation**: Novel use of AI for log analysis
- **Practical Impact**: Saves developers hours of manual log analysis
- **Completeness**: Full-stack solution with documentation

---

## 🎉 You're Ready!

You now have a fully functional AI-powered API monitoring system that:
- ✅ Integrates Elasticsearch with Google Cloud Vertex AI
- ✅ Provides conversational log analysis
- ✅ Detects anomalies automatically
- ✅ Suggests optimizations
- ✅ Scales to millions of logs

**Good luck with the hackathon! 🚀**

For detailed documentation, see [VERTEX_AI_SETUP.md](./VERTEX_AI_SETUP.md)
