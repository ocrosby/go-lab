// Package errorhandling demonstrates Go's error idioms on a tiny user
// store: sentinel errors, wrapping with fmt.Errorf %w, and a custom
// error type carrying structured information.
package errorhandling

import (
	"errors"
	"fmt"
)

// --- Sentinel errors ---------------------------------------------------
// Exported errors callers can match against with errors.Is. Naming
// convention: package-level var, name starts with Err.

var (
	// ErrUserNotFound is returned by Get when no user matches the ID.
	ErrUserNotFound = errors.New("user not found")

	// ErrDuplicateID is returned by Add when a user with the given ID
	// already exists.
	ErrDuplicateID = errors.New("duplicate user ID")
)

// --- Custom error type -------------------------------------------------
// Use a struct when the caller needs structured data about the failure —
// which field failed, what the invalid value was, an HTTP status code, etc.

// ValidationError carries the specific field and message so callers can
// programmatically show the right form-field error, log by field name,
// or map to a 422 response body.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface. Any type with Error() string IS
// an error — no explicit interface declaration required.
func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", v.Field, v.Message)
}

// --- The store ---------------------------------------------------------

type User struct {
	ID   string
	Name string
}

// Store is a tiny in-memory user store used to demonstrate the error
// idioms. Nothing here is production-quality — the goal is to show the
// error patterns.
type Store struct {
	users map[string]*User
}

func NewStore() *Store {
	return &Store{users: map[string]*User{}}
}

// Add inserts a user. Returns:
//   - *ValidationError wrapped with context if the input is invalid
//   - ErrDuplicateID wrapped with context if a user with that ID exists
//   - nil on success
func (s *Store) Add(u *User) error {
	if u.ID == "" {
		// Wrap a *ValidationError with fmt.Errorf %w — the caller can
		// still errors.As the wrapper into a *ValidationError below.
		return fmt.Errorf("Add: %w", &ValidationError{
			Field:   "ID",
			Message: "must be non-empty",
		})
	}
	if u.Name == "" {
		return fmt.Errorf("Add: %w", &ValidationError{
			Field:   "Name",
			Message: "must be non-empty",
		})
	}
	if _, exists := s.users[u.ID]; exists {
		// Wrap a sentinel with %w — errors.Is finds it through the wrap.
		return fmt.Errorf("Add: user %q: %w", u.ID, ErrDuplicateID)
	}
	s.users[u.ID] = u
	return nil
}

// Get returns the user for the given ID, or a wrapped ErrUserNotFound.
func (s *Store) Get(id string) (*User, error) {
	u, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("Get: user %q: %w", id, ErrUserNotFound)
	}
	return u, nil
}
