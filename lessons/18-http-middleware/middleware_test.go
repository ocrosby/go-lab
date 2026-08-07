package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newLogger() (*log.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return log.New(&buf, "", 0), &buf
}

func TestRequestID_GeneratesWhenAbsent(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := RequestIDFrom(r.Context()); id == "" {
			t.Error("expected request ID in context")
		}
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Request-ID"); got == "" {
		t.Error("expected X-Request-ID response header")
	}
}

func TestRequestID_TrustsUpstream(t *testing.T) {
	const upstreamID = "upstream-123"
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := RequestIDFrom(r.Context()); got != upstreamID {
			t.Errorf("context ID = %q, want %q", got, upstreamID)
		}
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", upstreamID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Request-ID"); got != upstreamID {
		t.Errorf("response X-Request-ID = %q, want %q", got, upstreamID)
	}
}

func TestLogging_RecordsStatusAndDuration(t *testing.T) {
	logger, buf := newLogger()
	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("short and stout"))
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := buf.String()
	if !strings.Contains(logLine, "status=418") {
		t.Errorf("log missing status=418: %s", logLine)
	}
	if !strings.Contains(logLine, "bytes=15") {
		t.Errorf("log missing bytes=15: %s", logLine)
	}
}

func TestRecover_ConvertsPanicTo500(t *testing.T) {
	handler := Recover(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("boom")
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "internal server error") {
		t.Errorf("body = %q, want JSON error envelope", rec.Body.String())
	}
}

func TestRecover_RethrowsAbortHandler(t *testing.T) {
	handler := Recover(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic(http.ErrAbortHandler)
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if rec := recover(); rec == nil {
			t.Error("expected Recover to rethrow http.ErrAbortHandler")
		}
	}()
	handler.ServeHTTP(rec, req)
}

func TestAuth_RejectsMissingToken(t *testing.T) {
	handler := Auth("s3cret")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("handler should not run")
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
	if got := rec.Header().Get("WWW-Authenticate"); !strings.Contains(got, "Bearer") {
		t.Errorf("WWW-Authenticate = %q, want to contain 'Bearer'", got)
	}
}

func TestAuth_AcceptsCorrectToken(t *testing.T) {
	handler := Auth("s3cret")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer s3cret")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestBodyLimit_RejectsOversizedBody(t *testing.T) {
	handler := BodyLimit(10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "too big", http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	body := strings.NewReader(strings.Repeat("x", 100))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", body)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want 413", rec.Code)
	}
}

func TestCORS_AnswersPreflight(t *testing.T) {
	handler := CORS("https://example.com")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("handler should not run for OPTIONS preflight")
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Errorf("Allow-Origin = %q", got)
	}
}

func TestChain_OrderIsFIFO(t *testing.T) {
	var order []string
	mw := func(label string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, label+" in")
				next.ServeHTTP(w, r)
				order = append(order, label+" out")
			})
		}
	}
	base := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		order = append(order, "base")
	})
	handler := Chain(base, mw("A"), mw("B"), mw("C"))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	want := []string{"A in", "B in", "C in", "base", "C out", "B out", "A out"}
	if strings.Join(order, ",") != strings.Join(want, ",") {
		t.Errorf("order = %v, want %v", order, want)
	}
}
