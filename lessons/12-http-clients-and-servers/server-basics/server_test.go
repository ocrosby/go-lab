package serverbasics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHandleRoot uses httptest.NewRecorder — the fastest way to test a
// handler. No socket, no port, no goroutine — just call the handler with a
// synthetic request/response pair and assert on the recorder.
func TestHandleRoot(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handleRoot(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/plain") {
		t.Errorf("Content-Type = %q, want text/plain*", got)
	}
	if !strings.Contains(rec.Body.String(), "go-lab HTTP server") {
		t.Errorf("body = %q, want it to mention go-lab", rec.Body.String())
	}
}

// TestHandleHello_DefaultName checks the query-parameter fallback.
func TestHandleHello_DefaultName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()

	handleHello(rec, req)

	var got Greeting
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Message != "hello, world" {
		t.Errorf("Message = %q, want %q", got.Message, "hello, world")
	}
}

// TestHandleHello_WithName verifies the query parameter is honored.
func TestHandleHello_WithName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?name=omar", nil)
	rec := httptest.NewRecorder()

	handleHello(rec, req)

	var got Greeting
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Message != "hello, omar" {
		t.Errorf("Message = %q, want %q", got.Message, "hello, omar")
	}
}

// TestServer_WiredEndToEnd uses httptest.NewServer — spins up a real listener
// on a random port, so the test exercises the ServeMux and both handlers as
// they'd behave in production. Slower than NewRecorder but higher-fidelity.
func TestServer_WiredEndToEnd(t *testing.T) {
	srv := httptest.NewServer(NewServer(":0").Handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/hello?name=go-lab")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	var got Greeting
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Message != "hello, go-lab" {
		t.Errorf("Message = %q, want %q", got.Message, "hello, go-lab")
	}
}
