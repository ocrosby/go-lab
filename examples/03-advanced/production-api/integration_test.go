package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/application"
	"github.com/ocrosby/go-lab/projects/api/internal/domain"
	httpAdapter "github.com/ocrosby/go-lab/projects/api/internal/infrastructure/adapters/http"
	"github.com/ocrosby/go-lab/projects/api/internal/infrastructure/adapters/repository"
	"github.com/ocrosby/go-lab/projects/api/pkg/health"
)

func TestIntegration_UserLifecycle(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryUserRepository()
	service := application.NewUserService(repo, logger)
	handler := httpAdapter.NewUserHandler(service, logger)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Test data
	createUserReq := httpAdapter.CreateUserRequest{
		Email: "integration@example.com",
		Name:  "Integration User",
	}
	updateUserReq := httpAdapter.UpdateUserRequest{
		Name: "Updated Integration User",
	}

	// 1. Create user
	t.Run("Create User", func(t *testing.T) {
		jsonBody, _ := json.Marshal(createUserReq)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var user domain.User
		err := json.NewDecoder(w.Body).Decode(&user)
		if err != nil {
			t.Errorf("Failed to decode user: %v", err)
		}

		if user.Email != createUserReq.Email {
			t.Errorf("Expected email %s, got %s", createUserReq.Email, user.Email)
		}
		if user.Name != createUserReq.Name {
			t.Errorf("Expected name %s, got %s", createUserReq.Name, user.Name)
		}
		if user.ID == "" {
			t.Error("Expected user ID, got empty string")
		}
	})

	// 2. List users (should contain our created user)
	var createdUserID string
	t.Run("List Users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var users []*domain.User
		err := json.NewDecoder(w.Body).Decode(&users)
		if err != nil {
			t.Errorf("Failed to decode users: %v", err)
		}

		if len(users) != 1 {
			t.Errorf("Expected 1 user, got %d", len(users))
		}

		if users[0].Email != createUserReq.Email {
			t.Errorf("Expected email %s, got %s", createUserReq.Email, users[0].Email)
		}

		createdUserID = users[0].ID
	})

	// 3. Get user by ID
	t.Run("Get User by ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/"+createdUserID, nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var user domain.User
		err := json.NewDecoder(w.Body).Decode(&user)
		if err != nil {
			t.Errorf("Failed to decode user: %v", err)
		}

		if user.ID != createdUserID {
			t.Errorf("Expected ID %s, got %s", createdUserID, user.ID)
		}
	})

	// 4. Update user
	t.Run("Update User", func(t *testing.T) {
		jsonBody, _ := json.Marshal(updateUserReq)
		req := httptest.NewRequest("PUT", "/users/"+createdUserID, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var user domain.User
		err := json.NewDecoder(w.Body).Decode(&user)
		if err != nil {
			t.Errorf("Failed to decode user: %v", err)
		}

		if user.Name != updateUserReq.Name {
			t.Errorf("Expected updated name %s, got %s", updateUserReq.Name, user.Name)
		}
	})

	// 5. Delete user
	t.Run("Delete User", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/users/"+createdUserID, nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})

	// 6. Verify user is deleted
	t.Run("Verify User Deleted", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/"+createdUserID, nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestIntegration_HealthEndpoints(t *testing.T) {
	// Setup health checker
	healthChecker := health.NewHealthChecker()
	healthChecker.AddCheck("basic", func(ctx context.Context) error {
		return nil
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.LivenessHandler(healthChecker))
	mux.HandleFunc("/readyz", health.ReadinessHandler(healthChecker))
	mux.HandleFunc("/startupz", health.StartupHandler(healthChecker))

	endpoints := []string{"/healthz", "/readyz", "/startupz"}

	for _, endpoint := range endpoints {
		t.Run("Health endpoint "+endpoint, func(t *testing.T) {
			req := httptest.NewRequest("GET", endpoint, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected content type application/json, got %s", contentType)
			}

			var response health.HealthStatus
			err := json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Errorf("Failed to decode health response: %v", err)
			}

			if response.Status != health.StatusUp {
				t.Errorf("Expected status UP, got %s", response.Status)
			}
		})
	}
}

func TestIntegration_ErrorHandling(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryUserRepository()
	service := application.NewUserService(repo, logger)
	handler := httpAdapter.NewUserHandler(service, logger)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		contentType    string
		expectedStatus int
	}{
		{
			name:           "Invalid JSON in create user",
			method:         "POST",
			url:            "/users",
			body:           "invalid json",
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Get non-existent user",
			method:         "GET",
			url:            "/users/non-existent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Update non-existent user",
			method:         "PUT",
			url:            "/users/non-existent",
			body:           `{"name": "Updated Name"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Delete non-existent user",
			method:         "DELETE",
			url:            "/users/non-existent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Method not allowed on /users",
			method:         "PATCH",
			url:            "/users",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Method not allowed on /users/{id}",
			method:         "PATCH",
			url:            "/users/some-id",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Empty user ID",
			method:         "GET",
			url:            "/users/",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, tt.url, nil)
			}

			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestIntegration_Pagination(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryUserRepository()
	service := application.NewUserService(repo, logger)
	handler := httpAdapter.NewUserHandler(service, logger)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Create multiple users
	for i := 0; i < 5; i++ {
		createReq := httpAdapter.CreateUserRequest{
			Email: fmt.Sprintf("user%d@example.com", i),
			Name:  fmt.Sprintf("User %d", i),
		}
		jsonBody, _ := json.Marshal(createReq)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Failed to create user %d: status %d", i, w.Code)
		}
	}

	// Test pagination
	t.Run("List with limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users?limit=3", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var users []*domain.User
		err := json.NewDecoder(w.Body).Decode(&users)
		if err != nil {
			t.Errorf("Failed to decode users: %v", err)
		}

		if len(users) != 3 {
			t.Errorf("Expected 3 users, got %d", len(users))
		}
	})

	t.Run("List with offset", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users?limit=10&offset=2", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var users []*domain.User
		err := json.NewDecoder(w.Body).Decode(&users)
		if err != nil {
			t.Errorf("Failed to decode users: %v", err)
		}

		if len(users) != 3 { // 5 total - 2 offset = 3
			t.Errorf("Expected 3 users, got %d", len(users))
		}
	})
}
