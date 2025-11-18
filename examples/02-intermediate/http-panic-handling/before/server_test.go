package before

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_Health(t *testing.T) {
	server := NewServer()
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

func TestServer_Panic(t *testing.T) {
	server := NewServer()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic to occur, but it didn't")
		}
	}()

	server.ServeHTTP(w, req)
}

func TestServer_SlowWithBackgroundPanic(t *testing.T) {
	server := NewServer()
	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	time.Sleep(200 * time.Millisecond)
}

func TestServer_JSON(t *testing.T) {
	server := NewServer()

	tests := []struct {
		name        string
		payload     map[string]interface{}
		expectPanic bool
		expectCode  int
	}{
		{
			name:        "valid JSON without panic",
			payload:     map[string]interface{}{"data": "test"},
			expectPanic: false,
			expectCode:  http.StatusOK,
		},
		{
			name:        "trigger panic",
			payload:     map[string]interface{}{"trigger_panic": true},
			expectPanic: true,
			expectCode:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic to occur, but it didn't")
					}
				}()
			}

			server.ServeHTTP(w, req)

			if !tt.expectPanic && w.Code != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, w.Code)
			}
		})
	}
}
