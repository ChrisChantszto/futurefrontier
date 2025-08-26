# Onetake CorpSite Backend

A Go + Fiber backend with MongoDB storage.

This guide shows how to run MongoDB via Docker and smoke test the API locally using Postman or curl.

## Prerequisites

- Go 1.22+
- Docker Desktop (for MongoDB container)
- Postman (recommended) or curl

## 1) Start MongoDB with Docker

Run MongoDB locally on port 27017:

```powershell
# Pull and start MongoDB 7
Docker pull mongo:7.0
Docker run --name mongo-local -p 27017:27017 -d mongo:7.0

# Verify the port is open
Test-NetConnection localhost -Port 27017
```

If Docker says the container name exists, remove or reuse it:

```powershell
Docker start mongo-local  # if it already exists
# or
Docker rm -f mongo-local; Docker run --name mongo-local -p 27017:27017 -d mongo:7.0
```

## 2) Configure environment

Create a `.env` at the project root (already supported in `main.go`):

```
PORT=8080
MONGODB_URI=mongodb://localhost:27017
DB_NAME=onetake
JWT_SECRET=dev-secret
JWT_REFRESH_SECRET=dev-refresh-secret
SECURE_COOKIES=false
```

Tip: Never commit your real secrets. `.env` is gitignored by default.

## 3) Install deps and run

From the project root:

```powershell
# Optional: let Go fetch a matching toolchain
$env:GOTOOLCHAIN = 'auto'

Go mod tidy
Go run .
```

You should see logs ending with:

```
listening on :8080
```

Health check:

```powershell
Curl http://localhost:8080/api/healthz
```

Expect: `{ "ok": true }`

## 4) Postman smoke test

Create a Postman environment:

- baseUrl = http://localhost:8080/api

Use Body = Raw + JSON for POST/PATCH. Content-Type = application/json.

1) Init first user (public)

- POST {{baseUrl}}/auth/init-first-user
- Body:
```json
{
  "Email": "admin@example.com",
  "Password": "P@ssw0rd!"
}
```

2) Login (public; sets cookies)

- POST {{baseUrl}}/auth/login
- Body:
```json
{
  "Email": "admin@example.com",
  "Password": "P@ssw0rd!"
}
```

3) Get settings (public)

- GET {{baseUrl}}/settings

4) Update settings (protected; requires login cookies)

- PATCH {{baseUrl}}/settings
- Body:
```json
{
  "theme": "light",
  "defaultLocale": "en"
}
```

5) Pages list (public)

- GET {{baseUrl}}/pages?limit=20

6) Create page (protected)

- POST {{baseUrl}}/pages
- Body:
```json
{
  "identifier": "home",
  "enabled": true,
  "meta": {
    "title": { "en": "Home" },
    "description": { "en": "Welcome home" },
    "keywords": ["corp","site"]
  },
  "grouping": { "label": { "en": "Default" } },
  "header": { "noticeBar": { "enabled": false } },
  "sections": []
}
```

7) Get one page (public)

- GET {{baseUrl}}/pages/home

8) Export page (public)

- GET {{baseUrl}}/pages/home/export

9) Patch page (protected)

- PATCH {{baseUrl}}/pages/home
- Body:
```json
{
  "meta": { "title": { "en": "Home Updated" } },
  "enabled": true
}
```

10) Import pages (protected)

- POST {{baseUrl}}/pages/import
- Body:
```json
{
  "pages": [
    {
      "identifier": "contact",
      "enabled": true,
      "meta": { "title": { "en": "Contact" }, "description": { "en": "" }, "keywords": [] },
      "grouping": { "label": { "en": "Default" } },
      "header": { "noticeBar": { "enabled": false } },
      "sections": []
    }
  ]
}
```

11) Delete page (protected)

- DELETE {{baseUrl}}/pages/home

Note: Protected routes require cookies set by the login response (Postman stores them automatically).

## 5) Curl examples (PowerShell)

```powershell
# Init first user
Curl -X POST http://localhost:8080/api/auth/init-first-user `
  -H "Content-Type: application/json" `
  -d '{"Email":"admin@example.com","Password":"P@ssw0rd!"}'

# Login and store cookies
Curl -i -c cookie.txt -b cookie.txt -X POST http://localhost:8080/api/auth/login `
  -H "Content-Type: application/json" `
  -d '{"Email":"admin@example.com","Password":"P@ssw0rd!"}'

# Create a page (protected)
Curl -i -c cookie.txt -b cookie.txt -X POST http://localhost:8080/api/pages `
  -H "Content-Type: application/json" `
  -d '{"identifier":"home","enabled":true,"meta":{"title":{"en":"Home"},"description":{"en":""},"keywords":[]},"grouping":{"label":{"en":"Default"}},"header":{"noticeBar":{"enabled":false}},"sections":[]}'

# List pages (public)
Curl http://localhost:8080/api/pages?limit=20
```

## Troubleshooting

- "mongo connect failed":
  - Ensure the Docker container is running and 27017 is open.
  - `Test-NetConnection localhost -Port 27017` should be True.
- "missing access token":
  - You’re calling a protected route before logging in, or you’re missing the cookies. First login, then retry.
- Change PORT in `.env` if 8080 is taken.
- If using MongoDB Atlas:
  - Set `MONGODB_URI` to your SRV connection string.
  - Whitelist your IP in Atlas Network Access.
  - URL-encode special characters in the password.

## Project layout

- `main.go`: server entrypoint and middleware
- `internal/config/`: config loader
- `internal/service/`: business logic (pages, auth, settings)
- `internal/transport/http/handlers/`: HTTP handlers
- `internal/transport/http/middleware/`: auth middleware
- `internal/models/`: data models

---

If you need a ready-to-import Postman collection, open an issue or ask in chat and we’ll attach one.
