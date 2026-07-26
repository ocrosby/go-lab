# Generics

Type parameters, constraints, and the `slices` / `maps` standard-library packages that Go 1.18+ made possible. The last major language feature to learn — and the reason you no longer have to write `IntMax`, `Float64Max`, `StringMax` separately.

> **Recommended after lesson 06 (interfaces).** Constraints are interface types under the hood; understanding interfaces first makes generics land more cleanly.

## Why it matters

Go went 12 years without generics. Every "container" library either used `interface{}` (unsafe, needs type assertions) or generated code for each type. Since Go 1.18, real type-parameterized functions and types are a normal part of the language. Modern Go — the standard library, popular third-party packages — uses them heavily. If you learn Go without generics, you're learning 2020 Go on a 2026 syllabus.

## Prerequisites

- Lesson 02: functions.
- Lesson 06: interfaces (constraints are a kind of interface).
- Lesson 21: slices and maps (generic functions on them are where you'll use this most).

## Run it

```bash
go test ./lessons/25-generics
```

Expected: 6 passes.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`generics.go`](./generics.go) | `Map`/`Filter`/`Reduce`, a generic `Max`, a generic `Stack[T]` type, and use of the `comparable` and `~` constraints. |
| [`generics_test.go`](./generics_test.go) | Tests that exercise each generic function with two different type arguments. |

## Syntax in one page

Type parameters go in square brackets **before** the regular parameters. Inside the function, use the type parameter like any other type name.

```go
func Max[T int | float64 | string](a, b T) T {
    if a > b {
        return a
    }
    return b
}

Max(1, 2)         // T inferred as int
Max(1.5, 2.7)     // T inferred as float64
Max("a", "b")     // T inferred as string
Max[float64](1, 2) // explicit T if inference fails
```

`[T int | float64 | string]` is a **constraint** — the union of allowed types. `T` must be one of them.

## Constraints

The general form of a constraint is: an interface. Regular interfaces (methods) work as constraints; new "type-set" interfaces (allowed underlying types) also work.

### Union constraints

```go
type Number interface {
    int | int32 | int64 | float32 | float64
}

func Sum[T Number](xs []T) T {
    var total T
    for _, x := range xs {
        total += x
    }
    return total
}
```

### The `~` operator: "the type or any type whose underlying type is this"

Without `~`, a named type wouldn't satisfy the constraint even though it's built on the same underlying type:

```go
type Celsius float64
Sum([]Celsius{1, 2, 3}) // FAILS without ~
```

With `~`:

```go
type Number interface {
    ~int | ~float64  // now accepts Celsius, MyInt, etc.
}
```

Use `~` on numeric constraints. Omit it for exact-type matches.

### The built-in `comparable` constraint

`comparable` covers "any type that supports `==` and `!=`" — most types except slices, maps, and functions. Needed for map keys and equality checks:

```go
func Contains[T comparable](xs []T, target T) bool {
    for _, x := range xs {
        if x == target {
            return true
        }
    }
    return false
}
```

### The `any` constraint (alias for `interface{}`)

`any` means "any type at all." Useful for containers where you don't need `==` or an ordering:

```go
type Stack[T any] struct {
    items []T
}
```

## The `slices` and `maps` standard-library packages

Go 1.21 added `slices` and `maps` — generic implementations of every operation you used to write by hand:

```go
import "slices"

slices.Contains([]int{1, 2, 3}, 2)     // true
slices.Sort([]string{"b", "a", "c"})   // sorts in place
slices.Reverse(xs)
slices.Index(xs, target)
slices.Min(xs) / slices.Max(xs)        // Go 1.21+
slices.SortFunc(xs, func(a, b int) int { return a - b })
```

```go
import "maps"

maps.Keys(m)     // returns an iter.Seq[K] (Go 1.23+)
maps.Values(m)   // returns an iter.Seq[V]
maps.Clone(m)    // shallow copy
maps.Equal(a, b) // deep equality
```

Learn these first. You'll rarely need to write a generic `Map`/`Filter` yourself once you know what the standard library provides.

## Try it yourself

1. Write a generic `Contains[T comparable]([]T, T) bool`. Compare with `slices.Contains`.
2. Add a `Peek()` method to `Stack[T]` that returns the top item without popping. Add a test.
3. Try to instantiate `Stack[func()]` (function value). Does it work? Why or why not?
4. Rewrite `Max` in `generics.go` to use `cmp.Ordered` from the `cmp` standard-library package (Go 1.21+) instead of the manual union constraint. Which reads better?

## Common pitfalls

- **Over-generic types.** If you find yourself writing `func Foo[T any](x T) T` where the body doesn't actually need type parameters (or does the same thing regardless of type), you don't need generics. Just take `T` as an `any`.
- **Constraints too tight.** `func Sum[T int](xs []T)` — you locked yourself into `int`. Use `~int` (or a union) so callers with a named int type can call it too.
- **Method sets on type parameters.** You can't call arbitrary methods on `T` unless the constraint promises them. `func Foo[T any](x T) { x.String() }` — doesn't compile; `any` doesn't include `String()`.
- **Generic types with pointer receivers.** Same rules as non-generic pointer receivers (lesson 22) — but the mental picture takes a beat longer. Read carefully.
- **Cost.** Generic code is not slower at runtime in the general case — Go monomorphizes some instantiations and dictionary-dispatches others. Not a reason to avoid generics.

## You've understood this lesson when...

- You can write a generic function with a union constraint and know when to use `~`.
- You can name three functions from `slices` you'd use in a real program.
- You know why `comparable` and `any` are the two most-common constraints.
- You can define a generic container type and know which operations require which constraints.

## Next

You've now covered every language-level feature Go offers a beginner. The rest of the syllabus (concurrency, HTTP, testing frameworks, hexagonal architecture) is application of these primitives — no new language keywords ahead.

- **If you started with the fundamentals track**, cycle back to [05-composition](../05-composition/) and continue through 06, 07, ... up to 18.
- **If you're ready for a project**, jump to [20-production-api](../20-production-api/) and build.
