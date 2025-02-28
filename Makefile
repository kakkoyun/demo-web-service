.PHONY: build run test clean lint fix fmt help tidy embedmd-check embedmd shellcheck shfmt yamlfmt install-linters update-help

BINARY_NAME=demo-web-service
MAIN_PATH=./cmd/api
SHELL_FILES := $(shell find . -name "*.sh" -not -path "./vendor/*")
YAML_FILES := $(shell find . -name "*.yml" -o -name "*.yaml" -not -path "./vendor/*")

GOFLAGS := GOFLAGS="${GOFLAGS} '-toolexec=orchestrion toolexec'"
DATADOG_ENV_VARS := DD_ENV=kakkoyun/local DD_SERVICE=demo-web-service DD_VERSION=0.0.0 DD_TAGS=env:local,version:0.0.0
DATADOG_DEBUG_ENV_VARS := DD_TRACE_DEBUG=true DD_RUNTIME_METRICS_ENABLED=true DD_PROFILING_ENABLED=true DD_DOGSTATSD_PORT=8135 DD_TRACE_AGENT_PORT=8136

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make test-cover   - Run tests with coverage report"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make lint         - Run all linters (Go, Shell, YAML)"
	@echo "  make fmt          - Format all source files (Go, Shell, YAML)"
	@echo "  make tidy         - Run go mod tidy"
	@echo "  make embedmd      - Update code snippets in the README"
	@echo "  make embedmd-check - Check if README code snippets are up-to-date"
	@echo "  make shellcheck   - Run shellcheck on shell scripts"
	@echo "  make shfmt        - Format shell scripts"
	@echo "  make yamlfmt      - Format YAML files with yamlfmt"
	@echo "  make install-linters - Install Go-based linting tools"
	@echo "  make help         - Show this help message"

build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)

build-instrumented:
	@echo "Building instrumented..."
	# @go build -toolexec "go tool errtrace toolexec" -o $(BINARY_NAME)-instrumented $(MAIN_PATH)
	go build -toolexec 'orchestrion toolexec' -o $(BINARY_NAME)-instrumented $(MAIN_PATH)

run: build
	@echo "Running..."
	@./$(BINARY_NAME)

dev:
	@echo "Running in development mode with Datadog tracing..."
	$(DATADOG_ENV_VARS) $(DATADOG_DEBUG_ENV_VARS) $(GOFLAGS) go run $(MAIN_PATH)

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


update-help:
	@echo "Updating help text"
	@mkdir -p tmp
	@$(MAKE) help > tmp/help.txt 2>/dev/null


embedmd: update-help
	@echo "Updating code snippets in README.md..."
	@go tool embedmd -w README.md

embedmd-check: update-help
	@echo "Checking if README code snippets are up-to-date..."
	@go tool embedmd -d README.md

lint: lint-go shellcheck
	@echo "All linters completed"

lint-all:
	@echo "Running all linters and formatters..."
	@./tools/lint-all.sh

fix-all:
	@echo "Running all linters and formatters in fix mode..."
	@./tools/lint-all.sh --fix

lint-go:
	@echo "Running Go linter..."
	@go tool golangci-lint run ./...

shellcheck:
	@echo "Running shellcheck on shell scripts..."
	@if command -v shellcheck >/dev/null 2>&1; then \
		shellcheck $(SHELL_FILES); \
	else \
		echo "shellcheck not found. Install with 'brew install shellcheck' or see shellcheck.net"; \
		exit 1; \
	fi

shfmt:
	@echo "Formatting shell scripts..."
	@go tool shfmt -i 2 -ci -w $(SHELL_FILES)

yamlfmt:
	@echo "Formatting YAML files..."
	@go tool yamlfmt $(YAML_FILES)

fmt: fmt-go shfmt yamlfmt
	@echo "Go, shell, and YAML formatting completed"

fmt-go:
	@echo "Formatting Go code..."
	@go fmt ./...

govulncheck:
	@echo "Running vulnerability scanner..."
	@go tool govulncheck ./...

fix:
	@echo "Running linter fix..."
	@go tool golangci-lint run ./... --fix

install-linters:
	@# Note about shellcheck
	@echo "Note: shellcheck is recommended but not installed automatically."
	@echo "      Install shellcheck with your OS package manager."
	@echo "      On macOS: brew install shellcheck"
	@echo "      On Linux: sudo apt-get install shellcheck"
	@echo "All Go-based linting tools installed successfully!"

tidy:
	@echo "Tidying up module dependencies..."
	@go mod tidy

inject-error-traces:
	@echo "Injecting error traces..."
	@go tool errtrace -w ./...

# Default target
.DEFAULT_GOAL := help