# Functions and packages

How to break a Go program into named, testable pieces — and what "exported" means in practice.

## Why it matters

Real programs are hundreds or thousands of functions across many files. Go organises them into **packages**: named collections of code that can be imported. Whether a function is visible outside its own package depends on a single, weird rule: **the first letter of its name.** Learning that rule early saves confusion later.

## Prerequisites

- Lesson 01: how to run a Go program.

## Run it

```bash
go test ./lessons/02-functions-and-packages
```

Expected output (last line):

```text
ok  	github.com/ocrosby/go-lab/lessons/02-functions-and-packages	0.005s
```

That's Go saying "the tests I found in this package passed."

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`math.go`](./math.go) | Four small functions (`Add`, `Subtract`, `Divide`, `Multiply`) in a package called `math`. |
| [`math_test.go`](./math_test.go) | Tests for those functions — a first look at `go test`. |

## The naming rule

Look at `math.go`:

```go
package math

func Add(x, y int) int { ... }
```

`Add` starts with a **capital letter**. That's Go's rule for **exported** — visible to code outside this package. If we'd written `add` (lowercase), it would be usable inside `math.go` and `math_test.go` but *not* from any other package.

Rule: **Capital first letter = public. Lowercase first letter = private to the package.**

That's the whole rule. No `public`, `private`, or `protected` keywords — just capitalization.

## Function shape

```go
func Add(x, y int) int {
    return x + y
}
```

Read it as: "function named `Add`, taking two `int` parameters called `x` and `y`, returning an `int`." Notice:

- **Type comes *after* the name.** `x int`, not `int x`. (This looks backwards if you're coming from C/Java. You'll adjust in an hour.)
- **Consecutive parameters of the same type share a type name.** `func Add(x, y int)` is shorthand for `func Add(x int, y int)`.
- **Return type comes last.** No `return type` keyword — just the type.

## Try it yourself

1. Add a new function `Square(x int) int` that returns `x * x`. Add a test for it that verifies `Square(3) == 9`. Run `go test`.
2. Rename `Add` to `add` (lowercase). Run `go test` again. What error do you get, and why?
3. Look at `Divide` — it silently returns 0 when `y == 0`. That's a bug (arguably). Rewrite it to return an `error` instead, and update the test.

## Common pitfalls

- **Forgetting to export a function** you want callers to use. If your test file (same package) can see it but no other package can, check the capitalization.
- **`Divide` returns a `float64` but uses integer division.** `Divide(7, 2)` returns `3.0`, not `3.5`, because `x / y` is `int / int` — the result is truncated *before* the cast. This is a real Go gotcha; the test file doesn't currently exercise it. See if you can write a test that catches it.
- **Package name doesn't have to match the folder name** but almost always does. This package is named `math`; the folder is `02-functions-and-packages`. Prefer matching them in real code.

## You've understood this lesson when...

- You can explain in one sentence why `Add` is callable from other packages but `add` is not.
- You can write a new function and test it with `go test`.
- You know where the `int` goes in a Go function signature (after the name, not before).

## Next

- **Next lesson:** [03-testing-basics](../03-testing-basics/) — a closer look at Go's testing framework, including table-driven tests.
