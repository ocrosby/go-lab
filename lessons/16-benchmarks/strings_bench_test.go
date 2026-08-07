package benchmarks

import (
	"strings"
	"testing"
)

// Three ways to concatenate a slice of strings. The benchmarks below let you
// measure the cost of each so you see, in numbers, why the standard library
// nudges you toward strings.Builder and strings.Join.

var words = []string{"the", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}

// Naive: build the result with += in a loop. O(n^2) allocations because each
// += creates a new string.
func concatPlus(items []string) string {
	result := ""
	for _, s := range items {
		result += s
	}
	return result
}

// Builder: strings.Builder amortises to O(n) — it grows its internal buffer
// like a slice, doubling when full.
func concatBuilder(items []string) string {
	var b strings.Builder
	for _, s := range items {
		b.WriteString(s)
	}
	return b.String()
}

// Join: strings.Join walks the slice once, sums lengths, allocates exactly
// once, then copies. Usually the fastest and clearest.
func concatJoin(items []string) string {
	return strings.Join(items, "")
}

func BenchmarkConcatPlus(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = concatPlus(words)
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = concatBuilder(words)
	}
}

func BenchmarkConcatJoin(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = concatJoin(words)
	}
}
