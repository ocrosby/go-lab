package toyota_test

import (
	"testing"

	mocking "github.com/ocrosby/go-lab/lessons/07-interfaces-and-mocking"
	"github.com/ocrosby/go-lab/lessons/07-interfaces-and-mocking/trucks/toyota"
)

// Tundra tests are black-box — public methods only, no reach-ins into
// tundra.year or any other unexported field. This is intentional: the
// previous Ginkgo version wrote `tundra.year = 2021` in its Arrange
// sections, which pinned the tests to the struct's internal layout.

func TestTundra_ConstantsMatchMakeAndModel(t *testing.T) {
	tr := toyota.NewTundra()

	if got := tr.GetMake(); got != "Toyota" {
		t.Errorf("make = %q, want Toyota", got)
	}
	if got := tr.GetModel(); got != "Tundra" {
		t.Errorf("model = %q, want Tundra", got)
	}
	if got := tr.GetWheelCount(); got != 4 {
		t.Errorf("wheels = %d, want 4", got)
	}
}

func TestTundra_StartsWithZeroYear(t *testing.T) {
	tr := toyota.NewTundra()
	if got := tr.GetYear(); got != 0 {
		t.Errorf("year = %d, want 0", got)
	}
}

func TestTundra_SetYearIsReflectedInGetYear(t *testing.T) {
	tr := toyota.NewTundra()

	tr.SetYear(2021)

	if got := tr.GetYear(); got != 2021 {
		t.Errorf("year after SetYear(2021) = %d, want 2021", got)
	}
}

// Compile-time check that *toyota.Tundra satisfies mocking.Truck.
var _ mocking.Truck = (*toyota.Tundra)(nil)
