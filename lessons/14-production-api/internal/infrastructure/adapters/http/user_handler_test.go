package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ocrosby/go-lab/lessons/14-production-api/internal/application"
	"github.com/ocrosby/go-lab/lessons/14-production-api/internal/domain"
	httpAdapter "github.com/ocrosby/go-lab/lessons/14-production-api/internal/infrastructure/adapters/http"
	"github.com/ocrosby/go-lab/lessons/14-production-api/internal/testutil"
)

// Black-box tests (see rules/black-box-testing.md). Every request goes
// through the mux so the tests exercise the same routing the running server
// uses. The collaborators are real — application.NewUserService wired to a
// testutil.FakeUserRepository — so the tests survive any refactor of the
// service's internals or the handler's method decomposition.
//
// The mock UserService was removed intentionally; asserting on
// mockService.EXPECT().CreateUser(...).Return(...) pinned the interaction
// shape and broke on every service refactor.

type harness struct {
	mux  *http.ServeMux
	repo *testutil.FakeUserRepository
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	repo := testutil.NewFakeUserRepository()
	svc := application.NewUserService(repo, testutil.NewTestLogger())
	handler := httpAdapter.NewUserHandler(svc, testutil.NewTestLogger())
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return &harness{mux: mux, repo: repo}
}

func (h *harness) do(t *testing.T, method, target string, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body == "" {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader([]byte(body))
	}
	req := httptest.NewRequestWithContext(t.Context(), method, target, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.mux.ServeHTTP(rec, req)
	return rec
}

func decode[T any](t *testing.T, body []byte) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("decode: %v; body=%s", err, string(body))
	}
	return out
}

