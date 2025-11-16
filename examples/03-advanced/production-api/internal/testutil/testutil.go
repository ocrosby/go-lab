// Package testutil provides testing utilities and helpers for the User Management API.
package testutil

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/domain"
)

// CreateTestUser creates a test user with default values
func CreateTestUser(id, email, name string) *domain.User {
	now := time.Now()
	return &domain.User{
		ID:        id,
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateTestUsers creates multiple test users
func CreateTestUsers(count int) []*domain.User {
	users := make([]*domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = CreateTestUser(
			fmt.Sprintf("user_%d", i+1),
			fmt.Sprintf("user%d@example.com", i+1),
			fmt.Sprintf("User %d", i+1),
		)
	}
	return users
}

// NewTestLogger creates a no-op logger for testing
func NewTestLogger() *zap.Logger {
	return zap.NewNop()
}

// NewTestDevelopmentLogger creates a development logger for testing
func NewTestDevelopmentLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return zap.NewNop()
	}
	return logger
}

// AssertUserEqual asserts that two users are equal
func AssertUserEqual(t *testing.T, expected, actual *domain.User) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		t.Errorf("Expected one user to be nil, got expected=%v, actual=%v", expected, actual)
		return
	}

	if expected.ID != actual.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, actual.ID)
	}
	if expected.Email != actual.Email {
		t.Errorf("Expected Email %s, got %s", expected.Email, actual.Email)
	}
	if expected.Name != actual.Name {
		t.Errorf("Expected Name %s, got %s", expected.Name, actual.Name)
	}
}

// AssertUsersEqual asserts that two user slices are equal
func AssertUsersEqual(t *testing.T, expected, actual []*domain.User) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("Expected %d users, got %d", len(expected), len(actual))
		return
	}

	for i := range expected {
		AssertUserEqual(t, expected[i], actual[i])
	}
}

// AssertError asserts that an error matches the expected error
func AssertError(t *testing.T, expected, actual error) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		t.Errorf("Expected error mismatch: expected=%v, actual=%v", expected, actual)
		return
	}

	if expected.Error() != actual.Error() {
		t.Errorf("Expected error %q, got %q", expected.Error(), actual.Error())
	}
}
