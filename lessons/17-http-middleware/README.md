# HTTP middleware

The `func(http.Handler) http.Handler` pattern that every serious Go server uses to layer request-scoped concerns — logging, auth, panic recovery, CORS, body limits — around handlers without touching them.

## Why it matters

Middleware is the single most common pattern in `net/http` production code. Once you can read the shape, every Go HTTP codebase in the world becomes intelligible — from stdlib-only services to Chi and Echo and Gin. It's also how you keep handlers focused on their business logic instead of drowning in cross-cutting concerns.

## Prerequisites

- Lesson 11: `http.Handler`, `http.HandlerFunc`.
- Lesson 16: full REST handlers to wrap.
- Lesson 10: panic recovery basics (this lesson generalizes what you saw there).

## Run it

```bash
go test -race ./lessons/17-http-middleware
```

Expected: 10 tests pass — one per middleware plus a chain-order test.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`middleware.go`](./middleware.go) | Six real middlewares + `Chain` composer. |
| [`middleware_test.go`](./middleware_test.go) | One focused test per middleware, plus a chain-order verification. |

## The shape

```go
type Middleware func(http.Handler) http.Handler
```

A middleware takes a handler and returns a new handler. The returned handler runs code, delegates to the wrapped handler, then runs more code. That's it — the entire pattern.

Composition:

```go
final := Chain(base,
    RequestID,       // outermost — first in, last out
    Logging(logger),
    Recover,
    CORS("*"),
    BodyLimit(1<<20),
    Auth(token),     // innermost — closest to base handler
)
```

Order matters. `Chain` applies the middlewares so that the first-listed becomes the outermost wrapper — it sees the request first and the response last. Every middleware sits in a "diamond" around the next one:

```text
RequestID in → Logging in → Recover in → ... → base → ... → Recover out → Logging out → RequestID out
```

## The six middlewares in this lesson

| Middleware | What it does | Order note |
|---|---|---|
| `RequestID` | Adds a random ID (or trusts an upstream `X-Request-ID`), echoes it in the response, and stashes it in the context so downstream code and logs can reference it. | Outermost — everything else should be able to log the ID. |
| `Logging` | Structured line per request: method, path, status, size, duration, request ID. Wraps `ResponseWriter` to capture the status code. | After `RequestID` so log lines carry the ID. |
| `Recover` | Catches panics, converts to `500 Internal Server Error` with a JSON envelope, logs the stack. Rethrows `http.ErrAbortHandler` — that's the standard library's own signal, not a bug. | After `Logging` so the log line records the 500. |
| `CORS` | Minimal browser CORS: allowed origin, methods, headers; answers preflight `OPTIONS` with `204`. Use `rs/cors` for anything beyond a demo. | Before `Auth` so preflights aren't rejected. |
| `BodyLimit` | Caps request body via `http.MaxBytesReader`. Unbounded bodies are an OOM vector. | Before `Auth` is fine; before handler is required. |
| `Auth(token)` | Requires `Authorization: Bearer <token>`; responds `401` + `WWW-Authenticate: Bearer` if missing/wrong. Constant-time compare to avoid a timing side-channel. | Innermost — protects the handler. |

## The three tricks worth naming

**Wrapping `http.ResponseWriter`.** To capture the status code a handler wrote, you replace the `ResponseWriter` before calling `next.ServeHTTP` with your own type that embeds the original and overrides `WriteHeader`. See `statusWriter` in `middleware.go`. This is the pattern any Go logger, metrics recorder, or compression middleware uses.

**Context values for cross-cutting data.** `RequestID` stashes the ID in `r.Context()` using an unexported `ctxKey` type. Handlers pull it out via `RequestIDFrom(ctx)`. The unexported key type prevents accidental key collisions across packages — a well-known Go idiom.

**Middleware factories.** Middlewares that need configuration (`Logging(logger)`, `Auth(token)`, `BodyLimit(1<<20)`) return the middleware from a function. This keeps the middleware itself parameterless and the wiring in one place.

## Try it yourself

1. Write a `RateLimit(n int, per time.Duration)` middleware that returns `429 Too Many Requests` with a `Retry-After` header when the caller exceeds `n` requests in `per` time. (Hint: `time.Ticker` or `golang.org/x/time/rate`.)
2. Write a `Compress` middleware that gzip-encodes responses when the client sends `Accept-Encoding: gzip`. Beware: you need to wrap `Write` and set `Content-Encoding` before the handler writes the body.
3. Reorder the chain in `Handler()` so `Auth` runs *before* `BodyLimit`. Which requests break? Which are cheaper?
4. Replace the ad-hoc log line in `Logging` with `log/slog` (Go 1.21+) for structured JSON output.

## Common pitfalls

- **Middleware that swallows panics without logging** — worse than no middleware because it hides bugs. Always log the recovered value AND the stack.
- **Forgetting to call `next.ServeHTTP`.** A middleware that returns without delegating short-circuits the pipeline — sometimes what you want (auth failure), often not (accidentally).
- **Mutating `r` directly** instead of `r.WithContext(newCtx)`. `*http.Request` is intended to be treated as immutable per-hop; use `WithContext` to attach values.
- **Applying `Recover` innermost.** The point of `Recover` is to catch panics in the wrapped stack. Put it near the outside so it wraps as much as possible.
- **Chain-order confusion.** Draw the diamond. Every middleware's "in" code runs top-to-bottom; every "out" code runs bottom-to-top. If you're not sure, add a test like `TestChain_OrderIsFIFO` in this lesson.

## You've understood this lesson when...

- You can sketch the `Middleware` type from memory and write a "hello, world" middleware without looking.
- You can explain why `Logging` needs a wrapped `ResponseWriter`.
- You can predict the order of "in" and "out" log lines for a given chain.
- You know why `Recover` should rethrow `http.ErrAbortHandler`.

## Next

- **Next lesson:** [18-http-client-depth](../18-http-client-depth/) — the client-side equivalents of everything you just learned: `http.RoundTripper` is the client's middleware seam, custom `Transport` is the connection-pool tuning knob, `context.WithTimeout` is the per-request cancellation primitive.