func TestPOSTUsers_CreatesAndPersists(t *testing.T) {
	h := newHarness(t)

	rec := h.do(t, http.MethodPost, "/users", `{"email":"new@example.com","name":"New User"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body)
	}
	user := decode[domain.User](t, rec.Body.Bytes())
	if user.Email != "new@example.com" {
		t.Errorf("email = %q, want %q", user.Email, "new@example.com")
	}

	// Read-back through the public API confirms persistence — the assertion
	// the previous mock-based test could not make.
	get := h.do(t, http.MethodGet, "/users/"+user.ID, "")
	if get.Code != http.StatusOK {
		t.Errorf("GET after POST status = %d, want 200", get.Code)
	}
}

func TestPOSTUsers_RejectsInvalidJSON(t *testing.T) {
	h := newHarness(t)
	rec := h.do(t, http.MethodPost, "/users", `{ this is not json }`)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestPOSTUsers_ReturnsConflictOnDuplicateEmail(t *testing.T) {
	h := newHarness(t)

	first := h.do(t, http.MethodPost, "/users", `{"email":"dup@example.com","name":"First"}`)
	if first.Code != http.StatusCreated {
		t.Fatalf("first POST status = %d, want 201", first.Code)
	}
	second := h.do(t, http.MethodPost, "/users", `{"email":"dup@example.com","name":"Second"}`)

	if second.Code != http.StatusConflict {
		t.Errorf("second POST status = %d, want 409", second.Code)
	}
}

func TestGETUsersByID_ReturnsSeededUser(t *testing.T) {
	h := newHarness(t)
	h.repo.Seed(testutil.CreateTestUser("u-1", "u@example.com", "Seeded"))

	rec := h.do(t, http.MethodGet, "/users/u-1", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	user := decode[domain.User](t, rec.Body.Bytes())
	if user.ID != "u-1" {
		t.Errorf("id = %q, want u-1", user.ID)
	}
}

func TestGETUsersByID_ReturnsNotFound(t *testing.T) {
	h := newHarness(t)

	rec := h.do(t, http.MethodGet, "/users/nobody", "")

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestPUTUsers_UpdatesAndPersists(t *testing.T) {
	h := newHarness(t)
	h.repo.Seed(testutil.CreateTestUser("u-2", "u2@example.com", "Old Name"))

	rec := h.do(t, http.MethodPut, "/users/u-2", `{"name":"New Name"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
	}
	user := decode[domain.User](t, rec.Body.Bytes())
	if user.Name != "New Name" {
		t.Errorf("returned name = %q, want %q", user.Name, "New Name")
	}

	// Persistence check via a fresh GET.
	get := h.do(t, http.MethodGet, "/users/u-2", "")
	after := decode[domain.User](t, get.Body.Bytes())
	if after.Name != "New Name" {
		t.Errorf("persisted name = %q, want %q", after.Name, "New Name")
	}
}

func TestPUTUsers_RejectsInvalidJSON(t *testing.T) {
	h := newHarness(t)
	h.repo.Seed(testutil.CreateTestUser("u-3", "u3@example.com", "User"))

	rec := h.do(t, http.MethodPut, "/users/u-3", `not json`)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestDELETEUsers_RemovesAndReturns204(t *testing.T) {
	h := newHarness(t)
	h.repo.Seed(testutil.CreateTestUser("del-1", "d@example.com", "Doomed"))

	rec := h.do(t, http.MethodDelete, "/users/del-1", "")

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("204 must have empty body, got %d bytes", rec.Body.Len())
	}

	// Subsequent GET now 404s.
	get := h.do(t, http.MethodGet, "/users/del-1", "")
	if get.Code != http.StatusNotFound {
		t.Errorf("GET after DELETE status = %d, want 404", get.Code)
	}
}

func TestDELETEUsers_ReturnsNotFoundForAbsent(t *testing.T) {
	h := newHarness(t)

	rec := h.do(t, http.MethodDelete, "/users/nobody", "")

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

// listResponse mirrors the wire shape of GET /users — the handler returns a
// paginated envelope, not a bare array.
type listResponse struct {
	Users      []domain.User `json:"users"`
	Pagination struct {
		Limit   int  `json:"limit"`
		Offset  int  `json:"offset"`
		HasNext bool `json:"has_next"`
		HasPrev bool `json:"has_prev"`
	} `json:"pagination"`
}

func TestGETUsers_ReturnsSeededList(t *testing.T) {
	h := newHarness(t)
	h.repo.Seed(testutil.CreateTestUsers(3)...)

	rec := h.do(t, http.MethodGet, "/users", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	out := decode[listResponse](t, rec.Body.Bytes())
	if len(out.Users) != 3 {
		t.Errorf("got %d users, want 3", len(out.Users))
	}
}

func TestGETUsers_HonorsPagination(t *testing.T) {
	h := newHarness(t)
	h.repo.Seed(testutil.CreateTestUsers(5)...)

	rec := h.do(t, http.MethodGet, "/users?limit=2&offset=1", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	out := decode[listResponse](t, rec.Body.Bytes())
	if len(out.Users) != 2 {
		t.Errorf("got %d users, want 2 (limit=2)", len(out.Users))
	}
	if out.Pagination.Limit != 2 || out.Pagination.Offset != 1 {
		t.Errorf("pagination = %+v, want limit=2 offset=1", out.Pagination)
	}
}

func TestGETUsers_FallsBackWhenPaginationInvalid(t *testing.T) {
	h := newHarness(t)
	h.repo.Seed(testutil.CreateTestUsers(2)...)

	rec := h.do(t, http.MethodGet, "/users?limit=invalid&offset=nope", "")

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 (invalid pagination should default)", rec.Code)
	}
}

func TestMethodNotAllowedOnCollection(t *testing.T) {
	h := newHarness(t)

	rec := h.do(t, http.MethodPatch, "/users", "")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}

func TestMethodNotAllowedOnItem(t *testing.T) {
	h := newHarness(t)

	rec := h.do(t, http.MethodPatch, "/users/some-id", "")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}

func TestGETUsers_EmptyIDIsBadRequest(t *testing.T) {
	h := newHarness(t)

	rec := h.do(t, http.MethodGet, "/users/", "")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 (empty ID)", rec.Code)
	}
}

// Service-error mapping: each domain error type produces the correct HTTP
// status via an arranged repo state (or fault injection where the state
// alone can't produce the error). This replaces the previous white-box
// test that called handler.handleServiceError directly.
func TestServiceErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*harness)
		method     string
		target     string
		body       string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "not found → 404",
			setup:      func(_ *harness) {},
			method:     http.MethodGet,
			target:     "/users/absent",
			wantStatus: http.StatusNotFound,
			wantMsg:    "user not found",
		},
		{
			name: "already exists → 409",
			setup: func(h *harness) {
				h.repo.Seed(testutil.CreateTestUser("existing", "same@example.com", "Existing"))
			},
			method:     http.MethodPost,
			target:     "/users",
			body:       `{"email":"same@example.com","name":"Dupe"}`,
			wantStatus: http.StatusConflict,
			wantMsg:    "user already exists",
		},
		{
			name:       "invalid input → 400",
			setup:      func(_ *harness) {},
			method:     http.MethodPost,
			target:     "/users",
			body:       `{"email":"","name":""}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "invalid input",
		},
		{
			name: "internal error → 500 (via fault injection)",
			setup: func(h *harness) {
				h.repo.FailNextGetByID = domain.ErrInternalError
			},
			method:     http.MethodGet,
			target:     "/users/any",
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness(t)
			tt.setup(h)

			rec := h.do(t, tt.method, tt.target, tt.body)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tt.wantStatus, rec.Body)
			}
			var envelope struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if envelope.Error != tt.wantMsg {
				t.Errorf("error = %q, want %q", envelope.Error, tt.wantMsg)
			}
		})
	}
}

// TestRoutesAreRegistered verifies each route is wired to a handler by
// hitting it and asserting the response is NOT the ServeMux "no route"
// default (which is the literal body "404 page not found\n").
// Handler-produced 404s (resource-not-found) return a JSON body and pass.
func TestRoutesAreRegistered(t *testing.T) {
	h := newHarness(t)

	routes := []struct{ method, path string }{
		{http.MethodGet, "/users"},
		{http.MethodPost, "/users"},
		{http.MethodGet, "/users/123"},
		{http.MethodPut, "/users/123"},
		{http.MethodDelete, "/users/123"},
	}

	for _, r := range routes {
		t.Run(r.method+"_"+r.path, func(t *testing.T) {
			body := ""
			if r.method == http.MethodPost || r.method == http.MethodPut {
				body = "{}"
			}
			rec := h.do(t, r.method, r.path, body)

			if rec.Code == http.StatusNotFound && strings.TrimSpace(rec.Body.String()) == "404 page not found" {
				t.Errorf("route %s %s not registered", r.method, r.path)
			}
		})
	}
}
