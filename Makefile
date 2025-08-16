# Makefile for Trains CLI

# Variables
APP_NAME = trains
CMD_DIR = ./cmd/trains
BUILD_DIR = .

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

# Run tests
.PHONY: test
test:
	go test ./...

# Run tests with verbose output
.PHONY: test-verbose
test-verbose:
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BUILD_DIR)/$(APP_NAME)

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	go vet ./...

# Download dependencies
.PHONY: deps
deps:
	go mod download

# Update dependencies
.PHONY: deps-update
deps-update:
	go get -u ./...
	go mod tidy

# Install the binary
.PHONY: install
install: build
	cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  test-verbose - Run tests with verbose output"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  deps         - Download dependencies"
	@echo "  deps-update  - Update dependencies"
	@echo "  install      - Install binary to /usr/local/bin"
	@echo "  help         - Show this help"