# Defer and cleanup

`defer` runs a function call when the surrounding function returns. It's Go's answer to Java's `finally` and C++'s destructors — and the primary tool for reliable cleanup even when errors happen.

> **Recommended before lessons 11 (panic/recover) and 11 (HTTP).** Both use `defer` prominently — closing HTTP response bodies, unlocking mutexes, recovering from panics — without ever explaining the timing rules.

## Why it matters

Half of "how do I handle errors without leaking resources?" in Go is answered by `defer`. Open a file? `defer f.Close()`. Acquire a lock? `defer mu.Unlock()`. Started a transaction? `defer tx.Rollback()` (before you commit). It's simple, universal, and has a couple of timing rules that catch newcomers.

## Prerequisites

- Lesson 21: control flow (`for`, `return`).
- Lesson 24: error handling (why you need cleanup in error paths).

## Run it

```bash
go test ./lessons/25-defer-and-cleanup
```

Expected output includes the LIFO-order demonstration.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`defer_demo.go`](./defer_demo.go) | The three rules of `defer` (LIFO order, argument evaluation at defer time, execution just before return) and common cleanup patterns. |
| [`defer_demo_test.go`](./defer_demo_test.go) | Tests that pin each timing rule so you can experiment. |

## The three rules

### 1. Deferred calls run in LIFO order

```go
func main() {
    defer fmt.Println("first")
    defer fmt.Println("second")
    defer fmt.Println("third")
}
// Output:
// third
// second
// first
```

Last deferred, first executed. Read the function bottom-up to see cleanup order.

### 2. Arguments are evaluated at the `defer` statement, not at execution

```go
x := 1
defer fmt.Println(x) // captures x=1 NOW
x = 2
// prints "1", not "2"
```

If you want the value at defer-execution time, close over it (use `func()` wrapper):

```go
x := 1
defer func() { fmt.Println(x) }() // captures x by reference
x = 2
// prints "2"
```

### 3. Deferred calls run just before the enclosing function returns

Including panics, including named-return-value mutations. This is what makes `defer` reliable — cleanup runs no matter how the function exits.

## The common cleanup patterns

**File close:**

```go
f, err := os.Open(path)
if err != nil {
    return err
}
defer f.Close()
// ... use f
```

Put `defer f.Close()` on the line right after the successful open. Never before — if `Open` returned nil for `f`, `defer f.Close()` would panic.

**Mutex unlock:**

```go
mu.Lock()
defer mu.Unlock()
// ... critical section
```

Idiomatic Go. Once you write `mu.Lock()`, the *very next line* should almost always be `defer mu.Unlock()`.

**HTTP response body close:**

```go
resp, err := client.Do(req)
if err != nil {
    return err
}
defer resp.Body.Close()
// ... read resp.Body
```

Missing this leaks a connection. `defer` immediately after the error check catches every exit path.

**Panic recovery:**

```go
func safe() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    // ... code that might panic
    return nil
}
```

`recover` only works inside a deferred function. This is the pattern lesson 11 (panic-and-recover) covers in depth.

**Named-return-value manipulation:**

```go
func work() (result int, err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("work: %w", err)
        }
    }()
    // ... code that sets err
    return
}
```

Deferred functions can inspect and modify named return values before the function actually returns. This is how libraries wrap every error path in one place.

## Try it yourself

1. Add a `defer fmt.Println(i)` inside a `for i := 0; i < 3; i++` loop. What order do the numbers print in? Why?
2. Rewrite `defer fmt.Println(x)` (with argument captured at defer) as `defer func() { fmt.Println(x) }()` (with closure). Change `x` between the defer and the return. Which one shows the updated value?
3. Write a function that opens two files and closes both with `defer`. Verify the closes happen in reverse order.
4. Combine `defer` with `recover` to turn a panic into an error return. Compare with lesson 11's HTTP-handler and worker-pool examples.

## Common pitfalls

- **`defer` inside a loop.** `for _, f := range files { defer f.Close() }` doesn't close each file after one iteration — it accumulates them and closes them all at function exit. Usually not what you want. Wrap the iteration body in a function so each iteration's defer fires immediately:
  ```go
  for _, path := range paths {
      func() {
          f, _ := os.Open(path)
          defer f.Close()
          // ... use f
      }()
  }
  ```
- **Deferring before the error check.** `defer f.Close()` before `if err != nil { return }` is a nil-pointer panic waiting to happen when `Open` failed and returned `nil` for `f`. Always error-check first.
- **Expecting an argument to be evaluated late.** `defer log.Println(counter)` captures `counter`'s value NOW. If you want the value later, use a closure.
- **Deferring a `recover()` call that isn't wrapped in a func.** `defer recover()` runs `recover` at `defer` time (returning nil, doing nothing). It must be `defer func() { recover() }()`.
- **Cost.** Every `defer` used to be measurably slower than a plain call. Go 1.14 dropped the cost to near-zero. Don't avoid `defer` for micro-perf reasons in application code.

## You've understood this lesson when...

- You can predict what a stack of three `defer` statements will print, in order.
- You know why `defer fmt.Println(x)` and `defer func() { fmt.Println(x) }()` behave differently when `x` changes.
- You can name three cleanup situations where `defer` is the idiomatic fix (files, locks, HTTP body).
- You can spot the "defer inside a for loop" bug in a code review.

## Next

- **Next lesson (recommended):** [06-composition](../06-composition/) — you now have every fundamental you need to start the main syllabus properly.
