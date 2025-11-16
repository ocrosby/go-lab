// Package health provides health check functionality for Kubernetes probes and monitoring.
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ocrosby/go-lab/projects/api/internal/config"
)

type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

type Check func(ctx context.Context) error

type HealthChecker interface {
	AddCheck(name string, check Check)
	RemoveCheck(name string)
	CheckHealth(ctx context.Context) *HealthStatus
}

type HealthStatus struct {
	Status    Status                  `json:"status"`
	Timestamp time.Time               `json:"timestamp"`
	Checks    map[string]*CheckStatus `json:"checks"`
}

type CheckStatus struct {
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}

type healthChecker struct {
	checks map[string]Check
	mutex  sync.RWMutex
}

func NewHealthChecker() HealthChecker {
	return &healthChecker{
		checks: make(map[string]Check),
	}
}

func (h *healthChecker) AddCheck(name string, check Check) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.checks[name] = check
}

func (h *healthChecker) RemoveCheck(name string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.checks, name)
}

func (h *healthChecker) CheckHealth(ctx context.Context) *HealthStatus {
	h.mutex.RLock()
	checks := make(map[string]Check, len(h.checks))
	for name, check := range h.checks {
		checks[name] = check
	}
	h.mutex.RUnlock()

	status := &HealthStatus{
		Status:    StatusUp,
		Timestamp: time.Now(),
		Checks:    make(map[string]*CheckStatus),
	}

	for name, check := range checks {
		checkStatus := &CheckStatus{Status: StatusUp}

		if err := check(ctx); err != nil {
			checkStatus.Status = StatusDown
			checkStatus.Error = err.Error()
			status.Status = StatusDown
		}

		status.Checks[name] = checkStatus
	}

	return status
}

func LivenessHandler(checker HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), config.DefaultHealthCheckTimeout)
		defer cancel()

		health := checker.CheckHealth(ctx)

		w.Header().Set("Content-Type", "application/json")
		if health.Status == StatusDown {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err := json.NewEncoder(w).Encode(health); err != nil {
			http.Error(w, "Failed to encode health status", http.StatusInternalServerError)
		}
	}
}

func ReadinessHandler(checker HealthChecker) http.HandlerFunc {
	return LivenessHandler(checker)
}

func StartupHandler(checker HealthChecker) http.HandlerFunc {
	return LivenessHandler(checker)
}
