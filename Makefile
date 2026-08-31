# HydraForge Development & Build Automation Makefile

.PHONY: dev build test clean help

# Default target
all: dev

# Run HydraForge in Live Development Mode (Go Control Plane + React Frontend + Python Worker)
dev:
	@echo "🚀 Starting HydraForge Control Plane & Training Studio..."
	@echo "💡 Web UI available on http://localhost:8081"
	go run ./cmd/hydraforge

# Build production binaries and web assets
build:
	@echo "📦 Building HydraForge Go Control Plane binary..."
	mkdir -p bin
	go build -o bin/hydraforge ./cmd/hydraforge

# Run test suite
test:
	@echo "🧪 Running Go unit tests..."
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/ build/
