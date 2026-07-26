# Type assertions and type switches

`x.(T)` — the assertion form — and `switch v := x.(type)` — the switch form — extract concrete types out of an interface value. Together with `errors.As` (a specific application), they're how you unwrap `any` into something you can call methods on.

> **Recommended after lesson 06 (interfaces).** Interface values are the input; type assertions are how you get back to concrete types.

## Why it matters

`any` and `interface{}` values hold ANYTHING. To do useful work with one, you have to extract the concrete type at runtime. The two operators that do this — the assertion and the type switch — appear in adapters, error handling, JSON round-trips, and any code that talks to `reflect`. Once you know both syntaxes, unwrapping any interface value becomes routine.

## Prerequisites

- Lesson 06: interfaces (satisfy them structurally).
- Lesson 23: error handling (already saw `errors.As` — this generalizes the pattern).

## Run it

```bash
go test -race ./lessons/31-type-assertions
```

Expected: 7 passes.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`assertions.go`](./assertions.go) | The two assertion forms, the type switch, an "any → structured" adapter, and the compile-time interface-satisfaction check idiom. |
| [`assertions_test.go`](./assertions_test.go) | Tests for each form. |

## Two syntaxes

### The single-return form (panics on mismatch)

```go
var x any = "hello"
s := x.(string)          // s is "hello"
n := x.(int)             // PANICS — x is a string, not an int
```

The single-return form panics if the assertion fails. Use only when the type is guaranteed (rare) or you truly want to crash on a bug.

### The comma-ok form (safe)

```go
s, ok := x.(string)      // s == "hello", ok == true
n, ok := x.(int)         // n == 0, ok == false — no panic
```

Same shape as map lookup and channel receive. This is what you use 99% of the time.

## The type switch

When you want to dispatch on which of several types an interface value holds:

```go
switch v := x.(type) {
case int:
    fmt.Println("int:", v)
case string:
    fmt.Println("string:", v)
case []byte:
    fmt.Println("bytes:", string(v))
case nil:
    fmt.Println("nil")
default:
    fmt.Printf("unknown type %T\n", v)
}
```

Inside each case, `v` has that case's specific type — you can call type-specific methods (`v.String()` if `v` is `time.Time`, etc.) without a further cast.

Case-multiple-types is allowed but `v` becomes the interface type:

```go
switch v := x.(type) {
case int, int32, int64:
    // v is still `any` here — you have to assert further
    fmt.Println("some int:", v)
}
```

Include a `default` case in production code — new types get added to systems and silent fallthrough hides bugs.

## `errors.As` — a specialized type assertion

From lesson 23. `errors.As(err, &target)` walks the error's wrap chain looking for something the target's type. It's a type assertion that follows wraps:

```go
var verr *ValidationError
if errors.As(err, &verr) {
    fmt.Println(verr.Field)
}
```

Prefer `errors.As` over `err.(*ValidationError)` — the latter misses wrapped errors.

## Compile-time interface-satisfaction check

A very common idiom in Go codebases:

```go
var _ Vehicle = (*Accord)(nil)
```

Read as: "assign `nil` (of type `*Accord`) to a discarded variable of type `Vehicle`." If `*Accord` doesn't satisfy `Vehicle`, this line fails to compile. Costs nothing at runtime; catches missing-method bugs where the type is defined, not where callers use it.

You'll see this at the top of adapter files throughout this repo — e.g. `lessons/06-interfaces-and-mocking` uses it in the test files.

## Try it yourself

1. Write a `Describe(x any) string` that returns different strings for `int`, `string`, `bool`, and any other type. Use a type switch.
2. Try `x.(int)` on a `*int`. Does it work? Why not? (Answer: `*int` and `int` are different types.)
3. Add `case int, int64:` to a type switch. Notice `v` becomes `any` in that case, not `int`.
4. Add a compile-time check `var _ error = (*ValidationError)(nil)` from lesson 23's ValidationError to its file. Verify it catches missing-method bugs (remove `Error() string` and watch the check fail).

## Common pitfalls

- **Single-return assertion without a guard.** `x.(string)` panics if `x` isn't a string. Use the comma-ok form.
- **Asserting on a pointer when the value is stored.** `v.(User)` won't match if the interface holds a `*User`. Two different types.
- **Missing `default` in a type switch.** New types get added; old switches silently fall through. Add a `default` that logs or errors.
- **Assertion after `nil` check missed.** A `nil` interface is different from an interface *holding* a nil pointer (the "typed nil" hazard from `rules/go-conventions.md`). `x != nil` doesn't tell you the wrapped value is non-nil.
- **`x.(any)` — always ok.** `x.(any)` succeeds for every non-nil interface. Rarely what you want.
- **Confusing type assertions with type conversions.** `int(3.14)` is a **conversion** — same operator syntax as C-style casts, works at compile time between compatible types. `x.(int)` is an **assertion** — runtime check on an interface value. Different mechanisms.

## You've understood this lesson when...

- You can write both forms of type assertion and know when to use each.
- You can spot the difference between a type switch and a regular switch.
- You know why `errors.As` is preferred over a direct assertion on an error.
- You can add a compile-time interface-satisfaction check to your own adapter type.

## Next

- **Next lesson:** [32-file-io-and-cli](../32-file-io-and-cli/) — `os.Args`, `flag`, `os.Open`, `io.Reader`, `bufio.Scanner` — the last set of standard-library basics you need to build a real command-line tool.
