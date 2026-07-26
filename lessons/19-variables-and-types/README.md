# Variables and types

The three ways to declare a variable in Go, the built-in types you'll use every day, zero values, constants, and `iota` (Go's answer to enums).

> **Recommended before lesson 02.** Lesson 02 introduces functions like `func Add(x, y int) int` without ever explaining what `int` is or why the type comes *after* the name. This lesson fills that gap.

## Why it matters

Go is a statically-typed language: every variable has a type known at compile time. That gives you speed and safety, but it means you have to learn the type system before anything else makes sense. Once you can name the six primitive types, the three ways to declare a variable, and the concept of a *zero value*, every other lesson in the syllabus reads more clearly.

## Prerequisites

- Lesson 01: how to run a Go program.

## Run it

```bash
go run ./lessons/19-variables-and-types
```

Expected output:

```text
--- primitive zero values ---
int:     0
float64: 0
bool:    false
string:  ""

--- declared with := ---
name: Ada, age: 36, weight: 62.3, licensed: true

--- constants ---
Pi ≈ 3.14159
--- iota ---
Sunday=0 Monday=1 Tuesday=2

--- type conversion ---
i=42, f=42, s=42
```

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`main.go`](./main.go) | The three declaration forms, zero values, `const`, `iota`, and explicit type conversions — all in one runnable program. |

## The three declaration forms

Go gives you three ways to declare a variable, and they overlap. Pick the shortest one that works.

```go
var x int = 42     // Full form — explicit type + value
var x = 42         // Type inferred from the value; still uses var
x := 42            // Short form — only inside a function
```

The `:=` short form is by far the most common. Use `var` at package level (outside any function) or when you want to declare-without-initializing (`var x int` — x is zero).

## Zero values

Go initializes every declared variable. There is no "uninitialized" garbage value. The default for each type:

| Type | Zero value |
|---|---|
| `int`, `int8`..`int64`, `uint`, etc. | `0` |
| `float32`, `float64` | `0` |
| `bool` | `false` |
| `string` | `""` (empty string, not `nil`) |
| pointer, slice, map, channel, function, interface | `nil` |
| struct | struct with every field at its own zero value |

`var x int` is a valid, useful declaration — you get `x == 0` guaranteed.

## The Go type cheat sheet

**Numbers:**
- `int`, `uint` — platform-sized integer (32 or 64 bit depending on OS). Use these by default.
- `int8`, `int16`, `int32`, `int64` — explicit sizes. `int32` is aliased as `rune`, used for Unicode code points.
- `uint8` is aliased as `byte`, used for raw bytes.
- `float32`, `float64` — IEEE-754. Prefer `float64` unless you're memory-constrained.
- `complex64`, `complex128` — rarely used.

**Text and booleans:**
- `string` — immutable sequence of bytes. Indexing gives you a `byte`, not a character (see the runes-vs-bytes pitfall below).
- `bool` — `true` or `false`. No implicit int-to-bool conversion; `if 1 { ... }` does not compile.

**Compound types (covered in later lessons):**
- `[]T` — slice of T (lesson 21).
- `map[K]V` — map of K to V (lesson 21).
- `*T` — pointer to T (lesson 22).
- `struct { ... }` — record (lesson 11).

## Constants and `iota`

Constants are values known at compile time. Declared with `const`.

```go
const Pi = 3.14159
const Greeting = "hello"

const (
    Small  = 1
    Medium = 2
    Large  = 3
)
```

For sequential integer constants — Go's answer to enums — use `iota`. Inside a `const` block, `iota` starts at 0 and increments by one per line:

```go
const (
    Sunday    = iota  // 0
    Monday            // 1
    Tuesday           // 2
    Wednesday         // 3
    Thursday          // 4
    Friday            // 5
    Saturday          // 6
)
```

`iota` also supports expressions — `1 << iota` for bit flags, `iota * 2` for even numbers, etc.

