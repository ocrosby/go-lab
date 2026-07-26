# Hello

Your first Go program. Ten lines of code that print `Hello World!` to your screen.

## Why it matters

Every Go program has the same skeleton — a `package` line, an `import` block, and a `main` function. Once you understand the seven lines of `hello.go`, you can read the first-half of *any* Go program you'll ever see. This is your foundation.

## Prerequisites

None. If `go version` prints a number, you're ready.

## Run it

```bash
go run ./lessons/01-hello
```

Expected output:

```text
Hello World!
```

If you got that, you just ran your first Go program. Everything from here builds on this.

## What's in this folder

| File | What it does |
|---|---|
| `hello.go` | The whole program — ten lines. |

## Anatomy of the file

Open [`hello.go`](./hello.go) and match each line to what it does:

```go
package main            // "This file is part of the 'main' package — a runnable program."

import (
    "fmt"               // "I need the standard-library 'fmt' package for printing."
)

func main() {           // "This is where the program starts."
    fmt.Println("Hello World!")   // Print the message and a newline.
}
```

Three ideas to notice:

1. **`package main` means runnable.** Only files in a package named `main` can be run as programs. Everything else (`package fmt`, `package strings`, packages you'll write yourself) is a library that gets imported by something else.
2. **`func main()` is the entry point.** When you run the program, Go calls this function first. When it returns, the program ends.
3. **You have to import what you use.** `fmt` isn't automatic — you asked for it. Try deleting the `import` line and see what happens.

## Try it yourself

1. Change `"Hello World!"` to your own name. Re-run — the output should match.
2. Add a second `fmt.Println("...")` line inside `main`. Re-run — both lines should print.
3. Delete the `import "fmt"` line. Run again. What error does Go give you? (This is your first Go compiler error — a useful one to see.)
4. Rename `main` to `Main` (capital M). Run again. What changes?

## Common pitfalls

- **Running from the wrong directory.** All commands in this repo are meant to be run from the *root* of the repo (the folder containing `go.mod`). If you `cd` into `lessons/01-hello/` and run `go run hello.go`, that also works — Go finds its own way. If you're getting "cannot find package" errors, check where you're standing.
- **Editing the file but forgetting to save.** Go reads the file from disk each time you run — an unsaved change won't take effect. Save first.
- **Typos in `func main()` or `package main`.** These strings are magic — Go looks for them by name. `func Main()` is not the same thing.

## You've understood this lesson when...

- You can name the three parts of a minimal Go program (package clause, imports, `main` function).
- You can explain in one sentence what `package main` means versus, say, `package fmt`.
- You can predict what the program will do if you change the string passed to `fmt.Println`.

## Next

- **Next lesson:** [02-functions-and-packages](../02-functions-and-packages/) — how to break code into functions and packages of your own.
