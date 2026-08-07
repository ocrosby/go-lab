# Channel patterns

Named patterns that recur every time you write concurrent Go: **done/quit signalling**, **pipelines**, and the moves that give them their names — **fan-out** and **fan-in**.

## Why it matters

Lesson 08 taught the three primitives (goroutines, channels, `select`). Real code composes them into repeatable shapes. Once you can name a shape ("this is a fan-in") you can spot it in someone else's code, and you can design your own concurrent code by picking a shape rather than inventing one from scratch.

## Prerequisites

- Lesson 08: goroutines, channels, `select`, closing channels.

## Run it

The pattern demos are standalone files marked `//go:build ignore` — run each with `go run` on the file itself:

```bash
go run ./lessons/09-channel-patterns/done.go
go run ./lessons/09-channel-patterns/pipeline.go
```

The `channels/` sub-package has runnable Ginkgo tests:

```bash
go test ./lessons/09-channel-patterns/channels/...
```

## What's in this folder

| Path | Pattern |
|---|---|
| [`done.go`](./done.go) | **Done channel** — cooperative cancellation. Producer checks `<-done` in its `select` and exits when signalled. |
| [`pipeline.go`](./pipeline.go) | **Pipeline** — stages connected by channels. Each stage is a goroutine reading its input channel and writing its output. |
| [`channels/`](./channels/) | Small Ginkgo tests exercising channel behaviour. |

## The patterns

### Done / quit signalling

```go
func doWork(done <-chan bool) {
    for {
        select {
        case <-done:
            return               // asked to stop
        default:
            // one unit of work
        }
    }
}
```

The `done` channel is **receive-only** to the worker — it can listen but not signal. The caller holds the send side and closes it (or sends on it) when the worker should exit. Closing is usually cleaner than sending because multiple workers can all wake up from one `close`.

In real code, prefer `context.Context` over a hand-rolled `done` channel — `ctx.Done()` returns a channel with the same shape and carries cancellation reasons and deadlines. But learn the raw pattern first; `context` is a wrapper around it.

### Pipeline

```go
sliceToChannel(nums) → sq → sq → main
```

Three stages, each a goroutine. Values flow left to right through channels. Each stage:

1. Takes a receive-only input channel.
2. Returns a receive-only output channel.
3. Spawns a goroutine that reads input, transforms, writes output, and closes output when input is drained.

Directional channels (`<-chan int` for receive-only, `chan<- int` for send-only) are documentation that compiles — the compiler stops you from writing on a receive-only channel.

### Fan-out and fan-in

Not in this lesson's files yet — but named because they are the natural extensions:

- **Fan-out** — one input channel, N workers all reading from it. Distributes work.
- **Fan-in** — N input channels, one output channel. One forwarding goroutine per input, all writing to the same output. Merges streams.

Lesson 10 (`10-worker-pools`) is a full fan-out example. `docs/csp-and-go-concurrency.md` diagrams both.

## Try it yourself

1. Modify [`done.go`](./done.go) to use `context.WithTimeout` instead of a raw `done` channel. The `case <-done` becomes `case <-ctx.Done()`. Does the rest of the function change?
2. Add a third stage to [`pipeline.go`](./pipeline.go) that adds 1 to each value. Compose it as `sq → addOne` and print the results.
3. In `pipeline.go`, remove the `close(out)` in one of the stages. Run it. What happens, and why? (Hint: the consumer's `for v := range ch` never terminates.)

## Common pitfalls

- **Forgetting to `close` a channel.** A `for v := range ch` loop only exits when the channel is closed *and* drained. Forgetting the close is a goroutine leak.
- **Closing from the receiver.** Only the sender should close. Closing from the receiver risks a panic if the sender is still writing.
- **Two senders, one channel, close called twice.** `close` panics on an already-closed channel. If multiple goroutines share the send side, coordinate the close (often with a `sync.WaitGroup`).
- **Confusing `close(ch)` with sending `nil`.** They are different. Close signals "no more values ever"; sending nil (only valid for pointer/interface channels) sends a value that happens to be nil.

## You've understood this lesson when...

- You can sketch a pipeline with three stages and directional channels from memory.
- You can explain in one sentence why `context.Context` is preferred over a raw done channel in production code.
- Given a snippet, you can spot which pattern (done, pipeline, fan-out, fan-in) is being used.

## Related deep-dive

- [`docs/csp-and-go-concurrency.md`](../../docs/csp-and-go-concurrency.md) — where these patterns come from (Hoare's CSP) and how each maps to Go.

## Next

- **Next lesson:** [10-worker-pools](../10-worker-pools/) — a fan-out pattern applied end-to-end, with bounded concurrency.
