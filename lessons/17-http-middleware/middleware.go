// Package middleware demonstrates the standard net/http middleware pattern:
// a function that takes an http.Handler and returns an http.Handler,
// composed into a chain around a base handler.
//
//	base := http.HandlerFunc(myHandler)
//	final := Chain(base, RequestID, Logging, Recover, Auth("secret"),
//	              BodyLimit(1<<20), CORS)
//
// Each middleware sees the request on the way in and the response on the way
// out. Order matters — Recover should wrap Auth (so a panic during auth is
// still logged) but be wrapped by RequestID (so panic logs include the ID).
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// Middleware is the alias almost every Go HTTP codebase uses. Once you know
// this shape, you can read any router or framework in the ecosystem.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares to h in the order given: the first entry is the
// outermost wrapper (runs first on the way in, last on the way out).
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	// Walk backwards so the first-listed middleware ends up outermost.
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// --- Request ID ------------------------------------------------------------

type ctxKey int

const requestIDKey ctxKey = iota

// RequestID attaches a random request ID to every request. If the client
// sent X-Request-ID, that value is trusted (common when a load balancer or
// gateway sets it upstream); otherwise generate a new one. The ID is echoed
// back in the response header and stashed in the context for downstream
// handlers and log lines.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = newRequestID()
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFrom pulls the ID out of the context. Returns "" if the middleware
// didn't run for this request — handlers should tolerate that gracefully.
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

func newRequestID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// --- Logging ---------------------------------------------------------------

// statusWriter wraps http.ResponseWriter to capture the status code the
// handler wrote. Handlers that never call WriteHeader default to 200.
type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	if sw.status == 0 {
		sw.status = http.StatusOK
	}
	n, err := sw.ResponseWriter.Write(b)
	sw.bytes += n
	return n, err
}

// Logging emits a structured line per request: method, path, status, size,
// duration, and (if present) request ID. Real code would use log/slog.
func Logging(logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusWriter{ResponseWriter: w}
			next.ServeHTTP(sw, r)
			logger.Printf("method=%s path=%s status=%d bytes=%d duration=%s request_id=%s",
				r.Method, r.URL.Path, sw.status, sw.bytes,
				time.Since(start), RequestIDFrom(r.Context()))
		})
	}
}

// --- Recover ---------------------------------------------------------------

// Recover catches panics that would otherwise crash the process (lesson 10)
// and returns 500 to the client. Rethrows the special http.ErrAbortHandler
// which the net/http package uses to signal "abort quietly" — recovering
// that would swallow legitimate server-side aborts.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(rec)
				}
				log.Printf("panic recovered: %v\n%s", rec, debug.Stack())
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// --- Auth ------------------------------------------------------------------

// Auth requires a valid Bearer token in the Authorization header. This is the
// simplest possible token check — real code validates a JWT or looks up
// an opaque token in a store. Uses constant-time comparison to avoid the
// timing side-channel that plain == exposes.
func Auth(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			presented, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || !constantTimeStringEqual(presented, token) {
				w.Header().Set("WWW-Authenticate", `Bearer realm="api"`)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// constantTimeStringEqual compares two strings in time proportional to the
// longer one rather than to the first mismatch. Wrap subtle.ConstantTimeCompare
// but with length equalization so different-length inputs don't leak length.
func constantTimeStringEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var diff byte
	for i := 0; i < len(a); i++ {
		diff |= a[i] ^ b[i]
	}
	return diff == 0
}

// --- Body limit ------------------------------------------------------------

// BodyLimit caps the number of bytes any handler can read from r.Body via
// http.MaxBytesReader. An unbounded body is an OOM vector and a real REST
// API sets a limit on every request. Applies globally when used as
// middleware; use it per-handler with more targeted limits (see lesson 16's
// decodeJSON) for endpoints with different needs.
func BodyLimit(max int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, max)
			next.ServeHTTP(w, r)
		})
	}
}

// --- CORS ------------------------------------------------------------------

// CORS adds the minimum headers needed for a browser client from
// allowedOrigin to talk to this API. Real production CORS handles origin
// allowlists, credentialed requests, preflight OPTIONS, and Vary headers —
// use a library (rs/cors) for anything beyond this demo.
func CORS(allowedOrigin string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Vary", "Origin")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// --- Demo assembly ---------------------------------------------------------

// Handler returns a fully-wired handler suitable for demos and tests: five
// middleware wrapped around a base that responds with the request ID.
func Handler(logger *log.Logger, token string) http.Handler {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "handled request %s\n", RequestIDFrom(r.Context()))
	})
	return Chain(base,
		RequestID,
		Logging(logger),
		Recover,
		CORS("*"),
		BodyLimit(1<<20),
		Auth(token),
	)
}
