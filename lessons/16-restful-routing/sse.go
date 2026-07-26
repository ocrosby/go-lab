package restfulrouting

import (
	"fmt"
	"net/http"
	"time"
)

// SSEHandler streams ticks to the client using Server-Sent Events, showcasing
// http.Flusher — the interface an http.ResponseWriter implements when the
// underlying protocol supports flushing (HTTP/1.1 chunked, HTTP/2). Without
// Flush, buffered writes would sit in the write buffer until the handler
// returned, defeating the point of streaming.
//
// The `text/event-stream` content type is the SSE contract. Each event is
// one or more "field: value" lines terminated by a blank line.
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Bound the number of events so this doesn't stream forever in a
	// misconfigured deployment. Real code would cap by time or client
	// disconnect only.
	const maxTicks = 5

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i < maxTicks; i++ {
		select {
		case <-r.Context().Done():
			// Client disconnected — stop sending. Every SSE handler must
			// respect r.Context().Done().
			return
		case t := <-ticker.C:
			fmt.Fprintf(w, "id: %d\n", i+1)
			fmt.Fprintf(w, "event: tick\n")
			fmt.Fprintf(w, "data: %s\n\n", t.UTC().Format(time.RFC3339))
			flusher.Flush()
		}
	}
}
