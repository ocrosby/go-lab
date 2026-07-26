// Package deferdemo shows the three timing rules of `defer` — LIFO order,
// argument evaluation at defer time, and execution just before return —
// plus the common cleanup patterns.
package deferdemo

import (
	"fmt"
	"strings"
)

// LIFO writes labels to sb in defer order. Since defers run LIFO, the
// labels appear in reverse of the order they were deferred.
func LIFO(sb *strings.Builder) {
	defer sb.WriteString("first ")  // runs LAST
	defer sb.WriteString("second ") // runs SECOND
	defer sb.WriteString("third ")  // runs FIRST
	sb.WriteString("body ")
}

// CaptureAtDefer demonstrates that a defer's arguments are evaluated at
// the DEFER statement, not at execution time. `x` was 1 when the defer
// captured it, so 1 is what gets written even though we then set x = 99.
func CaptureAtDefer(sb *strings.Builder) {
	x := 1
	defer sb.WriteString(fmt.Sprintf("captured=%d ", x))
	x = 99
	sb.WriteString(fmt.Sprintf("current=%d ", x))
}

// CaptureViaClosure wraps the defer in a func literal so `x` is captured
// by reference. When the deferred function actually runs, it reads the
// current value — 99, not 1.
func CaptureViaClosure(sb *strings.Builder) {
	x := 1
	defer func() { sb.WriteString(fmt.Sprintf("closure=%d ", x)) }()
	x = 99
	sb.WriteString(fmt.Sprintf("current=%d ", x))
}

// WrapErrorViaDefer uses defer + a named return value to wrap every error
// exit in one place. This is a common Go pattern for adding context to
// errors without repeating the wrap at each return.
func WrapErrorViaDefer(fail bool) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("WrapErrorViaDefer: %w", err)
		}
	}()

	if fail {
		return fmt.Errorf("underlying failure")
	}
	return nil
}
