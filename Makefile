.PHONY: help test test-verbose test-coverage build clean fmt vet lint tidy install

# Default target
help:
	@echo "Available targets:"
	@echo "  test           - Run tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  build          - Build the project"
	@echo "  clean          - Clean build artifacts and cache"
	@echo "  fmt            - Format code with gofmt"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run golangci-lint (if installed)"
	@echo "  tidy           - Tidy go modules"
	@echo "  install        - Download dependencies"

# Run tests
test:
	cd v3 && go test

# Run tests with verbose output
test-verbose:
	cd v3 && go test -v

# Run tests with coverage
test-coverage:
	cd v3 && go test -coverprofile=coverage.out
	cd v3 && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: v3/coverage.html"

# Build (validate the code compiles)
build:
	cd v3 && go build

# Clean build artifacts and cache
clean:
	cd v3 && go clean
	cd v3 && rm -f coverage.out coverage.html

# Format code
fmt:
	cd v3 && go fmt ./...

# Run go vet
vet:
	cd v3 && go vet ./...

# Run golangci-lint if available
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	cd v3 && golangci-lint run

# Tidy dependencies
tidy:
	cd v3 && go mod tidy

# Download dependencies
install:
	cd v3 && go mod download
