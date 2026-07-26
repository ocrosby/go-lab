# Concurrency

A beginner-friendly tour of Go's concurrency primitives: **goroutines**, **channels**, and the **`select` statement**. These three building blocks compose every higher-level concurrency pattern in Go.

## Mental Model

Go's concurrency follows the CSP (Communicating Sequential Processes) model: independent units of work coordinate by passing messages through channels rather than by sharing memory and locking.

> Don't communicate by sharing memory; share memory by communicating.
> — Go proverb

A channel is a typed FIFO queue. A goroutine is an independently-running function. `select` lets one goroutine wait on multiple channels at once. That is the entire model.

## Goroutines

A **goroutine** is a function running independently and concurrently with the caller. Spawn one by prefixing a function call with `go`:

```go
go doWork()
```

Key properties for beginners:

- **Lightweight.** Each goroutine starts with a tiny stack (~2 KB) that grows as needed. A program can run hundreds of thousands of them without exhausting memory.
- **Multiplexed onto OS threads.** The Go runtime schedules goroutines onto a small pool of OS threads — you do not pay one OS thread per goroutine.
- **`main` is a goroutine too.** When `main` returns the program exits, and any still-running goroutines are dropped on the floor.
- **No return value.** A goroutine returns nothing to its caller. Use a channel (or a shared variable + synchronization) to send results back.

See [`primitives3.go`](./primitives3.go) for a long-running goroutine driven by `select`, and [`done.go`](./done.go) for cooperative cancellation.

## Channels

A **channel** is a typed, FIFO queue used to pass values between goroutines. Create one with `make`:

```go
ch := make(chan int)        // unbuffered
ch := make(chan int, 10)    // buffered, capacity 10
```

Send and receive use the `<-` operator:

```go
ch <- 42        // send 42 into ch
x := <-ch       // receive a value from ch
```

### Unbuffered vs buffered

| Kind | Send blocks until... | Use when |
|---|---|---|
| Unbuffered (`make(chan T)`) | Another goroutine is ready to receive | You want **synchronization** — sender and receiver meet at the channel |
| Buffered (`make(chan T, N)`) | The buffer is full | You want **decoupling** — the sender can outpace the receiver by up to `N` items |

Communication on a buffered channel is asynchronous up to its capacity; once the buffer fills, sends block just like on an unbuffered channel.

### Direction

A channel parameter can be restricted to send-only or receive-only, and the compiler enforces it:

```go
func producer(out chan<- int) { /* can only send */ }
func consumer(in <-chan int)  { /* can only receive */ }
```

This is documentation that compiles. See [`pipeline.go`](./pipeline.go) for stages wired together with directional channels.

### Closing

The sender closes a channel when no more values will arrive:

```go
close(ch)
```

After closing:

- Receives still work — they drain remaining buffered values, then return the zero value.
- A `for v := range ch` loop exits when the channel is closed and drained.
- **Sending on a closed channel panics.** Only the sender should close, and only once.
- The two-value receive form distinguishes "got a value" from "channel is closed and empty":

```go
v, ok := <-ch
if !ok {
    // channel is closed and drained
}
```

## The `select` Statement

`select` waits on multiple channel operations at once and proceeds with whichever is ready first. It is the multiplexer for goroutine communication:

```go
select {
case msg := <-ch1:
    fmt.Println("from ch1:", msg)
case ch2 <- value:
    fmt.Println("sent to ch2")
case <-time.After(time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("nothing ready right now")
}
```

Behavior:

- If multiple cases are ready, one is chosen at random — never assume an order.
- If no case is ready and a `default` exists, `default` runs (non-blocking select).
- If no case is ready and no `default`, the select **blocks** until one becomes ready.

Common uses: timeouts, cancellation (`<-done`), fan-in from multiple producers, and non-blocking sends or receives.

See [`primitives.go`](./primitives.go) for `select` choosing between two channels.

## Beyond the Three Primitives

The three above are the language-level primitives. The standard library adds higher-level coordination tools you'll meet quickly:

- **`sync.WaitGroup`** — wait for a known number of goroutines to finish.
- **`sync.Mutex` / `sync.RWMutex`** — protect shared state when channels would be awkward.
- **`sync.Once`** — run an initializer exactly once.
- **`context.Context`** — propagate cancellation and deadlines across goroutine trees. In production code, prefer `context` over a hand-rolled `done` channel.
- **`errgroup.Group`** (`golang.org/x/sync/errgroup`) — like `WaitGroup` but with first-error short-circuit and context cancellation.

## Files in This Directory

| File | Demonstrates |
|---|---|
| [`primitives.go`](./primitives.go) | Two goroutines, two channels, `select` chooses whichever delivers first |
| [`primitives2.go`](./primitives2.go) | Buffered channel filled, closed, then drained with `range` |
| [`primitives3.go`](./primitives3.go) | Long-running goroutine driven by `select` |
| [`done.go`](./done.go) | Cooperative cancellation via a done channel |
| [`pipeline.go`](./pipeline.go) | Multi-stage pipeline with directional channels |

Each file declares its own `package main` and `main()`, so run them individually:

```bash
go run ./lessons/07-goroutines-and-channels/primitives.go
go run ./lessons/07-goroutines-and-channels/primitives2.go
```

The `done`- and `pipeline`-based primitives live in the next lesson:

```bash
go run ./lessons/08-channel-patterns/done.go
go run ./lessons/08-channel-patterns/pipeline.go
```

## Common Pitfalls

- **Goroutine leaks.** A goroutine blocked on a send or receive that never happens lives forever. Always give it a path to exit — `close`, `context`, or a done channel.
- **Sending on a closed channel panics.** Only the sender closes, and only once. If multiple senders share a channel, coordinate the close (often via a separate done signal).
- **`main` exits without waiting.** A program with goroutines still running ends when `main` returns. Use `WaitGroup`, channels, or `errgroup` to wait for completion.
- **Races on shared memory.** Run `go test -race` and `go run -race` to catch data races. Channels prevent most races by design — reach for `sync.Mutex` only when the data really is shared mutable state.

## Next Steps

- Whole programs that wire these primitives together: [`09-worker-pools`](../09-worker-pools/)
- Measuring concurrent code: [`15-benchmarks`](../15-benchmarks/)
