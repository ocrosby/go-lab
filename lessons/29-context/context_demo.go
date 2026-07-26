// Package contextdemo demonstrates context.WithCancel, WithTimeout, and
// WithValue via a small piece of "work" that respects cancellation.
package contextdemo

import (
	"context"
	"time"
)

// requestIDKey is the unexported key type used to attach and retrieve
// a request ID from a context. Using a private type here means no other
// package can accidentally pass the same key.
type requestIDKey struct{}

// WithRequestID returns a child context carrying the request ID.
func WithRequestID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, requestIDKey{}, id)
}

// RequestIDFrom retrieves the request ID from ctx, or "" if none is
// attached. Comma-ok pattern on the type assertion protects against a
// non-string value under the same key (shouldn't happen, but defensive).
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// DoWork simulates work that takes `duration` to complete, but respects
// context cancellation. This is the "long-running operation that can
// be cancelled early" pattern.
//
// Returns nil if the work completed on time, or ctx.Err() if the context
// cancelled first (context.Canceled or context.DeadlineExceeded).
func DoWork(ctx context.Context, duration time.Duration) error {
	// A ticker to simulate incremental work. A real function might be
	// reading from a channel or doing a chunk of computation per tick.
	select {
	case <-ctx.Done():
		// Cancelled before work started.
		return ctx.Err()
	case <-time.After(duration):
		// Work finished on time.
		return nil
	}
}

// StreamCount emits integers from 1 to n on the returned channel, one
// per tick of interval, stopping early if the context is cancelled.
// The channel is closed when the stream ends (successfully or via
// cancellation), so callers can `for v := range ch` safely.
func StreamCount(ctx context.Context, n int, interval time.Duration) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for i := 1; i <= n; i++ {
			select {
			case <-ctx.Done():
				// Caller cancelled — stop sending. Deferred close still
				// fires so the range loop terminates.
				return
			case <-ticker.C:
				select {
				case out <- i:
				case <-ctx.Done():
					// Even the send needs to respect cancellation, in
					// case no one is receiving.
					return
				}
			}
		}
	}()
	return out
}
