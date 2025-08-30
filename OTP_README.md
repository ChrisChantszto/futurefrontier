# OTP Authentication System

This document describes the One-Time Password (OTP) authentication system implemented for the Onetake backend.

## Overview

The OTP system provides secure, email-based authentication with the following features:
- 6-digit OTP codes generated using crypto/rand
- HMAC-SHA256 for secure code storage
- Configurable TTL, max attempts, and resend cooldown
- Support for multiple purposes (login, forgetPassword)
- Dual SMTP configuration (custom + company fallback)
- MongoDB-based storage with automatic cleanup

## Environment Variables

### Required OTP Configuration
```bash
APP_HMAC_SECRET=ZGV2LXNlY3JldC1obWFjLWtleS0zMi1ieXRlcy1taW5pbXVt  # Base64/hex encoded, ≥32 bytes
OTP_TTL_MINUTES=5                    # OTP expiration time
OTP_MAX_ATTEMPTS=5                   # Maximum verification attempts
OTP_RESEND_COOLDOWN_SECONDS=60       # Cooldown between resend requests
```

### SMTP Configuration (Optional - Custom)
```bash
SMTP_HOST=                           # Custom SMTP host
SMTP_PORT=465                        # SMTP port (default: 465)
SMTP_USER=                           # SMTP username
SMTP_PASS=                           # SMTP password
SMTP_FROM_EMAIL=                     # From email address
SMTP_FROM_NAME=                      # From display name
SMTP_REPLY_TO=                       # Reply-to address (optional)
SMTP_BCC=                           # BCC address (optional)
```

### Company SMTP (Fallback)
```bash
ONETAKE_SMTP_HOST=smtpdm-ap-southeast-1.aliyun.com
ONETAKE_SMTP_PORT=465
ONETAKE_SMTP_USER=system@onetakesolutions.com.hk
ONETAKE_SMTP_PASS=OneTake07102022
ONETAKE_SMTP_FROM_EMAIL=system@onetakesolutions.com.hk
ONETAKE_SMTP_FROM_NAME=OTS system
ONETAKE_SMTP_REPLY_TO=info@onetakesolutions.com.hk
ONETAKE_SMTP_BCC=tim.ho@onetakesolutions.com.hk
ONETAKE_SMTP_TLS=true
```

## API Endpoints

### 1. Request OTP Code
**POST** `/api/auth/otp/request`

Request a new OTP code to be sent via email.

**Request Body:**
```json
{
  "email": "user@example.com",
  "purpose": "login"  // or "forgetPassword"
}
```

**Response:**
```json
{
  "requestId": "uuid-string",
  "message": "If an account exists, we've sent a code to your email"
}
```

**Features:**
- Enforces resend cooldown to prevent spam
- Normalizes email (trim + lowercase)
- Generates cryptographically secure 6-digit code
- Stores HMAC of code (never stores plaintext)
- Automatically marks previous active OTPs as consumed

### 2. Verify OTP Code
**POST** `/api/auth/otp/verify`

Verify an OTP code and optionally authenticate the user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "code": "123456",
  "requestId": "uuid-string",  // optional but recommended
  "purpose": "login"
}
```

**Response (Login Success):**
```json
{
  "ok": true,
  "message": "Login successful"
}
```

**Response (Other Purposes):**
```json
{
  "ok": true,
  "message": "OTP verified successfully",
  "requestId": "uuid-string"
}
```

**Features:**
- Constant-time HMAC comparison to prevent timing attacks
- Increments attempt counter on failure
- Auto-consumes OTP after max attempts reached
- For login purpose: creates JWT tokens and sets HTTP-only cookies
- Marks OTP as consumed after successful verification

### 3. Resend OTP Code
**POST** `/api/auth/otp/resend`

Request a new OTP code, replacing any existing active OTP.

**Request Body:**
```json
{
  "email": "user@example.com",
  "purpose": "login"
}
```

**Response:**
```json
{
  "requestId": "uuid-string",
  "message": "If an account exists, we've sent a new code to your email"
}
```

## MongoDB Schema

### Collection: `otp_codes`

```javascript
{
  _id: ObjectId,
  email: String,           // Normalized (lowercase, trimmed)
  requestId: String,       // UUID for request tracking
  purpose: String,         // "login" | "forgetPassword"
  codeHmac: String,        // HMAC-SHA256 of the OTP code
  expiresAt: Date,         // UTC expiration timestamp
  attempts: Number,        // Current attempt count
  maxAttempts: Number,     // Maximum allowed attempts
  createdAt: Date,         // UTC creation timestamp
  consumedAt: Date | null, // UTC consumption timestamp
  lastSentAt: Date         // UTC last sent timestamp
}
```

### Indexes

1. **Compound Index**: `email + purpose + consumedAt + expiresAt` (for finding active OTPs)
2. **Unique Index**: `requestId` (for exact lookups)
3. **TTL Index**: `expiresAt` (automatic cleanup 1 hour after expiry)
4. **Descending Index**: `createdAt` (for sorting)
5. **Compound Index**: `email + purpose + lastSentAt` (for cooldown checks)

## Security Features

### Code Generation
- Uses `crypto/rand` for cryptographically secure random number generation
- 6-digit codes provide 1,000,000 possible combinations
- Codes are never stored in plaintext

### HMAC Protection
- Codes are hashed using HMAC-SHA256 with a secret key
- Constant-time comparison prevents timing attacks
- Secret key should be ≥32 bytes and base64/hex encoded

### Rate Limiting
- Configurable resend cooldown prevents spam
- Maximum attempt limit prevents brute force
- Automatic OTP consumption after max attempts

### Email Security
- Generic response messages prevent user enumeration
- Support for BCC to internal monitoring addresses
- TLS encryption for SMTP connections (port 465)

## Usage Examples

### Basic Login Flow
```bash
# 1. Request OTP
curl -X POST http://localhost:8081/api/auth/otp/request \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "purpose": "login"}'

# 2. Verify OTP (user enters code from email)
curl -X POST http://localhost:8081/api/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "code": "123456", "purpose": "login"}'
```

### Password Reset Flow
```bash
# 1. Request OTP for password reset
curl -X POST http://localhost:8081/api/auth/otp/request \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "purpose": "forgetPassword"}'

# 2. Verify OTP
curl -X POST http://localhost:8081/api/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "code": "123456", "purpose": "forgetPassword"}'

# 3. Use the verified session to reset password
```

## Error Handling

### Common Error Responses

**400 Bad Request:**
- Invalid JSON body
- Missing required fields
- Invalid purpose value
- Invalid email format

**401 Unauthorized:**
- Invalid or expired OTP code
- Maximum attempts exceeded
- User not found (for login purpose)

**429 Too Many Requests:**
- Resend cooldown period active

### Generic Responses
All endpoints return generic success messages to prevent user enumeration:
- "If an account exists, we've sent a code to your email"
- "If an account exists, we've sent a new code to your email"

## Configuration Notes

### SMTP Priority
1. If `ALLOW_CUSTOM_SMTP=true` and custom SMTP is configured → use custom SMTP
2. Otherwise → use company SMTP (Aliyun)

### Purpose Values
- `"login"`: For user authentication
- `"forgetPassword"`: For password reset flows

Both purposes use the same enum type and can be extended in the future.

## Monitoring and Logging

The system logs important events:
- OTP request failures
- OTP verification failures
- SMTP sending errors
- Database operation errors

All logs include relevant context (email, request ID) for debugging while maintaining security.
