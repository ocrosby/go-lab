# The Go toolchain

`go` is one binary with many subcommands. Below are the ones you'll use daily, monthly, and once-in-a-blue-moon, plus the three community tools every serious Go project runs.

## Daily — you'll type these more than any other

| Command | What it does |
|---|---|
| `go run ./path/to/pkg` | Compile and run a package. Great for scripts and REPL-ish exploration. |
| `go run ./path/to/pkg -flag value` | Same, with flags passed to the program. |
| `go build ./...` | Compile every package in the module. Errors out on the first failure. |
| `go test ./...` | Run every test in the module. Add `-v` for verbose, `-run TestFoo` for one, `-race` for the race detector. |
| `go test -race -count=1 ./...` | The full "confidence" run — race detector on, no cache. |
| `go mod tidy` | Add missing dependencies, remove unused. Run before every commit that touches `go.mod`. |
| `go fmt ./...` | Rewrite every `.go` file in canonical Go format. |

## Weekly — you'll want these regularly

| Command | What it does |
|---|---|
| `go vet ./...` | Static analysis for common bugs (nil derefs, printf mismatches, unreachable code). Free — runs in seconds. |
| `go doc <symbol>` | Print the docs for a symbol without leaving the terminal. `go doc net/http.Handler` or `go doc fmt.Errorf`. |
| `go get -u ./...` | Update every direct dependency to the latest minor version. |
| `go mod why <package>` | Explain which chain of imports pulled a dependency in. |
| `go install github.com/foo/bar@latest` | Install a tool globally to `$GOPATH/bin` (or `$GOBIN`). |
| `go test -bench=. ./...` | Run every benchmark. |
| `go test -cover ./...` | Test with coverage summary. Add `-coverprofile=cov.out` + `go tool cover -html=cov.out` for a browser view. |
| `go tool pprof <profile>` | Interactive profile explorer. Lesson 15 and `docs/benchmarking.md` cover this. |

## Once in a blue moon

| Command | What it does |
|---|---|
| `go work init ./mod1 ./mod2` | Bootstrap a `go.work` for multi-module local development. See `docs/go-workspaces.md`. |
| `go generate ./...` | Run every `//go:generate` directive in source files. `docs/go-build-directives.md` has the full story. |
| `go mod vendor` | Copy dependencies into `vendor/`. Some enterprises require it; most projects don't. |
| `go env` | Print the Go environment (GOPATH, GOROOT, GOOS, GOARCH, GOMODCACHE, etc.). |
| `go clean -modcache` | Nuke the module download cache. Fix for genuinely corrupted downloads. Slow to rebuild. |
| `go tool objdump -s <symbol> <binary>` | Disassemble compiled code. You'll never need this. Included for completeness. |

## Community tools

Not shipped in the Go toolchain, but every serious Go project uses at least the first two.

### `gofmt` and `goimports`

`gofmt` is part of the toolchain (`go fmt` calls it). `goimports` (`golang.org/x/tools/cmd/goimports`) is the strict superset every editor should use — it does what `gofmt` does, plus adds/removes/organizes imports.

Install:

```bash
go install golang.org/x/tools/cmd/goimports@latest
```

Wire into your editor (every Go plugin has a "run goimports on save" option — turn it on).

### `golangci-lint`

The meta-linter. Runs 50+ individual analyzers in one pass, respects a `.golangci.yml` config, has fast incremental modes for CI. This repo's `.golangci.yml` runs `errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`, plus `gofmt`/`goimports`.

Install:

```bash
brew install golangci-lint                        # macOS
# or:
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
```

Run:

```bash
golangci-lint run
golangci-lint run --fix    # apply auto-fixes for rules that support them
```

Two rules from the parent claude-config repo apply everywhere:

- **Every suppression needs a reason.** `//nolint:errcheck — test setup, error not actionable`. Bare `//nolint` is a Must Fix in code review (see `rules/lint-suppression.md`).
- **golangci-lint v2 config format.** Modern versions require `version: "2"` at the top of `.golangci.yml`. Migrate v1 configs with `golangci-lint migrate`.

### `govulncheck`

Standard-library vulnerability scanner. This repo runs it in CI:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

Reports known CVEs in your dependencies and (usefully) tells you whether your code actually calls the vulnerable function — false positives are rare.

## The `-race` flag

Every test run should use it when correctness matters. The race detector instruments memory accesses and reports data races at runtime. Costs 2-10× CPU and memory. Don't ship with `-race`, but always test with it:

```bash
go test -race ./...
go run -race ./cmd/server    # in dev
```

## The `-cover` flag

Line coverage. Two shapes:

```bash
# Just the number
go test -cover ./...

# Detailed HTML report
go test -coverprofile=cov.out ./...
go tool cover -html=cov.out
```

Treat coverage as a **detector**, not a target — see `rules/black-box-testing.md`. 100% coverage of a suite that only spies on internals is worse than 60% coverage of a suite that verifies outcomes.

## Repo shortcuts

This repo's `Makefile` wraps the common commands so you type less:

```bash
make hello       # run lesson 01
make test        # run every lesson's tests (skips intentional-panic packages)
make test-all    # includes the intentional-panic packages
make build       # compile every lesson
make vet         # go vet ./...
make lint        # golangci-lint run
make vulncheck   # govulncheck ./...
make clean       # remove test caches and coverage artifacts
```

Prefer the Make targets for daily work — they encode the correct filters (like skipping the intentional-panic packages).

## See also

- [`docs/standard-library-tour.md`](standard-library-tour.md) — the packages you'll use with these tools.
- [`docs/go-build-directives.md`](go-build-directives.md) — `//go:` compiler directives, useful with `go generate` and `//go:build` constraints.
- [`docs/benchmarking.md`](benchmarking.md) — the `-bench` and `pprof` deep dive.
- [go.dev/doc/](https://go.dev/doc/) — official documentation.
