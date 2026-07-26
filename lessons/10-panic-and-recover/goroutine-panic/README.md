# Goroutine panic

A worker pool where jobs can panic. Shows the `before/` unsafe pattern (silent goroutine death) and the `after/` fix (per-goroutine `defer recover()`).

## Prerequisites

- Lesson 07: goroutines and channels.
- Lesson 10 parent README.

## Run it

```bash
# Correct pattern — should pass
go test ./lessons/10-panic-and-recover/goroutine-panic/after/

# Broken pattern — supposed to demonstrate the failure
go test ./lessons/10-panic-and-recover/goroutine-panic/before/
```

The `before/` run will panic and terminate. That is the point.

## What's in this folder

| Path | What it demonstrates |
|---|---|
| [`before/worker.go`](./before/worker.go) | Worker pool with no panic recovery — a single panicking job kills its worker goroutine silently. |
| [`before/worker_test.go`](./before/worker_test.go) | Tests that surface the resulting hang / crash. |
| [`after/worker.go`](./after/worker.go) | Same pool wrapped in `defer func() { if r := recover(); r != nil { ... } }()` per worker. |
| [`after/worker_test.go`](./after/worker_test.go) | Tests that verify the pool keeps working after a panic. |

## The bug in `before/`

Each worker looks something like:

```go
go func(workerID int) {
    for {
        select {
        case <-ctx.Done():
            return
        case job := <-wp.jobs:
            process(job)   // may panic
        }
    }
}(i)
```

If `process(job)` panics, the goroutine unwinds and exits. The rest of the pool doesn't notice. The next request to the pool queues fine, but there's one fewer worker to pick it up — eventually the pool is silently empty and jobs queue up forever.

## The fix in `after/`

```go
go func(workerID int) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("worker %d recovered from panic: %v", workerID, r)
            // optionally: restart or signal the pool
        }
    }()

    for { /* same select as before */ }
}(i)
```

Two things about this:

1. **The `defer` must be inside the goroutine's function body.** A `defer recover()` in the caller of `go func(...)` does *nothing* — it runs in the wrong goroutine.
2. **Just recovering isn't enough.** The worker's goroutine exits after the deferred function returns. The `after/` implementation goes further — it uses a `PanicHandler` callback and (optionally) restarts the worker so the pool stays at full capacity.

## Try it yourself

1. Move the `defer recover()` from inside the goroutine to *outside* the `go func(...)` call. Watch it stop working.
2. Add a job that panics with a specific value (say, `panic("bad job")`). Verify the `after/` pool logs it and keeps processing subsequent jobs.
3. Extend `after/` to send the recovered panic as a `JobResult` back to the caller (see `PanicInfo` field on `JobResult`) — the caller then knows *which* job blew up.

## Common pitfalls

- **Assuming the parent's `defer recover()` catches child goroutines.** It doesn't. Every `go` you write needs to consider whether its goroutine needs its own recovery.
- **Recovering and just returning.** The pool now has one fewer worker. If your workload matters, restart or signal.
- **Recovering in a tight retry loop.** If a job panics *every* time it's retried, you burn CPU. Cap the retry count or move the bad job aside.
