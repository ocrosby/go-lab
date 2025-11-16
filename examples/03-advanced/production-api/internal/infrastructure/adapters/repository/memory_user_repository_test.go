package repository

import (
	"context"
	"testing"
	"time"

	"github.com/ocrosby/go-lab/projects/api/internal/domain"
)

func TestMemoryUserRepository_Create(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := &domain.User{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Try to create the same user again (should fail)
	err = repo.Create(ctx, user)
	if err != domain.ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestMemoryUserRepository_GetByID(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := &domain.User{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user first
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test getting existing user
	retrievedUser, err := repo.GetByID(ctx, "test-id")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedUser == nil {
		t.Error("Expected user, got nil")
		return
	}
	if retrievedUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrievedUser.ID)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}

	// Test getting non-existent user
	_, err = repo.GetByID(ctx, "non-existent")
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryUserRepository_GetByEmail(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := &domain.User{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user first
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test getting existing user by email
	retrievedUser, err := repo.GetByEmail(ctx, "test@example.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedUser == nil {
		t.Error("Expected user, got nil")
		return
	}
	if retrievedUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrievedUser.ID)
	}

	// Test getting non-existent user by email
	_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryUserRepository_Update(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := &domain.User{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user first
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Update user
	user.Name = "Updated User"
	user.UpdatedAt = time.Now()

	err = repo.Update(ctx, user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	updatedUser, err := repo.GetByID(ctx, "test-id")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updatedUser.Name != "Updated User" {
		t.Errorf("Expected name 'Updated User', got %s", updatedUser.Name)
	}

	// Test updating non-existent user
	nonExistentUser := &domain.User{
		ID:    "non-existent",
		Email: "test@example.com",
		Name:  "Test",
	}
	err = repo.Update(ctx, nonExistentUser)
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryUserRepository_Delete(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := &domain.User{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user first
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Delete user
	err = repo.Delete(ctx, "test-id")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, "test-id")
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound after deletion, got %v", err)
	}

	// Test deleting non-existent user
	err = repo.Delete(ctx, "non-existent")
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryUserRepository_List(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Create test users
	users := []*domain.User{
		{ID: "1", Email: "user1@example.com", Name: "User 1"},
		{ID: "2", Email: "user2@example.com", Name: "User 2"},
		{ID: "3", Email: "user3@example.com", Name: "User 3"},
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user %s: %v", user.ID, err)
		}
	}

	// Test listing all users
	result, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 users, got %d", len(result))
	}

	// Test pagination - limit
	result, err = repo.List(ctx, 2, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 users with limit, got %d", len(result))
	}

	// Test pagination - offset
	result, err = repo.List(ctx, 10, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 user with offset, got %d", len(result))
	}

	// Test empty result with high offset
	result, err = repo.List(ctx, 10, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 users with high offset, got %d", len(result))
	}
}

func TestMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := &domain.User{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test concurrent reads
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			_, err := repo.GetByID(ctx, "test-id")
			if err != nil {
				t.Errorf("Concurrent read failed: %v", err)
			}
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_, err := repo.GetByEmail(ctx, "test@example.com")
			if err != nil {
				t.Errorf("Concurrent read failed: %v", err)
			}
		}
		done <- true
	}()

	// Wait for goroutines to complete
	<-done
	<-done
}

func TestMemoryUserRepository_DataIsolation(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := &domain.User{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Original Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get user and modify returned instance
	retrievedUser, err := repo.GetByID(ctx, "test-id")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Modify the retrieved user
	retrievedUser.Name = "Modified Name"

	// Get user again to ensure data isolation
	retrievedUser2, err := repo.GetByID(ctx, "test-id")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if retrievedUser2.Name != "Original Name" {
		t.Errorf("Data isolation failed: expected 'Original Name', got %s", retrievedUser2.Name)
	}
}
