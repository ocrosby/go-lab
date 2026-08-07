package errorhandling_test

import (
	"errors"
	"strings"
	"testing"

	eh "github.com/ocrosby/go-lab/lessons/24-error-handling"
)

func TestAdd_Success(t *testing.T) {
	s := eh.NewStore()
	err := s.Add(&eh.User{ID: "u1", Name: "Ada"})
	if err != nil {
		t.Fatalf("Add err = %v, want nil", err)
	}
}

func TestGet_UnknownReturnsWrappedSentinel(t *testing.T) {
	s := eh.NewStore()

	_, err := s.Get("nobody")

	// errors.Is walks the wrap chain — this catches the sentinel even
	// though errors.go wrapped it with fmt.Errorf.
	if !errors.Is(err, eh.ErrUserNotFound) {
		t.Errorf("err = %v, want to wrap ErrUserNotFound", err)
	}
	// The wrapper added context; the message includes both.
	if !strings.Contains(err.Error(), "nobody") {
		t.Errorf("err message %q, want to include the id", err.Error())
	}
}

func TestAdd_EmptyIDReturnsValidationError(t *testing.T) {
	s := eh.NewStore()
	err := s.Add(&eh.User{ID: "", Name: "Ada"})

	// errors.As extracts the typed error from the wrap chain into the
	// pointer we pass.
	var verr *eh.ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("err = %v, want a *ValidationError", err)
	}
	if verr.Field != "ID" {
		t.Errorf("Field = %q, want %q", verr.Field, "ID")
	}
}

func TestAdd_DuplicateIDReturnsWrappedSentinel(t *testing.T) {
	s := eh.NewStore()
	_ = s.Add(&eh.User{ID: "dup", Name: "First"})

	err := s.Add(&eh.User{ID: "dup", Name: "Second"})

	if !errors.Is(err, eh.ErrDuplicateID) {
		t.Errorf("err = %v, want to wrap ErrDuplicateID", err)
	}
}

func TestErrorsIs_ReturnsFalseForNil(t *testing.T) {
	// errors.Is(nil, target) is safe — returns false. This is why the
	// idiomatic check pattern doesn't need an outer nil guard.
	if errors.Is(nil, eh.ErrUserNotFound) {
		t.Error("errors.Is(nil, ...) returned true")
	}
}

func TestErrorsAs_ReturnsFalseForNil(t *testing.T) {
	var verr *eh.ValidationError
	if errors.As(nil, &verr) {
		t.Error("errors.As(nil, ...) returned true")
	}
}

// The next two tests demonstrate the ==-vs-errors.Is difference.
// Direct comparison misses wrapped sentinels; errors.Is finds them.

func TestDirectComparison_MissesWrapped(t *testing.T) {
	s := eh.NewStore()
	_, err := s.Get("nobody")

	if err == eh.ErrUserNotFound { //nolint:errorlint — demonstrating the anti-pattern
		t.Error("direct == comparison shouldn't match a wrapped error — got true")
	}
}

func TestErrorsIs_FindsWrapped(t *testing.T) {
	s := eh.NewStore()
	_, err := s.Get("nobody")

	if !errors.Is(err, eh.ErrUserNotFound) {
		t.Error("errors.Is should match through the wrap — got false")
	}
}
