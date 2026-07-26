package generics_test

import (
	"testing"

	generics "github.com/ocrosby/go-lab/lessons/25-generics"
)

func TestMap_TransformsEachElement(t *testing.T) {
	got := generics.Map([]int{1, 2, 3}, func(n int) int { return n * n })
	want := []int{1, 4, 9}

	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestMap_ChangesElementType(t *testing.T) {
	// Map lets the input and output types differ — []int → []string.
	got := generics.Map([]int{1, 2, 3}, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	want := []string{"odd", "even", "odd"}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFilter_KeepsOnlyMatching(t *testing.T) {
	got := generics.Filter([]int{1, 2, 3, 4, 5}, func(n int) bool {
		return n%2 == 0
	})
	if len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Errorf("got %v, want [2 4]", got)
	}
}

func TestReduce_FoldsIntoDifferentType(t *testing.T) {
	// Sum the lengths of a slice of strings into an int. Element type
	// (string) and accumulator type (int) differ.
	total := generics.Reduce([]string{"hi", "there", "friend"}, 0, func(acc int, s string) int {
		return acc + len(s)
	})
	if total != 13 {
		t.Errorf("total = %d, want 13", total)
	}
}

func TestMax_WorksAcrossTypes(t *testing.T) {
	if got := generics.Max(2, 5); got != 5 {
		t.Errorf("Max(2, 5) = %d, want 5", got)
	}
	if got := generics.Max("apple", "banana"); got != "banana" {
		t.Errorf("Max(\"apple\", \"banana\") = %q, want \"banana\"", got)
	}
	if got := generics.Max(3.14, 2.71); got != 3.14 {
		t.Errorf("Max(3.14, 2.71) = %g, want 3.14", got)
	}
}

func TestStack_PushPopIsLIFO(t *testing.T) {
	var s generics.Stack[int]
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if s.Len() != 3 {
		t.Errorf("Len() = %d, want 3", s.Len())
	}

	for _, want := range []int{3, 2, 1} {
		got, ok := s.Pop()
		if !ok {
			t.Fatal("Pop returned ok=false on non-empty stack")
		}
		if got != want {
			t.Errorf("Pop() = %d, want %d", got, want)
		}
	}
	if _, ok := s.Pop(); ok {
		t.Error("Pop on empty stack returned ok=true")
	}
}
