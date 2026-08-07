// Package restfulrouting demonstrates a small REST API for a "tasks"
// resource using the Go 1.22+ ServeMux: method-in-pattern routing and
// path parameters via r.PathValue.
package restfulrouting

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// Task is the resource this API exposes.
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// store is a tiny in-memory repository. Real code would inject an interface
// (see lesson 13) and back it with a database.
type store struct {
	mu     sync.RWMutex
	nextID int
	tasks  map[int]Task
}

func newStore() *store {
	return &store{nextID: 1, tasks: make(map[int]Task)}
}

// server wires the store to the ServeMux.
type server struct {
	store *store
}

// NewMux returns a ServeMux configured with the REST endpoints. Split out so
// tests can wire the mux without booting a full http.Server.
func NewMux() *http.ServeMux {
	srv := &server{store: newStore()}
	mux := http.NewServeMux()

	// Go 1.22+ ServeMux: method and path in one pattern. `{id}` is a path
	// parameter, extracted via r.PathValue("id").
	mux.HandleFunc("GET /tasks", srv.listTasks)
	mux.HandleFunc("POST /tasks", srv.createTask)
	mux.HandleFunc("GET /tasks/{id}", srv.getTask)
	mux.HandleFunc("PUT /tasks/{id}", srv.replaceTask)
	mux.HandleFunc("DELETE /tasks/{id}", srv.deleteTask)

	return mux
}

// listTasks: GET /tasks[?done=true|false]
// Demonstrates query-parameter filtering.
func (s *server) listTasks(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("done")

	s.store.mu.RLock()
	defer s.store.mu.RUnlock()

	out := make([]Task, 0, len(s.store.tasks))
	for _, t := range s.store.tasks {
		switch filter {
		case "true":
			if !t.Done {
				continue
			}
		case "false":
			if t.Done {
				continue
			}
		}
		out = append(out, t)
	}
	writeJSON(w, http.StatusOK, out)
}

// createTask: POST /tasks
// Demonstrates safe body reading (size cap + strict decoder) and the
// 201 Created + Location header response shape from RFC 7231.
func (s *server) createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Title == "" {
		writeError(w, http.StatusUnprocessableEntity, "title is required")
		return
	}

	s.store.mu.Lock()
	id := s.store.nextID
	s.store.nextID++
	task := Task{ID: id, Title: input.Title}
	s.store.tasks[id] = task
	s.store.mu.Unlock()

	w.Header().Set("Location", fmt.Sprintf("/tasks/%d", id))
	writeJSON(w, http.StatusCreated, task)
}

// getTask: GET /tasks/{id}
// Demonstrates r.PathValue for path parameters, and the 404 shape.
func (s *server) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be an integer")
		return
	}
	s.store.mu.RLock()
	task, ok := s.store.tasks[id]
	s.store.mu.RUnlock()
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

// replaceTask: PUT /tasks/{id}
// PUT is idempotent — full replacement. Demonstrates 204 No Content when the
// server has no body to return, and the requirement that 204 have no body.
func (s *server) replaceTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be an integer")
		return
	}
	var input Task
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	input.ID = id

	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, ok := s.store.tasks[id]; !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	s.store.tasks[id] = input
	w.WriteHeader(http.StatusNoContent)
}

// deleteTask: DELETE /tasks/{id}
// DELETE is idempotent. 204 No Content whether the resource existed or not —
// the observable end state is the same.
func (s *server) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be an integer")
		return
	}
	s.store.mu.Lock()
	delete(s.store.tasks, id)
	s.store.mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

// decodeJSON is the safe body-reader every REST endpoint should use:
// - MaxBytesReader caps request body to prevent OOM
// - DisallowUnknownFields rejects payloads with extra fields (contract clarity)
// - errors on more-than-one-JSON-value in the stream
func decodeJSON(r *http.Request, dst any) error {
	const maxBody = 1 << 20 // 1 MiB
	r.Body = http.MaxBytesReader(nil, r.Body, maxBody)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if dec.More() {
		return errors.New("body must contain exactly one JSON value")
	}
	return nil
}

// writeJSON marshals v as JSON and writes it with the given status.
// Sets Content-Type before writing the status (required — headers can't
// be modified after WriteHeader).
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError is the API's uniform error envelope. Real code often adds a
// stable error code and a request ID for correlation.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
