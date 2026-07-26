# Getting started with Go Laboratory

A step-by-step guide from "no Go installed" to "you've run your first Go program and understood what happened." No prior programming experience assumed.

## Prerequisites

- A computer running macOS, Linux, or Windows.
- The ability to open a terminal (see [Step 1](#step-1-open-a-terminal)).
- A text editor. Notepad, TextEdit, VS Code, or anything you already have works — Go doesn't require a specific one.

## Step 1: Open a terminal

The terminal is a text-only way to talk to your computer. Every command below is typed into it.

- **macOS**: press ⌘+Space, type "Terminal", press Enter.
- **Windows**: press the Windows key, type "PowerShell", press Enter. (Command Prompt also works.)
- **Linux**: your distro's default terminal — usually Ctrl+Alt+T.

You should see a prompt ending in `$` or `%` or `>` — that means it's ready. Type `echo hello` and press Enter; it should print `hello` on the next line. If it did, you're set.

## Step 2: Install Go

### Option A: Official installer (works on every OS)

1. Visit [go.dev/dl](https://go.dev/dl/).
2. Download the installer for your operating system.
3. Run it and accept the defaults.

### Option B: Package manager (faster if you already have one)

**macOS (Homebrew):**
```bash
brew install go
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install golang-go
```

**Windows (Chocolatey):**
```bash
choco install golang
```

### Verify the install

In your terminal, type:

```bash
go version
```

Expected output (your patch number may differ):

```text
go version go1.26.5 darwin/amd64
```

If you see a version number, Go is installed. If you see "command not found", close and reopen your terminal — some installers only take effect in new terminal windows.

## Step 3: Install git

Git is how you download this repository. Check first:

```bash
git --version
```

If that prints a version, skip to Step 4. Otherwise:

- **macOS**: `git --version` will offer to install the Xcode Command Line Tools. Accept.
- **Windows**: download and install [Git for Windows](https://gitforwindows.org/).
- **Linux (Debian/Ubuntu)**: `sudo apt install git`.

## Step 4: Clone this repository

```bash
git clone https://github.com/ocrosby/go-lab.git
cd go-lab
go mod download
```

The first command downloads the repo. The second moves into its folder. The third fetches the external Go packages this repo uses — you only need to run it once.

## Step 5: Run your first Go program

From inside the `go-lab` folder:

```bash
go run ./lessons/01-hello
```

Expected output:

```text
Hello World!
```

If you saw that, congratulations — you just ran your first Go program. Open [`lessons/01-hello/README.md`](../../lessons/01-hello/README.md) to see what the seven lines of code actually do.

## Step 6: Run the tests

```bash
make test
```

You should see a wall of green `ok` lines. That means every lesson's tests pass on your machine — a nice sanity check before you dive in.

If you don't have `make` (Windows without WSL, for example), the equivalent is:

```bash
go list ./... | grep -v "10-panic-and-recover/.*/before$" | xargs go test
```

## Step 7: Follow the syllabus

The 18-lesson syllabus lives in the [root README's Syllabus section](../../README.md#syllabus). Work through the lessons in order — each folder under [`lessons/`](../../lessons/) is a self-contained mini-lesson with its own README, runnable code, and small "Try it yourself" exercises.

The first three lessons are the fastest way in:

1. [`lessons/01-hello`](../../lessons/01-hello/) — anatomy of a Go program.
2. [`lessons/02-functions-and-packages`](../../lessons/02-functions-and-packages/) — functions, packages, the "capital letter = public" rule.
3. [`lessons/03-testing-basics`](../../lessons/03-testing-basics/) — `go test` and table-driven tests.

## Step 8: Set up your editor (optional but recommended)

An editor with Go support gives you autocomplete, jump-to-definition, and inline error messages — worth the ten minutes of setup.

- **VS Code**: install the [Go extension by Google](https://marketplace.visualstudio.com/items?itemName=golang.Go). It will offer to install a handful of helper tools; accept them all.
- **GoLand** (JetBrains): Go support is built in.
- **Neovim / Vim**: install [`gopls`](https://pkg.go.dev/golang.org/x/tools/gopls) as your language server; every popular Vim plugin manager has a Go integration.

## Common commands (bookmark this)

### Running code

```bash
go run ./lessons/01-hello   # run a program in a folder
go run some_file.go         # run a single file
go build ./...              # compile everything
```

### Testing

```bash
go test ./lessons/03-testing-basics   # run tests in one lesson
go test -v ./...                      # verbose, all packages
go test -race ./...                   # with the race detector
go test -bench=. ./lessons/15-benchmarks   # run benchmarks
```

### Modules

```bash
go mod tidy       # add missing / remove unused dependencies
go mod download   # download declared dependencies
```

## Troubleshooting

### `go: command not found`

- Close and reopen your terminal — installers usually only take effect in new sessions.
- Check `which go` (macOS/Linux) or `where go` (Windows). If empty, Go isn't on your PATH.
- Reinstall from [go.dev/dl](https://go.dev/dl/) and follow the OS-specific installation notes.

### Tests fail on a fresh clone

- Try `go mod download` again.
- Two `lessons/10-panic-and-recover/*/before/` packages intentionally demonstrate panics — use `make test` (which skips them) rather than raw `go test ./...`.

### Import errors when running lesson files

- Run commands from the repo root (the folder containing `go.mod`), not from inside a lesson folder.
- If you `cd` into a lesson, `go run .` (with a dot) will still work — Go finds the current package.

## Getting help

- [Go documentation](https://go.dev/doc/) — the official reference.
- [Effective Go](https://go.dev/doc/effective_go) — one afternoon read that pays for itself.
- [Go by Example](https://gobyexample.com/) — bite-sized snippets.
- [r/golang](https://www.reddit.com/r/golang/) and the [Gophers Slack](https://gophers.slack.com/) — active communities.

## Next

Head to [`lessons/01-hello`](../../lessons/01-hello/) and follow the syllabus. Each lesson's README ends with a "You've understood this lesson when..." checklist so you know when you're ready for the next one.
