package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kakkoyun/demo-web-service/handlers"
	"github.com/kakkoyun/demo-web-service/models"
)

// setupAPITest creates a test server with the application's routes
func setupAPITest() *httptest.Server {
	// Set up the routes similar to how main.go does it
	mux := http.NewServeMux()

	// Set up routes with Go 1.22 pattern syntax (via handler mapping)
	mux.HandleFunc("GET /", handlers.HomeHandler)
	mux.HandleFunc("GET /api/health", handlers.HealthCheckHandler)
	mux.HandleFunc("GET /api/users", handlers.GetUsersHandler)
	mux.HandleFunc("POST /api/users", handlers.CreateUserHandler)
	mux.HandleFunc("GET /api/users/{id}", handlers.GetUserHandler)

	// Apply middleware
	var handler http.Handler = mux
	handler = handlers.LoggingMiddleware(handler)

	// Create a test server
	return httptest.NewServer(handler)
}

func TestAPIEndpoints(t *testing.T) {
	// Set up the test server
	server := setupAPITest()
	defer server.Close()

	// Test case 1: Get all users
	t.Run("Get All Users", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/users")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		// Check content type
		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected content type application/json, got %v", contentType)
		}

		// Parse response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var response models.UserResponse
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// Verify response data
		if response.Status != "success" {
			t.Errorf("Expected status 'success', got %v", response.Status)
		}

		if len(response.Users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(response.Users))
		}
	})

	// Test case 2: Get a specific user
	t.Run("Get User by ID", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/users/1")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		// Parse response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var response models.UserResponse
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// Verify response data
		if response.Status != "success" {
			t.Errorf("Expected status 'success', got %v", response.Status)
		}

		if response.User == nil {
			t.Fatal("No user returned in response")
		}

		if response.User.ID != 1 {
			t.Errorf("Expected user ID 1, got %d", response.User.ID)
		}

		if response.User.Name != "User 1" {
			t.Errorf("Expected user name 'User 1', got %s", response.User.Name)
		}
	})

	// Test case 3: Get a user with invalid ID
	t.Run("Get User with Invalid ID", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/users/invalid")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status Bad Request, got %v", resp.Status)
		}

		// Parse response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var response map[string]string
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// Verify response data
		if response["status"] != "error" {
			t.Errorf("Expected status 'error', got %v", response["status"])
		}
	})

	// Test case 4: Create a new user
	t.Run("Create User", func(t *testing.T) {
		reqBody := `{"name":"Test User"}`
		resp, err := http.Post(
			server.URL+"/api/users",
			"application/json",
			strings.NewReader(reqBody),
		)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status Created, got %v", resp.Status)
		}

		// Parse response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var response models.UserResponse
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// Verify response data
		if response.Status != "success" {
			t.Errorf("Expected status 'success', got %v", response.Status)
		}

		if response.Message != "User created successfully" {
			t.Errorf("Expected message 'User created successfully', got %v", response.Message)
		}

		if response.User == nil {
			t.Fatal("No user returned in response")
		}

		if response.User.ID != 3 {
			t.Errorf("Expected user ID 3, got %d", response.User.ID)
		}
	})

	// Test case 5: Health check
	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/health")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		// Parse response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var response map[string]string
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// Verify response data
		if response["status"] != "healthy" {
			t.Errorf("Expected status 'healthy', got %v", response["status"])
		}
	})
}
