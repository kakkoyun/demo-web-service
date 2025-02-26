package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kakkoyun/demo-web-service/models"
)

// HomeHandler handles the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling home request", "path", r.URL.Path, "method", r.Method)

	response := map[string]string{
		"message": "Welcome to the API",
	}

	jsonResponse(w, http.StatusOK, response)
}

// HealthCheckHandler returns the API health status
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Health check requested", "remote_addr", r.RemoteAddr)

	response := map[string]string{
		"status": "healthy",
	}

	jsonResponse(w, http.StatusOK, response)
}

// GetUsersHandler returns a list of users
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Getting all users", "path", r.URL.Path)

	// In a real application, we would get these from a database
	users := []models.User{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Smith"},
	}

	response := models.UserResponse{
		Status: "success",
		Users:  users,
	}

	jsonResponse(w, http.StatusOK, response)
}

// CreateUserHandler creates a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Creating new user", "path", r.URL.Path)

	var user models.User = models.User{
		ID:   3,
		Name: "New User",
	}

	response := models.UserResponse{
		Status:  "success",
		Message: "User created successfully",
		User:    &user,
	}

	jsonResponse(w, http.StatusCreated, response)
}

// GetUserHandler returns a specific user by ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID from path parameter using Go 1.22's PathValue method
	idStr := r.PathValue("id")

	slog.Info("Getting user by ID", "id", idStr, "path", r.URL.Path)

	id, err := strconv.Atoi(idStr)

	if err != nil {
		slog.Error("Invalid user ID", "id", idStr, "error", err)
		errorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// In a real application, we would get this from a database
	user := models.User{
		ID:   id,
		Name: fmt.Sprintf("User %d", id),
	}

	response := models.UserResponse{
		Status: "success",
		User:   &user,
	}

	jsonResponse(w, http.StatusOK, response)
}

// jsonResponse sends a JSON response
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// errorResponse sends an error response
func errorResponse(w http.ResponseWriter, status int, message string) {
	slog.Warn("Sending error response", "status", status, "message", message)

	response := map[string]string{
		"status":  "error",
		"message": message,
	}

	jsonResponse(w, status, response)
}
