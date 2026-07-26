# go-lab

A Go tutorial organized as a linear, 25-lesson syllabus. Each lesson is a self-contained folder with runnable code and its own tests — work through them in order.

![Quality Check & Learning Validation](https://github.com/ocrosby/go-lab/actions/workflows/quality-check.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> **Brand new to programming?** Start at [docs/tutorials/getting-started.md](docs/tutorials/getting-started.md). It walks you through installing Go, opening a terminal, and running your first program, with no prior experience assumed.

## Table of contents

- [What is Go, in one sentence](#what-is-go-in-one-sentence)
- [Requirements](#requirements)
- [Installation](#installation)
- [Your first Go program](#your-first-go-program)
- [Running tests](#running-tests)
- [Syllabus](#syllabus)
- [Repository layout](#repository-layout)
- [Notes on the code you'll see](#notes-on-the-code-youll-see)
- [Contributing](#contributing)
- [License](#license)

## What is Go, in one sentence

Go is a small, fast programming language from Google, designed for building server programs and command-line tools that need to do many things at once. If you can read this README you can start writing Go this afternoon.

## Requirements

- **Go 1.26 or newer.** To check, open a terminal (Terminal on macOS, Command Prompt or PowerShell on Windows, or your favourite terminal on Linux) and type `go version`. If it prints a version number, you're set. If it says "command not found," follow the getting-started guide linked at the top.
- **Git.** Same check: `git --version` in a terminal. macOS and most Linux distros ship it; on Windows install [Git for Windows](https://gitforwindows.org/).
- **Optional tools** — only needed if you plan to regenerate mocks or run the Ginkgo test suites directly rather than through `go test`:

  ```bash
  go install github.com/onsi/ginkgo/v2/ginkgo@latest
  go install go.uber.org/mock/mockgen@latest
  ```

## Installation

```bash
git clone https://github.com/ocrosby/go-lab.git
cd go-lab
go mod download
```

That last command reads the list of external Go packages this repo uses and downloads them into your local cache. It only needs to run once (or after someone adds a new dependency).

## Your first Go program

```bash
go run ./lessons/01-hello
```

You should see:

```
Hello World!
```

If you did, congratulations — you just ran your first Go program. Head to [`lessons/01-hello/README.md`](lessons/01-hello/README.md) to look at what actually happened.

## Running tests

The friendliest way is:

```bash
make test
```

That runs every lesson's tests and shows a green result. If you don't have `make`, the equivalent Go command is:

```bash
go list ./... | grep -v '10-panic-and-recover/.*/before$' | xargs go test
```

You have to skip the `10-panic-and-recover/*/before` packages because those are **intentional demonstrations** of crashes — lesson 10 shows what happens *without* panic recovery, and their tests are supposed to blow up. The `after` packages in the same lesson show how to fix it, and those pass. See [`lessons/10-panic-and-recover/README.md`](lessons/10-panic-and-recover/README.md) for the whole picture.

Other useful Make targets:

```bash
make help        # list every target
make hello       # run the first lesson
make build       # compile every lesson
make lint        # run the linter
```

## Syllabus

Work through the lessons in the order listed below. The **Order** column is the recommended reading order — the **Folder** column shows the actual directory name (numbers reflect the order lessons were added, not the pedagogical sequence). Each row links to the lesson folder — open its `README.md` for the concept, how to run it, and a small "try it yourself" exercise.

**If you're brand new to Go**, start at 1 and work down. **If you already know the language basics**, you can skim the language fundamentals (marked *fundamentals*) and start at "Functions and packages."

| Order | Folder | Concept |
|---|---|---|
| 1 | [01-hello](lessons/01-hello) | Running a Go program, `package main`, `fmt` |
| 2 | [19-variables-and-types](lessons/19-variables-and-types) | *fundamentals* — `var`/`:=`, primitives, zero values, `const`, `iota` |
| 3 | [20-control-flow](lessons/20-control-flow) | *fundamentals* — `if`/`switch`/`for` (Go's only loop) |
| 4 | [02-functions-and-packages](lessons/02-functions-and-packages) | Functions, packages, exported vs. unexported names |
| 5 | [03-testing-basics](lessons/03-testing-basics) | `go test`, table-driven tests |
| 6 | [04-test-suites-and-refactor](lessons/04-test-suites-and-refactor) | Ginkgo/Gomega, refactoring under a green suite (v1 → v2) |
| 7 | [21-slices-and-maps](lessons/21-slices-and-maps) | *fundamentals* — the two workhorse collection types |
| 8 | [22-pointers](lessons/22-pointers) | *fundamentals* — `&`/`*`, value vs pointer receivers, `new` vs `&T{}` |
| 9 | [23-error-handling](lessons/23-error-handling) | *fundamentals* — `if err != nil`, `%w` wrap, `errors.Is`/`errors.As` |
| 10 | [24-defer-and-cleanup](lessons/24-defer-and-cleanup) | *fundamentals* — LIFO order, argument evaluation, common patterns |
| 11 | [05-composition](lessons/05-composition) | Struct embedding, method promotion |
| 12 | [06-interfaces-and-mocking](lessons/06-interfaces-and-mocking) | Small interfaces at the consumer, fakes vs. mocks |
| 13 | [07-goroutines-and-channels](lessons/07-goroutines-and-channels) | `go`, buffered/unbuffered channels, `WaitGroup`, `Mutex` |
| 14 | [08-channel-patterns](lessons/08-channel-patterns) | Pipeline, done/quit, fan-in, fan-out |
| 15 | [09-worker-pools](lessons/09-worker-pools) | Bounded parallelism with a fixed pool |
| 16 | [10-panic-and-recover](lessons/10-panic-and-recover) | Controlled panic recovery in goroutines and HTTP handlers |
| 17 | [11-http-clients-and-servers](lessons/11-http-clients-and-servers) | `net/http`, JSON, testing handlers |
| 18 | [12-dependency-injection](lessons/12-dependency-injection) | Constructor injection, interface seams |
| 19 | [13-design-patterns](lessons/13-design-patterns) | Builder, Prototype, Singleton, Adapter — in idiomatic Go |
| 20 | [14-production-api](lessons/14-production-api) | Hexagonal architecture, config, health checks, integration tests |
| 21 | [15-benchmarks](lessons/15-benchmarks) | `testing.B`, `benchstat`, reading pprof |
| 22 | [16-restful-routing](lessons/16-restful-routing) | Go 1.22+ ServeMux, path params, status codes, static files, SSE |
| 23 | [17-http-middleware](lessons/17-http-middleware) | Middleware chain: request ID, logging, recover, auth, body limit, CORS |
| 24 | [18-http-client-depth](lessons/18-http-client-depth) | `http.Client` tuning, retries via `RoundTripper`, testing seams |
| 25 | [25-generics](lessons/25-generics) | Type parameters, `~` and `comparable` constraints, `slices`/`maps` packages |

## Repository layout

```
go-lab/
├── lessons/          # The 25-lesson syllabus (this is the tutorial)
├── docs/             # Deep-dive reference material and tutorials
├── deployment/       # Deployment scaffolding (used by lesson 14)
├── templates/        # Project templates for starting your own Go project
├── scripts/          # Helper scripts
├── Makefile          # `make help` lists targets
├── go.mod            # Single Go module for the whole repo
├── CLAUDE.md         # Conventions for Claude Code sessions in this repo
└── README.md         # You are here
```

## Notes on the code you'll see

**Two testing frameworks.** Most lessons use the Go standard library's `testing` package — simple functions like `TestFoo(t *testing.T)`. A few lessons (04, 06, 08) use [**Ginkgo** and **Gomega**](https://onsi.github.io/ginkgo/), a behaviour-driven style with `Describe`/`It` blocks. Both are common in real Go codebases, so the tutorial shows you both. Start with the standard library style; you'll pick up Ginkgo naturally when you meet it.

**`//go:build ignore` at the top of some files.** This line tells Go's build tool "do not include this file when compiling the package — only run it when the user explicitly asks with `go run <file>`." It's used for standalone example files that live inside a package but shouldn't participate in that package's build. When you see it, run the file directly:

```bash
go run ./lessons/07-goroutines-and-channels/primitives.go
```

For the full story on Go's build directives, see [docs/go-build-directives.md](docs/go-build-directives.md).

**Two Go modules.** The repo is one big module with one exception: `lessons/14-production-api/` has its own `go.mod` because it depends on packages (Cobra, Viper, Zap, Swagger) that shouldn't leak into every simple lesson. `go test ./...` at the repo root only tests the root module — the workflow tests submodule 14 separately.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bug reports and lesson-quality suggestions are welcome — open an issue before starting substantial work.

## License

MIT. See [LICENSE](LICENSE).
