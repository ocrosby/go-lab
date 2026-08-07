# go-lab — CLAUDE.md

This file records the audience, style, and structural conventions for the `go-lab` repository. It is read at the start of every Claude Code session in this directory and takes precedence over general habits.

## Audience — read this first

Every piece of writing and every lesson in this repo has to serve **two audiences at the same time**:

1. **A 6th-grade nephew** who is curious about programming, has typed maybe a hundred lines of code in their life, and has never used a terminal seriously. They do not know what `git` is, what a "module" means, or the difference between "the language" and "the compiler." Their goal is to see `Hello World!` appear on their screen and feel like a programmer.
2. **A senior Go engineer** who has shipped production Go for years, opens the repo to check a pattern, and needs to find the interesting bit inside five seconds. They already know what CSP is; they want to see how you wired the pipeline.

If a change makes it harder for either audience, it is the wrong change. In practice the tensions are smaller than they sound — beginner-friendly writing is almost always *also* clearer for pros. The place both audiences meet is: **short opening sentences, precise vocabulary, runnable examples, one concept per lesson.**

### The 6th-grader test

Before shipping any change to a lesson, a lesson README, or an entry-point doc, read it back and ask:

- Could a curious 6th grader, sitting at a fresh laptop with nothing installed, follow the *first* command on the page without googling any word?
- If they hit an error, does the doc say what to check?
- Is there something concrete they can do (run a command, change a value, see it change)?

If the answer to any of these is no, the doc is not done.

### The pro test

- Can a senior Go dev find the runnable code in this directory in under 5 seconds?
- Is the *interesting* thing about this lesson stated in the first paragraph, or buried under narrative?
- If they wanted to skip to the point, can they?

## Testing shape

