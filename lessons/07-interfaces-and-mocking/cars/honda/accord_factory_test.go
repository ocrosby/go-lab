package honda_test

import (
	"testing"

	"github.com/ocrosby/go-lab/lessons/07-interfaces-and-mocking/cars/honda"
)

// AccordFactory tests: wired end-to-end with a real *AccordBuilder as its
// collaborator. The previous version of this file mocked the builder and
// asserted mockedBuilder.EXPECT().Build() (behaviour verification against a
// same-team collaborator). That coupled every test to the exact sequence of
// builder calls the factory made — refactoring the factory to, say, cache
// the builder's output would have broken every test even though the
// observable behaviour was identical. Using the real builder makes the
// tests survive any such refactor.
//
// State-based assertions throughout: the tests check what the factory
// returns via the Accord's own public methods, not what it did to the
// builder along the way.

func newFactory() *honda.AccordFactory {
	return honda.NewAccordFactory(honda.NewAccordBuilder())
}

func TestFactory_CreateReturnsParkedAccord(t *testing.T) {
	f := newFactory()

	got := f.Create()

	if got == nil {
		t.Fatal("Create returned nil")
	}
	if state := got.GetState(); state != "parked" {
		t.Errorf("state = %q, want parked", state)
	}
}

func TestFactory_CreateWithStateSetsState(t *testing.T) {
	f := newFactory()

	got, err := f.CreateWithState("on")

	if err != nil {
		t.Fatalf("CreateWithState err = %v", err)
	}
	if state := got.GetState(); state != "on" {
		t.Errorf("state = %q, want on", state)
	}
}

func TestFactory_CreateWithYearSetsYear(t *testing.T) {
	f := newFactory()

	got, err := f.CreateWithYear(2020)

	if err != nil {
		t.Fatalf("CreateWithYear err = %v", err)
	}
	if year := got.GetYear(); year != 2020 {
		t.Errorf("year = %d, want 2020", year)
	}
}

func TestFactory_CreateWithStateAndYearSetsBoth(t *testing.T) {
	f := newFactory()

	got, err := f.CreateWithStateAndYear("on", 2020)

	if err != nil {
		t.Fatalf("CreateWithStateAndYear err = %v", err)
	}
	if state := got.GetState(); state != "on" {
		t.Errorf("state = %q, want on", state)
	}
	if year := got.GetYear(); year != 2020 {
		t.Errorf("year = %d, want 2020", year)
	}
}

// Independence: each factory call returns a fresh Accord. If the builder
// were leaky (sharing state across calls) this would catch it.
func TestFactory_ReturnsIndependentInstances(t *testing.T) {
	f := newFactory()

	first, _ := f.CreateWithYear(2020)
	second, _ := f.CreateWithYear(2024)

	if first.GetYear() != 2020 {
		t.Errorf("first.year = %d, want 2020 (second call clobbered first)", first.GetYear())
	}
	if second.GetYear() != 2024 {
		t.Errorf("second.year = %d, want 2024", second.GetYear())
	}
}
