# Functions and closures

Anonymous functions, closures that capture variables from their enclosing scope, functions as first-class values (parameters and returns), variadic parameters, named return values, and `init()`.

> **Recommended before lesson 08 (goroutines).** Every `go func() { ... }()` in the concurrency lessons is a closure capturing outer variables. This lesson explains that idiom.

## Why it matters

Half the interesting code in the syllabus uses anonymous functions and closures without ever introducing them. `go func(id int) { ... }(i)` is opaque unless you know what "anonymous function taking `id`, called with argument `i`" means. Once the pattern clicks, every concurrency lesson reads more clearly.

## Prerequisites

- Lesson 02: functions and packages.
- Lesson 20: variables and types.

## Run it

```bash
go test -race ./lessons/28-functions-and-closures
```

Expected: 8 passes.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`funcs.go`](./funcs.go) | Anonymous funcs, closures, higher-order helpers, variadic parameters, named return values, `init()`. |
| [`funcs_test.go`](./funcs_test.go) | Tests for each pattern above. |

## Anonymous functions

Functions without a name, declared inline:

```go
func() {
    fmt.Println("I have no name")
}()  // the trailing () CALLS the function immediately
```

Store one in a variable to call later:

```go
greet := func(name string) {
    fmt.Println("hi,", name)
}
greet("Ada")  // hi, Ada
```

The syntax mirrors regular functions minus the name.

## Closures — capturing variables

A closure is a function value that references variables from **its enclosing scope**. The variables aren't copied — the closure holds a reference to them, so mutations from either side are visible:

```go
counter := 0
inc := func() {
    counter++
}
inc()
inc()
fmt.Println(counter)  // 2 — the closure mutated the outer counter
```

This is the foundation of every "generator" or "stateful callback" you'll see in Go.

## Functions as first-class values

Functions have a type: `func(int) int`, `func(context.Context) error`, etc. You can:

- **Pass a function as an argument:**
  ```go
  slices.SortFunc(users, func(a, b User) int { return cmp.Compare(a.Age, b.Age) })
  ```
- **Return a function from a function** (also a closure — captures state):
  ```go
  func makeCounter() func() int {
      n := 0
      return func() int {
          n++
          return n
      }
  }
  c := makeCounter()
  fmt.Println(c(), c(), c())  // 1 2 3
  ```
- **Store functions in maps or slices** to build a dispatch table:
  ```go
  ops := map[string]func(int, int) int{
      "+": func(a, b int) int { return a + b },
      "-": func(a, b int) int { return a - b },
  }
  ```

## Variadic parameters (`...T`)

Accept any number of arguments of the same type:

```go
func Sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

Sum(1, 2, 3)           // total = 6
Sum()                   // total = 0
Sum([]int{1, 2, 3}...)  // spread a slice into individual args
```

`fmt.Println(args ...any)` is the variadic pattern you use every day.

## Named return values

Function returns can be named. Two effects:

1. **Documentation** — the reader sees what each return means.
2. **Automatic zero-value pre-declaration** — the named returns exist as local variables, so `return` with no args returns them.

```go
func divide(a, b int) (quotient int, err error) {
    if b == 0 {
        err = errors.New("divide by zero")
        return  // "naked return" — returns quotient=0, err=... as declared
    }
    quotient = a / b
    return
}
```

Use named returns sparingly. They shine for small helper functions (especially with a deferred wrap — see lesson 25) but hurt readability in long ones where you can lose track of what got assigned.

## `init()` — the package-init hook

Every package can have one or more `init()` functions. Go calls each `init()` **once**, at program startup, before `main()` runs:

```go
package config

var Settings map[string]string

func init() {
    Settings = loadFromEnv()
}
```

Rules:

- Each file can have any number of `init()` functions.
- They run in the order they're declared in a file, and package-init runs after all variable initializers in that package.
- Between packages, dependencies init first — the standard library initializes before your `main` package.
- **`init()` takes no arguments and returns nothing.**

Use `init()` for registering plugins with a registry, priming caches from disk, checking required environment variables. Avoid it for anything the caller might want to skip or reorder — it runs before `main` and can't be turned off.

## Try it yourself

1. Write a `makeAdder(n int) func(int) int` — returns a closure that adds `n` to whatever it's given. Verify each returned adder has its own `n`.
2. Fix the classic pre-Go-1.22 loop-variable-capture bug: `for i := 0; i < 3; i++ { go func() { fmt.Println(i) }() }`. This code on Go 1.21 or earlier prints 3 three times; on Go 1.22+ prints 0, 1, 2 (in some order). Try it — verify Go 1.26 gives you the fixed behaviour.
3. Write a `NewCounter() (inc func(), value func() int)` that returns two closures sharing the same private counter. What happens if you call `NewCounter()` twice — do the two counters share state?
4. Add a variadic `Concat(sep string, parts ...string) string` and call it both with individual args and with a slice-spread.
5. Add an `init()` to your file that prints "package initialized". Where in the program output does it appear?

## Common pitfalls

- **Loop-variable capture (Go ≤ 1.21).** Old code: `for i := 0; i < N; i++ { go func() { fmt.Println(i) }() }` — every goroutine printed the SAME `i` (the final value). Go 1.22+ fixed this by making the loop variable fresh per iteration. This repo pins Go 1.26 so you're safe, but the pattern still appears in older codebases.
- **Naked returns in long functions.** `return` in a 60-line function relies on the reader remembering what each named return got set to. Prefer explicit returns except in very short functions.
- **`init()` that panics.** A panic in `init` crashes the program before `main` runs, which is often surprising. If startup can fail, check for the error condition in `main()` where the caller can handle it.
- **`init()` order across files.** Files in the same package can have their init order shuffled between builds. Do not rely on it.
- **Capturing a pointer to a range-loop element:**
  ```go
  var ptrs []*User
  for _, u := range users {
      ptrs = append(ptrs, &u)  // WRONG (pre-1.22) — all same address
  }
  ```
  Go 1.22+ fixes this too. In older code, take `&users[i]` instead.

## You've understood this lesson when...

- You can write an anonymous function that immediately runs (an IIFE).
- You can write `makeCounter()` from scratch.
- You know the difference between a closure that captures a variable and a function that just takes it as an argument.
- You can spot the loop-variable-capture bug in pre-Go-1.22 code.
- You can name a legitimate use of `init()` and one abuse to avoid.

## Next

- **Next lesson (recommended):** [06-composition](../06-composition/) — you now have every language mechanic you need for the main syllabus. Structs, methods, and interfaces open the door to composition, polymorphism, and everything after.
