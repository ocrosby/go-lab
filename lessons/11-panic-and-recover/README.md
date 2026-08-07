# Panic and recover

What `panic` and `recover` are, when to reach for them, and — more importantly — when *not* to.

## Why it matters

Every other lesson so far has used `error` return values: expected failure paths, checked at each call site. `panic` is Go's escape hatch for the *unexpected* — the runtime signals it when something is genuinely broken (nil pointer dereference, out-of-bounds index, `close` on a closed channel), and your code can raise it deliberately with `panic("...")`.

`recover` is the only way to catch a panic and keep running. There are two places you *must* know how to use it: **long-running goroutines** and **HTTP handlers**. Everywhere else, prefer `error`.

## Prerequisites

- Lesson 08: goroutines and channels.
- Lesson 12: HTTP handlers (helpful, not required — this lesson also has HTTP examples).

## Run it

This lesson has two sub-lessons, each with a `before/` (broken) version and an `after/` (fixed) version. **The `before/` packages are supposed to fail** — they demonstrate uncontrolled panics.

Run each sub-lesson's tests:

```bash
# The "after" packages — these should pass
go test ./lessons/11-panic-and-recover/goroutine-panic/after/
go test ./lessons/11-panic-and-recover/http-panic/after/

# The "before" packages — these WILL panic; they're the demonstration
go test ./lessons/11-panic-and-recover/goroutine-panic/before/
go test ./lessons/11-panic-and-recover/http-panic/before/
```

For the whole-repo test command, use `make test` — it skips the `before/` packages automatically.

## What's in this folder

| Path | What it demonstrates |
|---|---|
| [`goroutine-panic/before/`](./goroutine-panic/before/) | A worker pool where a panicking job silently kills its worker. |
| [`goroutine-panic/after/`](./goroutine-panic/after/) | The same pool with `defer recover()` in each worker. |
| [`http-panic/before/`](./http-panic/before/) | An HTTP server where a handler panic crashes the whole server. |
| [`http-panic/after/`](./http-panic/after/) | The same server with panic-recovery middleware. |

Each sub-lesson has its own README with the specifics.

## Mental model

```text
                normal error path
    caller ────────────────────────→ callee
              (return err → nil)

                panic path
    caller ←═════════════════════════ callee
              (panic bubbles up)
```

- **Normal errors** flow down through returns. The caller decides what to do.
- **Panics** flow up through the call stack, unwinding as they go, until either:
  1. A deferred function calls `recover()` and catches it, or
  2. The goroutine's stack is fully unwound — in which case **the whole program crashes**.

The critical, non-obvious rule: **`recover` in goroutine A does not catch a panic in goroutine B.** Each goroutine's panic must be caught inside that goroutine. This is exactly what breaks in the `goroutine-panic/before/` example.

## The recovery idiom

```go
func safe() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("recovered: %v", r)
        }
    }()

    // ... code that might panic
}
```

Three things to notice:

1. **`recover` must be called inside a deferred function** — anywhere else it returns `nil` and does nothing.
2. **`recover` returns the value passed to `panic`** — usually a string or an `error`.
3. **After `recover` returns, execution continues *after* the deferred function** — the panicking function does not resume, but the caller of `safe()` does.

## When to use panic/recover

**Use them:**
- At the top of every long-running goroutine (worker pools, background jobs, request handlers running in their own goroutines) to prevent one bad request from crashing the whole process.
- At the top of every HTTP handler (usually via middleware) so one panicking handler returns a 500 instead of taking down the server.
- In library code, at package boundaries, to convert internal panics to `error` returns.

**Do not use them:**
- As a general control-flow mechanism ("try/catch"). Go prefers explicit `error` returns.
- To skip past invalid input. Validate at the boundary and return an `error`.
- To signal "expected" failure modes. If it's expected, it's not a panic.

## Try it yourself

1. Run the `before/` and `after/` tests side by side. Notice how the `before/` test crashes the whole test process, but `after/` reports the panic and keeps running.
2. Write a small program that panics inside a goroutine started by `go func() { panic("boom") }()`. Add `recover` to the outer function and observe that it does not catch the panic. Then add `recover` inside the goroutine and see it work.
3. Change the `after/` HTTP handler to also log the stack trace using `runtime/debug.Stack()`. That's what production servers do.

## Common pitfalls

- **`recover()` outside a deferred function.** Returns `nil`. Silently does nothing. Very easy to write.
- **Assuming a parent goroutine's recover covers its children.** It doesn't. Every goroutine needs its own `defer recover()`.
- **Recovering and continuing as if nothing happened.** A panic usually means invariants are broken — after recovering, at minimum log the stack trace and the recovered value. In many cases the right response is to log and *still* exit the request/worker.
- **`http.ErrAbortHandler`.** The `net/http` package uses this specific panic value as a signal — code that recovers panics in HTTP handlers must re-panic if the recovered value is `http.ErrAbortHandler`. See the `http-panic/after/` example.

## You've understood this lesson when...

- You can name at least two places where `defer recover()` belongs and at least two where it doesn't.
- You know why `recover` in the parent goroutine can't catch a child goroutine's panic.
- You can write a panic-recovery middleware for `net/http` from memory.

## Next

- **Next lesson:** [12-http-clients-and-servers](../12-http-clients-and-servers/) — the HTTP fundamentals that lesson 11's `http-panic` sub-lesson builds on.
