# Simple Go Web Service

A RESTful API web service built with Go, featuring user management endpoints.

[![CI/CD Pipeline](https://github.com/kakkoyun/demo-web-service/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/kakkoyun/demo-web-service/actions/workflows/ci-cd.yml)
[![Security Scan](https://github.com/kakkoyun/demo-web-service/actions/workflows/security.yml/badge.svg)](https://github.com/kakkoyun/demo-web-service/actions/workflows/security.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue)](https://golang.org/doc/go1.24)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/kakkoyun/demo-web-service?include_prereleases)](https://github.com/kakkoyun/demo-web-service/releases)

## Features

- RESTful API endpoints for user management
- JSON responses
- Health check endpoint
- Environment-based configuration
- Graceful shutdown

## Requirements

- Go 1.24 or higher

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/demo-web-service.git
cd demo-web-service
```

2. Install dependencies:

```bash
go mod tidy
```

## Running the Service

Using Make:
```bash
make run     # Build and run the application
make dev     # Run without building (for development)
```

Or manually:

```bash
go run cmd/api/main.go
```

The server will start on port 8080 by default.

## Development Commands

The project includes a Makefile with common tasks:

```txt
Available commands:
  make build        - Build the application
  make run          - Run the application
  make test         - Run tests
  make test-cover   - Run tests with coverage report
  make clean        - Remove build artifacts
  make lint         - Run all linters (Go, Shell, YAML)
  make fmt          - Format all source files (Go, Shell, YAML)
  make tidy         - Run go mod tidy
  make embedmd      - Update code snippets in the README
  make embedmd-check - Check if README code snippets are up-to-date
  make shellcheck   - Run shellcheck on shell scripts
  make shfmt        - Format shell scripts
  make yamlfmt      - Format YAML files with yamlfmt
  make install-linters - Install Go-based linting tools
  make help         - Show this help message
```

## Linting

The project uses several linters to ensure code quality:

### Go Code Linting

Go code is linted using golangci-lint, which includes multiple linters in one tool:

```bash
make lint-go
```

### Shell Script Linting

Shell scripts are linted using shellcheck, a static analysis tool that gives warnings and suggestions:

```bash
make shellcheck
```

To format shell scripts according to a consistent style using the Go-based `shfmt` tool:

```bash
make shfmt
```

### YAML Linting

YAML files are formatted and validated using Google's yamlfmt, a Go-based YAML formatter with validation capabilities:

```bash
make yamlfmt
```

### Run All Linters

To run all linters at once:

```bash
make lint
```

To run linters and formatters with auto-fixing enabled:

```bash
make fix-all
```

### Installing Linting Tools

Install all required Go-based linting tools:

```bash
make install-linters
```

This will install `yamlfmt` and other Go-based tools. Note that shellcheck must be installed separately using your OS package manager:
- macOS: `brew install shellcheck`
- Ubuntu/Debian: `sudo apt-get install shellcheck`

## Documentation

The README is kept in sync with the actual codebase using [embedmd](https://github.com/campoy/embedmd). 
After making changes to the code, run the following command to update the code snippets in the README:

```bash
make embedmd
```

You can also check if the README is in sync with the code without modifying it:

```bash
make embedmd-check
```

## Configuration

The service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Port the server listens on | 8080 |
| READ_TIMEOUT | HTTP read timeout | 15s |
| WRITE_TIMEOUT | HTTP write timeout | 15s |
| IDLE_TIMEOUT | HTTP idle timeout | 60s |
| ALLOWED_ORIGINS | CORS allowed origins (comma-separated) | http://localhost:3000,http://localhost:8080 |

## Code Examples

### Configuration

Here's how the application loads its configuration:

[embedmd]:# (config/config.go /func LoadConfig/ /^}/)
```go
func LoadConfig() *Config {
	return &Config{
		ServerPort:     env("SERVER_PORT", "8080"),
		ReadTimeout:    durationEnv("READ_TIMEOUT", "15s"),
		WriteTimeout:   durationEnv("WRITE_TIMEOUT", "15s"),
		IdleTimeout:    durationEnv("IDLE_TIMEOUT", "60s"),
		AllowedOrigins: sliceEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080"),
	}
}
```

### API Handlers

Example of an API handler:

[embedmd]:# (handlers/handlers.go /func HomeHandler/ /^}/)
```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling home request", "path", r.URL.Path, "method", r.Method)

	// Randomly generate an error 10% of the time (but not in test mode)
	if !TestMode && rand.Intn(10) == 0 {
		slog.Error("Random error in home handler", "error", "random service unavailable")
		errorResponse(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
		return
	}

	response := map[string]string{
		"message": "Welcome to the API",
	}

	jsonResponse(w, http.StatusOK, response)
}
```

Health check endpoint:

[embedmd]:# (handlers/handlers.go /func HealthCheckHandler/ /^}/)
```go
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Health check requested", "remote_addr", r.RemoteAddr)

	response := map[string]string{
		"status": "healthy",
	}

	jsonResponse(w, http.StatusOK, response)
}
```

### Models

User model definition:

[embedmd]:# (models/user.go /type User/ /^}/)
```go
type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
```

## API Endpoints

### Base URL

`http://localhost:8080`

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | / | Home page - Welcome message |
| GET | /api/health | Health check endpoint |
| GET | /api/users | Get all users |
| POST | /api/users | Create a new user |
| GET | /api/users/{id} | Get user by ID |

### Example Requests

#### Get all users

```bash
curl http://localhost:8080/api/users
```

#### Get a specific user

```bash
curl http://localhost:8080/api/users/1
```

#### Create a user

```bash
curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json" -d '{"name":"New User"}'
```

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── config/
│   └── config.go            # Configuration handling
├── handlers/
│   └── handlers.go          # HTTP request handlers
├── models/
│   └── user.go              # Data models
├── .golangci.yml            # Golangci-lint configuration
├── Makefile                 # Build automation
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
└── README.md                # This file
```

## Testing

This project includes a comprehensive test suite to ensure all components function as expected:

### Unit Tests

The handlers package includes unit tests for all API endpoints, which verify:
- Correct HTTP status codes are returned
- Response body structure and content is valid
- Error handling works as expected

Run the unit tests with:

```bash
go test ./handlers -v
```

### Integration Tests

The tests directory contains integration tests that verify the complete API flow using an HTTP test server. These tests:
- Set up a full HTTP server with all routes
- Test each endpoint with real HTTP requests
- Verify responses, including status codes and JSON payloads
- Ensure middleware is applied correctly

Run the integration tests with:

```bash
go test ./tests -v
```

### Load Testing

The project includes a comprehensive load testing setup using k6, InfluxDB, and Grafana for metrics visualization. The load testing components are organized in the `tests` directory.

The load tests perform a variety of scenarios:

- Health check endpoint performance
- Version endpoint checks
- Home page response time
- Listing all users (GET /api/users)  
- Getting users by ID, with both valid and invalid IDs
- Creating new users with random data

Each endpoint is monitored for:
- Response times (p95, p99 percentiles)
- Success/failure rates
- Correct response formats and status codes

To run the load tests:

1. Ensure your API is running (by default on port 8080)
2. Execute the load test script:

```bash
cd tests
./run-loadtest.sh
```

The script will:

- Start InfluxDB and Grafana containers for metrics collection and visualization
- Run k6 load tests against your API endpoints
- Store the results in InfluxDB
- Provide access to visualize performance metrics in Grafana (http://localhost:3000)

Detailed instructions for load testing can be found in [tests/README.md](tests/README.md).

## CI/CD Pipeline

This project uses GitHub Actions for continuous integration and deployment. The status of these workflows is displayed as badges at the top of this README:

- **CI/CD Pipeline**: Shows the status of tests, code verification, and builds
- **Security Scan**: Displays the status of security scans that check for vulnerabilities
- **Go Version**: Indicates the Go version used by the project
- **License**: Shows the project's license type
- **Release**: Shows the latest release version

The CI/CD pipeline includes:

- Automated testing with coverage reporting
- Code quality checks using golangci-lint
- Documentation verification with embedmd
- Automated builds for multiple platforms
- Security scanning with Gosec and govulncheck

## License

This project is licensed under the MIT License. 
