package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ocrosby/go-lab/lessons/14-production-api/internal/application"
	"github.com/ocrosby/go-lab/lessons/14-production-api/internal/domain"
	"github.com/ocrosby/go-lab/lessons/14-production-api/internal/testutil"
)

// These tests are black-box (see rules/black-box-testing.md):
// they use a real in-memory testutil.FakeUserRepository as the collaborator
// instead of mocking domain.UserRepository. Assertions are on observable
// outcomes — the returned value and the persisted state — not on which
// repository methods were called with which arguments.

func newService(t *testing.T) (domain.UserService, *testutil.FakeUserRepository) {
	t.Helper()
	repo := testutil.NewFakeUserRepository()
	svc := application.NewUserService(repo, testutil.NewTestLogger())
	return svc, repo
}

func TestCreateUser_PersistsAndReturnsUser(t *testing.T) {
	// Arrange
	svc, repo := newService(t)

	// Act
	user, err := svc.CreateUser(context.Background(), "new@example.com", "New User")

	// Assert — service return
	if err != nil {
		t.Fatalf("CreateUser err = %v, want nil", err)
	}
	if user.Email != "new@example.com" || user.Name != "New User" {
		t.Errorf("CreateUser returned user = %+v", user)
	}
	if user.ID == "" {
		t.Error("CreateUser returned user with empty ID")
	}

	// Assert — observable persistence
	persisted, err := svc.GetUser(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetUser after Create err = %v", err)
	}
	if persisted.Email != "new@example.com" {
		t.Errorf("persisted email = %q, want %q", persisted.Email, "new@example.com")
	}

	// Assert — fake state (redundant with GetUser but useful when the SUT
	// might route reads through a cache in a real implementation).
	if len(repo.Users()) != 1 {
		t.Errorf("repo has %d users, want 1", len(repo.Users()))
	}
}

func TestCreateUser_RejectsInvalidInput(t *testing.T) {
	svc, _ := newService(t)

	tests := []struct {
		name, email, userName string
	}{
		{"empty email", "", "Test User"},
		{"empty name", "test@example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.CreateUser(context.Background(), tt.email, tt.userName)
			if !errors.Is(err, domain.ErrInvalidInput) {
				t.Errorf("err = %v, want ErrInvalidInput", err)
			}
			if user != nil {
				t.Errorf("user = %+v, want nil on invalid input", user)
			}
		})
	}
}

func TestCreateUser_RejectsDuplicateEmail(t *testing.T) {
	// Arrange
	svc, _ := newService(t)
	_, err := svc.CreateUser(context.Background(), "same@example.com", "First")
	if err != nil {
		t.Fatalf("first CreateUser err = %v", err)
	}

	// Act
	user, err := svc.CreateUser(context.Background(), "same@example.com", "Second")

	// Assert
	if !errors.Is(err, domain.ErrUserAlreadyExists) {
		t.Errorf("err = %v, want ErrUserAlreadyExists", err)
	}
	if user != nil {
		t.Errorf("user = %+v, want nil when duplicate", user)
	}
}

func TestCreateUser_SurfacesRepositoryError(t *testing.T) {
	svc, repo := newService(t)
	repo.FailNextCreate = domain.ErrInternalError

	user, err := svc.CreateUser(context.Background(), "test@example.com", "Test User")

	if !errors.Is(err, domain.ErrInternalError) {
		t.Errorf("err = %v, want ErrInternalError", err)
	}
	if user != nil {
		t.Errorf("user = %+v, want nil on repository error", user)
	}
}

func TestGetUser_ReturnsSeededUser(t *testing.T) {
	svc, repo := newService(t)
	repo.Seed(testutil.CreateTestUser("seeded-1", "seeded@example.com", "Seeded User"))

	user, err := svc.GetUser(context.Background(), "seeded-1")
	if err != nil {
		t.Fatalf("GetUser err = %v", err)
	}
	if user.Email != "seeded@example.com" {
		t.Errorf("email = %q, want %q", user.Email, "seeded@example.com")
	}
}

func TestGetUser_RejectsInvalidInput(t *testing.T) {
	svc, _ := newService(t)

	user, err := svc.GetUser(context.Background(), "")

	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("err = %v, want ErrInvalidInput", err)
	}
	if user != nil {
		t.Errorf("user = %+v, want nil", user)
	}
}

