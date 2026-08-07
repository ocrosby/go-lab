# Control flow

Go's `if`, `switch`, and `for`. **`for` is Go's only loop** — no `while`, no `do-while`, no `foreach`. Once you know its three forms, you can loop over anything.

> **Recommended before lesson 02.** Lesson 02 introduces functions and tests without ever showing you how to write `if` or `for`. This lesson fixes that.

## Why it matters

Go's control flow is deliberately small. Instead of many keywords, you get a few flexible ones. `for` alone covers what other languages split into `for`/`while`/`do-while`/`foreach`. Once you internalize the three forms of `for` and the "expression-less" `switch`, you can read any Go code you meet.

## Prerequisites

- Lesson 20: variables and types.

## Run it

```bash
go run ./lessons/21-control-flow
```

Expected output:

```text
--- if / else ---
5 is positive
0 is zero

--- if with init ---
found 42 in map

--- switch ---
Monday is a weekday
Saturday is a weekend

--- expression-less switch ---
age 25 → adult

--- for: classic C-style ---
0 1 2 3 4

--- for: condition-only (Go's 'while') ---
countdown: 3 2 1 done

--- for: infinite ---
break after 3

--- for range: slice ---
0:apple 1:banana 2:cherry

--- for range: map ---
(order may vary)
```

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`main.go`](./main.go) | Every control-flow form Go has, in one runnable file. |

## `if` — with a twist

Standard shape:

```go
if x > 0 {
    fmt.Println("positive")
} else if x < 0 {
    fmt.Println("negative")
} else {
    fmt.Println("zero")
}
```

Go's twist: `if` can have an **init statement**. This is idiomatic for functions that return `(value, error)`:

```go
if v, ok := lookup(key); ok {
    fmt.Println("found", v)
}
// v and ok are out of scope here
```

`v` and `ok` are scoped to the `if` block, which is exactly what you usually want — you use the return values inside the branch and don't need them after.

## `switch` — pattern-y without the ceremony

Cleaner than nested `if`s:

```go
switch day {
case "Sat", "Sun":
    fmt.Println("weekend")
case "Mon", "Tue", "Wed", "Thu", "Fri":
    fmt.Println("weekday")
default:
    fmt.Println("unknown")
}
```

Key differences from C/Java `switch`:

- **No `break` needed.** Cases don't fall through by default. Use explicit `fallthrough` if you want it (rare).
- **Multiple values per case.** `case "Sat", "Sun":` matches either.
- **Any type.** Switch on strings, ints, structs — anything comparable.

### Expression-less switch

Omit the expression and each `case` is a full boolean. Cleaner than a long `if/else if` chain:

```go
switch {
case age < 13:
    return "child"
case age < 20:
    return "teen"
case age < 65:
    return "adult"
default:
    return "senior"
}
```

Type switches (`switch v := x.(type)`) come up in interface work — you'll see them in lesson 13.

## `for` — the only loop

Three shapes, one keyword.

### Classic three-part

```go
for i := 0; i < 10; i++ {
    fmt.Println(i)
}
```

### Condition-only (like `while`)

```go
for count > 0 {
    fmt.Println(count)
    count--
}
```

### Infinite (`for {}`)

```go
for {
    if done { break }
    // ... do work
}
```

`break` exits the loop. `continue` skips to the next iteration.

### `for range` — the "foreach"

Iterate a slice, string, map, or channel:

```go
fruits := []string{"apple", "banana", "cherry"}
for i, v := range fruits {
    fmt.Println(i, v)
}

for _, v := range fruits { ... }   // discard the index
for i := range fruits { ... }      // discard the value
```

Ranging over a map gives you `key, value` (in **random** order per iteration — see pitfalls). Ranging over a string gives you `byteIndex, rune`.

Since Go 1.22, `for range n` (where `n` is an int) iterates from 0 to n-1 — the shortest way to do a count-N loop:

```go
for i := range 10 {  // Go 1.22+
    fmt.Println(i)
}
```

## Try it yourself

1. Rewrite the classic `for i := 0; i < 5; i++` in lesson code as `for i := range 5`. Same output?
2. Change a `switch` case to have `fallthrough`. What happens on that case?
3. Write a function `classify(n int) string` that returns "negative", "zero", or "positive" using an expression-less switch.
4. Iterate a map twice in the same program and print the keys. Notice the order changes (or may) between iterations — this is by design.

## Common pitfalls

- **Loop variable capture in goroutines (pre-Go-1.22).** Old Go code has the classic `for i := 0; i < 3; i++ { go func() { fmt.Println(i) }() }` bug where all goroutines print `3`. **Go 1.22+ fixed this**: the loop variable is fresh per iteration. This repo pins Go 1.26, so you're safe. But if you read older Go code, watch for it.
- **Map iteration order is random.** Do not rely on it. Sort the keys first if you need deterministic order.
- **Missing `fallthrough`.** Go switches don't fall through. If you write two cases hoping to run both, they don't.
- **`continue` in a range loop.** `continue` skips to the next iteration of the innermost `for`. If you're nested, use a labelled `continue OUTER`.

## You've understood this lesson when...

- You can write a countdown loop three different ways (classic, condition-only, `for range n`).
- You can predict whether a `switch` case falls through.
- You know why iterating a map twice may print keys in different orders.
- You can use an `if` init statement to keep short-lived variables tightly scoped.

## Next

- **Next lesson (recommended):** [02-functions-and-packages](../02-functions-and-packages/) — put what you learned into a testable function.
