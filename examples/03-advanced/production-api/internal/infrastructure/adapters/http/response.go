// Package http provides HTTP adapter implementation for the User Management API.
package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/domain"
)

// ResponseWriter provides standardized HTTP response handling
type ResponseWriter struct {
	logger *zap.Logger
}

// NewResponseWriter creates a new response writer with logger
func NewResponseWriter(logger *zap.Logger) *ResponseWriter {
	return &ResponseWriter{logger: logger}
}

// WriteJSON writes a JSON response with the given status code and data
func (rw *ResponseWriter) WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		rw.logger.Error("failed to encode JSON response", zap.Error(err))
		// If we can't encode the original data, try to send a generic error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// WriteError writes a standardized error response
func (rw *ResponseWriter) WriteError(w http.ResponseWriter, statusCode int, message string) {
	rw.WriteJSON(w, statusCode, ErrorResponse{Error: message})
}

// WriteServiceError maps domain errors to appropriate HTTP status codes and messages
func (rw *ResponseWriter) WriteServiceError(w http.ResponseWriter, err error) {
	// Use errors.Is to handle wrapped errors
	if domain.IsUserNotFoundError(err) {
		rw.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	if domain.IsUserAlreadyExistsError(err) {
		rw.WriteError(w, http.StatusConflict, "user already exists")
		return
	}

	if domain.IsInvalidInputError(err) {
		rw.WriteError(w, http.StatusBadRequest, "invalid input")
		return
	}

	// Default to internal server error for any other errors
	rw.logger.Error("internal error", zap.Error(err))
	rw.WriteError(w, http.StatusInternalServerError, "internal server error")
}

// WriteCreated writes a successful creation response (201)
func (rw *ResponseWriter) WriteCreated(w http.ResponseWriter, data any) {
	rw.WriteJSON(w, http.StatusCreated, data)
}

// WriteSuccess writes a successful response (200)
func (rw *ResponseWriter) WriteSuccess(w http.ResponseWriter, data any) {
	rw.WriteJSON(w, http.StatusOK, data)
}

// WriteNoContent writes a successful response with no content (204)
func (rw *ResponseWriter) WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
