# Benchmarks

Measuring the performance of Go code with the built-in benchmark tool. Two small examples show, in numbers, why the standard library nudges you toward certain patterns.

## Why it matters

"This is faster" is worth roughly nothing without a measurement. Go's `testing` package includes benchmarks that iterate a piece of code enough times to produce statistically-useful numbers (nanoseconds per operation, bytes allocated, allocations per operation). Once you can measure, you can defend a choice ("Builder is 2× faster than +=") or catch a regression before shipping it.

This lesson is deliberately small. The [`docs/benchmarking.md`](../../docs/benchmarking.md) deep-dive covers profiling, `benchstat`, `pprof`, and production performance work in depth — come back for that when the basics feel comfortable.

## Prerequisites

- Lesson 03: `go test` basics.
- Enough Go to read a `for` loop and a slice.

## Run it

```bash
go test -bench=. -benchmem ./lessons/15-benchmarks
```

Expected shape of the output (numbers will differ on your machine):

```text
BenchmarkSliceGrow-24        	   28730	     45780 ns/op	  357626 B/op	      19 allocs/op
BenchmarkSlicePrealloc-24    	  126596	      9477 ns/op	   81920 B/op	       1 allocs/op
BenchmarkConcatPlus-24       	 6576954	       180.9 ns/op	     208 B/op	       8 allocs/op
BenchmarkConcatBuilder-24    	12717758	        93.58 ns/op	     120 B/op	       4 allocs/op
BenchmarkConcatJoin-24       	18135825	        66.66 ns/op	      48 B/op	       1 allocs/op
```

Read left to right:

1. Benchmark name (with the number of CPU cores after the `-`).
2. Iterations Go ran — chosen automatically to reach a stable measurement.
3. **`ns/op`** — nanoseconds per operation. Lower is better.
4. **`B/op`** — bytes allocated per operation. Lower is better.
5. **`allocs/op`** — number of heap allocations per operation. Fewer is better.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`strings_bench_test.go`](./strings_bench_test.go) | Three ways to concatenate strings: `+=`, `strings.Builder`, `strings.Join`. |
| [`slice_bench_test.go`](./slice_bench_test.go) | Growing a slice with `append` versus pre-allocating capacity. |

## Anatomy of a benchmark

```go
func BenchmarkConcatJoin(b *testing.B) {
    b.ReportAllocs()               // include B/op and allocs/op in output
    for i := 0; i < b.N; i++ {     // Go picks b.N to reach a stable time
        _ = concatJoin(words)
    }
}
```

Rules:

1. **Filename ends in `_test.go`** — same rule as regular tests.
2. **Function name starts with `Benchmark`** (not `Test`) and takes `*testing.B`.
3. **The measured work goes in a loop** from `0` to `b.N`. Go runs the loop long enough to hit a stable measurement, then reports per-iteration numbers.
4. **`b.ReportAllocs()`** turns on memory reporting. Almost always worth calling.
5. **Assign to `_`** to prevent the compiler from optimising the whole call away.

## Try it yourself

1. Add a `BenchmarkConcatFprintf` that uses `fmt.Sprintf("%s%s%s...")` for the same three strings. Compare its numbers to `ConcatJoin`.
2. Change `N` in `slice_bench_test.go` from `10_000` to `1_000_000`. Re-run. How does the gap between `Grow` and `Prealloc` change?
3. Run the benchmark five times with `-count=5`. Are the numbers consistent, or noisy? What happens if you run under load (open Chrome, play a video)?

## Common pitfalls

- **Measuring setup time.** If your benchmark has expensive setup, call `b.ResetTimer()` right after it — otherwise the setup counts toward every iteration.
- **Optimizing an already-fast function.** `ns/op` in the single digits is often noise. Focus on the bottleneck your profile shows, not on the operation that already runs in 3 ns.
- **Benchmarking without `-count=N`.** Single runs are noisy. Use `-count=5` or `-count=10` and check whether the results move.
- **Forgetting `_ =`.** The Go compiler can, and will, eliminate calls whose results are unused. `_ = f()` keeps the call live.

## You've understood this lesson when...

- You can write a new `Benchmark*` function without looking one up.
- You can read the four columns of benchmark output and explain each.
- You can name one situation where `strings.Builder` is the wrong choice and `+=` is fine. (Hint: how many strings?)

## Related deep-dive

- [`docs/benchmarking.md`](../../docs/benchmarking.md) — profiling with `pprof`, comparing runs with `benchstat`, parallel benchmarks (`b.RunParallel`), production performance monitoring, and load-testing tools.