## Type conversions are explicit

Go does not automatically convert between numeric types. Even `int` and `int64` need an explicit conversion:

```go
var i int32 = 42
var j int64 = int64(i)   // explicit — this compiles

var k int64 = i          // does NOT compile: type mismatch
```

Same for numbers and strings — `int(42)` is not `"42"`. Use `strconv.Itoa`/`strconv.Atoi` from the standard library (see [`docs/standard-library-tour.md`](../../docs/standard-library-tour.md#strconv--string-conversions)).

## Runes, bytes, and strings — the Unicode gotcha

A Go `string` is an **immutable sequence of bytes**, not characters. This distinction bites beginners the first time they index into a string containing anything beyond ASCII.

```go
s := "héllo"       // 6 bytes on disk (é is 2 bytes in UTF-8)
fmt.Println(len(s))   // 6 — bytes, not characters
fmt.Println(s[1])     // 195 — first byte of é, not 'é'
```

Three related types:

| Type | What it is | Written as |
|---|---|---|
| `byte` | Alias for `uint8`. Raw byte. | `'a'` if ASCII, `0x41` in general |
| `rune` | Alias for `int32`. One Unicode code point. | `'é'`, `'🚀'` |
| `string` | Immutable UTF-8 byte sequence. | `"hello"` |

To iterate over characters (code points), not bytes, use `for range`:

```go
s := "héllo"
for i, r := range s {
    fmt.Printf("byte %d: rune %q\n", i, r)
}
// byte 0: rune 'h'
// byte 1: rune 'é'    ← note: byte index jumps by 2 (é is 2 bytes)
// byte 3: rune 'l'
// byte 4: rune 'l'
// byte 5: rune 'o'
```

To get the count of code points (not bytes):

```go
import "unicode/utf8"

n := utf8.RuneCountInString("héllo")  // 5
```

Or convert to a `[]rune` (allocates):

```go
rs := []rune("héllo")
fmt.Println(len(rs))   // 5
fmt.Println(string(rs[1]))  // "é"
```

**Rule of thumb**: use `string` for text you're passing around unchanged. Convert to `[]rune` only when you need to index or modify individual characters. Almost every string function in the standard library (`strings.Contains`, `strings.Split`, `strings.ReplaceAll`) is Unicode-safe — you rarely need to think about this.

## Try it yourself

1. Add a line to `main.go` that declares `var count int` and prints it. What does it print? (Answer: `0` — the zero value.)
2. Try `var b bool = 1`. Does it compile? Why not? Now try `var b bool = true`.
3. Add a `Direction` type using `iota` for `North`, `East`, `South`, `West`. Print each.
4. Try assigning an `int32` to an `int64` without a conversion. Read the compiler error carefully — Go's error messages are educational.

## Common pitfalls

- **Runes vs bytes vs strings.** `s := "hi"; fmt.Println(s[0])` prints `104` (the byte value for `'h'`), not `"h"`. If you want characters, convert to `[]rune(s)`.
- **`:=` inside vs `var` outside.** The short form only works inside a function. At package level use `var name = value`.
- **Redeclaring with `:=`.** `x := 5; x := 6` doesn't compile — the second is a redeclaration. Use `x = 6` to reassign. `:=` requires at least one *new* variable on the left.
- **Precision loss on conversion.** `int64(3.9)` gives you `3`, silently. Convert with intent.
- **Constant expressions must be constant.** `const t = time.Now()` doesn't compile — `time.Now()` is a runtime call.

## You've understood this lesson when...

- You can list the six most-common Go primitive types from memory and name the zero value for each.
- You know when to use `var` vs `:=`.
- You can write a `const` block that uses `iota` to define a set of enum-like values.
- You can explain why `s[0]` on a string returns a number, not a character.

## Next

- **Next lesson:** [20-control-flow](../20-control-flow/) — Go's `if`, `switch`, and `for` (the only loop construct in the language).
