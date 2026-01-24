# Build stage
FROM golang:1.25-alpine AS builder

# Install git for go modules that require it
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files and download dependencies first (for better layer caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-w -s" -o bot ./cmd/bot

# Run stage - use distroless or scratch as base for minimal attack surface
FROM alpine:3.19

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install ca-certificates if needed for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bot .

# Copy goose binary from builder stage
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migrations directory (adjust path if different)
COPY internal/migrations/ ./internal/migrations/

# Switch to non-root user
USER appuser

# Run migrations first, then start the bot
# Assumes DATABASE_URL is set in Dokku environment variables
CMD ["sh", "-c", "goose up && ./bot"]
