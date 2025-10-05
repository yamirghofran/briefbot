# Multi-stage build for Go backend
# Stage 1: Build the application
FROM golang:1.25.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Stage 2: Runtime image
FROM alpine:latest

# Install runtime dependencies including ffmpeg for audio processing
RUN apk add --no-cache ca-certificates postgresql-client bash ffmpeg

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/server /app/server

# Copy goose binary
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migration files
COPY sql/migrations /app/sql/migrations

# Copy scripts
COPY scripts /app/scripts

# Make scripts executable
RUN chmod +x /app/scripts/*.sh

# Create migration script
RUN echo '#!/bin/sh' > /app/migrate.sh && \
    echo 'cd /app' >> /app/migrate.sh && \
    echo 'echo "Running migrations with goose..."' >> /app/migrate.sh && \
    echo 'goose up' >> /app/migrate.sh && \
    echo 'echo "Migrations completed successfully"' >> /app/migrate.sh && \
    chmod +x /app/migrate.sh

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/server"]
