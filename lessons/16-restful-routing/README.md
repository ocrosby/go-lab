# RESTful routing

Full REST-shaped routing for a `tasks` resource using the Go 1.22+ ServeMux — method-in-pattern (`GET /tasks/{id}`), path parameters via `r.PathValue`, safe JSON body reading, static-file serving, and Server-Sent Events. Everything a real REST API needs from `net/http`, in one lesson.

## Why it matters

Go 1.22 rewrote `http.ServeMux` to understand HTTP methods and path parameters. Before 1.22, REST-shaped routing meant either a giant `switch r.Method` inside a prefix handler (the pattern lesson 14's production API still uses) or reaching for a third-party router (chi, gorilla/mux, gin). Since 1.22, the standard library is enough for most APIs.

If you're writing a REST API in Go **today**, this is the routing style you should be using.

## Prerequisites

- Lesson 11: `http.Server`, `http.ServeMux`, `http.HandlerFunc`.
- Lesson 06: interfaces (the store is a value here, but any real API replaces it with an interface).
- Go 1.22 or newer (this repo pins Go 1.26).

## Run it

```bash
go test -race ./lessons/16-restful-routing
```

Expected: 8 tests pass, including the SSE streaming test.

Optional: wire it up as a real server and hit it with `curl`.

```go
package main
import (
    "net/http"
    "github.com/ocrosby/go-lab/lessons/16-restful-routing"
)
func main() { _ = http.ListenAndServe(":8080", restfulrouting.NewMux()) }
```

```bash
curl -i -X POST -d '{"title":"read a book"}' http://localhost:8080/tasks
curl -i http://localhost:8080/tasks
curl -i http://localhost:8080/tasks/1
curl -i -X PUT -d '{"title":"read a book","done":true}' http://localhost:8080/tasks/1
curl -i -X DELETE http://localhost:8080/tasks/1
```

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`tasks.go`](./tasks.go) | REST endpoints for `/tasks` and `/tasks/{id}`, safe JSON reading, uniform error envelope. |
| [`tasks_test.go`](./tasks_test.go) | Handler tests using `httptest.NewRecorder`. |
| [`static.go`](./static.go) | Static-file serving via `http.FileServer` + `http.StripPrefix`, and single-file serving via `http.ServeContent` (which handles `Range` requests for free). |
| [`sse.go`](./sse.go) | Server-Sent Events streaming with `http.Flusher`. |

## The 1.22+ ServeMux, in one table

| Pattern | Matches |
|---|---|
| `mux.HandleFunc("/tasks", h)` | Any method to `/tasks` |
| `mux.HandleFunc("GET /tasks", h)` | Only `GET /tasks` |
| `mux.HandleFunc("GET /tasks/{id}", h)` | Only `GET /tasks/{anything}` — capture with `r.PathValue("id")` |
| `mux.HandleFunc("GET /tasks/{id}/attachments/{name}", h)` | Multiple path params work |
| `mux.HandleFunc("GET /tasks/", h)` | Prefix (trailing slash) |
| `mux.HandleFunc("GET example.com/foo", h)` | Host-scoped |
| `mux.HandleFunc("GET /files/{path...}", h)` | Wildcard trailing segments |

Method mismatch returns **`405 Method Not Allowed`** with an `Allow` header automatically — you don't write that logic.

## The status code contract (the abbreviated version)

REST assigns semantic meaning to status codes. `tasks.go` uses these deliberately:

| Situation | Status | Extra headers |
|---|---|---|
| GET succeeded | `200 OK` | `Content-Type` |
| POST created a resource | `201 Created` | `Location: /path/to/{id}` |
| PUT/DELETE succeeded, no body | `204 No Content` | (must not have a body) |
| Body didn't parse | `400 Bad Request` | |
| Body parsed but was semantically invalid | `422 Unprocessable Entity` | |
| Resource does not exist | `404 Not Found` | |
| Method not allowed on this route | `405 Method Not Allowed` | `Allow: GET, POST, ...` |
| Auth required or wrong | `401 Unauthorized` | `WWW-Authenticate` |
| Auth OK, action forbidden | `403 Forbidden` | |
| Concurrent-modification conflict | `409 Conflict` | |
| Rate limited | `429 Too Many Requests` | `Retry-After` |