**Tests in this repo are black-box.** The diagnostic — from [Testing behavior, not implementation](https://omarcrosby.com/posts/testing-behavior-not-implementation/) — is:

> *If I change how the code is written without changing what it does, will the test still pass?*

If the answer is no, the test is measuring implementation shape, not behaviour. Fix the seam.

Concretely for this repo:

- **Mock at the edges of your system, not at the edges of your classes.** Databases, HTTP clients, message queues, the clock, filesystems — legitimate edges. Same-team collaborators inside your own module — not edges.
- **Prefer a small in-memory fake to a generated mock** at stateful boundaries. `internal/testutil/fake_repository.go` in lesson 15 is the reference example: an in-memory `domain.UserRepository` with optional fault-injection knobs (`FailNextCreate = err`) for tests that need the repo to fail.
- **Assert on outcomes, not call traces.** `svc.GetUser(ctx, id)` after `svc.CreateUser(...)` returning the user is the right assertion. `mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)` is not.
- **Do not reach into unexported fields.** `accord.state = "on"` in an arrange section pins the test to the struct's internal layout. Use the public API to set up the state you want.
- **Route through the public surface.** HTTP handler tests go through the `ServeMux` via `httptest.NewRecorder`/`NewServer`, not by calling unexported handler methods directly. See `lessons/15-production-api/internal/infrastructure/adapters/http/user_handler_test.go` for the reference pattern.
- **Test packages use `_test` suffix by default** (`package foo_test` in `foo/xxx_test.go`), which enforces access through the exported API. Use the internal test package (`package foo`) only when a test genuinely needs unexported symbols and the design cannot be reworked to make the seam public.

**Reference lessons** for testing shape:

- `lessons/19-http-client-depth` — a stub `http.RoundTripper` as the standard-library test seam, replacing the custom `IHttpClient` interface pattern.
- `lessons/07-interfaces-and-mocking` — the anti-pattern-to-refactor pattern, with both a real-collaborator test (`accord_factory_test.go`) and an edge-contract test with a fake (`accord_test.go`) side by side.
- `lessons/15-production-api/internal/testutil/fake_repository.go` — canonical in-memory fake with fault-injection.

**Anti-references** — what to avoid:

- Any test that asserts on `mockX.EXPECT().Method(...)` for a same-team `X`.
- Any test whose arrange section writes to an unexported field.
- Any test named `TestFoo_MethodName_Success` — that shape is method-per-test decomposition, not behaviour-per-test.

## Authoritative documents

There are **two** entry points, and no others. Everything else is reference material or archive.

| Purpose | File | Audience |
|---|---|---|
| Front door + syllabus | `README.md` | Both — a pro should see the syllabus in 5 seconds; a beginner should see the "brand new? start here" link in 5 seconds. |
| Ground-up setup | `docs/tutorials/getting-started.md` | Beginner-first — installs Go, opens a terminal, runs the first program. |

When adding a new "map" or "index" doc, **stop and consider whether it can be a section in one of the two above** first. The repo has a history of accumulating multiple overlapping roadmaps that drift apart from each other; do not add to that pile.

Documents that exist but should *not* be linked from the README as entry points:

- `docs/INDEX.md`, `docs/LEARNING_ROADMAP.md`, `docs/EXECUTIVE_SUMMARY.md`, `docs/project.md` — legacy overviews that predate the current `lessons/` layout. Do not edit these to add new content; if you have new content, put it in one of the two authoritative docs.

## Lesson README template

**Every lesson under `lessons/NN-*/` must have a `README.md`.** No exceptions. A lesson without a README is a broken lesson.

The template — apply consistently, omit sections only when they genuinely do not apply:

```markdown
# <Lesson name>

<One sentence: what this lesson teaches, in words a curious beginner understands.>

## Why it matters

<One short paragraph — what problem this concept solves, or where you meet it
in real Go code. This is the pro-onboarding section too: it names the concept
in the vocabulary they already have.>

## Prerequisites

- Lesson <NN>: <what it teaches>
- (Any tool that must be installed beyond the base Go toolchain.)

## Run it

```bash
go run ./lessons/NN-<name>
```

Expected output:

```
<literal output the reader should see>
```

## What's in this folder

| File | What it demonstrates |
|---|---|
| `foo.go` | <one line> |
| `foo_test.go` | <one line> |

## Mental model

<Optional but recommended for hard topics. A short prose or diagram that
gives the reader the shape of the idea before diving into syntax. See
`lessons/08-goroutines-and-channels/README.md` for the reference example.>

## Try it yourself

1. <A small, concrete change to make.>
2. <Another small change.>
3. <What to observe / what should change in the output.>

## Common pitfalls

- <A mistake beginners make with this concept, and how to recognize it.>

## You've understood this lesson when...

- You can explain <the core concept> in one sentence.
- You can predict what a small variation on the example will do before running it.
- <Any concrete skill that unlocks the next lesson.>

## Next

- Next lesson: NN-{name} (link to `../NN-{name}/`)
- Related deep-dive doc (optional): e.g. link to `../../docs/csp-and-go-concurrency.md`
```

**Reference example**: `lessons/08-goroutines-and-channels/README.md`. When in doubt about tone or depth, match that file.

**Anti-examples**:

- `lessons/01-hello/README.md` (three-line stub — do not ship anything this thin)
- `lessons/16-benchmarks/README.md` (500+ lines of narrative with no runnable code in the directory — either move the doc content out to `docs/` and add real code, or truncate)

## Terminology and tone

- **Introduce every jargon term the first time it appears.** "A *goroutine* is a function running independently and concurrently with the caller" is fine on first mention. "Spawn a goroutine" without introduction is not.
- **Assume no shell fluency.** Say "open a terminal (Terminal on macOS, Command Prompt or PowerShell on Windows)" the first time. Do not say "run this at the shell" without saying what a shell is.
- **Do not assume the reader knows what a package, a module, or a Go path is.** These are all Go-specific concepts that need one-line intros the first time they show up.
- **Prefer imperative sentences with expected output.** "Run `go run ./lessons/01-hello`. You should see `Hello World!` printed." beats "The command below executes the program."
- **Emoji sparingly, and only if the repo already uses them consistently.** The current repo uses them heavily in the older docs; the newer lesson READMEs (07) do not. Prefer none in new writing.
- **Sentence case for headings**, following `rules/docs-principles.md`. The historic docs use title case; do not chase them, but do not introduce more.

## Standards for entry-point commands

The **first runnable command** the reader is asked to type must succeed on a clean checkout of `main`.

- `go run ./lessons/01-hello` — must print `Hello World!`.
- `go build ./...` — must succeed everywhere.
- The "run all tests" command the README recommends must **not** produce `FAIL`. Two `lessons/11-panic-and-recover/*/before/` packages intentionally demonstrate panics; when recommending a bulk test command, use the same filter the CI workflow uses:

  ```bash
  go list ./... | grep -v '11-panic-and-recover/.*/before$' | xargs go test
  ```

  Or wrap it in a Makefile / Taskfile target so beginners type one word.

## Structure conventions

- **One concept per lesson.** If a lesson teaches two things, split it or make one a subfolder inside the lesson.
- **Sub-lessons live inside the numbered folder**, not at the top level (e.g. `lessons/11-panic-and-recover/goroutine-panic/` and `.../http-panic/`). Do not add a `lessons/11b-*/`.
- **Standalone `//go:build ignore` example files** are welcome inside a lesson but must be mentioned in the lesson README's "What's in this folder" table with the note that they run via `go run <file>` (not `go test`). A beginner who opens one of these files and sees `//go:build ignore` should already have read what it means.
- **Submodules** (`lessons/*/go.mod`) exist only when the sub-tree needs its own dependency graph. `lessons/15-production-api` legitimately does; `lessons/14-design-patterns` may not — challenge new submodule additions.

## Content that does not belong in a lesson

- **Deep-dive reference material** (`docs/go-build-directives.md`, `docs/csp-and-go-concurrency.md`) — link from the lesson, do not inline.
- **500+ line theory dumps** — see the `lessons/16-benchmarks` anti-example above.
- **Advice about production deployment** — that lives under `deployment/`, not in a lesson README.

## Interaction and change protocol

- **A lesson without a working example is a bug.** Prefer fixing the example over patching the README to explain why the example is broken.
- **A README that documents behaviour the code no longer has is worse than a missing README.** When you change code, check the lesson README in the same commit.
- **When the entry docs disagree with each other**, the root `README.md` wins. Update the others or delete them; do not leave the disagreement.
- **Do not add a new "roadmap" or "index" doc.** See "Authoritative documents" above.

## What Claude should do proactively in this repo

- When editing any lesson's code, check whether that lesson's README still accurately describes what's in the folder — and update it in the same change if not.
- When adding a new lesson, create the README from the template in this file *before* writing the code. The template forces you to name the concept in one sentence, which is the hardest and most valuable part.
- When editing the root `README.md`, verify every command and every link before saving. This file is the entry point; broken commands here compound the fastest.
- When touching `docs/`, ask whether the new content belongs in one of the two authoritative docs instead of a new file.

## Known drift to watch for

- **Legacy roadmaps** (`docs/INDEX.md`, `docs/LEARNING_ROADMAP.md`, `docs/EXECUTIVE_SUMMARY.md`) still reference the pre-restructure layout (`learning/`, `examples/`, `testing/`). They should be brought into line or retired; do not extend them.
- **`docs/oop.md`** is a 0-byte placeholder — delete or write.
- **`maps.go`** at the repo root is unattached to any lesson. It should live inside `lessons/02-functions-and-packages/` or be deleted.
- **Two testing frameworks (stdlib and Ginkgo/Gomega) coexist without a stated rationale** — the root README should include a one-paragraph note when that is added.

## Cross-references

- `docs/tutorials/getting-started.md` — beginner install / first-run path
- `lessons/08-goroutines-and-channels/README.md` — reference lesson-README quality bar
- `docs/go-build-directives.md`, `docs/csp-and-go-concurrency.md` — deep-dive references for lessons to link to, not duplicate
- `rules/docs-principles.md` (in the user's global config) — Write-the-Docs conventions this repo inherits
