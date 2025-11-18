package after

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type PanicHandler func(recovered interface{}, r *http.Request)

type Server struct {
	mux          *http.ServeMux
	panicHandler PanicHandler
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}

func NewServer(panicHandler PanicHandler) *Server {
	if panicHandler == nil {
		panicHandler = defaultPanicHandler
	}

	s := &Server{
		mux:          http.NewServeMux(),
		panicHandler: panicHandler,
	}
	s.setupRoutes()
	return s
}

func defaultPanicHandler(recovered interface{}, r *http.Request) {
	log.Printf("PANIC RECOVERED in %s %s: %v\n%s", r.Method, r.URL.Path, recovered, debug.Stack())
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", s.withPanicRecovery(s.handleHealth))
	s.mux.HandleFunc("/panic", s.withPanicRecovery(s.handlePanic))
	s.mux.HandleFunc("/slow", s.withPanicRecovery(s.handleSlow))
	s.mux.HandleFunc("/json", s.withPanicRecovery(s.handleJSON))
	s.mux.HandleFunc("/abort", s.withPanicRecovery(s.handleAbort))
}

func (s *Server) withPanicRecovery(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				panic(rec)
			}

			s.panicHandler(rec, r)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error:  "Internal Server Error",
				Status: "error",
			})
		}()

		handler(w, r)
	}
}

func safeGo(fn func(), panicHandler func(interface{})) {
	go func() {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				panic(rec)
			}

			if panicHandler != nil {
				panicHandler(rec)
				return
			}

			log.Printf("PANIC RECOVERED in goroutine: %v\n%s", rec, debug.Stack())
		}()
		fn()
	}()
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
	safeGo(func() {
		time.Sleep(100 * time.Millisecond)
		panic("panic in background goroutine")
	}, func(rec interface{}) {
		log.Printf("Background goroutine panic: %v", rec)
	})

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

func (s *Server) handleAbort(w http.ResponseWriter, r *http.Request) {
	panic(http.ErrAbortHandler)
}

func (s *Server) Start(addr string) error {
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
