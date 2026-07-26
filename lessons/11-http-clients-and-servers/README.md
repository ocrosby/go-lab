# HTTP clients and servers

A first look at `net/http` on both sides of the wire: making requests to a remote API, and answering requests from your own server. Two sub-lessons plus a bigger worked example.

## Why it matters

`net/http` is Go's flagship standard-library package. A working REST API touches maybe a dozen of its types and interfaces — this lesson introduces the smallest useful ones, and lessons 16–18 go into depth on routing, middleware, and production client patterns.

## Prerequisites

- Lesson 06: interfaces (mocking an HTTP client requires them).
- Lesson 03: `go test` basics.

## Run it

Each sub-lesson runs independently:

```bash
go test ./lessons/11-http-clients-and-servers/...
```

Expected: three sub-packages pass.

## What's in this folder

| Path | What it demonstrates |
|---|---|
| [`client-basics/`](./client-basics/) | The smallest useful client — one function, one endpoint, one JSON decode. First look at `http.Get` and `net/http/httptest`. |
| [`server-basics/`](./server-basics/) | The smallest useful server — `http.Server`, `ServeMux`, `HandlerFunc`, timeouts. First look at `httptest.NewRecorder` and `httptest.NewServer` for tests. |
| [`jsonplaceholder/`](./jsonplaceholder/) | A bigger, structured HTTP client that hits a real public API, with a service layer and a mockable transport. Ties together lessons 06 and 11. |

## Where to go next

The three add-on lessons that turn this foundation into a REST-ready toolkit:

- **[16-restful-routing](../16-restful-routing/)** — REST verbs and path parameters using the Go 1.22+ ServeMux.
- **[17-http-middleware](../17-http-middleware/)** — the middleware chain every server needs.
- **[18-http-client-depth](../18-http-client-depth/)** — production HTTP client: custom transport, timeouts, retries, testing.

## Common pitfalls

- **Forgetting `resp.Body.Close()`** — leaks a connection every call. Always `defer resp.Body.Close()` right after checking the error.
- **Using the default `http.Client`** for production requests — no timeout, follows up to 10 redirects. Fine for scripts; not fine for services. See lesson 18.
- **Confusing "server" with "endpoint I hit"** — the code in `client-basics/` sends requests, so it's a client, not a server. This lesson previously had this mislabelled and it caused real confusion.

## Related deep-dive

- [`docs/csp-and-go-concurrency.md`](../../docs/csp-and-go-concurrency.md) — the HTTP handlers each run in their own goroutine; understanding Go's concurrency model helps you reason about handler code.
