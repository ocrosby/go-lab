package main

import "testing"

// These tests pin the two most-surprising slice/map behaviours so a reader
// can experiment (change a value, rerun) and see which behaviours are
// contractual vs incidental.

func TestSliceAppend_ReassignmentIsMandatory(t *testing.T) {
	// This test documents the "always reassign append" rule. Without
	// reassignment, appending past capacity would silently drop the value.
	xs := make([]int, 0, 2)
	xs = append(xs, 1)
	xs = append(xs, 2)
	xs = append(xs, 3) // triggers reallocation

	if len(xs) != 3 {
		t.Errorf("len(xs) = %d, want 3", len(xs))
	}
	if got := xs[2]; got != 3 {
		t.Errorf("xs[2] = %d, want 3", got)
	}
}

func TestSlice_SlicingSharesBackingArray(t *testing.T) {
	// Slicing does NOT copy. Mutating the sub-slice mutates the parent.
	// This is the most common slice surprise.
	xs := []int{1, 2, 3, 4, 5}
	ys := xs[1:4] // ys = [2, 3, 4], shares backing with xs

	ys[0] = 999

	if xs[1] != 999 {
		t.Errorf("xs[1] = %d, want 999 (shared backing array)", xs[1])
	}
}

func TestMap_MissingKeyReturnsZeroValue(t *testing.T) {
	m := map[string]int{"answer": 42}

	if got := m["missing"]; got != 0 {
		t.Errorf("m[missing] = %d, want 0 (zero value for int)", got)
	}
}

func TestMap_CommaOkDistinguishesMissingFromZero(t *testing.T) {
	m := map[string]int{"zero-value-stored": 0}

	// A plain lookup can't tell "not there" from "value is 0".
	// The comma-ok form can.
	if _, ok := m["missing"]; ok {
		t.Error("missing key reported ok=true")
	}
	if _, ok := m["zero-value-stored"]; !ok {
		t.Error("stored-zero key reported ok=false")
	}
}

func TestNilMap_ReadOk_WritePanics(t *testing.T) {
	var m map[string]int // nil

	// Reads from a nil map return the zero value — no panic.
	if got := m["anything"]; got != 0 {
		t.Errorf("read from nil map = %d, want 0", got)
	}

	// Writes to a nil map panic. Recover so the test can assert on it.
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected write to nil map to panic")
		}
	}()
	m["oops"] = 1
}
