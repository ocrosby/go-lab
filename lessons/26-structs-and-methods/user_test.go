package user_test

import (
	"testing"

	user "github.com/ocrosby/go-lab/lessons/26-structs-and-methods"
)

func TestNew_ConstructsUserWithFields(t *testing.T) {
	u := user.New("u1", "Ada", "ada@example.com")

	if u.ID != "u1" || u.Name != "Ada" {
		t.Errorf("ID=%q Name=%q, want u1/Ada", u.ID, u.Name)
	}
	if got := u.Email(); got != "ada@example.com" {
		t.Errorf("Email() = %q, want ada@example.com", got)
	}
}

func TestSetEmail_MutatesTheReceiver(t *testing.T) {
	u := user.New("u1", "Ada", "old@example.com")

	if err := u.SetEmail("new@example.com"); err != nil {
		t.Fatalf("SetEmail err = %v", err)
	}

	// The pointer receiver means the mutation stuck — read back through
	// the exported accessor confirms.
	if got := u.Email(); got != "new@example.com" {
		t.Errorf("Email after SetEmail = %q, want new@example.com", got)
	}
}

func TestSetEmail_RejectsEmpty(t *testing.T) {
	u := user.New("u1", "Ada", "ada@example.com")

	err := u.SetEmail("")

	if err == nil {
		t.Error("SetEmail(\"\") returned nil, want error")
	}
	// Email should be unchanged after the failed setter.
	if got := u.Email(); got != "ada@example.com" {
		t.Errorf("Email after failed SetEmail = %q, want ada@example.com", got)
	}
}

func TestAddRole_AppendsToSlice(t *testing.T) {
	u := user.New("u1", "Ada", "ada@example.com")

	u.AddRole("admin")
	u.AddRole("billing")

	if len(u.Roles) != 2 || u.Roles[0] != "admin" || u.Roles[1] != "billing" {
		t.Errorf("Roles = %v, want [admin billing]", u.Roles)
	}
}

func TestString_ImplicitlyImplementsStringer(t *testing.T) {
	u := user.New("u1", "Ada", "ada@example.com")

	got := u.String()
	want := "User(u1: Ada)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

// The InventoryEntry tests exercise methods on a NON-struct named type.
// Same syntax, same rules — methods aren't struct-exclusive.

func TestInventoryEntry_LowStockIsHighPriority(t *testing.T) {
	if got := user.InventoryEntry(5).Priority(); got != "high" {
		t.Errorf("Priority() for count 5 = %q, want high", got)
	}
}

func TestInventoryEntry_HealthyStockIsLowPriority(t *testing.T) {
	if got := user.InventoryEntry(500).Priority(); got != "low" {
		t.Errorf("Priority() for count 500 = %q, want low", got)
	}
}
