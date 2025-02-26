.PHONY: build run test clean lint fix fmt help tidy embedmd-check embedmd

BINARY_NAME=demo-web-service
MAIN_PATH=./cmd/api

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make test-cover   - Run tests with coverage report"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make lint         - Run golangci-lint"
	@echo "  make fmt          - Run go fmt on all source files"
	@echo "  make tidy         - Run go mod tidy"
	@echo "  make embedmd      - Update code snippets in the README"
	@echo "  make embedmd-check - Check if README code snippets are up-to-date"
	@echo "  make help         - Show this help message"

build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	@echo "Running..."
	@./$(BINARY_NAME)

dev:
	@echo "Running in development mode..."
	@go run $(MAIN_PATH)

test:
	@echo "Running tests..."
	@go test -v ./...

test-cover:
	@echo "Running tests with coverage..."
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html


tmp/help.txt:
	@echo "Updating help text"
	@mkdir -p tmp
	make help > tmp/help.txt 2>&1

embedmd: tmp/help.txt
	@echo "Updating code snippets in README.md..."
	@go tool embedmd -w README.md

embedmd-check:
	@echo "Checking if README code snippets are up-to-date..."
	@go tool embedmd -d README.md

lint:
	@echo "Running linter..."
	@go tool golangci-lint run ./...

fix:
	@echo "Running linter fix..."
	@go tool golangci-lint run ./... --fix

fmt:
	@echo "Formatting code..."
	@go fmt ./...

tidy:
	@echo "Tidying up module dependencies..."
	@go mod tidy

# Default target
.DEFAULT_GOAL := help 