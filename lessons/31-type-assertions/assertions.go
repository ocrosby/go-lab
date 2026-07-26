// Package assertions demonstrates the two type-assertion forms, the
// type switch, and the compile-time interface-satisfaction check idiom.
package assertions

import "fmt"

// Stringable is a small interface — any type with a String() method.
type Stringable interface {
	String() string
}

// UserID is a domain type that satisfies Stringable.
type UserID string

func (u UserID) String() string { return "UserID(" + string(u) + ")" }

// Compile-time interface-satisfaction check. If UserID stops satisfying
// Stringable (e.g. someone removes the String method), this line fails
// to compile — you catch the bug at the type's definition site, not at
// some caller far away. Costs zero at runtime.
var _ Stringable = UserID("")

// AsString safely extracts a string from x using the comma-ok assertion.
// Returns "", false if x isn't a string. This is the pattern you use
// 99% of the time — the single-return form (which panics on mismatch)
// is almost always wrong.
func AsString(x any) (string, bool) {
	s, ok := x.(string)
	return s, ok
}

// Describe uses a type switch to dispatch on the concrete type inside x.
// Inside each case, `v` has that case's specific type — you can call
// methods on it without a further cast.
func Describe(x any) string {
	switch v := x.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int %d", v)
	case string:
		return fmt.Sprintf("string %q", v)
	case bool:
		return fmt.Sprintf("bool %v", v)
	case Stringable:
		// Any type satisfying Stringable — could be UserID, could be
		// something else. v has type Stringable, so v.String() works.
		return "stringable " + v.String()
	default:
		return fmt.Sprintf("unknown type %T", v)
	}
}

// PromoteToString wraps `any` inputs, converting numbers to their
// decimal string form and passing strings through. Useful demo of
// selecting on multiple concrete types with different actions per case.
func PromoteToString(x any) string {
	switch v := x.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%g", v)
	default:
		return ""
	}
}
