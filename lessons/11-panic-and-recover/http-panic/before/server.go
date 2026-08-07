package before

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	mux *http.ServeMux
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/panic", s.handlePanic)
	s.mux.HandleFunc("/slow", s.handleSlow)
	s.mux.HandleFunc("/json", s.handleJSON)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Server is healthy",
		Status:  "ok",
	})
}

func (s *Server) handlePanic(w http.ResponseWriter, r *http.Request) {
	panic("intentional panic in handler")
}

func (s *Server) handleSlow(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		panic("panic in background goroutine")
	}()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Started background job",
		Status:  "ok",
	})
}

func (s *Server) handleJSON(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if trigger, ok := data["trigger_panic"].(bool); ok && trigger {
		panic(fmt.Sprintf("panic triggered by request: %v", data))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Data processed successfully",
		Status:  "ok",
	})
}

func (s *Server) Start(addr string) error {
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
