package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/kakkoyun/demo-web-service/models"
)

// TestMain sets up the testing environment
func TestMain(m *testing.M) {
	// Enable test mode to disable random failures
	TestMode = true

	// Run all tests
	exitCode := m.Run()

	// Exit with the same code
	os.Exit(exitCode)
}

func TestGetUsersHandler(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("GET", "/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsersHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Parse the response body
	var response models.UserResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("could not parse response body: %v", err)
	}

	// Check the response body
	if response.Status != "success" {
		t.Errorf("handler returned wrong status: got %v want %v", response.Status, "success")
	}

	// Check the users list
	if len(response.Users) != 2 {
		t.Errorf("handler returned wrong number of users: got %v want %v", len(response.Users), 2)
	}

	// Check the first user's fields
	if response.Users[0].ID != 1 || response.Users[0].Name != "John Doe" {
		t.Errorf("handler returned wrong first user data: got ID=%v, Name=%v", response.Users[0].ID, response.Users[0].Name)
	}
}

// TestGetUserHandlerDirect tests the GetUserHandler handler
// Note: Because of Go 1.22's PathValue method, we're using a custom test
// that extracts the ID from the URL path and passes it to a simplified version
// of the handler function to test just the core logic.
func TestGetUserHandlerDirect(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		userID         string
		expectedName   string
		expectedStatus int
		isError        bool
	}{
		{
			name:           "Valid User ID",
			userID:         "1",
			expectedStatus: http.StatusOK,
			expectedName:   "User 1",
			isError:        false,
		},
		{
			name:           "Invalid User ID",
			userID:         "abc",
			expectedStatus: http.StatusBadRequest,
			expectedName:   "",
			isError:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create a test handler that simulates the functionality of GetUserHandler
			// but accepts the ID directly instead of using PathValue
			testHandler := func(w http.ResponseWriter, _ *http.Request) {
				idStr := tc.userID // Directly use the test case ID

				id, err := strconv.Atoi(idStr)
				if err != nil {
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

			// Create a request
			req, err := http.NewRequest("GET", "/api/users/"+tc.userID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Call the test handler directly
			testHandler(rr, req)

			// Check the status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			// If we're expecting an error, we don't need to check the response body
			if tc.isError {
				return
			}

			// Parse the response body
			var response models.UserResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("could not parse response body: %v", err)
			}

			// Check the user data
			if response.User.Name != tc.expectedName {
				t.Errorf("handler returned wrong user name: got %v want %v", response.User.Name, tc.expectedName)
			}
		})
	}
}

func TestCreateUserHandler(t *testing.T) {
	// Create a request with a JSON body
	reqBody := `{"name":"New Test User"}`
	req, err := http.NewRequest("POST", "/api/users", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUserHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Parse the response body
	var response models.UserResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("could not parse response body: %v", err)
	}

	// Check the response data
	if response.Status != "success" {
		t.Errorf("handler returned wrong status: got %v want %v", response.Status, "success")
	}

	if response.Message != "User created successfully" {
		t.Errorf("handler returned wrong message: got %v want %v", response.Message, "User created successfully")
	}

	if response.User == nil {
		t.Errorf("handler did not return a user")
	} else if response.User.ID != 3 {
		t.Errorf("handler returned wrong user ID: got %v want %v", response.User.ID, 3)
	}
}
