# Concurrency Examples

End-to-end runnable programs demonstrating Go concurrency patterns.

## Prerequisites

- Completed [07-goroutines-and-channels](../07-goroutines-and-channels/) and [08-channel-patterns](../08-channel-patterns/) or equivalent experience with goroutines, channels, and `select`
- Comfort with `context.Context` cancellation semantics

## Scope

This directory holds **whole programs** that exercise concurrency end-to-end (input → workers → output, with cancellation and error handling). For single-concept snippets that illustrate one primitive at a time (channels, `sync.WaitGroup`, `done` channel, pipeline stages), see the preceding lessons.

| Need | Location |
|---|---|
| Learn one primitive in isolation | [`07-goroutines-and-channels`](../07-goroutines-and-channels/), [`08-channel-patterns`](../08-channel-patterns/) |
| See a full program wire several primitives together | this lesson |

## Examples

_None yet — this section is a placeholder for future programs._

Planned topics:

- Worker pool with bounded concurrency and graceful shutdown
- Fan-out / fan-in pipeline with per-stage cancellation
- Rate limiter using `time.Ticker` and a token bucket
- `errgroup`-coordinated parallel fetches with first-error cancellation

## Layout Convention

Each example lives in its own subdirectory with a runnable `main.go`:

```
concurrency/
└── <example-name>/
    ├── main.go
    └── README.md
```

Run any example with:

```bash
go run ./lessons/09-worker-pools
```

## Related

- [07-goroutines-and-channels](../07-goroutines-and-channels/) and [08-channel-patterns](../08-channel-patterns/) — concept-by-concept primitives
- [15-benchmarks](../15-benchmarks/) — measuring concurrent code
