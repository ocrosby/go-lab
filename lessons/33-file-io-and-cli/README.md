# File I/O and CLI

Reading and writing files with `os` and `bufio`, streaming through `io.Reader` / `io.Writer`, and building a command-line tool with `flag` + `os.Args` + `os.Getenv`.

## Why it matters

Once you can read a file, write a file, and take command-line arguments, you can build every classic UNIX tool: a log filter, a config validator, a report generator, a data-migration script. Real programs also spend a lot of time with `io.Reader` / `io.Writer` — the two-method interfaces that most standard-library types satisfy, and that let you compose file/network/buffer/gzip readers freely.

## Prerequisites

- Lesson 20: types.
- Lesson 24: error handling (every I/O call returns an error).
- Lesson 25: `defer` for closing files.

## Run it

```bash
go test -race ./lessons/33-file-io-and-cli
```

Expected: 7 passes.

You can also run the CLI demo directly:

```bash
go run ./lessons/33-file-io-and-cli/cmd/wcish -lines /etc/hosts
go run ./lessons/33-file-io-and-cli/cmd/wcish -words -bytes README.md
```

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`fileio.go`](./fileio.go) | Reading and writing whole files, streaming line-by-line with `bufio.Scanner`, using `io.Reader`/`io.Writer` as function parameters. |
| [`fileio_test.go`](./fileio_test.go) | Tests using `t.TempDir()` for filesystem isolation. |
| [`cmd/wcish/main.go`](./cmd/wcish/main.go) | A tiny `wc`-style tool that counts lines, words, and bytes. Parses flags with the `flag` package. |

## Reading and writing files — the shortcuts

For small files that fit in memory:

```go
data, err := os.ReadFile("config.json")  // whole file → []byte
if err != nil {
    return err
}

err := os.WriteFile("out.txt", data, 0644)  // []byte → whole file
```

`0644` is the file mode (owner rw, group r, other r — standard for a regular file). Prefer these one-liners when the file is small.

## Streaming with `os.Open` and `bufio.Scanner`

For anything you don't want to load entirely into memory — large log files, streams from stdin, arbitrarily-long inputs — use `os.Open` + a scanner:

```go
f, err := os.Open(path)
if err != nil {
    return err
}
defer f.Close()

scanner := bufio.NewScanner(f)
for scanner.Scan() {
    line := scanner.Text()
    // ... process
}
if err := scanner.Err(); err != nil {
    return err  // scanner errors don't come from Scan() — they come from Err() at the end
}
```

Notice `defer f.Close()` right after the error check (lesson 25). `bufio.Scanner` splits on newlines by default; call `scanner.Split(bufio.ScanWords)` for word-splitting, `ScanBytes` for byte-by-byte, `ScanRunes` for Unicode code points.

## `io.Reader` and `io.Writer` — the two-method contract

Almost every I/O-shaped standard-library type satisfies one of these:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

A file is a `Reader` and a `Writer`. A `bytes.Buffer` is both. An HTTP request body is a `Reader`. Compression, encryption, and buffering all layer as `Reader` wrapping `Reader`.

Writing your functions to take `io.Reader` (instead of `*os.File` or `string`) makes them testable with any input:

```go
func CountLines(r io.Reader) (int, error) { ... }

// Production:
f, _ := os.Open(path); defer f.Close()
CountLines(f)

// Test:
CountLines(strings.NewReader("line1\nline2\n"))
```

Same function, both call sites. This is the "accept interfaces, return structs" rule from `rules/go-conventions.md`.

## `os.Args` — the raw arguments

Every Go program has `os.Args` — a slice where `[0]` is the program name and `[1:]` is everything the user typed:

```bash
$ mytool -v --config app.yaml one two
```

Gives you `os.Args = ["mytool", "-v", "--config", "app.yaml", "one", "two"]`.

Use `os.Args` directly only for tiny scripts. For anything more, use `flag`.

## The `flag` package

`flag` is the standard-library way to parse `-v`, `--name value`, `-n 10` style arguments. Minimal example:

```go
var verbose bool
var count int

flag.BoolVar(&verbose, "v", false, "verbose output")
flag.IntVar(&count, "n", 10, "number of items")

flag.Parse()

// After Parse(): verbose and count are set.
// flag.Args() gives you any positional args left over.
```

Then:

```bash
mytool -v -n 20 file1 file2
# verbose == true, count == 20, flag.Args() == ["file1", "file2"]
```

`flag` prints usage automatically if the user passes `-h` or gives bad input. For more complex CLIs (subcommands, colored output, shell completion), reach for `spf13/cobra`.

## `os.Getenv` and `os.LookupEnv`

Environment variables:

```go
port := os.Getenv("PORT")           // "" if unset — silently
port, ok := os.LookupEnv("PORT")    // ok == false if unset
```

Prefer `LookupEnv` when the difference between "not set" and "set to empty string" matters. `Getenv` is fine for optional settings with a default.

## Try it yourself

1. Modify `wcish` to also print a total when given multiple files (like the real `wc`).
2. Add a `-o outfile` flag that writes the counts to a file instead of stdout.
3. Rewrite `CountLines` in `fileio.go` to also return the number of bytes. Update the tests.
4. Write a function that reads a file and returns the top-10 most-frequent words. Use `map[string]int` (lesson 22) and `slices.SortFunc` (lesson 26).
5. Try `os.Open` on a file that doesn't exist. What error type do you get? Check with `errors.Is(err, os.ErrNotExist)`.

## Common pitfalls

- **`defer f.Close()` before the error check.** If `Open` failed, `f` is nil and `Close` panics. Always error-check first, then defer.
- **`os.ReadFile` on a huge file.** Loads everything into memory. Use `os.Open` + streaming for anything unbounded.
- **Ignoring `scanner.Err()`.** `bufio.Scanner` returns errors via `Err()` after the loop, not through `Scan()`. Missing this hides read errors.
- **`os.WriteFile` with mode 0777.** Overly permissive. Regular files should be 0644, executables 0755, secrets 0600.
- **Assuming `os.Args[0]` is the tool name.** It's whatever the user invoked — could be `/full/path/to/mytool` or `mytool`. Use `filepath.Base(os.Args[0])` for the short form.
- **`flag.Parse()` called before flags declared.** The declarations must run first. Idiomatic pattern: declare all flags at package top with `var flagFoo = flag.String(...)`, then `flag.Parse()` at the top of `main`.

## You've understood this lesson when...

- You can read a file, both the "whole thing at once" and "line by line" ways.
- You can spot the "defer before error check" bug.
- You can write a function that takes `io.Reader` and test it with `strings.NewReader`.
- You can build a small CLI with `flag` from memory.
- You know when to prefer `LookupEnv` over `Getenv`.

## Next

You've now covered every core standard-library and language concept a beginner needs. From here:

- **Cycle back to lesson 06** and continue the main syllabus if you started on the fundamentals track.
- **Try building a small tool of your own** in a fresh module (`go mod init github.com/you/name`) using everything from lessons 20-32. Log tailer, config validator, whatever you'd find useful.
