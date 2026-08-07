# Pointers

`&x` gives you the address, `*p` follows it. Understanding pointers unlocks methods with mutating receivers, interface implementation rules, and half the standard library's function signatures.

> **Recommended before lesson 06 (composition).** Lesson 06 uses `*Accord` receivers throughout without ever explaining what the `*` means.

## Why it matters

Go's pointers are simpler than C's — no pointer arithmetic, no unions, no `void*` — but they're everywhere. Every function that mutates its argument, every method on a struct that needs to change fields, every "optional" parameter (`*int` where `nil` means "not set"), every interface value that holds a struct — pointers under the hood.

Once you can read `func (u *User) Rename(...)` and know why the `*` is there, most Go signatures stop being mysterious.

## Prerequisites

- Lesson 20: variables and types.
- Lesson 22: slices and maps (which are already reference-y and have prepared you for the idea).

## Run it

```bash
go run ./lessons/23-pointers
go test ./lessons/23-pointers
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

## Value semantics vs pointer semantics — the deep dive

The "when to use which" table above is the summary. This section is the *why* — the material that trips up nearly everyone learning Go and that the rest of the syllabus assumes you've internalized.

### What "value" actually copies

When you pass a struct by value (as a function argument or receiver), Go copies **every field** into a new struct. That copy is shallow — copies of primitives are independent, but copies of **reference-y fields** (`[]T`, `map[K]V`, `chan T`, `func`, and `*T`) share the underlying storage:

```go
type Bag struct {
    Label string   // primitive — copy is independent
    Items []string // slice — copy shares the backing array!
}

original := Bag{Label: "mine", Items: []string{"a", "b"}}
copy := original

copy.Label = "yours"           // affects only the copy
copy.Items[0] = "MUTATED"      // affects BOTH — same backing array

fmt.Println(original.Label)    // "mine"   — copy was independent
fmt.Println(original.Items[0]) // "MUTATED" — shared backing array
```

This is the single biggest "wait, what?" moment for people coming to Go from Java, Python, or JavaScript. **Copying a struct doesn't deep-copy its slices, maps, or channels.** Reason: those types are already implemented as headers pointing at a backing array/hashtable/queue. Copying the header gives you a second header pointing at the *same* backing storage.

Rule of thumb: **treat any struct with slice/map/channel fields as "shallow-copyable"** — you can hand it around cheaply, but mutations through the reference-y fields are visible to every other holder of the "copy."

### The reference types (`[]T`, `map[K]V`, `chan T`, `func`)

Because these types are already reference-y, you almost never need a *pointer to* them:

```go
// Rarely useful — []int is already a "pointer to backing array"
func addOne(xs *[]int) { ... }

// Idiomatic — the caller's slice is mutable through this parameter
func addOne(xs []int) { ... }
```

Two situations DO justify `*[]T`:

1. **You need to replace the slice itself** (change its length or capacity via `append`) so the change is visible to the caller. `func(*[]int)` lets you reassign. `func([]int)` only lets you mutate elements in place.
2. **The nil-vs-empty distinction matters** and you want to convey "unset."

For maps, channels, and funcs: essentially never take a pointer to them.

### Addressability — when method calls on values just work, and when they don't

Go automatically takes the address of a value when you call a pointer-receiver method on it — *if* the value is addressable:

```go
c := Counter{}
c.Increment()     // Go rewrites this to (&c).Increment(). Works.
```

But some values are **not addressable**, and pointer-receiver calls on them fail to compile:

```go
// Map values — NOT addressable
counters := map[string]Counter{"a": {n: 0}}
counters["a"].Increment()   // Does NOT compile

// Function return values — NOT addressable
NewCounter().Increment()    // Does NOT compile if NewCounter returns Counter (value)

// Interface values holding a value type — NOT addressable through the interface
var i Incrementer = Counter{}
i.Increment()               // Does NOT compile if Increment has a pointer receiver
```

**Fix pattern**: either return `*T` from constructors (`func NewCounter() *Counter`), or extract to a local variable first:

```go
c := counters["a"]
c.Increment()
counters["a"] = c   // reassign — the map value was a copy anyway

