# Go Build Directives

Go source files can carry special comments that instruct the compiler, linker, and build tooling. These are called **compiler directives** (also "pragmas" in older documentation and "build directives" colloquially). They look like ordinary comments but are recognized by the toolchain because of a strict syntax: they begin with `//go:` with **no space** between `//` and `go:`.

This document covers what they are, why they exist, how the toolchain recognizes them, and provides a detailed catalog of every directive available in modern Go (as of Go 1.24).

## Table of contents

- [What they are](#what-they-are)
- [Why they exist](#why-they-exist)
- [How they work](#how-they-work)
- [Build constraints](#build-constraints)
- [Code generation and embedding](#code-generation-and-embedding)
- [Compiler directives](#compiler-directives)
- [Runtime and linker directives](#runtime-and-linker-directives)
- [CGo directives](#cgo-directives)
- [WebAssembly directives](#webassembly-directives)
- [Runtime debug directives](#runtime-debug-directives)
- [The `//line` directive](#the-line-directive)
- [Common pitfalls](#common-pitfalls)
- [References](#references)

## What they are

A Go build directive is a specially formatted single-line comment that carries an instruction to some part of the Go toolchain. Unlike documentation comments — which are ignored by the compiler but consumed by `godoc` — directives change how the code is compiled, linked, generated, or included in a build.

Three properties distinguish a directive from a regular comment:

1. **No space after `//`.** A directive is `//go:noinline`, never `// go:noinline`. The space is what tells the compiler "this is a comment for humans" rather than "this is a directive for me."
2. **Recognized prefix.** The compiler understands a fixed set of prefixes (`//go:`, `//line`, `// +build`). Other prefixes are treated as comments.
3. **Positional rules.** Most directives must sit immediately above the declaration they affect — a blank line or another comment between the directive and its target usually breaks the association.

## Why they exist

Go deliberately keeps the language surface small. Directives are the escape hatch for concerns that do not belong in the language itself:

- **Build selection** — a file should be compiled only on Linux, or only when a race detector is enabled, or only when a specific tag is passed.
- **Code generation** — a tool needs to be re-run whenever the source changes.
- **Static assets** — a file's contents should be embedded into the compiled binary.
- **Low-level runtime work** — a function must not be inlined, or must not call anything that allocates on the heap, because it is part of the scheduler or garbage collector.
- **FFI and interop** — cgo needs to describe how Go and C symbols connect.

By putting these concerns in structured comments, Go keeps the language grammar clean while still exposing the knobs that systems programming demands.

## How they work

### Syntax rules

- Directives are always **`//`-style** line comments. Block comments (`/* ... */`) are never recognized as directives.
- No space between `//` and the directive name. `//go:build` works; `// go:build` is silently ignored.
- One directive per line.
- The directive must appear in a position appropriate to its category (see each section below).

### Placement rules

- **Build constraints** appear near the top of the file, before the `package` clause, followed by a blank line.
- **Function-level directives** (like `//go:noinline`) appear on the line immediately above the function declaration — no blank line between them.
- **`//go:generate`** may appear anywhere in the file; `go generate` collects them in source order.
- **`//go:embed`** appears immediately above a package-level `var` declaration of type `string`, `[]byte`, or `embed.FS`.

### Requesting recognition

Most directives require the file to import `unsafe` — this is Go's convention that the file is doing something that steps outside normal safety guarantees. The compiler refuses to honor several directives (notably `//go:linkname`) unless `import "unsafe"` is present.

## Build constraints

Build constraints control whether a file is included in a particular build. They are the most commonly encountered directive.

### `//go:build`

Introduced in Go 1.17, `//go:build` is the modern, expression-based build constraint. It replaces the older `// +build` form.

```go
//go:build linux && amd64

package foo
```

The expression supports `&&`, `||`, `!`, and parentheses. Common tags:

- **GOOS**: `linux`, `darwin`, `windows`, `freebsd`, `netbsd`, `openbsd`, `dragonfly`, `solaris`, `android`, `ios`, `plan9`, `js`, `wasip1`
- **GOARCH**: `amd64`, `arm64`, `386`, `arm`, `mips`, `mips64`, `ppc64`, `riscv64`, `s390x`, `wasm`
- **Compiler**: `gc`, `gccgo`
- **Umbrella tags** (Go 1.19+): `unix` (matches all Unix-like GOOS values)
- **Go version**: `go1.18`, `go1.21`, etc. (matches this version *and later*)
- **Race detector**: `race`
- **CGo enabled**: `cgo`
- **User-defined tags**: any identifier passed via `go build -tags=<name>`

The constraint must be followed by a blank line before the `package` clause, or `gofmt` will treat it as a package doc comment.

```go
//go:build (linux || darwin) && !race

package cache
```

### `// +build` (legacy)

Before Go 1.17, build constraints used space-separated OR, comma-separated AND, and `!` for NOT:

```go
// +build linux darwin
// +build !race

package cache
```

`gofmt` in Go 1.17+ automatically inserts a matching `//go:build` line above any `// +build` line, so both forms often appear together in older code. New code should use `//go:build` exclusively; the `// +build` form is preserved only for backward compatibility with Go 1.16 and earlier.

### Filename-based constraints

The build system also treats certain filename suffixes as implicit constraints. A file named `foo_linux.go` is compiled only on Linux; `foo_amd64.go` only on amd64; `foo_linux_amd64.go` only on the intersection of both.

Separately, any file whose name ends in `_test.go` is compiled only under `go test`. This is an independent rule — a file named `handler_linux_test.go` is subject to *both* constraints (Linux, test-only).

### The `ignore` convention

There is no built-in `ignore` tag, but `//go:build ignore` is a widely-used convention for scratch files and standalone examples that live inside a package but should not participate in its build:

```go
//go:build ignore

package main

// A standalone example. Run with: go run this_file.go
```

Because nothing sets the `ignore` tag, the file is excluded from every normal build.

## Code generation and embedding

### `//go:generate`

Runs an arbitrary command when `go generate` is invoked. The command is not run by `go build` or `go test`.

```go
//go:generate stringer -type=Color

type Color int

const (
    Red Color = iota
    Green
    Blue
)
```

`go generate ./...` walks the module, finds every `//go:generate` line, and runs the commands in source order in the file's directory. Common uses: `stringer`, `mockgen`, `protoc`, `sqlc`, and any project-specific code generator.

The command sees these environment variables: `$GOFILE`, `$GOLINE`, `$GOPACKAGE`, `$GOARCH`, `$GOOS`, `$GOROOT`, `$DOLLAR` (a literal `$`).

### `//go:embed`

Introduced in Go 1.16. Embeds files or directories from the module into the compiled binary.

```go
import _ "embed"

//go:embed version.txt
var version string

//go:embed schema.sql
var schema []byte

//go:embed static/*
var staticFS embed.FS
```

Rules:

- The `embed` package must be imported (blank import if only the directive is needed).
- The `var` immediately below the directive must be `string`, `[]byte`, or `embed.FS`.
- Patterns are relative to the source file's directory and cannot use `..` to escape it.
- Multiple patterns can appear on one line: `//go:embed a.txt b.txt`.
- Files starting with `_` or `.` are excluded unless matched explicitly by pattern.

## Compiler directives

These directives change how the Go compiler translates a function or type. Most are intended for the standard library and runtime; using them in application code should be a conscious decision.

### `//go:noinline`

Prevents the compiler from inlining a function. Applied to the function declaration immediately below.

```go
//go:noinline
func slowPath(x int) int {
    return x * 2
}
```

Useful when benchmarking (to observe true call overhead), when a function's stack frame must exist for a stack trace, or when a runtime-support routine must remain a distinct symbol.

### `//go:noescape`

Declares that the arguments passed to a function do not escape to the heap. The compiler trusts this without proof — it is the caller's contract to honor it. Only valid on function *declarations* without bodies (typically assembly implementations).

```go
//go:noescape
func indexByte(b []byte, c byte) int
```

Wrong use causes memory corruption. Reserved for stdlib and hand-written assembly.

### `//go:nosplit`

Skips the usual stack-overflow check on function entry. The compiler still allocates stack space, but it does not check whether the goroutine has enough stack before allocating.

```go
//go:nosplit
func onSystemStack() {
    // ...
}
```

Used inside the runtime for functions that must run on a tiny fixed stack (e.g. signal handlers, the scheduler). Application code should never need this.

### `//go:norace`

Disables race-detector instrumentation for the function.

```go
//go:norace
func fastPath() {
    // ...
}
```

Reserved for functions that manipulate shared memory in ways the race detector cannot model (typically inside the runtime).

### `//go:nocheckptr`

Disables the `-d=checkptr` runtime check that detects invalid `unsafe.Pointer` conversions.

```go
//go:nocheckptr
func packHeaders(p unsafe.Pointer) uintptr {
    // Doing pointer arithmetic that checkptr would falsely flag.
}
```

Used when a function performs a valid but unusual pointer operation that the checker cannot prove safe.

### `//go:linkname`

Links a local Go symbol name to a symbol defined in another package, potentially bypassing export rules. Requires `import "unsafe"`.

```go
import _ "unsafe"

//go:linkname nanotime runtime.nanotime
func nanotime() int64
```

This is the mechanism that lets `sync`, `time`, and `net` reach into `runtime` internals. It is powerful and fragile: the target symbol is not part of any stable API and can change between Go versions. Go 1.23 tightened the rules — packages outside the standard library can only linkname a symbol whose defining package explicitly opts in. See the [Go 1.23 release notes](https://go.dev/doc/go1.23#linker) for the exact conditions and migration path.

### `//go:uintptrescapes`

Marks all `uintptr` arguments to the function as escaping. Necessary for syscall wrappers where a `uintptr` is used to hold what is really a pointer, and the garbage collector must keep the backing object alive across the call.

```go
//go:uintptrescapes
func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
```

### `//go:notinheap` (restricted)

Marks a type as not allocatable from the garbage-collected heap. Used inside the runtime for structures that must live in special memory regions. Not intended for general use.

### `//go:systemstack`

The annotated function must run on the system (g0) stack rather than a goroutine's stack. Runtime-only.

### `//go:nowritebarrier`

The compiler emits an error if the annotated function contains a write barrier. Used in runtime code that runs when the GC's write barrier is not in a valid state.

### `//go:nowritebarrierrec`

Like `//go:nowritebarrier`, but applies transitively — the function and everything it calls must be write-barrier-free.

### `//go:yeswritebarrierrec`

The companion to `//go:nowritebarrierrec`: marks a function as an acceptable stopping point for the recursive check.

### `//go:registerparams` (experimental / historical)

Marked the function as using the register-based calling convention during Go 1.17's ABI transition. No longer needed in modern Go; the register ABI is the default.

## Runtime and linker directives

### `//go:fix inline` (Go 1.24+)

Marks a function or constant as suggested for inlining by the `go fix` tool's inliner. When users run `go fix -inline`, the tool replaces call sites of the marked function with its body, and replaces uses of the marked constant with its value. Used to shepherd deprecations gently: keep the old symbol working, but let downstream code migrate automatically.

```go
//go:fix inline
func OldName(x int) int { return NewName(x) }
```

## CGo directives

CGo directives are emitted by the `cgo` tool into generated Go files. Hand-writing them is rare but possible for advanced FFI work. They control how Go symbols are exposed to C and how C symbols are pulled into Go.

### `//go:cgo_export_static <local> <exported>`

Exports a Go function to C using a static symbol name (visible at link time only).

### `//go:cgo_export_dynamic <local> <exported>`

Exports a Go function to C using a dynamic symbol (visible via the dynamic linker at runtime).

### `//go:cgo_import_static <name>`

Imports a C symbol by static linkage.

### `//go:cgo_import_dynamic <local> [<remote> ["<dynamic_linker>"]]`

Imports a C symbol dynamically. Used to weakly reference symbols from shared libraries.

### `//go:cgo_ldflag "<flag>"`

Passes a flag to the host linker when linking the final binary. Emitted by `cgo` from `#cgo LDFLAGS:` lines in Go source.

### `//go:cgo_dynamic_linker "<path>"`

Specifies the dynamic linker to be recorded in the output binary.

### `//go:cgo_unsafe_args`

Declares that a cgo-generated function passes arguments containing unsafe pointers.

## WebAssembly directives

### `//go:wasmimport <module> <name>` (Go 1.21+)

Declares that the annotated Go function is implemented by a host import in a WebAssembly module. Only valid on function declarations without bodies, and only when building for `GOOS=wasip1` or a WASM target.

```go
//go:wasmimport wasi_snapshot_preview1 fd_write
func fd_write(fd int32, iovs unsafe.Pointer, iovsLen int32, nwritten unsafe.Pointer) int32
```

### `//go:wasmexport <name>` (Go 1.22+)

Exports a Go function so it can be called by the WebAssembly host. Only meaningful when building a WASM library (`-buildmode=c-shared` on a wasm target).

```go
//go:wasmexport add
func Add(a, b int32) int32 { return a + b }
```

## Runtime debug directives

### `//go:debug <name>=<value>` (Go 1.21+)

Sets the default value of a `GODEBUG` setting for a program. Only allowed in `main` packages; must appear in a file that contains the `main` package clause, before that clause, in the doc comment or in `//go:debug` lines above it.

```go
//go:debug gotypesalias=1

package main
```

Users can still override the setting with the `GODEBUG` environment variable; the directive only changes the default when the environment does not specify a value. This is the mechanism the standard library uses to opt a program into or out of backward-incompatible behavior changes.

## The `//line` directive

Not a `//go:` directive, but part of the same family. `//line` overrides the file and line number the compiler reports for subsequent code.

```go
//line generated.y:42
func foo() { /* ... */ }
```

Used exclusively by code generators (yacc, protoc, sqlc, etc.) so that compiler errors point at the source file the human wrote rather than the generated intermediary. Two accepted forms:

- `//line file:line`
- `//line file:line:col`

Applies from the line *after* the directive.

## Common pitfalls

**Space after `//` silently disables the directive.** `// go:noinline` is a comment; `//go:noinline` is a directive. The compiler emits no warning either way.

**Blank line between directive and target breaks the binding.** For function-level directives, put the directive on the line immediately above the `func` keyword — no blank line, no intervening comment.

**`//go:build` without a following blank line becomes a package doc comment.** Always leave one blank line between the build constraint and the `package` clause.

**Filename constraints combine with `//go:build`.** A file named `foo_linux.go` with `//go:build darwin` is never built — the intersection is empty. This is a common cause of "why isn't this file being compiled" confusion.

**`//go:linkname` without `import "unsafe"` is rejected.** The compiler will complain about a missing unsafe import even though the file itself has no unsafe operations.

**`//go:embed` on an unexported var still works, but the imported files ship in the binary regardless of visibility.** Do not embed secrets — anyone with the binary can extract them.

**Legacy `// +build` constraints do not use `&&` or `||`.** Space is OR, comma is AND. Mixing modern and legacy syntax silently produces the wrong constraint.

## References

- [Official compiler directive documentation](https://pkg.go.dev/cmd/compile#hdr-Compiler_Directives) — the authoritative list of `//go:` directives with their exact semantics.
- [Build constraints reference](https://pkg.go.dev/cmd/go#hdr-Build_constraints) — full syntax and evaluation rules for `//go:build`.
- [`go generate` design and rationale](https://go.dev/blog/generate) — the Go Blog post that introduced `//go:generate`.
- [`embed` package documentation](https://pkg.go.dev/embed) — patterns, variable types, and directory-embedding rules for `//go:embed`.
- [GODEBUG documentation](https://go.dev/doc/godebug) — the runtime debug settings framework that `//go:debug` configures.
- [Go 1.23 release notes on linkname](https://go.dev/doc/go1.23#linker) — the tightened rules for cross-package `//go:linkname`.
