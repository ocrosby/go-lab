package contextdemo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	cd "github.com/ocrosby/go-lab/lessons/30-context"
)

func TestDoWork_CompletesWhenFast(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
	defer cancel()

	err := cd.DoWork(ctx, 10*time.Millisecond)

	if err != nil {
		t.Errorf("DoWork err = %v, want nil", err)
	}
}

func TestDoWork_ReturnsDeadlineExceededWhenSlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()

	err := cd.DoWork(ctx, 500*time.Millisecond)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
}

func TestDoWork_ReturnsCanceledWhenExplicit(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	// Cancel synchronously before starting the work.
	cancel()

	err := cd.DoWork(ctx, 500*time.Millisecond)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestStreamCount_EmitsRequestedValues(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
	defer cancel()

	out := cd.StreamCount(ctx, 3, 10*time.Millisecond)

	var got []int
	for v := range out {
		got = append(got, v)
	}

	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Errorf("stream = %v, want [1 2 3]", got)
	}
}

func TestStreamCount_StopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	// Cancel after two values.
	go func() {
		time.Sleep(25 * time.Millisecond)
		cancel()
	}()

	out := cd.StreamCount(ctx, 100, 10*time.Millisecond)

	count := 0
	for range out {
		count++
	}

	// We should get some but not all.
	if count == 0 {
		t.Error("received zero values before cancel")
	}
	if count == 100 {
		t.Error("received all 100 values — cancel was ignored")
	}
}

func TestWithRequestID_RoundTrips(t *testing.T) {
	parent := t.Context()

	ctx := cd.WithRequestID(parent, "req-42")

	if got := cd.RequestIDFrom(ctx); got != "req-42" {
		t.Errorf("RequestIDFrom = %q, want req-42", got)
	}
	// The parent context is untouched — WithValue returns a child.
	if got := cd.RequestIDFrom(parent); got != "" {
		t.Errorf("parent RequestIDFrom = %q, want empty (parent untouched)", got)
	}
}