The full table lives in `rules/rest-api-conventions.md` in the parent claude-config repo.

## Safe body reading

`json.NewDecoder(r.Body).Decode(&x)` is dangerous by default. `decodeJSON` in `tasks.go` shows the three guardrails every REST endpoint should have:

1. **`http.MaxBytesReader`** — cap the body. An attacker sending a 10 GB JSON payload otherwise burns your memory.
2. **`dec.DisallowUnknownFields()`** — reject payloads with extra fields. Catches typos (`titel` instead of `title`) that would otherwise silently drop data.
3. **Reject stream continuations** — a request body should contain exactly one JSON value. `dec.More()` catches the concatenated-JSON edge case.

## Static files, done right

```go
// Serve everything under ./public/ at /assets/*
mux.Handle("GET /assets/", FileServerHandler("./public", "/assets"))
```

Three primitives compose here — `http.Dir` (path → filesystem), `http.FileServer` (filesystem → handler), `http.StripPrefix` (remove URL prefix before lookup). See [`static.go`](./static.go).

For a single resource where you already have the bytes, use **`http.ServeContent`** instead. It handles `If-Modified-Since`, `ETag`/`If-None-Match`, and `Range` requests automatically — so a client asking for bytes 100–200 gets a proper `206 Partial Content` response with no extra code from you.

## Server-Sent Events

The `/events` endpoint in [`sse.go`](./sse.go) streams events one at a time using `http.Flusher` — the interface an `http.ResponseWriter` satisfies when the underlying protocol supports flushing (HTTP/1.1 chunked, HTTP/2). Without `Flush()`, buffered writes sit in the write buffer until the handler returns.

Every SSE handler should:
1. Set `Content-Type: text/event-stream`.
2. Assert `w.(http.Flusher)` and 500 if the assertion fails.
3. Select on `r.Context().Done()` in the send loop — the client can disconnect at any time.

## A note on WebSockets

`net/http` exposes `http.Hijacker` — the interface that lets a handler take over the raw TCP connection, which is what WebSockets need. In practice every real WebSocket implementation uses the [`gorilla/websocket`](https://github.com/gorilla/websocket) or [`nhooyr.io/websocket`](https://nhooyr.io/websocket) library on top of Hijacker rather than rolling one from scratch. If you need bidirectional streaming, reach for those libraries.

## Try it yourself

1. Add `PATCH /tasks/{id}` for partial updates (`{ "done": true }` alone should work). What status code should it return?
2. Change the response of `POST /tasks` to include a `Retry-After: 1` header when the store hits a size cap (say, 100 tasks) and return `429`. Add a test.
3. Add `HEAD /tasks/{id}` — the ServeMux does not synthesize this from `GET`. What should the response look like?
4. Change `static.go`'s `ServeReadmeContent` to serve from an `embed.FS` (Go 1.16+, see `docs/go-build-directives.md`) so the content is baked into the binary.

## Common pitfalls

- **Setting headers after `w.Write`.** Once the first byte is written, headers are committed. Do all header work before any body write.
- **Returning a body from a `204 No Content`.** Some clients silently strip it; some raise an error; some render garbage. `204` means empty — write nothing.
- **Forgetting `r.PathValue` returns `""` for missing segments** rather than an error. Always validate the value.
- **Method-specific patterns don't imply an OPTIONS handler.** If you need CORS preflight support, add explicit `OPTIONS /tasks` routes or use the middleware from lesson 17.
- **Overlapping patterns.** The 1.22 ServeMux is stricter about conflicts than the old one — two patterns that match the same request cause a startup-time panic. Read the error carefully; it usually names both patterns.

## You've understood this lesson when...

- You can write a full REST endpoint set for a new resource using method-in-pattern routing without looking anything up.
- You can name three things `decodeJSON` protects against that a bare `json.Decode` does not.
- You can explain when to use `http.ServeContent` vs `http.FileServer` vs `http.ServeFile`.
- You can sketch an SSE handler and know why `Flush()` matters.

## Related deep-dive

- `rules/rest-api-conventions.md` (parent claude-config repo) — the full REST status-code + method + URL table this lesson abbreviates.

## Next

- **Next lesson:** [17-http-middleware](../17-http-middleware/) — the middleware chain (logging, auth, CORS, body-limit, panic-recovery) every real server layers around handlers like these.
