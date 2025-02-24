# -------------------------------------------
# Stage 1: Build the Go application
# -------------------------------------------
FROM golang:1.24.0-alpine3.21 AS builder

# Set build arguments for versioning and security
ARG APP_VERSION="1.0.0"
ARG BUILD_TIMESTAMP
ARG COMMIT_HASH

# Set safe defaults for Go environment
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPRIVATE="github.com/GabrielEValenzuela/*"

WORKDIR /app

# Copy dependency files first for optimal caching
COPY go.mod go.sum ./
RUN go mod download -x && go mod verify

# Copy source code
COPY ./src ./src

# Build with security flags and version information
RUN go build -v -trimpath \
    -ldflags="-w -s \
        -X main.Version=${APP_VERSION} \
        -X main.BuildTime=${BUILD_TIMESTAMP} \
        -X main.CommitHash=${COMMIT_HASH}" \
    -o /app/bin/main ./src/cmd/main.go

# -------------------------------------------
# Stage 2: Minimal production image
# -------------------------------------------
FROM alpine:3.19

# Install security updates and required dependencies
RUN apk --no-cache upgrade && \
    apk add --no-cache \
        ca-certificates \
        tzdata \
        dumb-init \
        wget

# Create non-root user
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup -h /app

WORKDIR /app

# Copy necessary files from the builder
COPY --from=builder --chown=appuser:appgroup /app/bin/main /app/main
COPY ./config.yml /app/config.yml

# Set permissions and ownership
RUN chmod +x /app/main && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the service port
EXPOSE 9000

# Health check to ensure the service is responding
HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:9000/ || exit 1

# Use init system for proper signal handling
ENTRYPOINT ["dumb-init", "--"]

# Correctly define CMD with executable and arguments as separate array elements
CMD ["/app/main", "-config", "/app/config.yml"]
