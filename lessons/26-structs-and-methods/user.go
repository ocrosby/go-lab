// Package user demonstrates struct definition, exported vs unexported
// fields, value and pointer receivers, and the implicit fmt.Stringer
// interface.
package user

import "fmt"

// User is the domain type. Fields with capital first letters (ID, Name,
// Roles) are exported — accessible outside the package. `email` is
// package-private; access it only through Email() and SetEmail().
type User struct {
	ID    string
	Name  string
	email string // unexported — kept validated by the setter
	Roles []string
}

// New constructs a User with the given ID, name, and email. It's the
// only way callers outside this package can set `email`, which is what
// makes the "keep email validated" invariant enforceable.
func New(id, name, email string) *User {
	return &User{ID: id, Name: name, email: email}
}

// Email is a value-receiver method — u is a copy of the User. Fine for
// read-only accessors; the copy cost is negligible for a struct this small.
func (u User) Email() string {
	return u.email
}

// SetEmail is a POINTER-receiver method. It mutates the receiver's
// state, which only works through a pointer. If this were a value
// receiver, the caller's User would be unchanged after the call.
func (u *User) SetEmail(e string) error {
	if e == "" {
		return fmt.Errorf("email cannot be empty")
	}
	u.email = e
	return nil
}

// AddRole appends to the Roles slice on the receiver. Pointer receiver
// so the mutation is visible to the caller.
func (u *User) AddRole(r string) {
	u.Roles = append(u.Roles, r)
}

// String implements the fmt.Stringer interface implicitly — no
// declaration required. fmt looks for this method on any value it
// prints. Because User has it, fmt.Println(u) prints this string
// instead of the default {u1 Ada ada@example.com [admin]}.
func (u User) String() string {
	return fmt.Sprintf("User(%s: %s)", u.ID, u.Name)
}

// InventoryEntry demonstrates methods on a non-struct named type.
// Methods aren't struct-exclusive — you can attach them to any named
// type in the package.
type InventoryEntry int

// Priority returns "low", "medium", or "high" based on the count. A
// method on an int-based type; called as `e.Priority()`.
func (e InventoryEntry) Priority() string {
	switch {
	case e < 10:
		return "high"
	case e < 100:
		return "medium"
	default:
		return "low"
	}
}
