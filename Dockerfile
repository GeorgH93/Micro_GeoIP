# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o micro_geoip .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Create data directory for GeoIP database
RUN mkdir -p /data

# Copy the binary from builder stage
COPY --from=builder /app/micro_geoip .

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Change ownership of the data directory
RUN chown -R appuser:appgroup /data

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Environment variables
ENV PORT=8080
ENV HOST=0.0.0.0
ENV GEOIP_DB_PATH=/data/GeoLite2-Country.mmdb

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./micro_geoip"]