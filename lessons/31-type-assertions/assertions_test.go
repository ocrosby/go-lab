package assertions_test

import (
	"testing"

	assertions "github.com/ocrosby/go-lab/lessons/31-type-assertions"
)

func TestAsString_ExtractsStringSuccess(t *testing.T) {
	got, ok := assertions.AsString("hello")
	if !ok || got != "hello" {
		t.Errorf("AsString(\"hello\") = (%q, %v), want (\"hello\", true)", got, ok)
	}
}

func TestAsString_FailsGracefullyOnMismatch(t *testing.T) {
	got, ok := assertions.AsString(42)
	if ok {
		t.Error("AsString(42) reported ok=true")
	}
	if got != "" {
		t.Errorf("failed AsString returned %q, want empty", got)
	}
}

func TestDescribe_HandlesNil(t *testing.T) {
	if got := assertions.Describe(nil); got != "nil" {
		t.Errorf("Describe(nil) = %q, want nil", got)
	}
}

func TestDescribe_HandlesInt(t *testing.T) {
	if got := assertions.Describe(42); got != "int 42" {
		t.Errorf("Describe(42) = %q", got)
	}
}

func TestDescribe_HandlesString(t *testing.T) {
	if got := assertions.Describe("hi"); got != `string "hi"` {
		t.Errorf("Describe(\"hi\") = %q", got)
	}
}

func TestDescribe_HandlesStringable(t *testing.T) {
	// UserID satisfies the Stringable interface case in the switch —
	// the interface case matches after the concrete cases fail.
	uid := assertions.UserID("u-42")
	got := assertions.Describe(uid)
	want := "stringable UserID(u-42)"
	if got != want {
		t.Errorf("Describe(UserID) = %q, want %q", got, want)
	}
}

func TestPromoteToString_ConvertsNumericTypes(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{"already a string", "already a string"},
		{42, "42"},
		{3.14, "3.14"},
		{true, ""}, // bool has no case → default returns ""
	}

	for _, tt := range tests {
		got := assertions.PromoteToString(tt.in)
		if got != tt.want {
			t.Errorf("PromoteToString(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
