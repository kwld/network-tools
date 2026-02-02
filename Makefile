.PHONY: help build run clean docker-build docker-run docker-clean test

# Default target
help:
	@echo "Network Tools - Makefile Commands"
	@echo "=================================="
	@echo "  make build         - Build the scanner binary"
	@echo "  make run           - Run the scanner locally"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"
	@echo "  make docker-clean  - Clean Docker resources"
	@echo "  make test          - Run tests"
	@echo "  make deps          - Download dependencies"

# Build the scanner binary
build:
	@echo "Building scanner..."
	go build -o bin/scanner ./cmd/scanner

# Run the scanner locally
run: build
	@echo "Running scanner..."
	@mkdir -p output
	./bin/scanner

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf output/

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker-compose build

# Run with Docker Compose
docker-run:
	@echo "Running with Docker Compose..."
	@mkdir -p output
	docker-compose up

# Clean Docker resources
docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker-compose rm -f

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	go vet ./...
