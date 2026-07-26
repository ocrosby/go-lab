# Contributing to go-lab

Bug reports, lesson-quality suggestions, and new lessons are all welcome. This is a small guide — the goal is to get you from "I want to change something" to "I have an open PR" without ceremony.

## Before you start

**Open an issue first** for anything beyond a typo or a broken link. A one-paragraph description of what you want to change lets us tell you if someone's already on it, or if there's a better place for the change. Substantial changes without a linked issue often need to be reworked.

**Follow the code of conduct.** See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md). Be kind, especially to newcomers — this repo exists to help people learn.

## Set up your local copy

```bash
git clone https://github.com/<your-fork>/go-lab.git
cd go-lab
go mod download
make test
```

If `make test` is all green, you're ready.

Requirements: **Go 1.26 or newer**, **git**, and — for lint changes — **golangci-lint v2**.

## Making a change

1. Create a branch named after the intended commit — e.g. `fix/lesson-05-typo`, `docs/getting-started-windows`, `feat/lesson-16-generics`.

   ```bash
   git checkout -b <type>/<short-description>
   ```

2. Make your change. Run `make test` and `make lint` before pushing.

3. Commit with a [Conventional Commits](https://www.conventionalcommits.org/) subject line:

   ```
   docs(lesson-05): clarify method promotion example
   fix(workflow): pin actions to non-deprecated versions
   feat(lesson-12): add functional-options variant
   ```

   Common types in this repo: `docs`, `fix`, `feat`, `refactor`, `test`, `chore`, `ci`. One PR = one type + one scope. If your PR needs "and" in the description, split it.

4. Push and open a PR. Fill in what the PR does and how you verified it.

## Adding or editing a lesson

Every lesson under `lessons/NN-*/` follows the canonical template documented in [CLAUDE.md](CLAUDE.md#lesson-readme-template). The short version:

- One concept per lesson.
- Every lesson has a `README.md`. New lessons include one from the start.
- The first runnable command in the README must succeed on a clean checkout of `main`.
- Match the tone and depth of [`lessons/07-goroutines-and-channels/README.md`](lessons/07-goroutines-and-channels/README.md) — the reference example.

The **[6th-grader test](CLAUDE.md#the-6th-grader-test)** and **[pro test](CLAUDE.md#the-pro-test)** in `CLAUDE.md` describe the audience you're writing for. Both apply.

## Testing requirements

- New Go code needs tests. Bug fixes need a test that would have caught the bug.
- `make test` must pass locally before you push.
- If your change touches the CI workflow, wait for the CI check on your PR to go green before requesting review.

## Documentation

- Update the lesson's `README.md` in the same commit as any code change to that lesson.
- Broken links are a review-blocking issue. Verify every link you add.
- The two authoritative entry points are `README.md` and `docs/tutorials/getting-started.md` — do not add a third "index" or "roadmap" doc. See [CLAUDE.md](CLAUDE.md#authoritative-documents).

## Questions

- **Not sure if a change belongs?** Open an issue and ask.
- **Stuck on the setup?** Check [docs/troubleshooting/](docs/troubleshooting/) or open an issue.
- **Want to propose a big refactor?** Open an issue first — significant changes without prior discussion often need rework.
