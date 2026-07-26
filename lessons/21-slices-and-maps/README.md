# Slices and maps

Go's workhorse collection types. Slices for ordered lists, maps for key-value lookups. Both use `make`, both grow with `append`/assignment, both have a comma-ok idiom that trips up newcomers.

> **Recommended before lesson 05 (composition)** and every lesson after. Slices and maps appear on nearly every page from lesson 05 onward without ever being introduced formally.

## Why it matters

Almost every non-trivial Go program uses slices or maps or both. Beginners often confuse arrays and slices (they're different), pass slices "by value" and are surprised when the callee's mutations show up in the caller (they don't, but also they *do*, in a specific way), or use `if v := m[key]; v != 0` to test map presence and are burned when zero is a legitimate value (use the comma-ok form).

## Prerequisites

- Lesson 19: variables and types.
- Lesson 20: control flow (`for range`).

## Run it

```bash
go run ./lessons/21-slices-and-maps
```

And run the tests:

```bash
go test ./lessons/21-slices-and-maps
```

Expected output (last line):

```text
ok  	github.com/ocrosby/go-lab/lessons/21-slices-and-maps	...
```

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`main.go`](./main.go) | Slice creation, `append`, indexing, iteration, and the "shared backing array" surprise. Map creation, insert, delete, comma-ok lookup, iteration. |
| [`main_test.go`](./main_test.go) | Tests that pin the observable behaviour — for `append`'s growth semantics and for the comma-ok idiom. |

## Arrays vs slices

- **Array** — fixed-size, value type. `var a [3]int` is a distinct type from `[4]int`. Copying an array copies every element.
- **Slice** — dynamically-sized view into an underlying array. `[]int` is one type regardless of length. Copying a slice copies the *header* (pointer + len + cap), not the elements.

You almost never write arrays in application code. Slices are what you want.

## Slice essentials

Create:

```go
var xs []int             // nil slice — len 0, cap 0, no backing array yet
xs := []int{1, 2, 3}     // literal
xs := make([]int, 5)     // length 5, all zeros
xs := make([]int, 0, 10) // length 0, capacity 10 — pre-allocated
```

Grow:

```go
xs = append(xs, 42)         // reassign — append may allocate a new backing array
xs = append(xs, 1, 2, 3)    // append several at once
xs = append(xs, other...)   // append another slice (variadic spread)
```

**Rule**: always reassign the result of `append`. If the underlying array can't fit the new element, `append` allocates a bigger one and returns a slice pointing at it. Ignoring the return value silently loses the addition.

Slice a slice:

```go
xs[1:3]     // elements 1 and 2 (end exclusive)
xs[:2]      // first two elements
xs[3:]      // from index 3 to end
xs[:]       // the whole thing (rare — just use xs)
```

Slicing does **not** copy. The result shares the backing array. If you want a copy, use `copy()` or `slices.Clone()`.

## The "shared backing array" surprise

```go
xs := []int{1, 2, 3, 4, 5}
ys := xs[1:4]       // ys = [2, 3, 4], shared backing with xs
ys[0] = 999
fmt.Println(xs)     // [1, 999, 3, 4, 5] — surprising!
```

This is Go's biggest slice gotcha. If you're mutating a slice and don't want the caller to see the change, `slices.Clone(ys)` first.

## Maps

Create:

```go
var m map[string]int         // nil map — reading is OK, writing panics
m := map[string]int{}         // empty map, ready to write
m := make(map[string]int)     // same as above
m := make(map[string]int, 100) // pre-sized hint for the runtime
```

Assign, read, delete:

```go
m["answer"] = 42
v := m["answer"]        // 42
v := m["missing"]       // 0 — zero value for int
delete(m, "answer")
```

### The comma-ok idiom

`m[k]` always returns a value — the zero value if the key is missing. That means you can't distinguish "the key isn't there" from "the value happens to be zero." The two-value form does:

```go
v, ok := m["answer"]
if ok {
    fmt.Println("found:", v)
} else {
    fmt.Println("missing")
}
```

Idiomatically wrapped in an `if` init:

```go
if v, ok := m["answer"]; ok {
    fmt.Println("found:", v)
}
```

Use the two-value form whenever the zero value is a legitimate stored value.

## Sets via maps

Go has no built-in set. `map[T]struct{}` is the idiomatic sub — `struct{}` uses zero bytes of storage:

```go
seen := map[string]struct{}{}
seen["alice"] = struct{}{}
if _, ok := seen["alice"]; ok { ... }
```

`map[T]bool` also works and is less noisy at the cost of a byte per entry.

## Try it yourself

1. Take a slice `xs := []int{1, 2, 3}`. Assign a slice of it to `ys` and mutate `ys[0]`. What does `xs` look like? Now redo with `slices.Clone(xs)` first (Go 1.21+).
2. Create a map that counts word occurrences in a string. Hint: `for _, w := range strings.Fields(s) { counts[w]++ }`.
3. Delete a key from a map inside a `for range` loop over that map. Does it work safely? (Yes — Go's spec allows it.)
4. Try assigning to a nil map: `var m map[string]int; m["a"] = 1`. What happens?

## Common pitfalls

- **Forgetting to reassign `append`.** `append(xs, 1)` doesn't do what you want — you need `xs = append(xs, 1)`. This one bug per hour is normal for the first week.
- **Nil map write panics.** `var m map[K]V` gives you a nil map. Reads return the zero value, writes panic. Always `make` or `{}` before writing.
- **Map iteration order is random.** Deliberate — the runtime randomizes to prevent code from accidentally depending on order. Sort keys explicitly if you need a stable order.
- **Slice-of-slices from `append` capacity surprise.** Building a `[][]int` with append across iterations can produce weird results when the inner slice's capacity is shared. When in doubt, `slices.Clone` before appending the row.
- **`len(m)` is O(1), fine to call in loops.** `len(s)` for a slice is also O(1). But `len(str)` on a string returns *bytes*, not characters (see lesson 19).

## You've understood this lesson when...

- You can explain why `append(xs, 1)` sometimes returns a slice that points at a different backing array than `xs`.
- You know the difference between `xs[1:3]` and `slices.Clone(xs[1:3])`.
- You can write the comma-ok map lookup from memory.
- You know why writing to `var m map[string]int` panics.

## Next

- **Next lesson (recommended):** [22-pointers](../22-pointers/) — the last piece of Go you need before methods and interfaces start clicking.
