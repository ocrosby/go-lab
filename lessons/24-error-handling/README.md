# Error handling

Go's `if err != nil` idiom, `errors.New` and `fmt.Errorf %w`, sentinel errors, custom error types, and the `errors.Is` / `errors.As` inspection tools. This is the pattern that defines idiomatic Go code.

> **Recommended before lessons 08 (goroutines) and 11 (HTTP)** — both of which use `if err != nil` heavily without ever explaining it.

## Why it matters

Go does not have exceptions. Every function that can fail returns an `error` as its last return value. The caller checks it. If they don't want to handle it, they return it up the stack. That's the whole model. Once you've written a hundred `if err != nil { return fmt.Errorf("...: %w", err) }` lines, you'll wonder why anyone would want anything else.

`if err != nil` appears roughly 95 times in the rest of the lessons in this repo. This lesson is where you learn what it actually means.

## Prerequisites

- Lesson 20: types (the `error` interface).
- Lesson 23: pointers (custom error types often use pointer receivers).

## Run it

```bash
go test ./lessons/24-error-handling
```

Expected: 8 passes.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`errors.go`](./errors.go) | A tiny "user store" API with the four common error shapes: sentinel errors, wrapped errors, custom error types, and joined errors. |
| [`errors_test.go`](./errors_test.go) | Tests that exercise each error shape with `errors.Is` / `errors.As`. |

## The `error` interface

There's nothing special about `error`. It's a plain interface with one method:

```go
type error interface {
    Error() string
}
```

Any type that has an `Error() string` method IS an `error`. That's why you can define your own error types easily.

## The four shapes of Go errors

### 1. `errors.New` — the simplest form

For one-off error messages:

```go
if len(name) == 0 {
    return errors.New("name cannot be empty")
}
```

Fine for internal-to-a-function errors. Doesn't carry context, isn't inspectable — just a string.

### 2. Sentinel errors — comparable, package-exported

For errors the caller needs to check for and handle specially:

```go
// In the package that produces the error:
var ErrUserNotFound = errors.New("user not found")

// In the caller:
if errors.Is(err, users.ErrUserNotFound) {
    // handle "not found" specifically
}
```

The variable is `var`, not `const`, and it's exported (`ErrX` naming convention). `errors.Is` walks the wrap chain, so callers can match a sentinel even if it's been wrapped.

Common standard-library sentinels: `io.EOF`, `os.ErrNotExist`, `context.Canceled`, `sql.ErrNoRows`.

### 3. Wrapped errors — `fmt.Errorf %w`

For adding context to an error you're returning up the stack, without losing the original:

```go
if err := db.Query(...); err != nil {
    return fmt.Errorf("loading user %s: %w", id, err)
}
```

The `%w` verb wraps `err`. Callers can still `errors.Is(err, sql.ErrNoRows)` to check the underlying cause. Use `%w` in every wrapper — this is the single most important idiom in Go error handling.

**One `%w` per format string.** More than one gives you `errors.Join`-like behaviour (see below), but the tooling is friendlier if you stick to one.

### 4. Custom error types

When the caller needs structured information about the error (a field, a code, a URL), define a struct:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (v *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", v.Field, v.Message)
}
```

Callers use `errors.As` to extract the typed error:

```go
var verr *ValidationError
if errors.As(err, &verr) {
    fmt.Println("field:", verr.Field)
}
```

`errors.As` also walks the wrap chain, so wrapped custom errors are still findable.

## `errors.Is` vs `errors.As`

Confusingly-similar names, different jobs:

- **`errors.Is(err, target)`** — "does `err` match this sentinel?" Compares by identity (or a custom `Is` method). Use for sentinel checks.
- **`errors.As(err, &target)`** — "does `err` unwrap to something of this type?" Assigns the matched error to `target`. Use for custom error types.

Both walk the wrap chain. Neither is direct: prefer them over `err == ErrX` (misses wraps) or a type assertion (also misses wraps).

## `errors.Join` — Go 1.20+

Sometimes you want to return multiple errors from one function (e.g. every failing field in a validation). `errors.Join` wraps several into one:

```go
return errors.Join(
    fmt.Errorf("email: required"),
    fmt.Errorf("age: must be positive"),
)
```

`errors.Is` and `errors.As` see through joined errors too. The `Error()` string is the joined messages separated by newlines.

## Try it yourself

1. Add a `GetByEmail` method to the store in `errors.go` that returns `ErrUserNotFound` when there's no match. Add a test using `errors.Is`.
2. Add a `NotAuthorizedError` custom type carrying a `UserID` field. Return it from a new `Delete` method and test with `errors.As`.
3. Wrap a returned error with `fmt.Errorf("adding %s: %w", name, err)` in the `Add` method. Verify `errors.Is` still finds the underlying error.
4. Try `errors.Is(err, ErrUserNotFound)` when `err` is `nil`. What happens? (Answer: returns false, safely.)

## Common pitfalls

- **Comparing errors with `==`.** `if err == ErrX` misses any error that's been wrapped. Use `errors.Is`.
- **Type-asserting errors.** `if verr, ok := err.(*ValidationError); ok` misses wrapped ones. Use `errors.As`.
- **Missing `%w` in wraps.** `fmt.Errorf("failed: %v", err)` (using `%v` instead of `%w`) formats the error but doesn't wrap it. `errors.Is` and `errors.As` won't see through it. Use `%w`.
- **Overly-generic error messages.** `errors.New("failed")` tells the caller nothing. Include what was being attempted: `errors.New("open config file")` is much better.
- **Logging AND returning.** Pick one. Either log at the boundary where you can act on the error, or return it up. Logging every layer produces duplicated log lines and hides which one is the source.
- **Panicking instead of returning.** `panic` is for programmer errors (nil deref, out-of-bounds). Recoverable failures — bad input, network hiccups, missing files — should be errors, not panics.

## You've understood this lesson when...

- You can write a function that returns an error and know when to use `errors.New` vs `fmt.Errorf` vs a custom type.
- You can explain when to use `errors.Is` vs `errors.As`.
- You know why wrapping with `%w` matters and what breaks if you use `%v` instead.
- You can spot the "log AND return" anti-pattern in a code review.

## Next

- **Next lesson (recommended):** [25-defer-and-cleanup](../25-defer-and-cleanup/) — the `defer` keyword you'll want for reliable resource cleanup in error paths.
