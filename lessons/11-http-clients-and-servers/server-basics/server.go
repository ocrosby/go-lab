// Package serverbasics is the smallest useful HTTP server in Go: an
// http.Server bound to a ServeMux with two endpoints. Everything you need
// to answer HTTP requests, and no framework.
package serverbasics

import (
	"encoding/json"
	"net/http"
	"time"
)

// Greeting is the JSON shape returned by /hello.
type Greeting struct {
	Message string `json:"message"`
	At      string `json:"at"`
}

// NewServer wires the ServeMux and returns a fully-configured *http.Server.
// The caller decides when to call ListenAndServe or Shutdown.
//
// Timeouts matter — a server without them is a slowloris target. See the
// lesson README for the specific role each timeout plays.
func NewServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/hello", handleHello)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

// handleRoot answers the health check on `/`.
func handleRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("go-lab HTTP server — try GET /hello?name=you\n"))
}

// handleHello returns a JSON greeting. Demonstrates reading a query parameter,
// setting a response header, and encoding a struct as JSON.
func handleHello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}

	greeting := Greeting{
		Message: "hello, " + name,
		At:      time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(greeting)
}