// Or:
counters := map[string]*Counter{"a": {n: 0}}   // store *Counter, addressable through the pointer
counters["a"].Increment()
```

This is why almost every real Go codebase stores pointers in maps and returns pointers from constructors.

### The interface-satisfaction asymmetry

If a type `Counter` has these methods:

```go
func (c Counter) Value() int { return c.n }     // value receiver
func (c *Counter) Increment() { c.n++ }         // pointer receiver
```

The two "method sets" are:

- **`Counter`'s method set** = `{Value}` (only the value-receiver methods)
- **`*Counter`'s method set** = `{Value, Increment}` (both — Go automatically dereferences to call the value method)

So if you have an interface:

```go
type Adjuster interface {
    Value() int
    Increment()
}
```

Then `*Counter` satisfies `Adjuster`, but plain `Counter` does not. This is why:

```go
var a Adjuster = Counter{}    // COMPILE ERROR — Counter's method set is missing Increment
var a Adjuster = &Counter{}   // OK
```

And it's why **the consistency rule matters**: if you're going to use pointer receivers anywhere, use them everywhere on the type. Otherwise `T` and `*T` satisfy different sets of interfaces, and callers get surprised.

### The immutable-value idiom (when to prefer VALUE receivers)

Some types are designed to be treated as immutable values — passed and returned by value throughout their APIs, methods return NEW values instead of mutating in place. The standard-library reference examples:

- `time.Time` — `t.Add(d)` returns a new `time.Time`, doesn't mutate `t`.
- `time.Duration` — a typed `int64`, always by value.
- `net/url.URL` — usually held as `*URL`, but the accessor methods take value receivers.
- `netip.Addr` — the modern IP address type, all-value everywhere.
- `uuid.UUID` (google/uuid) — 16-byte array, always by value.

Common traits: **small (a couple of words), semantically a "point in a value space" (a time, a duration, an ID), and mutation would be confusing** ("what does it mean to mutate a UUID?"). If your type has those properties, use value receivers everywhere and hand instances around by value.

### Constructor convention

- **Return `*T`** when the type is stateful and mutable (`*bytes.Buffer`, `*http.Request`, `*sql.DB`). Callers expect to hold one and pass it around.
- **Return `T`** when the type is a small immutable value (`time.Now()` returns `time.Time`, `netip.MustParseAddr("1.2.3.4")` returns `netip.Addr`). Callers expect to compose with it, not mutate.
- If you're not sure: prefer `*T`. Mutability is more common in real applications than immutability.

### The decision table

| Purpose of the type / method | Receiver |
|---|---|
| Method mutates the receiver | **Pointer** |
| Type contains a `sync.Mutex`, `sync.WaitGroup`, or a channel | **Pointer** (copying breaks the invariant) |
| Type is large (rule of thumb: > ~100 bytes) | **Pointer** (avoid the copy) |
| Type is designed to be immutable (like `time.Time`) | **Value** |
| Type is a wrapper around a slice/map/channel and never mutates | Value works, pointer works — pick consistency |
| Type is a small named alias (`type UserID string`) | **Value** — no reason to pointer-ify a string |
| Method needs to satisfy an interface, and other methods on the type use pointer receivers | **Pointer** (see the asymmetry above) |
| You want the type usable as a map key or in equality checks | **Value** — and every field must be comparable (see lesson 27) |

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
5. Create a `Bag` struct with a `Label string` and `Items []string`. Copy it into `bag2`. Change `bag2.Label` and `bag2.Items[0]`. Print both. Predict what happens before you run.
6. Store `Counter` (value type) in a `map[string]Counter`. Try to call `.Increment()` on `m["key"]`. Read the compile error. Now switch to `map[string]*Counter`.
7. Define an interface with two methods. Give one method a value receiver, one a pointer receiver, on the same type. Assign a value and a pointer to a variable of the interface type. Which compiles?

## Common pitfalls

- **Nil pointer dereference.** The most common Go panic. Always check for nil before dereferencing a pointer that could be nil (function returns, struct fields, map values that are `*T`).
- **Pointer to a range-loop variable.** Pre-Go-1.22, `for _, u := range users { save(&u) }` gave you the SAME pointer every iteration. Go 1.22+ fixed this. But keep it in mind when reading older code.
- **Mixing value and pointer receivers.** If half your methods take `T` and half take `*T`, only `*T` satisfies interfaces that require both — the value type doesn't. This is one of the top-3 interface bugs.
- **Returning a pointer to a stack-local**. In Go, safe! The compiler does escape analysis and moves the variable to the heap. In C this would be a bug; in Go it just works.
- **Comparing pointers.** `p1 == p2` compares *addresses*, not the pointed-at values. If you want value equality, dereference: `*p1 == *p2`.
- **Thinking "copying a struct is expensive."** For most structs (under ~100 bytes) copying is essentially free. Don't reach for pointers just to avoid a copy — use them for *mutation* or for types that *shouldn't* be copied (mutexes, etc). Optimizing away a struct copy is nearly always premature.
- **Assuming value copy is deep.** A struct's slice/map/channel/func fields share their backing storage with the original. Copy is shallow. If you need a true deep copy, `slices.Clone` the slices and re-create the maps yourself.
- **Over-pointering reference types.** `*[]int` is almost never what you want — `[]int` is already reference-shaped. Same for maps, channels, funcs. Read `*[]int` in a signature as a hint that the callee will reassign the slice header itself, which is rare.
- **Calling a pointer-receiver method on a non-addressable value.** Map values, function returns, and interface-boxed values aren't addressable. Extract to a local variable, or return a pointer from the constructor.

## You've understood this lesson when...

- You can predict whether a function's changes to its argument will be visible in the caller, just by looking at the signature.
- You can name three cases where a method should use a pointer receiver AND one case where you should prefer a value receiver.
- You can explain why `var p *int; fmt.Println(*p)` panics but `var m map[string]int; fmt.Println(m["x"])` doesn't.
- You can convert a value-receiver-only type into a pointer-receiver-only type and update the call sites.
- You can predict what will happen when you copy a struct that has both an `int` field and a `[]string` field, then mutate both on the copy.
- You can look at a signature `NewX() X` vs `NewX() *X` and infer whether `X` is meant to be treated as an immutable value or a mutable object.
- You know why `m["key"].SomeMutatingMethod()` doesn't compile when `m` is `map[string]Counter`, and what the two fixes are.

## Next

- **Next lesson (recommended):** [24-error-handling](../24-error-handling/) — Go's `if err != nil` idiom, the pattern that defines idiomatic Go code.
