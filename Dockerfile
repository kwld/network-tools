# Multi-stage build for network-tools scanner
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scanner ./cmd/scanner

# Final stage
FROM alpine:latest

# Install runtime dependencies (GraphViz for diagram generation)
RUN apk add --no-cache ca-certificates graphviz ttf-dejavu openssh-client

# Create directories
RUN mkdir -p /config /output

# Copy binary from builder
COPY --from=builder /app/scanner /usr/local/bin/scanner

# Set working directory
WORKDIR /app

# Volume mounts
VOLUME ["/config", "/output"]

# Run the scanner
CMD ["scanner"]