func TestGetUser_ReturnsNotFoundWhenAbsent(t *testing.T) {
	svc, _ := newService(t)

	user, err := svc.GetUser(context.Background(), "nobody-home")

	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("err = %v, want ErrUserNotFound", err)
	}
	if user != nil {
		t.Errorf("user = %+v, want nil", user)
	}
}

func TestUpdateUser_ChangesNameAndPersists(t *testing.T) {
	svc, repo := newService(t)
	repo.Seed(testutil.CreateTestUser("upd-1", "u@example.com", "Old Name"))

	updated, err := svc.UpdateUser(context.Background(), "upd-1", "New Name")
	if err != nil {
		t.Fatalf("UpdateUser err = %v", err)
	}
	if updated.Name != "New Name" {
		t.Errorf("returned name = %q, want %q", updated.Name, "New Name")
	}

	// Assert — read back through the public API confirms persistence.
	after, err := svc.GetUser(context.Background(), "upd-1")
	if err != nil {
		t.Fatalf("GetUser after Update err = %v", err)
	}
	if after.Name != "New Name" {
		t.Errorf("persisted name = %q, want %q", after.Name, "New Name")
	}
}

func TestUpdateUser_RejectsInvalidInput(t *testing.T) {
	svc, _ := newService(t)

	tests := []struct {
		name, userID, newName string
	}{
		{"empty user ID", "", "New Name"},
		{"empty name", "user-id", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.UpdateUser(context.Background(), tt.userID, tt.newName)
			if !errors.Is(err, domain.ErrInvalidInput) {
				t.Errorf("err = %v, want ErrInvalidInput", err)
			}
			if user != nil {
				t.Errorf("user = %+v, want nil", user)
			}
		})
	}
}

func TestUpdateUser_ReturnsNotFoundWhenAbsent(t *testing.T) {
	svc, _ := newService(t)

	user, err := svc.UpdateUser(context.Background(), "nobody", "New Name")

	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("err = %v, want ErrUserNotFound", err)
	}
	if user != nil {
		t.Errorf("user = %+v, want nil", user)
	}
}

func TestDeleteUser_RemovesFromRepo(t *testing.T) {
	svc, repo := newService(t)
	repo.Seed(testutil.CreateTestUser("del-1", "d@example.com", "Doomed"))

	if err := svc.DeleteUser(context.Background(), "del-1"); err != nil {
		t.Fatalf("DeleteUser err = %v", err)
	}

	// Read-back confirms removal.
	_, err := svc.GetUser(context.Background(), "del-1")
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("GetUser after Delete err = %v, want ErrUserNotFound", err)
	}
}

func TestDeleteUser_RejectsInvalidInput(t *testing.T) {
	svc, _ := newService(t)

	err := svc.DeleteUser(context.Background(), "")

	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("err = %v, want ErrInvalidInput", err)
	}
}

func TestDeleteUser_SurfacesRepositoryError(t *testing.T) {
	svc, repo := newService(t)
	repo.Seed(testutil.CreateTestUser("del-2", "d2@example.com", "User"))
	repo.FailNextDelete = domain.ErrInternalError

	err := svc.DeleteUser(context.Background(), "del-2")

	if !errors.Is(err, domain.ErrInternalError) {
		t.Errorf("err = %v, want ErrInternalError", err)
	}
}

func TestListUsers_ReturnsSeededUsers(t *testing.T) {
	svc, repo := newService(t)
	repo.Seed(testutil.CreateTestUsers(3)...)

	users, err := svc.ListUsers(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("ListUsers err = %v", err)
	}
	if len(users) != 3 {
		t.Errorf("got %d users, want 3", len(users))
	}
}

func TestListUsers_DefaultsInvalidPagination(t *testing.T) {
	svc, _ := newService(t)

	users, err := svc.ListUsers(context.Background(), -1, -5)
	if err != nil {
		t.Fatalf("ListUsers err = %v", err)
	}
	if users == nil {
		t.Error("users = nil, want non-nil empty slice")
	}
}

func TestListUsers_SurfacesRepositoryError(t *testing.T) {
	svc, repo := newService(t)
	repo.FailNextList = domain.ErrInternalError

	users, err := svc.ListUsers(context.Background(), 10, 0)

	if !errors.Is(err, domain.ErrInternalError) {
		t.Errorf("err = %v, want ErrInternalError", err)
	}
	if users != nil {
		t.Errorf("users = %v, want nil", users)
	}
}
