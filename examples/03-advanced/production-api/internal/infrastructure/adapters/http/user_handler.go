// Package http provides HTTP adapter implementation for the User Management API.
package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/domain"
	"github.com/ocrosby/go-lab/projects/api/internal/utils"
)

type UserHandler struct {
	userService    domain.UserService
	logger         *zap.Logger
	responseWriter *ResponseWriter
}

func NewUserHandler(userService domain.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService:    userService,
		logger:         logger,
		responseWriter: NewResponseWriter(logger),
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/users/", h.handleUserByID)
}

func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		h.responseWriter.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *UserHandler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	if userID == "" {
		h.responseWriter.WriteError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUserByID(w, r, userID)
	case http.MethodPut:
		h.updateUser(w, r, userID)
	case http.MethodDelete:
		h.deleteUser(w, r, userID)
	default:
		h.responseWriter.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// createUser godoc
// @Summary Create a new user
// @Description Create a new user with email and name
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User creation request"
// @Success 201 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responseWriter.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.responseWriter.WriteCreated(w, user)
}

// getUserByID godoc
// @Summary Get user by ID
// @Description Get a user by their unique identifier
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request, userID string) {
	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.responseWriter.WriteSuccess(w, user)
}

// updateUser godoc
// @Summary Update user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateUserRequest true "User update request"
// @Success 200 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, userID string) {
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responseWriter.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), userID, req.Name)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.responseWriter.WriteSuccess(w, user)
}

// deleteUser godoc
// @Summary Delete user
// @Description Delete a user by their unique identifier
// @Tags users
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, userID string) {
	err := h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.responseWriter.WriteNoContent(w)
}

// listUsers godoc
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Produce json
// @Param limit query int false "Number of users to return (default: 10)"
// @Param offset query int false "Number of users to skip (default: 0)"
// @Success 200 {array} domain.User
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	pagination := utils.ParsePaginationFromQuery(limitStr, offsetStr)

	users, err := h.userService.ListUsers(r.Context(), pagination.Limit, pagination.Offset)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Create response with pagination metadata
	paginationResponse := utils.NewPaginationResponse(pagination, len(users), nil)
	response := map[string]interface{}{
		"users":      users,
		"pagination": paginationResponse,
	}

	h.responseWriter.WriteSuccess(w, response)
}

func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
	h.responseWriter.WriteServiceError(w, err)
}

type CreateUserRequest struct {
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"John Doe"`
}

type UpdateUserRequest struct {
	Name string `json:"name" example:"Jane Doe"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
