# API Testing Guide - OneTake CorpSite Backend

## Overview
This guide provides step-by-step instructions for testing the OneTake CorpSite Backend API using Postman or Apifox. The system includes centralized logging with Elasticsearch for monitoring all API requests.

## Prerequisites
- Postman or Apifox installed
- Backend server running on `http://localhost:8080`
- Elasticsearch logging server at `http://server.onetakesolutions.com.hk:9201`

## Setup Instructions

### 1. Environment Configuration

**For Postman:**
1. Create a new Environment called "OneTake Backend"
2. Add these variables:
   - `baseUrl`: `http://localhost:8080`
   - `elasticsearchUrl`: `http://server.onetakesolutions.com.hk:9201`

**For Apifox:**
1. Create a new Environment called "OneTake Backend"
2. Add these variables:
   - `baseUrl`: `http://localhost:8080`
   - `elasticsearchUrl`: `http://server.onetakesolutions.com.hk:9201`
3. Enable "Auto-save cookies" in Environment settings

## API Testing Steps

### Step 1: Test Basic Health Check

**Request:**
```
Method: GET
URL: {{baseUrl}}/api/healthz
```

**Expected Response:**
```json
{
  "ok": true
}
```

### Step 2: Request OTP for Login

**Request:**
```
Method: POST
URL: {{baseUrl}}/api/auth/otp/request
Headers: Content-Type: application/json
Body (JSON):
{
  "email": "chantszto.chris@gmail.com",
  "purpose": "login"
}
```

**Expected Response:**
```json
{
  "requestId": "abc123xyz",
  "message": "If an account exists, we've sent a code to your email"
}
```

**Important:** Save the `requestId` from the response for the next step.

### Step 3: Verify OTP and Login

**Request:**
```
Method: POST
URL: {{baseUrl}}/api/auth/otp/verify
Headers: Content-Type: application/json
Body (JSON):
{
  "email": "chantszto.chris@gmail.com",
  "code": "123456",
  "requestId": "abc123xyz",
  "purpose": "login"
}
```

**Note:** Replace `"123456"` with the actual OTP code (check server logs for the code during development).

**Expected Response:**
```json
{
  "ok": true,
  "message": "Login successful"
}
```

**Important:** This sets authentication cookies automatically.

### Step 4: Test Protected Endpoints

**Request:**
```
Method: GET
URL: {{baseUrl}}/api/settings/public
```

**Expected Response:**
```json
[
  // Settings data or empty array
]
```

### Step 5: Test Logging Statistics

**Request:**
```
Method: GET
URL: {{baseUrl}}/health/logging
```

**Expected Response:**
```json
{
  "logging_stats": {
    "failed_logs": 0,
    "processed_logs": 5,
    "queued_logs": 0
  }
}
```

## Verifying Elasticsearch Logs

### Check if Logs are Being Captured

**Request:**
```
Method: GET
URL: {{elasticsearchUrl}}/api-logs-*/_search?pretty&size=10
```

**Expected Response:**
```json
{
  "took": 5,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 7,
      "relation": "eq"
    },
    "hits": [
      {
        "_source": {
          "project_id": "onetake-corpsite-backend",
          "method": "GET",
          "path": "/api/healthz",
          "status_code": 200,
          "latency_ms": 15,
          "client_ip": "127.0.0.1",
          "timestamp": "2025-09-08T08:17:19.262Z"
        }
      }
    ]
  }
}
```

### Search for Specific API Calls

**Request:**
```
Method: GET
URL: {{elasticsearchUrl}}/api-logs-*/_search?pretty
Headers: Content-Type: application/json
Body (JSON):
{
  "query": {
    "bool": {
      "must": [
        {"match": {"path": "/api/auth/otp/request"}},
        {"range": {"timestamp": {"gte": "now-1h"}}}
      ]
    }
  },
  "sort": [{"timestamp": {"order": "desc"}}],
  "size": 5
}
```

## Complete Test Sequence

1. **Health Check** → Should return `{"ok": true}`
2. **Request OTP** → Should return `requestId`
3. **Verify OTP** → Should return login success and set cookies
4. **Test Protected Endpoint** → Should return data (not 401 error)
5. **Check Logging Stats** → Should show processed logs count
6. **Verify Elasticsearch** → Should show all API calls logged

## Expected Log Data

Each API call should generate a log entry containing:
- **timestamp**: When the request was made
- **method**: HTTP method (GET, POST, etc.)
- **path**: API endpoint called
- **status_code**: HTTP response code
- **latency_ms**: Response time in milliseconds
- **client_ip**: IP address of the client
- **user_agent**: Client application identifier
- **request_id**: Unique identifier for request tracking

## Troubleshooting

### Common Issues:

1. **401 Unauthorized on protected endpoints**
   - Ensure cookies are enabled in Postman/Apifox
   - Complete the OTP login flow first

2. **OTP verification fails**
   - Check server logs for the actual OTP code
   - Ensure the `requestId` matches the one from OTP request

3. **No logs in Elasticsearch**
   - Check if backend server shows "Successfully connected to Elasticsearch"
   - Verify `/health/logging` shows `processed_logs > 0`

4. **Connection timeouts**
   - Ensure backend server is running on port 8080
   - Check if Elasticsearch server is accessible

## Success Criteria

✅ All API endpoints respond correctly
✅ OTP authentication flow works
✅ Protected endpoints accessible after login
✅ Logging statistics show processed requests
✅ Elasticsearch contains detailed request logs
✅ Log entries include all required fields (timestamp, method, path, status, etc.)

## Contact

For technical issues or questions, contact the development team.
