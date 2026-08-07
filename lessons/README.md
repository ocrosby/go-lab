# Lessons

The 33-lesson syllabus for `go-lab`. Each folder is a self-contained mini-lesson with its own `README.md`, runnable code, and (where the concept warrants it) tests. Work through them in the recommended order — see the [syllabus table in the root README](../README.md#syllabus) for the full sequence.

## How this directory is organized

Every top-level folder here matches the pattern `NN-name/`, where `NN` is a two-digit number and `name` is a short kebab-case description. **Folder numbers reflect the order lessons were *authored*, not the order they should be *read*.** When the fundamentals track was added later, its folders got high numbers (20–33) to avoid renumbering older lessons and breaking cross-references. The root README's syllabus table shows the pedagogical order.

Each lesson folder generally contains:

- `README.md` — required. Concept, why it matters, how to run, "try it yourself", common pitfalls, and a "you've understood this lesson when..." checklist. Every lesson has this — CI enforces it.
- One or more `.go` source files with a `main.go` or a small library package.
- One or more `*_test.go` files exercising the concept through the public API. Not every lesson has tests — a few pure-demo lessons (variables, control flow, go-modules) are read-and-run-and-experiment rather than test-anything-observable.
- Sub-folders when a lesson has multiple worked examples (e.g. `11-panic-and-recover/goroutine-panic/`, `11-panic-and-recover/http-panic/`).

## Running a single lesson

Every lesson under this directory is runnable from the repo root. Three shapes:

```bash
# Run a program lesson — prints something to stdout
go run ./lessons/01-hello

# Run a test lesson — reports pass/fail
go test ./lessons/04-testing-basics

# Run everything in a lesson including sub-packages
go test ./lessons/07-interfaces-and-mocking/...
```

For lessons with a subcommand (like the wcish demo in lesson 33):

```bash
go run ./lessons/33-file-io-and-cli/cmd/wcish -lines README.md
```

The root [`Makefile`](../Makefile) wraps the whole-repo variants — `make test`, `make build`, `make lint`. Prefer those for the "run everything" case; they encode the correct filters (like skipping the intentional-panic packages under `11-panic-and-recover/*/before/`).

## Lessons that intentionally fail

Two sub-lessons under `11-panic-and-recover/` demonstrate what happens **without** panic recovery. Their tests are supposed to fail — that's the pedagogical point.

- `lessons/11-panic-and-recover/goroutine-panic/before/`
- `lessons/11-panic-and-recover/http-panic/before/`

The matching `after/` packages show the fix and pass. `make test` filters the `before/` packages out so a clean checkout produces green output; use `make test-all` if you want to see the intentional failures. See [`lessons/11-panic-and-recover/README.md`](./11-panic-and-recover/README.md) for the whole picture.

## Lessons that live in their own Go module

Two lessons have their own `go.mod` (visible with `find lessons -name go.mod`):

- **`lessons/14-design-patterns/`** — a separate module for historical reasons; contains no external dependencies. Would fold cleanly back into the root module in a future refactor.
- **`lessons/15-production-api/`** — a separate module because it depends on packages (Cobra, Viper, Zap, Swagger) that shouldn't leak into every simple lesson.

`go test ./...` at the repo root only exercises the root module. The CI workflow builds, lints, and vuln-scans the submodules independently (`Build submodules` and `Lint (lessons/15-production-api)` jobs). If you want to work inside one:

```bash
cd lessons/15-production-api
go test ./...
```

## Where to start

Depends on your background:

| You are | Start at |
|---|---|
| Brand new to programming | The [root README's "brand new?" call-out](../README.md), which sends you to `docs/tutorials/getting-started.md` first. |
| New to Go, know another language | Lesson [01-hello](./01-hello/), then follow the syllabus table in order — the fundamentals track (20–25, 27–33) is threaded in at the right dependency points. |
| Experienced Go dev, browsing | Any lesson that looks interesting. Concurrency: 08–10. HTTP: 12, 17–19. Testing shape: 07 and 31-json. |
| Looking for a specific concept | Grep the syllabus table in the root README, or `grep -r "concept" lessons/*/README.md`. |

## Related reference material

Deep-dive material that lessons link out to lives under [`docs/`](../docs/):

- [`docs/standard-library-tour.md`](../docs/standard-library-tour.md) — one-page reference for the twelve most-used stdlib packages.
- [`docs/tooling.md`](../docs/tooling.md) — the Go toolchain (`go run/build/test/mod`), `-race`, `-cover`, community tools (goimports, golangci-lint, govulncheck).
- [`docs/go-build-directives.md`](../docs/go-build-directives.md) — every `//go:` directive (build constraints, `//go:generate`, `//go:embed`, cgo, WebAssembly).
- [`docs/csp-and-go-concurrency.md`](../docs/csp-and-go-concurrency.md) — where goroutines and channels come from.
- [`docs/benchmarking.md`](../docs/benchmarking.md) — the `-bench` and `pprof` deep dive.
- [`docs/go-workspaces.md`](../docs/go-workspaces.md) — multi-module workspaces via `go.work`.

## Contributing a new lesson

The template every lesson follows is documented in [`CLAUDE.md`](../CLAUDE.md#lesson-readme-template) at the repo root. In short: `README.md` (following the template), one or more `.go` files, tests where applicable, an entry in the root README's syllabus table.

The `lesson-readmes` CI job enforces that every `lessons/NN-*/` folder has a `README.md`. If you add a lesson without one, CI fails with a specific missing-list.
