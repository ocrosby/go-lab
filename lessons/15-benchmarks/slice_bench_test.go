package benchmarks

import "testing"

// Two ways to build a slice of a known size. The pre-allocated version avoids
// the repeated grow-and-copy cycle that happens when append() runs past the
// slice's capacity.

const N = 10_000

func buildGrow() []int {
	var out []int
	for i := 0; i < N; i++ {
		out = append(out, i)
	}
	return out
}

func buildPrealloc() []int {
	out := make([]int, 0, N)
	for i := 0; i < N; i++ {
		out = append(out, i)
	}
	return out
}

func BenchmarkSliceGrow(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = buildGrow()
	}
}

func BenchmarkSlicePrealloc(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = buildPrealloc()
	}
}
