# Go modules

How a Go project is organized — `go.mod`, `go.sum`, `go mod init`, `go mod tidy`, `go get` — and how imports resolve. The plumbing you need to start your own project after the syllabus.

> This lesson is mostly reference material. There's no code to run; there's a step-by-step walkthrough you follow in a fresh directory outside this repo.

## Why it matters

Every Go project (this repo included) starts with `go mod init`. Every dependency addition is `go get`. Every clean-up is `go mod tidy`. You'll spend the rest of your Go life with these commands. Learning them once means never wondering later why `go build` can't find a package or what `go.sum` is protecting you from.

## Prerequisites

- Any of the syllabus lessons — you've been running code inside a Go module the whole time.

## Try it in a new directory

The whole point of this lesson is walking through creating a module from scratch. Do this in a **separate directory** outside `go-lab`:

```bash
mkdir /tmp/hello-mod
cd /tmp/hello-mod

# 1. Initialize a module. The path is the import path the module
#    publishes under. For a hobby project use "example.com/hello" or
#    your github path like "github.com/you/hello".
go mod init example.com/hello

# 2. Write a program that uses an external dependency.
cat > main.go <<'EOF'
package main

import (
    "fmt"
    "github.com/google/uuid"
)

func main() {
    fmt.Println("hello,", uuid.New())
}
EOF

# 3. Fetch dependencies and update go.mod / go.sum.
go mod tidy

# 4. Run it.
go run .
```

You'll see something like `hello, 3f28b9d0-6b8c-4ac9-b6f1-9c7e88f2b1a4`.

Then inspect what happened:

```bash
cat go.mod   # module path, Go version, one require line
cat go.sum   # cryptographic checksums for every module you depend on
```

## What each command does

| Command | What it does |
|---|---|
| `go mod init <path>` | Creates `go.mod` in the current directory. `<path>` is the import path the module publishes under. |
| `go get <package>` | Adds a dependency. Bumps `go.mod` and downloads to the module cache. |
| `go get <package>@<version>` | Same, pinned to a specific version. |
| `go get -u ./...` | Updates every direct dependency to the latest minor version. |
| `go mod tidy` | Removes dependencies you no longer import; adds any you do. Run before every commit. |
| `go mod download` | Fetches all declared dependencies into the module cache (no `go.mod` changes). |
| `go mod why <package>` | Explains why a package is included — what chain of imports pulled it in. |
| `go mod graph` | Prints the dependency graph. Handy piped through `grep`. |
| `go list -m all` | Lists every module in the build, direct and indirect. |

## `go.mod` anatomy

```go
module github.com/ocrosby/go-lab

go 1.26

require (
    github.com/onsi/ginkgo/v2 v2.12.0
    go.uber.org/mock v0.2.0
)

require (
    github.com/kr/text v0.2.0 // indirect
    // ... more indirects
)
```

Line by line:

- **`module ...`** — the import path this module publishes under. Every file in the tree that says `package foo` is available as `github.com/ocrosby/go-lab/path/to/foo`.
- **`go 1.26`** — the minimum Go language version this module requires. Bump it when you use features from a newer version (like `slices.Concat` from 1.22).
- **`require` blocks** — direct dependencies (with a version) and indirect ones (transitively pulled in — marked `// indirect`). `go mod tidy` maintains the split.

## `go.sum` anatomy

`go.sum` stores a **cryptographic checksum** for every version of every module you depend on. Its job: detect if the code you download is different from the code the person who added the dependency saw. Commit `go.sum`. Don't edit it by hand. If Go complains "checksum mismatch", something is genuinely wrong — someone tampered with the module registry, or your local cache is corrupted (delete `$GOPATH/pkg/mod/cache` and retry).

## Versioning: semver rules

Go modules use semantic versioning: `vMAJOR.MINOR.PATCH`.

- **Patch** (`v1.2.3` → `v1.2.4`) — bug fixes. Backwards compatible.
- **Minor** (`v1.2.3` → `v1.3.0`) — new features. Backwards compatible.
- **Major** (`v1.2.3` → `v2.0.0`) — breaking changes.

