# Simple Go Web Service

A RESTful API web service built with Go, featuring user management endpoints.

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

[embedmd]:# (tmp/help.txt)
```txt
Available commands:
  make build        - Build the application
  make run          - Run the application
  make test         - Run tests
  make test-cover   - Run tests with coverage report
  make clean        - Remove build artifacts
  make lint         - Run golangci-lint
  make fmt          - Run go fmt on all source files
  make tidy         - Run go mod tidy
  make embedmd      - Update code snippets in the README
  make embedmd-check - Check if README code snippets are up-to-date
  make help         - Show this help message
```

## Linting

The project uses golangci-lint for code quality checking. Install it with:

```bash
make lint-install
```

Then run the linter with:

```bash
make lint
```

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
func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{
		"message": "Welcome to the API",
	}

	jsonResponse(w, http.StatusOK, response)
}
```

Health check endpoint:

[embedmd]:# (handlers/handlers.go /func HealthCheckHandler/ /^}/)
```go
func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
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

## License

This project is licensed under the MIT License. 
