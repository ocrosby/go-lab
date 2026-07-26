package clientdepth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// TestGet_HappyPath uses httptest.NewServer for a fully-integrated test.
func TestGet_HappyPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}))
	defer srv.Close()

	body, status, err := Get(t.Context(), srv.Client(), srv.URL, 2*time.Second)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("status = %d, want 200", status)
	}
	if string(body) != "pong" {
		t.Errorf("body = %q, want pong", string(body))
	}
}

// TestGet_TimeoutViaContext demonstrates the per-request timeout pattern.
// The server sleeps longer than the timeout — the client must give up.
func TestGet_TimeoutViaContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte("late"))
	}))
	defer srv.Close()

	_, _, err := Get(t.Context(), srv.Client(), srv.URL, 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "deadline exceeded") &&
		!strings.Contains(err.Error(), "context deadline") {
		t.Errorf("error = %v, want a deadline error", err)
	}
}

// countingRT is a test-double RoundTripper — the pattern for stubbing
// HTTP responses without an httptest.Server. Way cheaper than spinning
// up a real listener, and gives you full control over what the client sees.
type countingRT struct {
	calls   atomic.Int32
	respond func(int32) (*http.Response, error)
}

func (c *countingRT) RoundTrip(req *http.Request) (*http.Response, error) {
	n := c.calls.Add(1)
	return c.respond(n)
}

func makeResp(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       http.NoBody,
		Header:     http.Header{},
		Request:    &http.Request{},
		// Body doesn't matter for retry-count tests but must be closable.
	}
}

func TestRetryTransport_RetriesOn5xx(t *testing.T) {
	rt := &countingRT{
		respond: func(n int32) (*http.Response, error) {
			if n < 3 {
				return makeResp(http.StatusInternalServerError, ""), nil
			}
			return makeResp(http.StatusOK, ""), nil
		},
	}
	client := &http.Client{Transport: &RetryTransport{Base: rt, Max: 5, Backoff: time.Millisecond}}
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://x", nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if got := rt.calls.Load(); got != 3 {
		t.Errorf("attempts = %d, want 3", got)
	}
}

func TestRetryTransport_DoesNotRetryPOST(t *testing.T) {
	rt := &countingRT{
		respond: func(int32) (*http.Response, error) {
			return makeResp(http.StatusInternalServerError, ""), nil
		},
	}
	client := &http.Client{Transport: &RetryTransport{Base: rt, Max: 5, Backoff: time.Millisecond}}
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://x", http.NoBody)

	_, _ = client.Do(req)

	if got := rt.calls.Load(); got != 1 {
		t.Errorf("POST attempts = %d, want 1 (POST is non-idempotent)", got)
	}
}

func TestRetryTransport_HonorsContextCancel(t *testing.T) {
	rt := &countingRT{
		respond: func(int32) (*http.Response, error) {
			return makeResp(http.StatusServiceUnavailable, ""), nil
		},
	}
	client := &http.Client{Transport: &RetryTransport{Base: rt, Max: 100, Backoff: 50 * time.Millisecond}}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://x", nil)
	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected cancellation error, got nil")
	}
}

func TestNewClient_ReturnsUseLastResponseOnRedirect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			w.Header().Set("Location", "/dest")
			w.WriteHeader(http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient()
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL+"/redirect", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("status = %d, want 302 (redirect not followed)", resp.StatusCode)
	}
	if got := resp.Header.Get("Location"); got != "/dest" {
		t.Errorf("Location = %q, want /dest", got)
	}
}
