package funcs_test

import (
	"errors"
	"testing"

	funcs "github.com/ocrosby/go-lab/lessons/27-functions-and-closures"
)

func TestInit_RanBeforeMain(t *testing.T) {
	// If init() ran, initMessage was set. This test wouldn't reach the
	// point of running if init hadn't completed, so a passing assertion
	// here is a positive confirmation of init timing.
	if got := funcs.InitMessage(); got != "package initialized" {
		t.Errorf("InitMessage() = %q, want %q", got, "package initialized")
	}
}

func TestMakeAdder_ClosesOverN(t *testing.T) {
	addFive := funcs.MakeAdder(5)
	addTen := funcs.MakeAdder(10)

	if got := addFive(3); got != 8 {
		t.Errorf("addFive(3) = %d, want 8", got)
	}
	if got := addTen(3); got != 13 {
		t.Errorf("addTen(3) = %d, want 13", got)
	}
	// The two closures don't share state — each has its own captured n.
	if got := addFive(0); got != 5 {
		t.Errorf("addFive(0) = %d, want 5 (state leaked from addTen)", got)
	}
}

func TestMakeCounter_HasIndependentState(t *testing.T) {
	c1 := funcs.MakeCounter()
	c2 := funcs.MakeCounter()

	c1()
	c1()
	c1()
	c2()

	if got := c1(); got != 4 {
		t.Errorf("c1 next = %d, want 4", got)
	}
	if got := c2(); got != 2 {
		t.Errorf("c2 next = %d, want 2 (state leaked from c1)", got)
	}
}

func TestSum_AcceptsVariadic(t *testing.T) {
	if got := funcs.Sum(); got != 0 {
		t.Errorf("Sum() = %d, want 0", got)
	}
	if got := funcs.Sum(1, 2, 3, 4); got != 10 {
		t.Errorf("Sum(1,2,3,4) = %d, want 10", got)
	}
}

func TestSum_AcceptsSliceSpread(t *testing.T) {
	nums := []int{10, 20, 30}
	if got := funcs.Sum(nums...); got != 60 {
		t.Errorf("Sum(nums...) = %d, want 60", got)
	}
}

func TestDivide_NamedReturnsWithNakedReturn(t *testing.T) {
	q, err := funcs.Divide(10, 2)
	if err != nil || q != 5 {
		t.Errorf("Divide(10,2) = (%d, %v), want (5, nil)", q, err)
	}

	q, err = funcs.Divide(10, 0)
	if !errors.Is(err, err) || err == nil {
		t.Errorf("Divide(10,0) err = %v, want non-nil", err)
	}
	if q != 0 {
		t.Errorf("Divide(10,0) q = %d, want 0 (zero value on error)", q)
	}
}

func TestApply_UsesFunctionValueAsParameter(t *testing.T) {
	got := funcs.Apply([]int{1, 2, 3}, func(x int) int { return x * x })
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

func TestDispatch_FindsRegisteredOperator(t *testing.T) {
	result, ok := funcs.Dispatch("+", 3, 4)
	if !ok {
		t.Fatal("Dispatch(\"+\") reported ok=false")
	}
	if result != 7 {
		t.Errorf("Dispatch(+, 3, 4) = %d, want 7", result)
	}

	_, ok = funcs.Dispatch("^", 1, 2)
	if ok {
		t.Error("unregistered operator reported ok=true")
	}
}
