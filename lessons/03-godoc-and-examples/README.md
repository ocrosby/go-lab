# GoDoc and examples

How Go teaches you to document code: doc comments that start with the identifier's name, and `Example*` functions that are executable *and* rendered into the docs.

## Why it matters

Go has one documentation convention and one tool that reads it. Every public library on [pkg.go.dev](https://pkg.go.dev) — including the standard library — is generated from the same rules you're about to learn. The pay-off is real: your own packages get browsable docs for free, and your example code cannot silently drift from the implementation because `go test` runs the examples and diffs their output.

Two things make this a great early lesson:

1. **The habit sets on day one.** Once you're writing exported names without doc comments, it's hard to go back and add them later. Learning the convention right after you learn what "exported" means (previous lesson) is the moment to install the habit.
2. **`Example*` functions are your first test *and* your first doc.** Before you meet `t.Run` or table-driven tests (next lesson), you can already write a self-verifying usage snippet.

## Prerequisites

- Lesson 02: functions and packages — you need to know what an exported name is.

## Run it

Run the examples like any other tests:

```bash
go test ./lessons/03-godoc-and-examples
```

Expected output:

```
ok  	github.com/ocrosby/go-lab/lessons/03-godoc-and-examples	0.187s
```

Then read the rendered docs from the terminal — no browser, no server:

```bash
go doc ./lessons/03-godoc-and-examples
go doc ./lessons/03-godoc-and-examples Rectangle
go doc -all ./lessons/03-godoc-and-examples
```

The `-all` form is the whole-package view, with every symbol, every doc comment, and every example laid out the way pkg.go.dev would render it.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`doc.go`](./doc.go) | Package-level doc comment (`// Package shapes …`). The `doc.go` file is the idiomatic home when a package has more than one source file. |
| [`rectangle.go`](./rectangle.go) | Doc comments on a type, a sentinel error, a constructor, and two methods — each starting with the identifier's name. |
| [`example_test.go`](./example_test.go) | Three `Example*` functions demonstrating the three godoc-recognized name shapes: `ExampleFunc`, `ExampleFunc_variant`, and `ExampleType_Method`. |

## The rules, in one place

Every rule below is enforced by convention, not by the compiler — but tooling assumes them. Skip them and your godoc renders wrong.

1. **Every exported identifier has a doc comment.** Types, functions, methods, variables, constants. The comment sits immediately above the declaration with no blank line between them.
2. **The comment starts with the identifier's name.** `// Rectangle is an axis-aligned rectangle …`, not `// This type represents …`. godoc uses the first sentence as the summary.
3. **Complete sentences, present tense, third person.** `// Area returns the rectangle's area.`
4. **Errors and concurrency go in the comment.** If a function returns a specific sentinel error, say so. If a type is safe for concurrent use (or isn't), say that too.
5. **Package doc lives in `doc.go`.** For a multi-file package, put `// Package name …` in a dedicated file with no code below it. For a one-file package, put it above the `package` line in that file.
6. **`Example*` functions document behavior and verify it.** Naming shapes godoc recognizes:
   - `ExampleFoo` — example for function/type `Foo`
   - `ExampleFoo_bar` — a second example for `Foo`, disambiguated with a lowercase suffix
   - `ExampleType_Method` — example for a method (note the uppercase method name)
   - Add `// Output:` at the bottom listing the expected `fmt.Print*` output. `go test` diffs it.

## Try it yourself

1. **Break an example.** Change the `// Output: 30` in `ExampleRectangle_Area` to `// Output: 42` and run `go test ./lessons/03-godoc-and-examples`. Note that the failure message shows the diff, exactly like a regular test.
2. **Add a new example.** Create `ExampleRectangle_Perimeter` in `example_test.go`. Give it an `// Output:` line and see it pass.
3. **Add a `Square` type** in a new `square.go` with a doc comment and a `NewSquare` constructor that reuses `NewRectangle` internally. Run `go doc -all ./lessons/03-godoc-and-examples` and confirm `Square` shows up in the package overview.
4. **Peek at the standard library.** From any shell:

   ```bash
   go doc strings.Builder
   go doc -all bytes.Buffer
   ```

   Everything you're reading is the same convention applied at scale.

## Common pitfalls

- **Blank line between the doc comment and the declaration.** godoc treats the comment as a floating comment and drops it. Keep them adjacent.
- **Comment doesn't start with the identifier's name.** `go doc` still shows the comment, but the auto-generated summary looks wrong on pkg.go.dev.
- **Missing `// Output:` in an `Example*` function.** Without it, `go test` compiles the example but does not run it — you get compilation coverage, not behavior coverage. Add it.
- **Whitespace in `// Output:`.** The comparison is line-by-line with trailing whitespace trimmed, but leading whitespace matters. Use `// Unordered output:` if the order isn't deterministic (e.g., map iteration).
- **Documenting *what* instead of *why*.** `// Add returns a + b` on an `Add` function is filler. Prefer `// Add returns the sum; overflow wraps per Go's untyped-arithmetic rules.` when there's a subtlety worth naming. If there isn't, keep it one line.

## You've understood this lesson when...

- You can write a doc comment for a new exported function without looking it up.
- You can predict what `go doc pkgname.Symbol` will print for a symbol you just wrote.
- You can write an `Example*` function with an `// Output:` block and know it will be run by `go test`.
- You reach for a doc comment reflexively whenever you write `func Name(...)` starting with an uppercase letter.

## Next

- Next lesson: [04-testing-basics](../04-testing-basics/) — `go test`, `t.Run`, table-driven tests. The `Example*` functions you wrote here are one shape of test; that lesson covers the rest.
- Related reference: the Godoc conventions section of Go's own [Effective Go](https://go.dev/doc/effective_go#commentary) and [Go Doc Comments](https://tip.golang.org/doc/comment).