For major version 2 and up, the import path itself changes: `github.com/foo/bar/v2`. That's why you see `github.com/onsi/ginkgo/v2` in this repo's `go.mod`.

## Common commands you'll actually run

**Starting a new project:**

```bash
mkdir myproject && cd myproject
go mod init github.com/you/myproject
```

**Adding a dependency:**

```bash
go get github.com/spf13/cobra           # latest
go get github.com/spf13/cobra@v1.8.0    # specific version
go mod tidy                             # clean up
```

**Upgrading:**

```bash
go get -u github.com/spf13/cobra   # to latest minor of the same major
go get github.com/spf13/cobra@v2   # to a specific major (breaking)
go mod tidy
```

**Removing:**

Delete the `import` in your code. Then `go mod tidy` removes the require line.

**Debugging "what pulled this in?":**

```bash
go mod why -m github.com/some/mod
```

## Multi-module workspaces (`go.work`)

If you're working on two modules together (e.g. a library and an app that uses it) and want changes to the library to be visible immediately without publishing, use a `go.work` file:

```bash
mkdir workspace && cd workspace
git clone https://github.com/you/mylib.git
git clone https://github.com/you/myapp.git
go work init ./mylib ./myapp
```

Now inside `myapp/`, `import "github.com/you/mylib"` uses your local `mylib` directory. See [`docs/go-workspaces.md`](../../docs/go-workspaces.md) for the full walkthrough.

This repo uses submodules (`lessons/13-design-patterns/go.mod`, `lessons/14-production-api/go.mod`) — those are separate modules within the same repo, not a workspace, but the same principles apply.

## Vendoring

`go mod vendor` copies every dependency into a local `vendor/` directory. Some organizations require it for reproducible builds without a network. Most projects skip it — the module cache and `go.sum` already give reproducibility. If you see a `vendor/` directory in someone else's project and wonder what it is, that's the answer.

## Try it yourself

1. Do the whole walkthrough above in `/tmp/hello-mod`. Read the resulting `go.mod` and `go.sum`.
2. `go get github.com/onsi/ginkgo/v2` and then `go get github.com/onsi/ginkgo/v2@v2.0.0` — check what `go.mod` shows after each.
3. Delete the `import "github.com/google/uuid"` in main.go and run `go mod tidy`. Where does the require line go?
4. Run `go mod why github.com/stretchr/testify` in this repo. Read the output — it shows the chain of imports.

## Common pitfalls

- **Forgetting `go mod tidy` before committing.** `go build` may work locally because your cache is warm, but CI hits "missing go.sum entry" errors. Always run `tidy` before pushing dependency changes.
- **Editing `go.sum` by hand.** Don't. It's a checksum database; `go mod tidy` maintains it.
- **Committing the module cache.** `$GOPATH/pkg/mod/` is a per-machine cache, not a per-project artifact. Don't commit it.
- **Manually editing `go.mod` version numbers.** Prefer `go get`; it validates the version exists and updates `go.sum` too.
- **Naming a module without a real path.** `go mod init myproject` compiles today but breaks when you try to publish or import from elsewhere. Use a real path like `github.com/you/myproject` even for private code.
- **Major version 2+ without the `/v2` path suffix.** `v2.x.x` modules must have `/v2` in the module path AND in every import. Break this and Go can't resolve the module correctly.

## You've understood this lesson when...

- You can bootstrap a new Go project from scratch and add one external dependency.
- You know what each column in `go.mod` means and what `go.sum` protects you from.
- You know which command to run before every commit (`go mod tidy`) and which to run when a dependency has a new release you want (`go get -u`).
- You can explain why `github.com/onsi/ginkgo/v2` has `/v2` in the path.

## Next

- **Next lesson:** [29-context](../29-context/) — `context.Context` for cancellation, deadlines, and request-scoped values. The last built-in-package concept you need before the concurrency and HTTP lessons make full sense.
