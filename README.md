# go-lab

A Go tutorial organized as a linear, 15-lesson syllabus. Each lesson is a self-contained folder with runnable code and its own tests — work through them in order.

![Quality Check & Learning Validation](https://github.com/ocrosby/go-lab/actions/workflows/quality-check.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Syllabus](#syllabus)
- [Configuration](#configuration)
- [Development](#development)
- [Repository Layout](#repository-layout)
- [Contributing](#contributing)
- [License](#license)

## Overview

`go-lab` teaches Go by walking through 15 numbered lessons, from `hello, world` to a hexagonal-architecture HTTP API with benchmarks. Each lesson concentrates on one concept — you can read the code, run it, and inspect its tests without cross-referencing the rest of the tree. Lessons that need multiple demos use subfolders inside the lesson (e.g. `10-panic-and-recover/goroutine-panic/`, `10-panic-and-recover/http-panic/`).

The intent is a lesson plan you can follow end-to-end, not a reference library to grep. If you just want to look up a pattern, `git grep` still works — but the syllabus is the recommended entry point.

## Features

- Linear 15-lesson syllabus covering fundamentals through production patterns
- Every lesson is a Go package with runnable code and tests
- Single Go module — `go test ./...` runs every lesson's tests
- Demonstrates Ginkgo/Gomega, `go.uber.org/mock`-generated mocks, `context`-driven concurrency, hexagonal architecture, and benchmarking

## Requirements

- Go 1.19 or newer
- (Optional) `ginkgo` and `mockgen` CLIs — only needed if you want to regenerate mocks or run Ginkgo suites directly rather than through `go test`:

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

## Usage

Run the first lesson:

```bash
go run ./lessons/01-hello
```

Run every lesson's tests:

```bash
go test ./...
```

Two lessons ship with intentional "before" demos of buggy behavior — see the [Development](#development) notes for the current status of those tests.

## Syllabus

Work through the lessons in order. Each row links to the lesson folder; open its `README.md` (where present) for setup notes.

| # | Lesson | Concept |
|---|---|---|
| 01 | [01-hello](lessons/01-hello) | Running a Go program, `package main`, `fmt` |
| 02 | [02-functions-and-packages](lessons/02-functions-and-packages) | Functions, packages, exported vs. unexported names |
| 03 | [03-testing-basics](lessons/03-testing-basics) | `go test`, table-driven tests |
| 04 | [04-test-suites-and-refactor](lessons/04-test-suites-and-refactor) | Ginkgo/Gomega, refactoring under a green suite (v1 → v2) |
| 05 | [05-composition](lessons/05-composition) | Struct embedding, method promotion |
| 06 | [06-interfaces-and-mocking](lessons/06-interfaces-and-mocking) | Small interfaces at the consumer, generated mocks |
| 07 | [07-goroutines-and-channels](lessons/07-goroutines-and-channels) | `go`, buffered/unbuffered channels, `WaitGroup`, `Mutex` |
| 08 | [08-channel-patterns](lessons/08-channel-patterns) | Pipeline, done/quit, fan-in, fan-out |
| 09 | [09-worker-pools](lessons/09-worker-pools) | Bounded parallelism with a fixed pool |
| 10 | [10-panic-and-recover](lessons/10-panic-and-recover) | Controlled panic recovery in goroutines and HTTP handlers |
| 11 | [11-http-clients-and-servers](lessons/11-http-clients-and-servers) | `net/http`, JSON, testing handlers |
| 12 | [12-dependency-injection](lessons/12-dependency-injection) | Constructor injection, interface seams |
| 13 | [13-design-patterns](lessons/13-design-patterns) | Builder, Prototype, Singleton, Adapter — in idiomatic Go |
| 14 | [14-production-api](lessons/14-production-api) | Hexagonal architecture, config, health checks, integration tests |
| 15 | [15-benchmarks](lessons/15-benchmarks) | `testing.B`, `benchstat`, reading pprof |

## Configuration

No configuration is required to run the lessons — they read no environment variables and open no external services. Lesson 14 (`14-production-api`) is the sole exception: it parses config from environment variables at startup. See `lessons/14-production-api/internal/config/config.go` for the defaults.

## Development

```bash
go build ./...    # compile every lesson
go test ./...     # run every lesson's tests
```

Two `before/` demo tests currently fail on purpose — they demonstrate uncontrolled panics for lessons that teach recovery. They are intentional but the tests do not yet wrap the panics in `assert.Panics`, so `go test ./...` exits non-zero. Tracked as a follow-up.

The repo's golangci-lint config (`.golangci.yml`) is in the v1 format and is not yet compatible with golangci-lint v2. Also a tracked follow-up.

## Repository Layout

```
go-lab/
├── lessons/          # The 15-lesson syllabus (this is the tutorial)
├── docs/             # Additional reference material
├── deployment/       # Deployment scaffolding (used by lesson 14)
├── templates/        # Project templates
├── scripts/          # Assorted helper scripts
├── go.mod            # Single Go module for the whole repo
└── README.md         # You are here
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bug reports and lesson-quality suggestions are welcome — open an issue before starting substantial work.

## License

MIT. See [LICENSE](LICENSE).
