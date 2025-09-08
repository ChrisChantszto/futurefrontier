# Middleware Logger Implementation

## Overview

This implementation provides a comprehensive middleware logger for your Go Fiber project that logs all API requests to Elasticsearch with MongoDB as a backup system. The system follows your requirements for centralized logging with data resilience.

## Architecture

```
API Request
    ↓
Middleware captures log data
    ↓ (async queue)
Background worker goroutines
    ↓──────→ Elasticsearch write fails ─────→ MongoDB backup
    ↓ (success)
Complete

Recovery worker checks MongoDB backup
    ↓
Retry write to Elasticsearch
    ↓ (success)
Delete from MongoDB backup
```

## Features Implemented

### 1. **Complete Log Data Capture**
- **ProjectID**: Configurable project identifier
- **Timestamp**: Precise request time
- **Request Method**: GET, POST, PUT, DELETE, etc.
- **Path & Route**: API endpoint information
- **Status Code**: HTTP response status
- **Latency**: Request processing time in milliseconds
- **Client IP & Port**: Real client IP (handles proxies/load balancers)
- **User-Agent**: Client browser/application info
- **Request ID**: Unique identifier for distributed tracing
- **Referer**: Request source page
- **Query Parameters**: URL query string data
- **Headers**: HTTP headers (sensitive ones excluded)
- **Request/Response Body**: Optional, with size limits
- **Bytes Sent/Received**: Data transfer metrics
- **Error Information**: Detailed error messages and types

### 2. **Dual Index System**
- **API Logs Index**: `api-logs-YYYY-MM` (monthly rotation)
- **Error Logs Index**: `error-logs-YYYY-MM` (monthly rotation)

### 3. **Async Processing with Queue**
- Non-blocking logging using Go channels
- Configurable queue size (default: 1000)
- Multiple worker goroutines (default: 3)
- Graceful handling of queue overflow

### 4. **Backup & Recovery System**
- MongoDB backup when Elasticsearch fails
- Automatic retry mechanism every 30 seconds
- TTL indexes for automatic cleanup (30 days)
- Retry count tracking with limits

### 5. **Security Features**
- Sensitive header filtering (Authorization, Cookie, etc.)
- Sensitive path body exclusion
- Configurable body size limits
- Generic error responses to prevent enumeration

### 6. **Local Development Mode**
- `LOG_LOCAL_MODE=true` for development
- Logs to structured zap logger instead of Elasticsearch
- No external dependencies needed for testing

## Configuration

### Environment Variables

```bash
# Logging Configuration
LOG_PROJECT_ID=onetake-corpsite-backend
ELASTICSEARCH_URL=http://localhost:9200
ELASTICSEARCH_USER=
ELASTICSEARCH_PASS=
LOG_ENABLE_REQUEST_BODY=false
LOG_ENABLE_RESPONSE_BODY=false
LOG_ENABLE_HEADERS=false
LOG_MAX_BODY_SIZE=1024
LOG_QUEUE_SIZE=1000
LOG_WORKER_COUNT=3
LOG_RETRY_ATTEMPTS=3
LOG_RETRY_INTERVAL=30
LOG_LOCAL_MODE=true
```

### Production Elasticsearch Setup

When your manager sets up Elasticsearch with Docker, update these settings:

```bash
LOG_LOCAL_MODE=false
ELASTICSEARCH_URL=http://your-elasticsearch-server:9200
ELASTICSEARCH_USER=your_username
ELASTICSEARCH_PASS=your_password
LOG_ENABLE_REQUEST_BODY=true
LOG_ENABLE_RESPONSE_BODY=true
LOG_ENABLE_HEADERS=true
```

## API Endpoints

### Health Check
- `GET /health/logging` - Returns logging statistics

Example response:
```json
{
  "logging_stats": {
    "processed_logs": 1250,
    "failed_logs": 5,
    "queued_logs": 0
  }
}
```

## Files Created/Modified

### New Files
1. `internal/models/log.go` - Log data structures
2. `internal/service/elasticsearch.go` - Elasticsearch client
3. `internal/service/logger.go` - Async logging service
4. `internal/service/mongo_backup.go` - MongoDB backup service
5. `internal/middleware/logger.go` - Fiber middleware

### Modified Files
1. `internal/config/config.go` - Added logging configuration
2. `main.go` - Integrated middleware and graceful shutdown
3. `.env` - Added logging environment variables
4. `go.mod` - Added Elasticsearch dependency

## Usage Examples

### Basic Request Logging
Every API request is automatically logged with basic information:

```json
{
  "project_id": "onetake-corpsite-backend",
  "timestamp": "2025-09-08T10:30:00Z",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/auth/otp/request",
  "status_code": 200,
  "latency_ms": 45,
  "client_ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "bytes_sent": 156,
  "bytes_received": 89
}
```

### Error Logging
5xx errors automatically create additional error log entries:

```json
{
  "project_id": "onetake-corpsite-backend",
  "timestamp": "2025-09-08T10:30:00Z",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/users",
  "status_code": 500,
  "error_message": "database connection failed",
  "error_type": "server_error"
}
```

## Testing

### 1. Start the Server
```bash
go run main.go
```

### 2. Make Test Requests
```bash
# Test API logging
curl -X GET http://localhost:8080/health/logging

# Test with request body
curl -X POST http://localhost:8080/auth/otp/request \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com"}'
```

### 3. Check Logs
In local mode, logs will appear in the console output with structured JSON data.

## Production Deployment

### 1. Update Environment
```bash
LOG_LOCAL_MODE=false
ELASTICSEARCH_URL=http://your-es-server:9200
```

### 2. Monitor Performance
- Check `/health/logging` endpoint regularly
- Monitor MongoDB backup collection sizes
- Watch Elasticsearch index growth

### 3. Maintenance
- Elasticsearch indexes rotate monthly automatically
- MongoDB backup has 30-day TTL
- Recovery worker runs every 30 seconds

## Troubleshooting

### Common Issues

1. **Queue Full Warnings**
   - Increase `LOG_QUEUE_SIZE`
   - Add more workers with `LOG_WORKER_COUNT`

2. **Elasticsearch Connection Fails**
   - Logs automatically backup to MongoDB
   - Check `ELASTICSEARCH_URL` configuration
   - Verify network connectivity

3. **High Memory Usage**
   - Reduce `LOG_MAX_BODY_SIZE`
   - Disable body logging: `LOG_ENABLE_REQUEST_BODY=false`

4. **MongoDB Backup Growing**
   - Check Elasticsearch connectivity
   - Monitor recovery worker logs
   - Verify retry mechanism is working

## Performance Considerations

- **Memory**: ~1MB per 1000 queued logs
- **CPU**: Minimal overhead with async processing
- **Network**: Batched writes to Elasticsearch
- **Storage**: Monthly index rotation prevents bloat

## Security Notes

- Sensitive headers are automatically filtered
- Request/Response bodies are optional and size-limited
- Authentication endpoints exclude body logging by default
- Request IDs enable secure distributed tracing
