# Build stage
FROM golang:1.24-alpine3.22 AS builder

# Install build dependencies
RUN apk add --no-cache git make protoc protobuf-dev

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy proto file and generate code
COPY sevalet.proto ./
RUN mkdir -p pb && \
    protoc --go_out=pb --go_opt=paths=source_relative \
           --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
           sevalet.proto

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags="-w -s -X 'github.com/zinrai/sevalet/cmd.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)'" \
    -o sevalet .

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

# Create non-root user and group
RUN addgroup -g 1000 -S sevalet && \
    adduser -u 1000 -S sevalet -G sevalet

# Create necessary directories
RUN mkdir -p /etc/sevalet /var/run && \
    chown sevalet:sevalet /var/run

# Copy binary from builder
COPY --from=builder /build/sevalet /usr/local/bin/sevalet
RUN chmod 755 /usr/local/bin/sevalet

# Copy default API configuration
COPY --chown=sevalet:sevalet configs/api.yaml /etc/sevalet/api.yaml

# Switch to non-root user
USER sevalet

# Expose HTTP port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set default command for API mode
ENTRYPOINT ["/usr/local/bin/sevalet"]
CMD ["api", "--config", "/etc/sevalet/api.yaml"]
