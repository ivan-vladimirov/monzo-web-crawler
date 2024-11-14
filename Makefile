# Variables
APP_NAME := monzo-web-crawler
GO := go
SRC := $(shell find . -name '*.go')
PKGS := $(shell $(GO) list ./...)

# Default target
all: build

# Build the application
build: $(SRC)
	@echo "Building $(APP_NAME)..."
	$(GO) build -o $(APP_NAME) ./cmd/main.go
	@echo "Build complete: $(APP_NAME)"

# Run the application with example flags
run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME) -url=https://monzo.com -max-depth=3 -delay=100ms

# Run tests with coverage
test:
	@echo "Running tests..."
	$(GO) test -v -cover ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Lint code (optional, requires golangci-lint)
lint:
	@echo "Running lint checks..."
	@golangci-lint run ./... || echo "Linting skipped: install golangci-lint for lint checks."

# Clean up generated files
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
	@echo "Cleanup complete."

# Display help
help:
	@echo "Monzo Coding Challenge Makefile"
	@echo "Available targets:"
	@echo "  build       Build the crawler application."
	@echo "  run         Build and run the crawler with example flags."
	@echo "  test        Run tests with coverage."
	@echo "  fmt         Format the code using 'go fmt'."
	@echo "  lint        Run lint checks (optional: requires golangci-lint)."
	@echo "  clean       Clean up generated files."
	@echo "  help        Display this help message."

# Phony targets
.PHONY: all build run test fmt lint clean help
