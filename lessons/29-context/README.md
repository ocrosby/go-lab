# Context

`context.Context` is Go's mechanism for propagating **cancellation**, **deadlines**, and **request-scoped values** across API boundaries. Every function that does I/O or blocks should take one as its first parameter.

> **Recommended before lesson 11 (HTTP) and lesson 14 (production API).** Both use `context.Context` on every method without ever introducing it.

## Why it matters

If a request-handler function calls a database query which calls an HTTP client which calls another service, and the original request times out, you want the whole chain to stop — not for each layer to finish its work and be discarded. `context.Context` is the standard-library primitive that makes this work. Every real Go server, every stdlib I/O operation added since Go 1.7, and every popular framework takes a `context.Context` for exactly this reason.

## Prerequisites

- Lesson 07: goroutines and channels — `ctx.Done()` returns a channel.
- Lesson 24: `defer` — every `WithCancel`/`WithTimeout` returns a cancel func you should `defer cancel()`.
- Lesson 27: functions and closures (goroutines started with `go func()` capture ctx).

## Run it

```bash
go test -race ./lessons/29-context
```

Expected: 6 passes including two that assert on cancellation behaviour.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`context_demo.go`](./context_demo.go) | Cancellable work with `context.WithCancel`, timeouts with `context.WithTimeout`, and passing a request ID via `context.WithValue`. |
| [`context_demo_test.go`](./context_demo_test.go) | Tests exercising each: normal completion, cancellation, timeout, and value propagation. |

## The three things `context.Context` carries

1. **A cancellation signal** — `ctx.Done()` returns a channel that's closed when the context is cancelled.
2. **A deadline** — `ctx.Deadline()` returns the time by which the operation should be complete.
3. **Values** — `ctx.Value(key)` retrieves a request-scoped value like a trace ID or a user ID. Use sparingly.

## The four ways to get a context

```go
ctx := context.Background()                       // root context — no deadline, no cancel, no values
ctx := context.TODO()                              // like Background(); marker for "I don't have one yet"
ctx, cancel := context.WithCancel(parent)          // child that can be cancelled manually
ctx, cancel := context.WithTimeout(parent, dur)    // child that cancels after dur elapses
ctx, cancel := context.WithDeadline(parent, time)  // child that cancels at a specific time
ctx := context.WithValue(parent, key, value)       // child that carries a value
```

Every `WithCancel`/`WithTimeout`/`WithDeadline` returns a `cancel` function. **You must call it** (via `defer cancel()`) to release the context's internal state. Forgetting is a memory leak — the Go vet tool will warn.

## The idiomatic pattern

Every function that does I/O or blocks takes a `context.Context` as its **first parameter**, conventionally named `ctx`:

```go
func FetchUser(ctx context.Context, id string) (*User, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", "/users/"+id, nil)
    // ...
}
```

The function propagates `ctx` down to whatever it calls (`http.NewRequestWithContext(ctx, ...)`, `db.QueryContext(ctx, ...)`, etc.). If the caller cancels the parent context, cancellation reaches every function in the chain.

Inside long-running loops or goroutines, check `ctx.Done()` in a `select`:

```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case work := <-in:
        process(work)
    }
}
```

## `ctx.Err()` — why did the context cancel?

After `<-ctx.Done()` fires, `ctx.Err()` tells you which flavor of cancellation happened:

- `context.Canceled` — someone explicitly called `cancel()`.
- `context.DeadlineExceeded` — a `WithTimeout` or `WithDeadline` fired.
- `nil` — the context is still live (this shouldn't happen if you just read from `Done()`).

Use `errors.Is(err, context.DeadlineExceeded)` to distinguish.

## `context.WithValue` — use sparingly

You can attach a value to a context and retrieve it downstream:

```go
type requestIDKey struct{}

ctx = context.WithValue(ctx, requestIDKey{}, "req-123")

// later:
id, _ := ctx.Value(requestIDKey{}).(string)
```

Two conventions:

- **Use an unexported key type** (like `requestIDKey struct{}` here) so no other package can accidentally collide with your key. Just `"request-id"` as a plain string is a footgun.
- **Only for request-scoped values that cross API boundaries** — request IDs, trace IDs, authenticated user IDs. Never for optional function parameters or configuration; those belong in explicit arguments.

## Try it yourself

1. Change the timeout in a test from 50ms to 500ms. Which tests still pass, which start failing?
2. Wrap a call to `time.Sleep` in a `select` with `ctx.Done()` so it can be cancelled early. This is exactly the pattern in `context_demo.go`'s `DoWork` function.
3. Try storing an `int` value in the context using a plain string key. What could go wrong if another package used the same key?
4. Chain `WithCancel` under `WithTimeout` (`ctx1, _ := WithTimeout(bg, 1s)`, then `ctx2, cancel2 := WithCancel(ctx1)`). Which context needs to fire for the child to cancel?

## Common pitfalls

- **Not calling `cancel()`.** Every `WithCancel`/`WithTimeout`/`WithDeadline` returns a cancel func. If you don't call it (usually via `defer cancel()`), the context and its internal goroutine leak until the parent cancels. `go vet` warns about this.
- **Blocking work in a select that doesn't check `ctx.Done()`.** A hot loop that doesn't see cancellation is unstoppable.
- **Storing a context in a struct field.** `context.Context` should flow through function arguments, not be attached to long-lived objects. If you need per-instance cancellation, store the `cancel` func, not the context.
- **`context.TODO()` shipped to production.** It's a marker for "I didn't know what context to use here" — fix before merging.
- **Using string keys with `WithValue`.** Two packages can collide silently. Always use an unexported struct type as the key.
- **Overusing `WithValue`.** If it's a normal function parameter, pass it as a parameter. Reserve context values for cross-boundary request-scoped data.

## You've understood this lesson when...

- You can name the three things a `context.Context` carries.
- You know why you must `defer cancel()` after `WithTimeout`.
- You can write a goroutine that respects context cancellation.
- You know why the convention is `ctx context.Context` as the first parameter of every I/O function.
- You can spot the "string key with `WithValue`" anti-pattern.

## Next

- **Next lesson:** [30-json-and-struct-tags](../30-json-and-struct-tags/) — how struct tags direct `encoding/json`, and the safe body-decoding pattern lesson 16 already uses.
