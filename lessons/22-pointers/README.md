# Pointers

`&x` gives you the address, `*p` follows it. Understanding pointers unlocks methods with mutating receivers, interface implementation rules, and half the standard library's function signatures.

> **Recommended before lesson 05 (composition).** Lesson 05 uses `*Accord` receivers throughout without ever explaining what the `*` means.

## Why it matters

Go's pointers are simpler than C's — no pointer arithmetic, no unions, no `void*` — but they're everywhere. Every function that mutates its argument, every method on a struct that needs to change fields, every "optional" parameter (`*int` where `nil` means "not set"), every interface value that holds a struct — pointers under the hood.

Once you can read `func (u *User) Rename(...)` and know why the `*` is there, most Go signatures stop being mysterious.

## Prerequisites

- Lesson 19: variables and types.
- Lesson 21: slices and maps (which are already reference-y and have prepared you for the idea).

## Run it

```bash
go run ./lessons/22-pointers
go test ./lessons/22-pointers
```

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`main.go`](./main.go) | `&`/`*` operators, nil pointers, value vs pointer receivers, `new` vs `&T{}`. |
| [`main_test.go`](./main_test.go) | Tests that pin the pass-by-value / pass-by-pointer difference. |

## The two operators

- **`&x`** — take the address of `x`. Yields a pointer to `x`.
- **`*p`** — follow the pointer `p`. Yields the value pointed at (or lets you assign to it).

```go
x := 42
p := &x         // p is a *int, pointing at x
fmt.Println(*p) // 42 — dereference to read
*p = 100        // dereference to write
fmt.Println(x)  // 100 — x was modified through p
```

**Zero value is `nil`.** An uninitialized pointer is `nil`. Dereferencing `nil` panics.

```go
var p *int
fmt.Println(p) // <nil>
// fmt.Println(*p) — panics!
if p != nil {
    fmt.Println(*p)
}
```

## Pass by value vs pass by pointer

Go passes function arguments **by value**. Every argument is copied. That means:

```go
func addOne(n int) { n++ }
x := 5
addOne(x)
fmt.Println(x) // 5 — unchanged
```

To let a function mutate the caller's variable, pass a pointer:

```go
func addOne(n *int) { *n++ }
x := 5
addOne(&x)
fmt.Println(x) // 6
```

This is why so many Go APIs take `*T`: mutation.

## Value receivers vs pointer receivers

The most common place you'll write `*` in Go: method receivers.

```go
type Counter struct { n int }

// Value receiver — copies the counter. Mutation to c.n is lost.
func (c Counter) IncrementValue() { c.n++ }

// Pointer receiver — modifies the original.
func (c *Counter) Increment() { c.n++ }

var c Counter
c.IncrementValue()
fmt.Println(c.n) // 0 — the increment was lost

c.Increment()
fmt.Println(c.n) // 1 — the pointer receiver worked
```

**When to use which:**

- **Pointer receiver (`*T`)** if the method mutates the receiver, if `T` is large (avoid copying), or if `T` contains a mutex/pointer/channel (copying would be a bug).
- **Value receiver (`T`)** otherwise — for small, immutable types.
- **Consistency matters.** If any method on `T` uses a pointer receiver, use pointer receivers for *all* methods on `T`. Mixing them is confusing and can cause interface-satisfaction surprises.

## `new(T)` vs `&T{}`

Two ways to allocate a struct on the heap and get a pointer to it:

```go
p := new(Counter)     // pointer to a zeroed Counter — fields all zero values
q := &Counter{n: 5}   // pointer to a Counter{n: 5}
```

Prefer `&T{...}` — it's more common in Go code, and it lets you initialize fields at the same time. Reserve `new` for types that don't have a struct literal (rare).

## Pointers to what?

Anything can have a pointer taken:

```go
&x        // pointer to a variable
&someStruct{...}  // pointer to a struct literal
&someArray[0]     // pointer to a slice/array element (rare — careful with slice regrowth)
```

You can also take the address of a struct field: `&user.Name` is a `*string` pointing at the `Name` field.

## Try it yourself

1. Change `IncrementValue` in `main.go` to a pointer receiver. What changes about how the tests behave?
2. Write a `Swap(a, b *int)` function that swaps the two values.
3. Try declaring `var p *int` and calling `*p`. See the panic. Now add `if p != nil` around it.
4. Change the `Counter` methods to use *both* value and pointer receivers. What does `go vet` say?

## Common pitfalls

- **Nil pointer dereference.** The most common Go panic. Always check for nil before dereferencing a pointer that could be nil (function returns, struct fields, map values that are `*T`).
- **Pointer to a range-loop variable.** Pre-Go-1.22, `for _, u := range users { save(&u) }` gave you the SAME pointer every iteration. Go 1.22+ fixed this. But keep it in mind when reading older code.
- **Mixing value and pointer receivers.** If half your methods take `T` and half take `*T`, only `*T` satisfies interfaces that require both — the value type doesn't. This is one of the top-3 interface bugs.
- **Returning a pointer to a stack-local**. In Go, safe! The compiler does escape analysis and moves the variable to the heap. In C this would be a bug; in Go it just works.
- **Comparing pointers.** `p1 == p2` compares *addresses*, not the pointed-at values. If you want value equality, dereference: `*p1 == *p2`.

## You've understood this lesson when...

- You can predict whether a function's changes to its argument will be visible in the caller, just by looking at the signature.
- You can name three cases where a method should use a pointer receiver.
- You can explain why `var p *int; fmt.Println(*p)` panics but `var m map[string]int; fmt.Println(m["x"])` doesn't.
- You can convert a value-receiver-only type into a pointer-receiver-only type and update the call sites.

## Next

- **Next lesson (recommended):** [23-error-handling](../23-error-handling/) — Go's `if err != nil` idiom, the pattern that defines idiomatic Go code.
