# Clockify MCP Server Makefile

# Variables
BINARY_NAME=clockify-mcp
GO=go
GOFLAGS=-v
BUILD_DIR=.

# Default target
.DEFAULT_GOAL := build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) ./...

# Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Run go fmt
.PHONY: fmt
fmt:
	@echo "Running go fmt..."
	$(GO) fmt ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Clean complete"

# Run all checks (fmt, vet, test, build)
.PHONY: check
check: fmt vet test build
	@echo "All checks passed!"

# Install dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Run the server (requires CLOCKIFY_API_KEY env var)
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build   - Build the binary (default)"
	@echo "  test    - Run tests"
	@echo "  vet     - Run go vet"
	@echo "  fmt     - Run go fmt"
	@echo "  clean   - Remove build artifacts"
	@echo "  check   - Run fmt, vet, test, and build"
	@echo "  deps    - Download and tidy dependencies"
	@echo "  run     - Build and run the server"
	@echo "  help    - Show this help message"
