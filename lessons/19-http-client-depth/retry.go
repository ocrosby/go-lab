package clientdepth

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// RetryTransport wraps an underlying RoundTripper to retry idempotent
// requests on transient failures (5xx and network errors). Non-idempotent
// requests (POST, PATCH) are never retried — a caller who wants that must
// send an Idempotency-Key header and use a Transport that respects it.
//
// Why RoundTripper? http.Client calls Transport.RoundTrip for every request,
// so wrapping the Transport is the way to add cross-cutting client behaviour
// (retries, tracing, auth-token injection, metrics) without changing every
// call site. It's the exact analogue of middleware on the server side.
type RetryTransport struct {
	Base    http.RoundTripper // underlying transport (nil ⇒ http.DefaultTransport)
	Max     int               // max attempts including the first try
	Backoff time.Duration     // initial backoff; doubles each attempt
}

// RoundTrip implements http.RoundTripper.
func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := rt.Base
	if base == nil {
		base = http.DefaultTransport
	}
	max := rt.Max
	if max <= 0 {
		max = 3
	}
	backoff := rt.Backoff
	if backoff <= 0 {
		backoff = 100 * time.Millisecond
	}

	var lastResp *http.Response
	var lastErr error

	for attempt := 0; attempt < max; attempt++ {
		// Clone the request to give each attempt a fresh body.
		clone := req.Clone(req.Context())
		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, err
			}
			clone.Body = body
		}

		resp, err := base.RoundTrip(clone)
		if !shouldRetry(req.Method, resp, err) {
			return resp, err
		}
		// Discard the retryable response body before the next attempt.
		if resp != nil {
			_ = resp.Body.Close()
		}
		lastResp, lastErr = resp, err

		// Wait, honoring context cancellation.
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(backoff):
		}
		backoff *= 2
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return lastResp, nil
}

// shouldRetry decides whether to make another attempt.
func shouldRetry(method string, resp *http.Response, err error) bool {
	if !isIdempotent(method) {
		return false
	}
	if err != nil {
		// Context-cancelled errors are not transient — the caller wants out.
		return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
	}
	// 5xx and 429 are retryable.
	if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
		return true
	}
	return false
}

// isIdempotent returns true for HTTP methods safe to retry. GET/HEAD/PUT/
// DELETE/OPTIONS are idempotent per RFC 7231; POST and PATCH are not.
func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut,
		http.MethodDelete, http.MethodOptions:
		return true
	}
	return false
}
