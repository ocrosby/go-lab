package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/domain"
	"github.com/ocrosby/go-lab/projects/api/mocks"
)

func TestUserHandler_CreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	expectedUser := &domain.User{
		ID:        "user_123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.EXPECT().
		CreateUser(gomock.Any(), "test@example.com", "Test User").
		Return(expectedUser, nil)

	reqBody := CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response domain.User
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Email != expectedUser.Email {
		t.Errorf("Expected email %s, got %s", expectedUser.Email, response.Email)
	}
}

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	req := httptest.NewRequest("POST", "/users", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Error != "invalid JSON" {
		t.Errorf("Expected error 'invalid JSON', got %s", response.Error)
	}
}

func TestUserHandler_CreateUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	mockService.EXPECT().
		CreateUser(gomock.Any(), "test@example.com", "Test User").
		Return(nil, domain.ErrUserAlreadyExists)

	reqBody := CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestUserHandler_GetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	expectedUser := &domain.User{
		ID:        "user_123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.EXPECT().
		GetUser(gomock.Any(), "user_123").
		Return(expectedUser, nil)

	req := httptest.NewRequest("GET", "/users/user_123", nil)
	w := httptest.NewRecorder()

	handler.getUserByID(w, req, "user_123")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response domain.User
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.ID != expectedUser.ID {
		t.Errorf("Expected ID %s, got %s", expectedUser.ID, response.ID)
	}
}

func TestUserHandler_GetUserByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	mockService.EXPECT().
		GetUser(gomock.Any(), "nonexistent").
		Return(nil, domain.ErrUserNotFound)

	req := httptest.NewRequest("GET", "/users/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.getUserByID(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	updatedUser := &domain.User{
		ID:        "user_123",
		Email:     "test@example.com",
		Name:      "Updated Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.EXPECT().
		UpdateUser(gomock.Any(), "user_123", "Updated Name").
		Return(updatedUser, nil)

	reqBody := UpdateUserRequest{
		Name: "Updated Name",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/users/user_123", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, req, "user_123")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response domain.User
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", response.Name)
	}
}

func TestUserHandler_UpdateUser_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	req := httptest.NewRequest("PUT", "/users/user_123", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, req, "user_123")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	mockService.EXPECT().
		DeleteUser(gomock.Any(), "user_123").
		Return(nil)

	req := httptest.NewRequest("DELETE", "/users/user_123", nil)
	w := httptest.NewRecorder()

	handler.deleteUser(w, req, "user_123")

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	mockService.EXPECT().
		DeleteUser(gomock.Any(), "nonexistent").
		Return(domain.ErrUserNotFound)

	req := httptest.NewRequest("DELETE", "/users/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.deleteUser(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUserHandler_ListUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	expectedUsers := []*domain.User{
		{ID: "1", Email: "user1@example.com", Name: "User 1"},
		{ID: "2", Email: "user2@example.com", Name: "User 2"},
	}

	mockService.EXPECT().
		ListUsers(gomock.Any(), 10, 0).
		Return(expectedUsers, nil)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.listUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []*domain.User
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 users, got %d", len(response))
	}
}

func TestUserHandler_ListUsers_WithPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	expectedUsers := []*domain.User{
		{ID: "3", Email: "user3@example.com", Name: "User 3"},
	}

	mockService.EXPECT().
		ListUsers(gomock.Any(), 5, 2).
		Return(expectedUsers, nil)

	req := httptest.NewRequest("GET", "/users?limit=5&offset=2", nil)
	w := httptest.NewRecorder()

	handler.listUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []*domain.User
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 user, got %d", len(response))
	}
}

func TestUserHandler_ListUsers_InvalidPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	mockService.EXPECT().
		ListUsers(gomock.Any(), 10, 0).
		Return([]*domain.User{}, nil)

	req := httptest.NewRequest("GET", "/users?limit=invalid&offset=invalid", nil)
	w := httptest.NewRecorder()

	handler.listUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserHandler_HandleUsers_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	req := httptest.NewRequest("PATCH", "/users", nil)
	w := httptest.NewRecorder()

	handler.handleUsers(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestUserHandler_HandleUserByID_EmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	req := httptest.NewRequest("GET", "/users/", nil)
	w := httptest.NewRecorder()

	handler.handleUserByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserHandler_HandleUserByID_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	req := httptest.NewRequest("PATCH", "/users/user_123", nil)
	w := httptest.NewRecorder()

	handler.handleUserByID(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestUserHandler_HandleServiceError_Coverage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	tests := []struct {
		name           string
		serviceError   error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "user not found",
			serviceError:   domain.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "user not found",
		},
		{
			name:           "user already exists",
			serviceError:   domain.ErrUserAlreadyExists,
			expectedStatus: http.StatusConflict,
			expectedMsg:    "user already exists",
		},
		{
			name:           "invalid input",
			serviceError:   domain.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid input",
		},
		{
			name:           "internal error",
			serviceError:   domain.ErrInternalError,
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.handleServiceError(w, tt.serviceError)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response ErrorResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Errorf("Failed to decode response: %v", err)
			}

			if response.Error != tt.expectedMsg {
				t.Errorf("Expected error message '%s', got '%s'", tt.expectedMsg, response.Error)
			}
		})
	}
}

func TestUserHandler_RegisterRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	logger, _ := zap.NewDevelopment()
	handler := NewUserHandler(mockService, logger)

	// Set up expectations for routes that will be called during testing
	mockService.EXPECT().
		ListUsers(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*domain.User{}, nil).
		AnyTimes()
	mockService.EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		Return(nil, domain.ErrUserNotFound).
		AnyTimes()
	mockService.EXPECT().
		UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, domain.ErrUserNotFound).
		AnyTimes()
	mockService.EXPECT().
		DeleteUser(gomock.Any(), gomock.Any()).
		Return(domain.ErrUserNotFound).
		AnyTimes()

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Test that routes are registered by checking they don't return 404
	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/users"},
		{"POST", "/users"},
		{"GET", "/users/123"},
		{"PUT", "/users/123"},
		{"DELETE", "/users/123"},
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.method == "POST" || tt.method == "PUT" {
				req.Header.Set("Content-Type", "application/json")
				// Add empty JSON body for POST/PUT
				req.Body = http.NoBody
			}
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			// Should not return 404 (route not found)
			if w.Code == http.StatusNotFound {
				t.Errorf("Route %s %s not registered", tt.method, tt.path)
			}
		})
	}
}
