# Build stage
FROM golang:1.25-alpine AS builder

# Install git for go modules that require it
RUN apk add --no-cache git

WORKDIR /app


# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -ldflags="-w -s" -o bot ./cmd/bot

# Run stage - use distroless or scratch as base for minimal attack surface
FROM alpine:3.19

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install ca-certificates if needed for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bot .

# Switch to non-root user
USER appuser

CMD ["./bot"]
