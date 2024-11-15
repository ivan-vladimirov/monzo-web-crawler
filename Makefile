# Variables
APP_NAME := monzo-web-crawler
GO := go
PKG := ./...
TEST_FLAGS := -v 

# Default values for custom flags
URL ?= http://monzo.com
MAX_DEPTH ?= 3
DELAY ?= 100ms
FILENAME ?= crawled.json

.PHONY: build run test clean help

# Build the monzo-web-crawler binary
build:
	@echo "Building the$(APP_NAME) binary..."
	go build -o monzo-web-crawler ./cmd

# Run the monzo-web-crawler with custom flags
run: build
	@echo "Running $(APP_NAME) with:"
	@echo "  URL: $(URL)"
	@echo "  Max Depth: $(MAX_DEPTH)"
	@echo "  Delay: $(DELAY)"
	@./monzo-web-crawler -url=$(URL) -max-depth=$(MAX_DEPTH) -delay=$(DELAY) -
# Run tests
test:
	@echo "Running tests..."
	go test $(PKG) $(TEST_FLAGS)
	
# Clean target to remove build artifacts
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf $(APP_NAME) logs/ coverage.out

# Help target to display usage
help:
	@echo "Available targets:"
	@echo "  build     - Build the $(APP_NAME) binary"
	@echo "  run       - Build and run the $(APP_NAME) (customizable with URL, MAX_DEPTH, and DELAY)"
	@echo "             e.g., make run URL=http://example.com MAX_DEPTH=5 DELAY=200ms"
	@echo "  clean     - Clean up build artifacts and logs"
	@echo "  help      - Show this help message"
