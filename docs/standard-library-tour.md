# Standard library tour

One-page reference for the twelve standard-library packages you'll use most days. Not a replacement for [pkg.go.dev](https://pkg.go.dev/); a "you probably want one of these" cheat sheet with the shape of each package's most-used APIs.

Lessons throughout the syllabus link back here when a package first appears.

## Table of contents

- [`fmt`](#fmt---formatted-io) — formatted I/O
- [`strings`](#strings---text-manipulation) — text manipulation
- [`strconv`](#strconv---string-conversions) — string ↔ number conversions
- [`os`](#os---process--environment) — process and environment
- [`io`](#io---the-reader--writer-interfaces) — the `Reader` and `Writer` interfaces
- [`bufio`](#bufio---buffered-io) — buffered I/O
- [`time`](#time---timers-durations-formatting) — timers, durations, formatting
- [`encoding/json`](#encodingjson---the-json-workhorse) — the JSON workhorse
- [`log/slog`](#logslog---structured-logging-go-121) — structured logging (Go 1.21+)
- [`sort`](#sort--slices---sorting) and [`slices`](#sort--slices---sorting) — sorting
- [`flag`](#flag---cli-argument-parsing) — CLI argument parsing
- [`regexp`](#regexp---pattern-matching) — pattern matching

---

## `fmt` — formatted I/O

The printf-family everyone uses. Verbs you'll see 90% of the time:

| Verb | Meaning |
|---|---|
| `%v` | Default representation. Handy for any value. |
| `%+v` | Struct with field names — `{Name: Ada}` instead of `{Ada}`. |
| `%#v` | Go-syntax representation — great for debug logs. |
| `%T` | The Go type of the value. |
| `%d` | Integer, decimal. |
| `%f`, `%g` | Float. `%g` picks the shorter of `%e`/`%f`. |
| `%s` | String, or anything implementing `Stringer`. |
| `%q` | Quoted string — safe for logs. |
| `%x` / `%X` | Hex — lower / upper case. |
| `%b`, `%o` | Binary / octal. |
| `%p` | Pointer address. |
| `%w` | Wrap an error (only in `fmt.Errorf`). |

Common functions:

```go
fmt.Println(args...)          // print + newline to stdout
fmt.Printf(format, args...)   // formatted, no newline
fmt.Fprintln(w, args...)      // write to any io.Writer
fmt.Sprintf(format, args...)  // return string, don't print
fmt.Errorf("...: %w", err)    // wrap an error
```

## `strings` — text manipulation

Working with strings without regex:

```go
strings.Contains(s, "sub")
strings.HasPrefix(s, "http://")
strings.HasSuffix(name, ".go")
strings.Index(s, "sub")           // -1 if not found
strings.Count(s, "sub")
strings.ToLower(s), strings.ToUpper(s)
strings.TrimSpace(s)              // strip whitespace
strings.Trim(s, "()")             // strip any of those chars
strings.Replace(s, "old", "new", n)  // n=-1 for all
strings.ReplaceAll(s, "old", "new")
strings.Split(s, ",")             // → []string
strings.Join([]string{"a","b"}, ",")
strings.Fields(s)                 // split on whitespace, drop empties
strings.Repeat("ab", 3)           // "ababab"

// Building strings efficiently — use Builder, not += in a loop
var b strings.Builder
b.WriteString("hello")
b.WriteRune('!')
result := b.String()
```

`strings.Cut` (Go 1.18+) is often what you want instead of `Split`:

```go
before, after, found := strings.Cut("name=Ada", "=")
// before="name", after="Ada", found=true
```

## `strconv` — string conversions

```go
n, err := strconv.Atoi("42")            // string → int
f, err := strconv.ParseFloat("3.14", 64) // string → float64
b, err := strconv.ParseBool("true")     // string → bool

s := strconv.Itoa(42)                    // int → string
s := strconv.FormatFloat(3.14, 'f', 2, 64) // 3.14 → "3.14"
s := strconv.FormatBool(true)            // → "true"

// Quote a string for safe display
q := strconv.Quote(`hello "world"`)      // `"hello \"world\""`
```

`Atoi`/`Itoa` are aliases for `ParseInt`/`FormatInt` with `base=10`, `bitSize=int`. Use them for the common case.

## `os` — process and environment

```go
os.Args                      // []string, args[0] is the program name
os.Getenv("HOME")            // env var, "" if unset
val, ok := os.LookupEnv("HOME") // ok distinguishes unset from ""
os.Setenv("KEY", "value")
os.Exit(1)                    // exit with status code (skips deferred funcs!)

// Files
os.ReadFile("path")           // whole file → []byte
os.WriteFile("path", data, 0644)
os.Open("path")               // *os.File; use with defer f.Close()
os.Create("path")             // create for writing; use with defer f.Close()
os.Remove("path"), os.RemoveAll("dir")
os.MkdirAll("a/b/c", 0755)

// Standard streams
os.Stdin, os.Stdout, os.Stderr  // *os.File values, all io.Reader/Writer
```

## `io` — the Reader / Writer interfaces

Two interfaces every I/O type in the standard library satisfies:

```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
```

Utilities:

```go
io.Copy(dst, src)          // stream src → dst
io.ReadAll(r)              // read entire reader → []byte
io.LimitReader(r, n)       // wrap r, cap at n bytes
io.MultiReader(r1, r2)     // concatenate readers
io.MultiWriter(w1, w2)     // fan out one write to many writers
io.EOF                      // sentinel error signalling end of stream
```

`io.EOF` is the value most streaming code checks for. `errors.Is(err, io.EOF)` is safer than `err == io.EOF` when errors might be wrapped.

## `bufio` — buffered I/O

Wraps an `io.Reader` or `io.Writer` to buffer bytes, and adds line/word scanning:

```go
scanner := bufio.NewScanner(f)
for scanner.Scan() {
    line := scanner.Text()
}
scanner.Err()                         // check after loop!

scanner.Split(bufio.ScanWords)        // switch to word splitting
scanner.Buffer(buf, 10*1024*1024)      // allow lines up to 10 MiB

writer := bufio.NewWriter(f)
writer.WriteString("hello\n")
writer.Flush()                         // MUST call before closing!
```

**Two gotchas:** always check `scanner.Err()` after the loop (errors come out there, not from `Scan()`), and always `Flush()` a `bufio.Writer` before closing the underlying file.

## `time` — timers, durations, formatting

```go
t := time.Now()                              // current time
t := time.Date(2026, time.July, 26, 12, 0, 0, 0, time.UTC)
t := time.Unix(1700000000, 0)                // from Unix timestamp

t.Year(), t.Month(), t.Day()                 // components
t.Format(time.RFC3339)                        // "2026-07-26T12:00:00Z"
t.Format("2006-01-02 15:04:05")               // custom (see below)

// Durations
5 * time.Second, 100 * time.Millisecond, 2 * time.Hour
d := time.Since(startTime)                    // now - startTime
d.Seconds(), d.Milliseconds()

// Timers and tickers
timer := time.NewTimer(5 * time.Second)
<-timer.C  // block until fires

ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for range ticker.C { ... }
```

**The time-format gotcha:** Go's format strings use the reference date `Mon Jan 2 15:04:05 MST 2006` (`01/02 03:04:05PM '06 -0700`) instead of `YYYY-MM-DD`. Weird but memorable — the sequence 1 2 3 4 5 6 7 (month, day, hour, minute, second, year, offset).

## `encoding/json` — the JSON workhorse

The syllabus's [lesson 31](../lessons/31-json-and-struct-tags/) has the full treatment. The short form:

```go
b, err := json.Marshal(v)           // struct → []byte
err := json.Unmarshal(b, &v)        // []byte → struct

// Streaming
enc := json.NewEncoder(w); enc.Encode(v)
dec := json.NewDecoder(r); dec.Decode(&v)
dec.DisallowUnknownFields()          // strict decoding
```

Struct tags: `` `json:"field_name,omitempty"` ``.

## `log/slog` — structured logging (Go 1.21+)

Modern Go logging. Key-value pairs, JSON or text output, log levels:

```go
import "log/slog"

slog.Info("user created", "id", u.ID, "email", u.Email)
slog.Warn("rate limit exceeded", "user", uid, "count", n)
slog.Error("db query failed", "err", err)

// A structured logger with attributes attached
logger := slog.With("service", "users", "version", "1.2.3")
logger.Info("started")

// JSON output for machine-readable logs
h := slog.NewJSONHandler(os.Stdout, nil)
slog.SetDefault(slog.New(h))
```

Prefer `log/slog` over the older `log` package in new code. The old `log.Printf` still works everywhere; `slog` is what queryable, searchable log pipelines expect.

## `sort` / `slices` — sorting

Since Go 1.21, `slices` covers most sort needs generically. Reach for `sort` only for interface-based sorting on non-slice types.

```go
import "slices"

slices.Sort(xs)                            // in-place, ascending
slices.SortFunc(xs, func(a, b User) int {
    return cmp.Compare(a.Name, b.Name)
})
slices.SortStableFunc(xs, cmpFn)           // preserves order of equals
slices.IsSorted(xs)
slices.BinarySearch(xs, target)            // returns (index, found)

// Also common:
slices.Reverse(xs)
slices.Contains(xs, x)
slices.Index(xs, x)
slices.Min(xs), slices.Max(xs)             // Go 1.21+
slices.Concat(a, b, c)                      // Go 1.22+
slices.Clone(xs)                            // fresh backing array
```

Lesson 25 covers generics if that whole `[T any]` prefix is unfamiliar.

## `flag` — CLI argument parsing

Lesson 32 covers this in depth. Minimal shape:

```go
verbose := flag.Bool("v", false, "verbose output")
count := flag.Int("n", 10, "number of items")
name := flag.String("name", "world", "who to greet")

flag.Parse()

// After Parse: *verbose, *count, *name are set
// flag.Args() = remaining positional args
```

For anything beyond flat flags (subcommands like `git commit`, coloured help, shell completion), reach for [`spf13/cobra`](https://github.com/spf13/cobra).

## `regexp` — pattern matching

RE2 syntax (no backreferences, no lookarounds — but linear-time guaranteed):

```go
r := regexp.MustCompile(`\d+`)   // panic if bad pattern (fine for constants)
r.MatchString("abc 42 def")       // true
r.FindString("abc 42 def")        // "42"
r.FindAllString("1, 2, 3", -1)    // ["1", "2", "3"]
r.ReplaceAllString(s, "[$0]")     // wrap matches in brackets

// With capture groups:
r := regexp.MustCompile(`(\w+)@(\w+)`)
m := r.FindStringSubmatch("ada@example.com")
// m[0] = whole match, m[1] = "ada", m[2] = "example"
```

Compile once (at package init or with `MustCompile` for constants), reuse many times.

---

## When to reach for the standard library vs a third-party package

Go's standard library is unusually rich. Before adding a dependency:

- **JSON parsing** — `encoding/json` is fine for 99% of use cases. Only reach for `json-iterator` or `sonic` after profiling shows it matters.
- **HTTP client/server** — `net/http` is production-grade. Third-party frameworks (Gin, Echo) add routing conveniences over it; needed less often since Go 1.22's ServeMux upgrade.
- **Logging** — `log/slog` since Go 1.21. External loggers (zap, logrus) still exist but the standard library is the default now.
- **Testing** — `testing` + `httptest` cover ~90%. Add `testify` (assertions/mocks) or Ginkgo (BDD) only when the standard library genuinely doesn't fit.
- **CLI** — `flag` is minimal; `spf13/cobra` is the community standard for anything more complex.

## See also

- [`docs/go-build-directives.md`](go-build-directives.md) — the full `//go:` catalog.
- [`docs/csp-and-go-concurrency.md`](csp-and-go-concurrency.md) — where goroutines and channels come from.
- [`docs/benchmarking.md`](benchmarking.md) — measuring performance.
- [pkg.go.dev](https://pkg.go.dev/) — canonical package documentation.
- [Go by Example](https://gobyexample.com/) — worked snippets for every stdlib topic.
