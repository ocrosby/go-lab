package restfulrouting

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// helper: create + GET a task by id.
func createTask(t *testing.T, mux *http.ServeMux, title string) Task {
	t.Helper()
	body := bytes.NewBufferString(`{"title":"` + title + `"}`)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/tasks", body)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /tasks status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	if loc := rec.Header().Get("Location"); !strings.HasPrefix(loc, "/tasks/") {
		t.Errorf("Location = %q, want /tasks/*", loc)
	}
	var got Task
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return got
}

func TestCreateAndGet(t *testing.T) {
	mux := NewMux()
	created := createTask(t, mux, "write tests")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/tasks/"+strconv.Itoa(created.ID), nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want 200", rec.Code)
	}
	var got Task
	_ = json.Unmarshal(rec.Body.Bytes(), &got)
	if got.Title != "write tests" {
		t.Errorf("Title = %q, want %q", got.Title, "write tests")
	}
}

func TestGetUnknownIsNotFound(t *testing.T) {
	mux := NewMux()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/tasks/999", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestCreateRejectsEmptyTitle(t *testing.T) {
	mux := NewMux()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/tasks", bytes.NewBufferString(`{"title":""}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", rec.Code)
	}
}

func TestCreateRejectsUnknownFields(t *testing.T) {
	mux := NewMux()
	// DisallowUnknownFields catches typos — "titel" instead of "title".
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/tasks", bytes.NewBufferString(`{"titel":"typo"}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestReplace(t *testing.T) {
	mux := NewMux()
	created := createTask(t, mux, "original")

	body := bytes.NewBufferString(`{"title":"revised","done":true}`)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/tasks/"+strconv.Itoa(created.ID), body)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("PUT status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("204 must not have a body, got %d bytes", rec.Body.Len())
	}
}

func TestDeleteIsIdempotent(t *testing.T) {
	mux := NewMux()
	created := createTask(t, mux, "delete me")

	for i := 0; i < 2; i++ {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/tasks/"+strconv.Itoa(created.ID), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("DELETE #%d status = %d, want 204", i+1, rec.Code)
		}
	}
}

func TestListFilter(t *testing.T) {
	mux := NewMux()
	a := createTask(t, mux, "a")
	_ = createTask(t, mux, "b")

	// Mark a as done via PUT
	body := bytes.NewBufferString(`{"title":"a","done":true}`)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/tasks/"+strconv.Itoa(a.ID), body)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// GET /tasks?done=true → 1 result
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/tasks?done=true", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var doneTasks []Task
	_ = json.Unmarshal(rec.Body.Bytes(), &doneTasks)
	if len(doneTasks) != 1 {
		t.Errorf("done=true count = %d, want 1", len(doneTasks))
	}
}

// SSE test: verify the endpoint streams at least one event before we
// cancel the request context to shut it down.
func TestSSE(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(SSEHandler))
	defer srv.Close()

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, nil)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if got := resp.Header.Get("Content-Type"); got != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", got)
	}
	buf := make([]byte, 256)
	n, _ := resp.Body.Read(buf)
	if !strings.Contains(string(buf[:n]), "event: tick") {
		t.Errorf("first chunk missing 'event: tick': %q", string(buf[:n]))
	}
}
