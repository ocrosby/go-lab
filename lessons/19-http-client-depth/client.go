// Package clientdepth demonstrates production-grade patterns for
// net/http clients: custom Transport, per-request timeouts via context,
// retries via a wrapping RoundTripper, and safe redirect policy.
//
// The default http.Client is fine for scripts and demos. For long-running
// services that talk to APIs, you want each of the knobs shown here.
package clientdepth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// NewClient returns an *http.Client configured for a real service:
//
//   - Explicit Transport so we can tune connection pooling. Sharing one
//     Transport across the process is what enables keep-alive reuse and
//     the connection pool — a Client-per-call defeats it.
//   - No Client.Timeout — per-request timeouts belong on the context.
//     Client.Timeout cancels ALL work including reads of a response body,
//     which is often wrong for streaming or large responses.
//   - CheckRedirect returns http.ErrUseLastResponse so redirects are
//     surfaced to the caller. Default follows up to 10 redirects, which
//     is usually not what an API client wants (redirects can pivot from
//     the target host to an attacker-controlled one).
func NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// Get is a convenience wrapper that adds a per-request timeout via context
// and reads the whole body. Use for small responses; for streaming responses,
// return the body and let the caller close it.
func Get(ctx context.Context, client *http.Client, url string, timeout time.Duration) ([]byte, int, error) {
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http Do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read body: %w", err)
	}
	return body, resp.StatusCode, nil
}

// ErrRetriesExhausted signals that all retry attempts failed.
var ErrRetriesExhausted = errors.New("retries exhausted")
