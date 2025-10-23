# Build stage: use Go 1.24 to satisfy genai requirements
FROM golang:1.24 AS build
WORKDIR /app

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Run stage: minimal image
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/main /app/main

# Railway provides PORT env var
ENV PORT=8080
EXPOSE 8080

CMD ["/app/main"]