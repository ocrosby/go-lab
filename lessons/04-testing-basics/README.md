# Testing basics

Go's built-in test framework — no libraries, just `go test`.

## Why it matters

Go treats testing as a first-class part of the language. Every file whose name ends in `_test.go` is a test file, and every function whose name starts with `Test` is a test. There are no annotations, no configuration files, no `pytest` to install. You write a function, you call `go test`, you see the result. That's it.

## Prerequisites

- Lesson 01: how to run a Go program.
- Lesson 02: what a package is.

## Run it

```bash
go test -v ./lessons/04-testing-basics
```

Expected output (last few lines):

```text
=== RUN   TestHello
--- PASS: TestHello (0.00s)
PASS
ok  	github.com/ocrosby/go-lab/lessons/04-testing-basics	0.003s
```

The `-v` flag ("verbose") prints one line per test. Drop it once you're comfortable — bare `go test` just says `ok` or `FAIL`.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`hello.go`](./hello.go) | A `Hello()` function that returns a string. |
| [`hello_test.go`](./hello_test.go) | A test that calls `Hello()` and checks the result. |

## Anatomy of a Go test

Look at [`hello_test.go`](./hello_test.go):

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello()
    want := "Hello, world"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

Three conventions Go enforces:

1. **Filename ends in `_test.go`** — that's how the test tool finds it, and it's excluded from your production build.
2. **Function name starts with `Test`** and takes one parameter: `t *testing.T`. That parameter is your handle for reporting failures.
3. **Failure means calling `t.Errorf` (or `t.Fatalf`)**, not `return`. A test that never calls one of those passes. `t.Errorf` records the failure and continues; `t.Fatalf` records it and stops the test.

The `got`/`want` variable names are a **Go convention** — it's not required, but nearly all Go tests use those names. Stick with them.

## Table-driven tests

Lesson 02's `math_test.go` has a hint of the next pattern you'll see everywhere:

```go
tests := []struct {
    x, y   int
    result int
}{
    {1, 2, 3},
    {3, 5, 8},
    {7, -4, 3},
}

for _, tt := range tests {
    // ... call the function, check the result
}
```

This is a **table-driven test** — one loop that runs many cases. It's the dominant style in Go. Add a new case by adding a line to the slice, not by writing another `TestX` function.

## Try it yourself

1. Change `Hello()` to return `"Hello, Go"`. Run `go test`. Read the failure message and note the `got %q want %q` format.
2. Rename `TestHello` to `CheckHello`. Run `go test` again. Does it run?
3. Add a second test function, `TestHelloIsNotEmpty`, that asserts `Hello()` returns a non-empty string.
4. Add a `t.Log("started")` line at the top of `TestHello`. Run with `go test -v`; run without `-v`. Notice when the log appears.

## Common pitfalls

- **File named `hello_tests.go`** (extra "s") — Go won't pick it up. The suffix must be exactly `_test.go`.
- **Test function named `testHello`** (lowercase `t`) — same problem. Test names must start with capital `T`.
- **Using `return` after a failed check** — the test won't run further assertions but the failure IS recorded because you called `t.Errorf`. If you want to bail immediately, use `t.Fatalf` instead.
- **Forgetting `-v`** and being confused about which tests actually ran. `go test -v` shows you.

## You've understood this lesson when...

- You can write a new test file and function without looking anything up.
- You know the difference between `t.Errorf` (record and continue) and `t.Fatalf` (record and stop).
- You can explain in one sentence why the file must end in `_test.go`.

## Next

- **Next lesson:** [05-test-suites-and-refactor](../05-test-suites-and-refactor/) — a second testing style (Ginkgo) and using tests as a safety net during refactors.
