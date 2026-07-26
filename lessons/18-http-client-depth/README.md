# HTTP client depth

Production patterns for `http.Client` and `http.Transport`: connection pool tuning, per-request timeouts via context, retries via a wrapping `RoundTripper`, safe redirect policy, and the `RoundTripper`-based testing seam that beats hand-rolled `IHttpClient` interfaces.

## Why it matters

The default `http.Client{}` and `http.Get` are fine for scripts. For a service that talks to APIs in production, every one of the knobs shown here matters — timeouts you don't set become open-ended waits, redirects you don't guard against become an SSRF vector, and retries you write inline become the retry logic you copy-paste badly into every call site.

This lesson replaces the client-side content of lessons 11 and 14 with production-ready shapes.

## Prerequisites

- Lesson 11: basic `http.Client` and `http.Get`.
- Lesson 07: goroutines and channels (for understanding what `context.Context` propagates through).
- Lesson 12: dependency injection (the `RoundTripper` pattern is DI applied to the HTTP client).

## Run it

```bash
go test -race ./lessons/18-http-client-depth
```

Expected: 6 tests pass, including the retry-count and context-cancellation cases.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`client.go`](./client.go) | `NewClient` with a tuned `Transport` and a safe `CheckRedirect`. `Get` shows the per-request-timeout-via-context pattern. |
| [`retry.go`](./retry.go) | `RetryTransport` — a wrapping `RoundTripper` that adds retries on transient failures for idempotent methods only. |
| [`client_test.go`](./client_test.go) | Two testing seams: `httptest.NewServer` for full-stack tests and a stub `RoundTripper` for pure-unit tests. |

## The seven client knobs that matter

**1. Share the `Transport`, not the `Client`.** The default `http.DefaultTransport` is shared globally, which is why `http.Get` works out of the box. If you construct your own client, hold onto the `Transport` and share *it* across every client that talks to the same set of hosts. A fresh `Transport` per call defeats connection pooling.

**2. `MaxIdleConnsPerHost`.** Default is 2. If you're talking heavily to one API, bump this to 10–100. Otherwise your service opens a new TCP connection almost every request.

**3. Per-connection dial and TLS timeouts.** `Transport.DialContext` and `TLSHandshakeTimeout` bound the connection-establishment work. Without them a slow DNS lookup or a hanging TLS handshake stalls the client indefinitely.

**4. `ResponseHeaderTimeout`.** Bounds the time from "request written" to "response headers received." Distinct from a body-read timeout — a server can hold the connection open indefinitely before sending a single response byte.

**5. Per-request `context.WithTimeout`, not `Client.Timeout`.** `Client.Timeout` covers the entire round-trip including body reads. That's fine for JSON APIs, wrong for streaming responses (which legitimately take arbitrary time). Timeout on the context of each request instead — you can then use a Client with `Timeout: 0` and let each caller pick.

**6. `CheckRedirect: return http.ErrUseLastResponse`.** The default `http.Client` follows up to 10 redirects. For API clients this is almost always wrong — a redirect from `api.example.com/users/1` to `attacker.com/users/1` (via a compromised `Location` header) becomes an SSRF-style credential leak. Return `http.ErrUseLastResponse` and let the caller decide.

**7. Wrap the `Transport` for cross-cutting concerns.** Retries, auth-token injection, tracing, metrics, request/response logging — every one of these is a `RoundTripper` that wraps another `RoundTripper`. This is middleware on the client side.

## The `RoundTripper` interface

```go
type RoundTripper interface {
    RoundTrip(req *Request) (*Response, error)
}
```

That's the whole thing. Every `http.Transport` implements it. So does anything you write. `http.Client.Transport` is a `RoundTripper` — your custom wrappers slot in transparently.

**The wrapping pattern:**

```go
client := &http.Client{
    Transport: &TracingTransport{
        Base: &AuthInjectingTransport{
            Base: &RetryTransport{
                Base: http.DefaultTransport,
            },
        },
    },
}
```

`retry.go` in this lesson shows the shape. Every layer is optional, testable in isolation, and swappable without changing the call site.

## The testing seam

Lesson 11's `jsonplaceholder/` defines a custom `IHttpClient` interface so a mock can be swapped in. That works, but it's noise — Go already has a perfect testing seam:

```go
type stubRT struct { respond func(*http.Request) (*http.Response, error) }
func (s *stubRT) RoundTrip(r *http.Request) (*http.Response, error) { return s.respond(r) }

// In test:
client := &http.Client{Transport: &stubRT{
    respond: func(*http.Request) (*http.Response, error) {
        return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
    },
}}
```

No wrapper interface, no `mockgen`, no `IHttpClient`. Just a `RoundTripper` — the standard library's own seam. `client_test.go` uses this pattern for the retry tests.

Use `httptest.NewServer` when the test needs to exercise the full protocol (TLS, redirects, chunked encoding). Use a stub `RoundTripper` when the test just needs to control what the client sees.

## What's idempotent (and therefore retriable)?

RFC 7231 defines these methods as **idempotent** — the result of two identical requests is the same as one, so a retry is safe:

- `GET`, `HEAD`, `PUT`, `DELETE`, `OPTIONS`

These are **not** idempotent — retrying can duplicate side effects:

- `POST`, `PATCH`

`RetryTransport.shouldRetry` refuses to retry non-idempotent methods. If your API needs retriable POST/PATCH, use an `Idempotency-Key` header (Stripe's pattern) and a Transport that respects it.

## Try it yourself

1. Add a `LoggingTransport` that prints method, URL, status, and duration for every request. Wrap it around `RetryTransport`. Which order gives you the most useful log line?
2. Add an `AuthTransport` that injects `Authorization: Bearer <token>`. Which order relative to `RetryTransport` matters? (Hint: what if the auth token expires mid-retry?)
3. Change `RetryTransport` to honor the `Retry-After` response header on `429` responses instead of just using exponential backoff.
4. Add a per-response `context.Cancel` so the client aborts if the response body is too large — like an in-flight `MaxBytesReader`.

## Common pitfalls

- **`Client.Timeout` for API calls.** Coarse. Prefer `context.WithTimeout` per request.
- **Not closing `resp.Body`.** Every `resp` you get from a client must have `Body.Close()` called, even on error paths. A leaked body is a leaked connection.
- **Retrying non-idempotent methods.** Duplicated POSTs create duplicate resources. `RetryTransport` will not do this and neither should you.
- **Following redirects by default.** Set `CheckRedirect` explicitly. Even if you decide to follow redirects, decide it deliberately.
- **Fresh `http.Client` per call.** Defeats pooling. Construct one per host-family, share it.
- **Mocking via a custom `IHttpClient` interface** when `http.RoundTripper` is right there. Every extra abstraction you invent is one your team has to maintain.

## You've understood this lesson when...

- You can name the seven knobs above and what each protects against.
- You can write a `RoundTripper` from scratch that adds one behaviour (auth, logging, metrics).
- You know why `Client.Timeout` and `context.WithTimeout` behave differently on a streaming response.
- You can explain why POST is not retried automatically.
- You can test an HTTP client without spinning up a server.

## Related deep-dive

- `rules/rest-api-conventions.md` (parent claude-config repo) — client-side behaviour on redirects, retries, and idempotency.

## Next

You've now completed the HTTP track. The remaining lessons — the production API in [14-production-api](../14-production-api/) — are worth revisiting with the lessons 16–18 lens: what routing style would you use today (1.22+ ServeMux), what middleware would you add (chain from 17), what would the client-side of that API look like (`NewClient` from 18)?
