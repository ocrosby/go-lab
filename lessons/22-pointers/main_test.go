package main

import "testing"

func TestAddOneValue_DoesNotMutateCaller(t *testing.T) {
	n := 5
	addOneValue(n)
	if n != 5 {
		t.Errorf("n = %d, want 5 (value receiver mustn't mutate)", n)
	}
}

func TestAddOne_MutatesCallerThroughPointer(t *testing.T) {
	n := 5
	addOne(&n)
	if n != 6 {
		t.Errorf("n = %d, want 6", n)
	}
}

func TestCounter_ValueReceiverLosesMutations(t *testing.T) {
	c := Counter{}
	c.IncrementValue()
	c.IncrementValue()
	if c.n != 0 {
		t.Errorf("c.n = %d, want 0 (value receiver mutations don't stick)", c.n)
	}
}

func TestCounter_PointerReceiverKeepsMutations(t *testing.T) {
	c := Counter{}
	c.Increment()
	c.Increment()
	c.Increment()
	if got := c.Value(); got != 3 {
		t.Errorf("c.Value() = %d, want 3", got)
	}
}

func TestNewAndStructLiteral_ProduceEquivalentPointers(t *testing.T) {
	a := new(Counter) // pointer to zero-value Counter
	b := &Counter{}   // also pointer to zero-value Counter
	if a.Value() != b.Value() {
		t.Errorf("new(Counter) and &Counter{} disagreed: %d vs %d", a.Value(), b.Value())
	}
}
