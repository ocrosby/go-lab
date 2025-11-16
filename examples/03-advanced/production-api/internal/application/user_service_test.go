package application

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/domain"
	"github.com/ocrosby/go-lab/projects/api/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	name := "Test User"

	// Mock expectations
	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, domain.ErrUserNotFound)

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	// Execute
	user, err := service.CreateUser(ctx, email, name)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Error("Expected user, got nil")
		return
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}

	if user.Name != name {
		t.Errorf("Expected name %s, got %s", name, user.Name)
	}
}

func TestUserService_CreateUser_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()

	// Test empty email
	user, err := service.CreateUser(ctx, "", "Test User")
	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
	if user != nil {
		t.Error("Expected nil user for invalid input")
	}

	// Test empty name
	user, err = service.CreateUser(ctx, "test@example.com", "")
	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
	if user != nil {
		t.Error("Expected nil user for invalid input")
	}
}

func TestUserService_CreateUser_UserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	email := "existing@example.com"
	name := "Test User"

	existingUser := &domain.User{
		ID:    "existing-id",
		Email: email,
		Name:  "Existing User",
	}

	// Mock expectations
	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(existingUser, nil)

	// Execute
	user, err := service.CreateUser(ctx, email, name)

	// Assertions
	if err != domain.ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}

	if user != nil {
		t.Error("Expected nil user when user already exists")
	}
}

func TestUserService_CreateUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	name := "Test User"

	// Mock expectations
	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, domain.ErrUserNotFound)

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(domain.ErrInternalError)

	// Execute
	user, err := service.CreateUser(ctx, email, name)

	// Assertions
	if err != domain.ErrInternalError {
		t.Errorf("Expected ErrInternalError, got %v", err)
	}
	if user != nil {
		t.Error("Expected nil user on repository error")
	}
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	userID := "test-id"
	expectedUser := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Mock expectations
	mockRepo.EXPECT().
		GetByID(ctx, userID).
		Return(expectedUser, nil)

	// Execute
	user, err := service.GetUser(ctx, userID)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
		return
	}
	if user.ID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, user.ID)
	}
}

func TestUserService_GetUser_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()

	// Execute
	user, err := service.GetUser(ctx, "")

	// Assertions
	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
	if user != nil {
		t.Error("Expected nil user for invalid input")
	}
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	userID := "non-existent"

	// Mock expectations
	mockRepo.EXPECT().
		GetByID(ctx, userID).
		Return(nil, domain.ErrUserNotFound)

	// Execute
	user, err := service.GetUser(ctx, userID)

	// Assertions
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
	if user != nil {
		t.Error("Expected nil user when not found")
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	userID := "test-id"
	newName := "Updated Name"
	existingUser := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Old Name",
	}

	// Mock expectations
	mockRepo.EXPECT().
		GetByID(ctx, userID).
		Return(existingUser, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *domain.User) error {
			if user.Name != newName {
				t.Errorf("Expected updated name %s, got %s", newName, user.Name)
			}
			return nil
		})

	// Execute
	user, err := service.UpdateUser(ctx, userID, newName)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
		return
	}
	if user.Name != newName {
		t.Errorf("Expected updated name %s, got %s", newName, user.Name)
	}
}

func TestUserService_UpdateUser_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		newName string
	}{
		{"empty user ID", "", "New Name"},
		{"empty name", "user-id", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.UpdateUser(ctx, tt.userID, tt.newName)
			if err != domain.ErrInvalidInput {
				t.Errorf("Expected ErrInvalidInput, got %v", err)
			}
			if user != nil {
				t.Error("Expected nil user for invalid input")
			}
		})
	}
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	userID := "non-existent"
	newName := "New Name"

	// Mock expectations
	mockRepo.EXPECT().
		GetByID(ctx, userID).
		Return(nil, domain.ErrUserNotFound)

	// Execute
	user, err := service.UpdateUser(ctx, userID, newName)

	// Assertions
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
	if user != nil {
		t.Error("Expected nil user when not found")
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	userID := "test-id"

	// Mock expectations
	mockRepo.EXPECT().
		Delete(ctx, userID).
		Return(nil)

	// Execute
	err := service.DeleteUser(ctx, userID)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUserService_DeleteUser_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()

	// Execute
	err := service.DeleteUser(ctx, "")

	// Assertions
	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
}

func TestUserService_DeleteUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	userID := "test-id"

	// Mock expectations
	mockRepo.EXPECT().
		Delete(ctx, userID).
		Return(domain.ErrInternalError)

	// Execute
	err := service.DeleteUser(ctx, userID)

	// Assertions
	if err != domain.ErrInternalError {
		t.Errorf("Expected ErrInternalError, got %v", err)
	}
}

func TestUserService_ListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	limit := 10
	offset := 0
	expectedUsers := []*domain.User{
		{ID: "1", Email: "user1@example.com", Name: "User 1"},
		{ID: "2", Email: "user2@example.com", Name: "User 2"},
	}

	// Mock expectations
	mockRepo.EXPECT().
		List(ctx, limit, offset).
		Return(expectedUsers, nil)

	// Execute
	users, err := service.ListUsers(ctx, limit, offset)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(users) != len(expectedUsers) {
		t.Errorf("Expected %d users, got %d", len(expectedUsers), len(users))
	}
}

func TestUserService_ListUsers_InvalidPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()

	// Test with invalid limit (should default to 10)
	mockRepo.EXPECT().
		List(ctx, 10, 0).
		Return([]*domain.User{}, nil)

	users, err := service.ListUsers(ctx, -1, -5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if users == nil {
		t.Error("Expected empty slice, got nil")
	}
}

func TestUserService_ListUsers_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	service := NewUserService(mockRepo, logger)

	ctx := context.Background()
	limit := 10
	offset := 0

	// Mock expectations
	mockRepo.EXPECT().
		List(ctx, limit, offset).
		Return(nil, domain.ErrInternalError)

	// Execute
	users, err := service.ListUsers(ctx, limit, offset)

	// Assertions
	if err != domain.ErrInternalError {
		t.Errorf("Expected ErrInternalError, got %v", err)
	}
	if users != nil {
		t.Error("Expected nil users on repository error")
	}
}
