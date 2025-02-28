package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"braces.dev/errtrace"

	"github.com/kakkoyun/demo-web-service/models"
)

// TestMode controls whether random errors are generated
// Set this to true in tests to disable random failures
var TestMode bool

// HomeHandler handles the root endpoint
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

	// Randomly generate an error 20% of the time (but not in test mode)
	if !TestMode && rand.Intn(5) == 0 {
		// Simple error handling - just log and return an error
		err := errors.New("database connection failed")
		slog.Error("Failed to get users", "error", err)
		errorResponse(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

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

	// Example of basic error checking
	if r.ContentLength == 0 {
		err := errors.New("empty request body")
		slog.Error("Failed to create user", "error", err)
		errorResponse(w, http.StatusBadRequest, "Empty request body")
		return
	}

	// Process the user data and handle any errors
	if err := validateAndCreateUser(r); err != nil {
		// Here we handle errors from our nested function
		statusCode := http.StatusBadRequest
		errMsg := err.Error()

		slog.Error("User creation failed",
			"error", err,
			"status", statusCode)
		errorResponse(w, statusCode, errMsg)
		return
	}

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

// Common validation errors
var (
	ErrValidation = errors.New("validation error")
)

// validateAndCreateUser demonstrates nested function calls with error wrapping
func validateAndCreateUser(_ *http.Request) error {
	// Randomly generate validation errors
	if !TestMode && rand.Intn(3) == 0 {
		return errtrace.Wrap(fmt.Errorf("%w: required fields missing", ErrValidation))
	}

	// Try to process the user data
	if err := processUserData(); err != nil {
		// Wrap the lower-level error
		return errtrace.Wrap(fmt.Errorf("user processing failed: %w", err))
	}

	return nil
}

// processUserData is a nested function that might return errors
func processUserData() error {
	// Randomly fail this operation (but not in test mode)
	if !TestMode && rand.Intn(4) == 0 {
		return errtrace.Wrap(errors.New("database constraint violation"))
	}

	// Simulate slow processing (minimal in test mode)
	if TestMode {
		time.Sleep(time.Millisecond)
	} else {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}

	return nil
}

// GetUserHandler returns a specific user by ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID from path parameter using Go 1.22's PathValue method
	idStr := r.PathValue("id")

	slog.Info("Getting user by ID", "id", idStr, "path", r.URL.Path)

	// Convert string ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Log the error with its stack trace for debugging
		stack := debug.Stack()
		wrappedErr := fmt.Errorf("%w: %s is not a valid integer", ErrInvalidUserID, idStr)

		slog.Error("Invalid user ID",
			"id", idStr,
			"error", wrappedErr,
			"stack", string(stack))

		errorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid user ID: %s", idStr))
		return
	}

	// Validate the ID
	if id <= 0 {
		wrappedErr := fmt.Errorf("%w: ID must be positive", ErrInvalidUserID)
		slog.Error("Invalid user ID value",
			"id", id,
			"error", wrappedErr)

		errorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid user ID: %d", id))
		return
	}

	// Simulate database query that might fail
	if err := queryDatabase(id); err != nil {
		slog.Error("Database query failed",
			"id", id,
			"error", err)
		errorResponse(w, http.StatusInternalServerError, "Failed to retrieve user data")
		return
	}

	// Randomly generate "not found" errors for valid IDs over 10 (but not in test mode)
	if !TestMode && id > 10 && rand.Intn(2) == 0 {
		notFoundErr := fmt.Errorf("%w: ID %d", ErrUserNotFound, id)
		slog.Error("User not found",
			"id", id,
			"error", notFoundErr)

		errorResponse(w, http.StatusNotFound, fmt.Sprintf("User with ID %d not found", id))
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

// Common user errors
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user ID")
)

// queryDatabase simulates a database query that might fail
func queryDatabase(id int) error {
	// Simulate different database errors (but not in test mode)
	if TestMode {
		return nil
	}

	// Use id in random error generation
	errorChance := rand.Intn(10)

	// IDs divisible by 5 have a higher chance of connection timeout
	if id%5 == 0 && errorChance < 3 {
		return errtrace.Wrap(errors.New("connection timeout"))
	}

	// IDs divisible by 3 have a higher chance of query execution failure
	if id%3 == 0 && errorChance < 3 {
		return errtrace.Wrap(errors.New("query execution failed"))
	}

	// Very high IDs might cause a constraint error
	if id > 1000 && errorChance < 2 {
		return errtrace.Wrap(errors.New("primary key constraint violation"))
	}

	return nil
}

// jsonResponse sends a JSON response
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
	}
}

// JSONResponse is an exported version of jsonResponse that can be used by other packages
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	jsonResponse(w, status, data)
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
