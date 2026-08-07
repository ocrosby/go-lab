# Server basics

The smallest useful HTTP server in Go: `http.Server`, `http.ServeMux`, two `http.HandlerFunc` endpoints, and correct timeouts.

## Why it matters

Every REST API you write in Go starts with these four types. Once you can wire them together, everything the ecosystem builds on top (Gin, Echo, Chi) is just sugar over the same shapes.

The four timeout knobs are also introduced here, in the smallest context where they make sense — they'll appear on every server in the rest of the syllabus.

## Prerequisites

- Lesson 01: how to run a Go program.
- Lesson 04: `go test`.

## Run it

Test the handlers without spinning up a real server:

```bash
go test ./lessons/12-http-clients-and-servers/server-basics
```

Expected output:

```text
ok  	github.com/ocrosby/go-lab/lessons/12-http-clients-and-servers/server-basics	...
```

To actually start the server, drop this two-line file next to `server.go`:

```go
package main

import "github.com/ocrosby/go-lab/lessons/12-http-clients-and-servers/server-basics"
func main() { _ = serverbasics.NewServer(":8080").ListenAndServe() }
```

Then in another terminal:

```bash
curl http://localhost:8080/
curl 'http://localhost:8080/hello?name=go-lab'
```

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`server.go`](./server.go) | `NewServer` wires `ServeMux` to two handlers and sets four timeouts. |
| [`server_test.go`](./server_test.go) | Two testing patterns from `net/http/httptest` — `NewRecorder` (fastest, handler-only) and `NewServer` (end-to-end, real socket). |

## The five things every Go HTTP server has

1. **`http.Handler` interface** — anything with a `ServeHTTP(w, r)` method. That's the whole contract.
2. **`http.HandlerFunc`** — a function `func(w http.ResponseWriter, r *http.Request)` that Go treats as an `http.Handler` for free (via a one-method type adapter).
3. **`http.ServeMux`** — the request router. Match URLs to handlers. Go 1.22 added method-in-pattern support; see lesson 17.
4. **`http.Server`** — the actual listening loop, plus timeouts and connection settings.
5. **The four timeouts.**

## The four timeout knobs

Any `http.Server` deployed to production sets these — leaving them at zero (the default) is an unbounded resource commitment per connection.

| Field | What it caps | Slowloris protection? |
|---|---|---|
| `ReadHeaderTimeout` | Time to read request headers | **Yes — this is the specific knob for it** |
| `ReadTimeout` | Time to read the entire request (headers + body) | Yes, but coarser than `ReadHeaderTimeout` |
| `WriteTimeout` | Time from end-of-request-read to end-of-response-write | Not directly |
| `IdleTimeout` | Time a keep-alive connection can sit unused | Frees up idle connections |

**Set `ReadHeaderTimeout` explicitly**, even if you also set `ReadTimeout` — `gosec` will flag its absence and it's cheap insurance against slow-header attacks.

## Testing with `httptest`

`net/http/httptest` gives you two ways to test handlers:

- **`httptest.NewRecorder`** — a fake `http.ResponseWriter` you pass to your handler directly. No socket, no port, no goroutine. Fastest. Best for unit tests.
- **`httptest.NewServer`** — a real HTTP server on a random port. Real socket, real protocol. Slower. Best for tests that must exercise the full stack (redirects, TLS, cross-handler behaviour).

Both are in [`server_test.go`](./server_test.go). Reach for `NewRecorder` first; graduate to `NewServer` only when you need it.

## Try it yourself

1. Add a third endpoint `/echo` that returns whatever `?msg=…` was passed in. Add a test for it.
2. Change `NewServer` to return an error if the port is already bound. (Hint: `net.Listen("tcp", addr)` first, then wrap in a Server.)
3. Add a `/panic` handler that calls `panic("test")`. Watch the whole process crash when you `curl` it. Then look at lesson 18 — that's what middleware is for.

## Common pitfalls

- **Writing headers after calling `Write`.** Once `Write` runs, the status code is committed. Set all headers before you write anything, and prefer `w.WriteHeader(status)` before `w.Write(body)` when the status isn't 200.
- **Missing `ReadHeaderTimeout`.** Silent security regression — Slowloris still works even when `ReadTimeout` is set, because attackers keep sending header bytes just below the deadline.
- **Using `http.ListenAndServe` at package level.** Fine for a hello-world; not fine for anything you want to test. `NewServer` returning a `*http.Server` lets tests hold the value without booting a real listener.

## You've understood this lesson when...

- You can name the four `http.Server` timeout fields and what each covers.
- You can write a handler test using both `httptest.NewRecorder` and `httptest.NewServer`, and explain when to reach for each.
- You know why setting headers after `Write` doesn't work.

## Next

- **Next sub-lesson:** [`../jsonplaceholder/`](../jsonplaceholder/) — a bigger, structured client that puts a service layer on top of `net/http`.
- **REST-specific routing:** [`../../17-restful-routing/`](../../17-restful-routing/).
- **Middleware:** [`../../18-http-middleware/`](../../18-http-middleware/).
