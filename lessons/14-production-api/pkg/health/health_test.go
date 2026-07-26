package health

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHealthChecker(t *testing.T) {
	checker := NewHealthChecker()
	if checker == nil {
		t.Error("Expected health checker, got nil")
	}
}

func TestHealthChecker_AddCheck(t *testing.T) {
	checker := NewHealthChecker()

	checkExecuted := false
	check := func(ctx context.Context) error {
		checkExecuted = true
		return nil
	}

	checker.AddCheck("test", check)

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	if !checkExecuted {
		t.Error("Expected check to be executed")
	}

	if status.Status != StatusUp {
		t.Errorf("Expected status UP, got %s", status.Status)
	}

	if len(status.Checks) != 1 {
		t.Errorf("Expected 1 check, got %d", len(status.Checks))
	}

	if status.Checks["test"].Status != StatusUp {
		t.Errorf("Expected check status UP, got %s", status.Checks["test"].Status)
	}
}

func TestHealthChecker_RemoveCheck(t *testing.T) {
	checker := NewHealthChecker()

	check := func(ctx context.Context) error {
		return nil
	}

	checker.AddCheck("test", check)
	checker.RemoveCheck("test")

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	if len(status.Checks) != 0 {
		t.Errorf("Expected 0 checks after removal, got %d", len(status.Checks))
	}
}

func TestHealthChecker_CheckHealth_Success(t *testing.T) {
	checker := NewHealthChecker()

	check1 := func(ctx context.Context) error { return nil }
	check2 := func(ctx context.Context) error { return nil }

	checker.AddCheck("check1", check1)
	checker.AddCheck("check2", check2)

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	if status.Status != StatusUp {
		t.Errorf("Expected status UP, got %s", status.Status)
	}

	if len(status.Checks) != 2 {
		t.Errorf("Expected 2 checks, got %d", len(status.Checks))
	}

	for name, check := range status.Checks {
		if check.Status != StatusUp {
			t.Errorf("Expected check %s status UP, got %s", name, check.Status)
		}
		if check.Error != "" {
			t.Errorf("Expected no error for check %s, got %s", name, check.Error)
		}
	}
}

func TestHealthChecker_CheckHealth_Failure(t *testing.T) {
	checker := NewHealthChecker()

	testError := errors.New("test error")
	check1 := func(ctx context.Context) error { return nil }
	check2 := func(ctx context.Context) error { return testError }

	checker.AddCheck("check1", check1)
	checker.AddCheck("check2", check2)

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	if status.Status != StatusDown {
		t.Errorf("Expected status DOWN, got %s", status.Status)
	}

	if status.Checks["check1"].Status != StatusUp {
		t.Errorf("Expected check1 status UP, got %s", status.Checks["check1"].Status)
	}

	if status.Checks["check2"].Status != StatusDown {
		t.Errorf("Expected check2 status DOWN, got %s", status.Checks["check2"].Status)
	}

	if status.Checks["check2"].Error != testError.Error() {
		t.Errorf("Expected check2 error %s, got %s", testError.Error(), status.Checks["check2"].Error)
	}
}

func TestHealthChecker_CheckHealth_Empty(t *testing.T) {
	checker := NewHealthChecker()

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	if status.Status != StatusUp {
		t.Errorf("Expected status UP for empty checks, got %s", status.Status)
	}

	if len(status.Checks) != 0 {
		t.Errorf("Expected 0 checks, got %d", len(status.Checks))
	}

	if status.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}

func TestHealthChecker_ConcurrentAccess(t *testing.T) {
	checker := NewHealthChecker()

	// Add some initial checks
	for i := 0; i < 5; i++ {
		checkName := fmt.Sprintf("check%d", i)
		checker.AddCheck(checkName, func(ctx context.Context) error { return nil })
	}

	done := make(chan bool, 3)

	// Goroutine 1: Add/Remove checks
	go func() {
		for i := 5; i < 10; i++ {
			checkName := fmt.Sprintf("check%d", i)
			checker.AddCheck(checkName, func(ctx context.Context) error { return nil })
			time.Sleep(time.Millisecond)
			checker.RemoveCheck(checkName)
		}
		done <- true
	}()

	// Goroutine 2: Check health
	go func() {
		for i := 0; i < 20; i++ {
			ctx := context.Background()
			checker.CheckHealth(ctx)
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Goroutine 3: Add/Remove checks
	go func() {
		for i := 10; i < 15; i++ {
			checkName := fmt.Sprintf("check%d", i)
			checker.AddCheck(checkName, func(ctx context.Context) error { return nil })
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestLivenessHandler_Success(t *testing.T) {
	checker := NewHealthChecker()
	checker.AddCheck("test", func(ctx context.Context) error { return nil })

	handler := LivenessHandler(checker)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", contentType)
	}

	var response HealthStatus
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Status != StatusUp {
		t.Errorf("Expected status UP, got %s", response.Status)
	}
}

func TestLivenessHandler_Failure(t *testing.T) {
	checker := NewHealthChecker()
	checker.AddCheck("test", func(ctx context.Context) error {
		return errors.New("test error")
	})

	handler := LivenessHandler(checker)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var response HealthStatus
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Status != StatusDown {
		t.Errorf("Expected status DOWN, got %s", response.Status)
	}
}

func TestReadinessHandler(t *testing.T) {
	checker := NewHealthChecker()
	readinessHandler := ReadinessHandler(checker)
	livenessHandler := LivenessHandler(checker)

	// Both handlers should behave identically
	req := httptest.NewRequest("GET", "/readyz", nil)

	w1 := httptest.NewRecorder()
	readinessHandler(w1, req)

	w2 := httptest.NewRecorder()
	livenessHandler(w2, req)

	if w1.Code != w2.Code {
		t.Errorf("Expected same status codes, got readiness=%d, liveness=%d", w1.Code, w2.Code)
	}
}

func TestStartupHandler(t *testing.T) {
	checker := NewHealthChecker()
	startupHandler := StartupHandler(checker)
	livenessHandler := LivenessHandler(checker)

	// Both handlers should behave identically
	req := httptest.NewRequest("GET", "/startupz", nil)

	w1 := httptest.NewRecorder()
	startupHandler(w1, req)

	w2 := httptest.NewRecorder()
	livenessHandler(w2, req)

	if w1.Code != w2.Code {
		t.Errorf("Expected same status codes, got startup=%d, liveness=%d", w1.Code, w2.Code)
	}
}
