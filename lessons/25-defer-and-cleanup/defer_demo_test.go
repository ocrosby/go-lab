package deferdemo

import (
	"errors"
	"strings"
	"testing"
)

func TestLIFO_OrderIsReversed(t *testing.T) {
	var sb strings.Builder
	LIFO(&sb)

	got := sb.String()
	want := "body third second first "
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCaptureAtDefer_UsesValueAtDeferTime(t *testing.T) {
	var sb strings.Builder
	CaptureAtDefer(&sb)

	got := sb.String()
	// The deferred call captured x=1 at defer time; the body set x=99.
	// current= shows the body ran with x=99, captured= shows the defer
	// saw x=1.
	want := "current=99 captured=1 "
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCaptureViaClosure_UsesValueAtExecTime(t *testing.T) {
	var sb strings.Builder
	CaptureViaClosure(&sb)

	got := sb.String()
	// The closure defers a func() that reads x at exec time — by then
	// x is 99.
	want := "current=99 closure=99 "
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWrapErrorViaDefer_WrapsOnError(t *testing.T) {
	err := WrapErrorViaDefer(true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "WrapErrorViaDefer:") {
		t.Errorf("err = %v, want wrapped with 'WrapErrorViaDefer:'", err)
	}
	// The original error is still there — errors.Unwrap can reach it.
	if unwrapped := errors.Unwrap(err); unwrapped == nil {
		t.Error("wrapped error should be unwrappable to the original")
	}
}

func TestWrapErrorViaDefer_LeavesNilAlone(t *testing.T) {
	err := WrapErrorViaDefer(false)
	if err != nil {
		t.Errorf("err = %v, want nil (successful path unaffected by deferred wrap)", err)
	}
}
