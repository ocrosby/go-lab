package honda_test

import (
	"testing"

	"github.com/ocrosby/go-lab/lessons/06-interfaces-and-mocking/cars/honda"
)

// AccordBuilder tests: exercised through its public API only. No test reaches
// into builder.instance directly — asserting on private struct state was the
// pattern the /docs/black-box-testing rule warns against, because any
// refactor of the builder's storage breaks tests whose behaviour it did not
// touch.

func TestBuilder_ProducesNilBeforeBuild(t *testing.T) {
	b := honda.NewAccordBuilder()
	if b.GetInstance() != nil {
		t.Error("GetInstance before Build = non-nil, want nil")
	}
}

func TestBuilder_BuildProducesAnAccord(t *testing.T) {
	b := honda.NewAccordBuilder()
	b.Build()

	got := b.GetInstance()
	if got == nil {
		t.Fatal("GetInstance after Build = nil, want *Accord")
	}
	if got.GetState() != "parked" {
		t.Errorf("state = %q, want parked", got.GetState())
	}
}

func TestBuilder_BuildStateSetsState(t *testing.T) {
	b := honda.NewAccordBuilder()
	b.Build()

	if err := b.BuildState("on"); err != nil {
		t.Fatalf("BuildState err = %v", err)
	}

	if got := b.GetInstance().GetState(); got != "on" {
		t.Errorf("state = %q, want on", got)
	}
}

func TestBuilder_BuildStateBeforeBuildErrors(t *testing.T) {
	b := honda.NewAccordBuilder()

	err := b.BuildState("on")

	if err == nil || err.Error() != "instance is nil" {
		t.Errorf("err = %v, want 'instance is nil'", err)
	}
}

func TestBuilder_BuildYearSetsYear(t *testing.T) {
	b := honda.NewAccordBuilder()
	b.Build()

	if err := b.BuildYear(2020); err != nil {
		t.Fatalf("BuildYear err = %v", err)
	}

	if got := b.GetInstance().GetYear(); got != 2020 {
		t.Errorf("year = %d, want 2020", got)
	}
}

func TestBuilder_BuildYearBeforeBuildErrors(t *testing.T) {
	b := honda.NewAccordBuilder()

	err := b.BuildYear(2020)

	if err == nil || err.Error() != "instance is nil" {
		t.Errorf("err = %v, want 'instance is nil'", err)
	}
}

func TestBuilder_FullFluentFlow(t *testing.T) {
	// End-to-end: build a car through the builder and confirm each attribute
	// via the returned Accord's own public methods.
	b := honda.NewAccordBuilder()
	b.Build()
	if err := b.BuildState("on"); err != nil {
		t.Fatalf("BuildState err = %v", err)
	}
	if err := b.BuildYear(2020); err != nil {
		t.Fatalf("BuildYear err = %v", err)
	}

	accord := b.GetInstance()
	if accord == nil {
		t.Fatal("GetInstance = nil")
	}
	if got := accord.GetState(); got != "on" {
		t.Errorf("state = %q, want on", got)
	}
	if got := accord.GetYear(); got != 2020 {
		t.Errorf("year = %d, want 2020", got)
	}
}
