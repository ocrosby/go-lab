# Design patterns

Classic Gang-of-Four design patterns implemented in idiomatic Go. This lesson holds four small, self-contained example programs — one per pattern.

## Why it matters

Design patterns are shared vocabulary. When someone says "I used a Builder for the config" you save a paragraph of explanation. Every pattern here came from object-oriented languages (Java, C++, Smalltalk), and each looks a little different in Go — usually smaller, sometimes simpler, occasionally missing entirely because a language feature obviates it.

The point of this lesson is not to memorise all 23 GoF patterns. It's to see how the four most useful ones translate to Go's smaller feature set, and to build the muscle for translating the rest yourself.

## Prerequisites

- Lesson 05: composition.
- Lesson 06: interfaces satisfied structurally.
- Lesson 12: dependency injection (the same idea shows up here in constructor form).

## Run it

Each pattern is a standalone `main` package. Because this lesson is its own Go module, `cd` into it first — or run the patterns directly:

```bash
go run ./lessons/13-design-patterns/creational/builder
go run ./lessons/13-design-patterns/creational/prototype
go run ./lessons/13-design-patterns/creational/singleton   # if present
go run ./lessons/13-design-patterns/structural/adpater     # note: folder is "adpater" (typo)
```

## What's in this folder

| Pattern | Category | Path |
|---|---|---|
| Builder | Creational | [`creational/builder/`](./creational/builder/) |
| Prototype | Creational | [`creational/prototype/`](./creational/prototype/) |
| Singleton | Creational | [`creational/singleton/`](./creational/singleton/) |
| Adapter | Structural | [`structural/adpater/`](./structural/adpater/) |

Each pattern folder has its own README explaining that specific pattern; open it before reading the code.

## Why this lesson is its own Go module

You'll notice a `go.mod` file in this folder. That normally means "this is a separate Go module." For this lesson it's mostly a legacy of an earlier layout — none of the patterns here need dependencies that would conflict with the root module.

Practical consequence: `go test ./...` at the repo root does *not* descend into this folder. To run this lesson's programs, use the explicit paths shown above, or `cd` into the lesson and use `./...` from there.

## Pattern notes

### Builder

Solves the "constructor with 15 optional parameters" problem. In Go, the pattern often shows up as a `Config` struct with defaults, or as a chain of `WithFoo(x)` setter methods — both are idiomatic Builder variants. The example uses the classic "Director + concrete builders" shape from the GoF book, which is heavier than most real Go code uses.

### Prototype

"Copy an existing object to make a new one." Go's built-in behaviour — passing a struct value copies it — already covers most of this. The example demonstrates a `Clone()` method for deep copies of a tree structure.

### Singleton

"Only one instance ever exists." In Go, use `sync.Once` for thread-safe lazy initialization, or a package-level variable initialised in `init()`. The pattern is often overused; before reaching for it, ask "would injecting the dependency work instead?" (See lesson 12.)

### Adapter

"Wrap a type with a different interface to match a contract you need." Very common in Go — think `sql.DB` wrapping a driver, or `http.Handler` wrapping a function. Small, useful, and lives right at the seam between subsystems.

## Try it yourself

1. In `creational/builder/`, refactor to use functional options (`WithWalls(...)`, `WithRoof(...)`) instead of the Director pattern. Which style reads more like typical Go?
2. In `creational/prototype/`, replace `Clone()` with a plain value copy. What breaks? (Look at whether any fields are pointers or slices.)
3. Sketch a Decorator pattern in Go on paper. Given lesson 06's material, how much code do you need? (Answer: less than you'd think.)

## Common pitfalls

- **Applying patterns because you know them, not because you need them.** Go's small feature set and standard library often solve GoF problems directly — reach for a pattern only when you have the problem it solves.
- **Assuming every GoF pattern maps 1:1 to Go.** Some don't. "Iterator" is just `for range` and channels. "Command" is a first-class function. "Template Method" is often just interface composition.
- **The typo in `structural/adpater/`.** It's in the folder name — a small drift from an older version. Kept for now to avoid breaking any external links; can be renamed in a future PR.

## You've understood this lesson when...

- You can name at least one place each of these patterns appears in the standard library.
- You can convert the Builder example to functional options without looking it up.
- You know one GoF pattern that has no natural equivalent in Go (and why).

## Next

- **Next lesson:** [14-production-api](../14-production-api/) — the patterns used together in a real service with hexagonal architecture.
