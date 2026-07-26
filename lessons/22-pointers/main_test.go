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

// The next two tests pin the "shallow copy" behavior — the single most
// important thing to internalize about value semantics.

func TestCopy_PrimitiveFieldIsIndependent(t *testing.T) {
	original := Bag{Label: "mine", Items: []string{"a", "b"}}

	_ = CopyAndMutate(original)

	// The copy's Label mutation didn't touch the original — Label is a
	// primitive field, so the copy got its own independent value.
	if original.Label != "mine" {
		t.Errorf("original.Label = %q, want %q (primitive copy should be independent)",
			original.Label, "mine")
	}
}

func TestCopy_SliceFieldSharesBackingArray(t *testing.T) {
	original := Bag{Label: "mine", Items: []string{"a", "b"}}

	_ = CopyAndMutate(original)

	// The copy's Items[0] mutation DID touch the original — slice fields
	// share their backing array between copies. This is the surprising
	// half of Go's value semantics.
	if original.Items[0] != "MUTATED" {
		t.Errorf("original.Items[0] = %q, want %q (slice backing array should be shared)",
			original.Items[0], "MUTATED")
	}
}
