package after

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestServer_Health(t *testing.T) {
	server := NewServer(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", resp.Status)
	}
}

func TestServer_PanicRecovery(t *testing.T) {
	var panicHandlerCalled atomic.Bool
	panicHandler := func(recovered interface{}, r *http.Request) {
		panicHandlerCalled.Store(true)
		if r.URL.Path != "/panic" {
			t.Errorf("Expected path '/panic', got '%s'", r.URL.Path)
		}
		if recovered != "intentional panic in handler" {
			t.Errorf("Expected panic message 'intentional panic in handler', got '%v'", recovered)
		}
	}

	server := NewServer(panicHandler)
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if !panicHandlerCalled.Load() {
		t.Error("Expected panic handler to be called")
	}

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", resp.Status)
	}
}

func TestServer_SlowWithBackgroundPanic(t *testing.T) {
	var backgroundPanicHandled atomic.Bool
	panicHandler := func(recovered interface{}, r *http.Request) {
		t.Logf("Panic handled: %v", recovered)
	}

	server := NewServer(panicHandler)
	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()

	done := make(chan bool)
	go func() {
		time.Sleep(200 * time.Millisecond)
		backgroundPanicHandled.Store(true)
		done <- true
	}()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	<-done
	if !backgroundPanicHandled.Load() {
		t.Error("Background goroutine should have completed")
	}
}

func TestServer_JSON(t *testing.T) {
	tests := []struct {
		name            string
		payload         map[string]interface{}
		expectPanic     bool
		expectCode      int
		expectErrorResp bool
	}{
		{
			name:            "valid JSON without panic",
			payload:         map[string]interface{}{"data": "test"},
			expectPanic:     false,
			expectCode:      http.StatusOK,
			expectErrorResp: false,
		},
		{
			name:            "trigger panic with recovery",
			payload:         map[string]interface{}{"trigger_panic": true},
			expectPanic:     true,
			expectCode:      http.StatusInternalServerError,
			expectErrorResp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var panicHandlerCalled atomic.Bool
			panicHandler := func(recovered interface{}, r *http.Request) {
				panicHandlerCalled.Store(true)
			}

			server := NewServer(panicHandler)
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if tt.expectPanic && !panicHandlerCalled.Load() {
				t.Error("Expected panic handler to be called")
			}

			if w.Code != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, w.Code)
			}

			if tt.expectErrorResp {
				var resp ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if resp.Status != "error" {
					t.Errorf("Expected error status, got '%s'", resp.Status)
				}
			}
		})
	}
}

func TestServer_AbortHandler(t *testing.T) {
	server := NewServer(nil)
	req := httptest.NewRequest(http.MethodGet, "/abort", nil)
	w := httptest.NewRecorder()

	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected http.ErrAbortHandler to be re-panicked")
			return
		}

		if err, ok := r.(error); !ok || err != http.ErrAbortHandler {
			t.Errorf("Expected http.ErrAbortHandler, got %v", r)
		}
	}()

	server.ServeHTTP(w, req)
}

func TestServer_MultipleConcurrentPanics(t *testing.T) {
	var panicCount atomic.Int32
	panicHandler := func(recovered interface{}, r *http.Request) {
		panicCount.Add(1)
	}

	server := NewServer(panicHandler)

	concurrency := 10
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/panic", nil)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)
			done <- true
		}()
	}

	for i := 0; i < concurrency; i++ {
		<-done
	}

	if panicCount.Load() != int32(concurrency) {
		t.Errorf("Expected %d panics to be handled, got %d", concurrency, panicCount.Load())
	}
}
