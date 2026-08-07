package honda_test

import (
	"errors"
	"testing"

	mocking "github.com/ocrosby/go-lab/lessons/07-interfaces-and-mocking"
	"github.com/ocrosby/go-lab/lessons/07-interfaces-and-mocking/cars/honda"
)

// These tests are black-box: they exercise Accord through its exported
// methods only, and assert on the values those methods return. No test
// reaches into unexported fields (accord.state = "on" was the old style
// and coupled every test to the struct's internal layout).
//
// For the Accord ⇄ TelemetryReporter edge — a genuine system boundary —
// the tests use a small in-memory recorder (below) and assert on the
// recorded interactions. Recording-and-asserting on an edge contract IS
// black-box: it verifies the wire behaviour the external service depends
// on. Compare with mocking the AccordBuilder in accord_factory_test.go,
// which would pin an internal collaboration.

// recordingReporter is a fake TelemetryReporter that stores each ReportIgnition
// call for later inspection. Prefer a fake to a mock when the assertion is
// naturally state-based ("we sent two events with these values") rather than
// order-based ("we called this then that").
type recordingReporter struct {
	events []ignitionEvent
	err    error // if set, returned from ReportIgnition
}

type ignitionEvent struct {
	vin string
	on  bool
}

func (r *recordingReporter) ReportIgnition(vin string, on bool) error {
	r.events = append(r.events, ignitionEvent{vin: vin, on: on})
	return r.err
}

func newAccord() *honda.Accord {
	return honda.NewAccord()
}

func TestNewAccord_StartsParked(t *testing.T) {
	a := newAccord()
	if a == nil {
		t.Fatal("NewAccord returned nil")
	}
	if got := a.GetState(); got != "parked" {
		t.Errorf("state = %q, want %q", got, "parked")
	}
	if got := a.GetYear(); got != 0 {
		t.Errorf("year = %d, want 0", got)
	}
}

func TestAccord_MakeAndModelAreConstant(t *testing.T) {
	a := newAccord()
	if got := a.GetMake(); got != "Honda" {
		t.Errorf("make = %q, want Honda", got)
	}
	if got := a.GetModel(); got != "Accord" {
		t.Errorf("model = %q, want Accord", got)
	}
	if got := a.GetWheelCount(); got != 4 {
		t.Errorf("wheels = %d, want 4", got)
	}
}

func TestSetYear_ReflectsInGetYear(t *testing.T) {
	a := newAccord()
	a.SetYear(2020)
	if got := a.GetYear(); got != 2020 {
		t.Errorf("year = %d, want 2020", got)
	}
}

func TestTurnOn_ChangesStateToOn(t *testing.T) {
	a := newAccord()
	if err := a.TurnOn(); err != nil {
		t.Fatalf("TurnOn err = %v", err)
	}
	if got := a.GetState(); got != "on" {
		t.Errorf("state after TurnOn = %q, want on", got)
	}
}

func TestTurnOn_ErrorsWhenAlreadyOn(t *testing.T) {
	a := newAccord()
	_ = a.TurnOn()

	err := a.TurnOn()
	if err == nil {
		t.Fatal("second TurnOn returned nil, want error")
	}
	if err.Error() != "car already on" {
		t.Errorf("err = %q, want %q", err.Error(), "car already on")
	}
}

func TestTurnOff_ErrorsWhenAlreadyOff(t *testing.T) {
	a := newAccord()
	_ = a.TurnOn()  // parked → on
	_ = a.TurnOff() // on → off

	err := a.TurnOff()
	if err == nil || err.Error() != "car already off" {
		t.Errorf("err = %v, want 'car already off'", err)
	}
}

func TestSetState_OverridesStateMachine(t *testing.T) {
	a := newAccord()
	if err := a.SetState("cruising"); err != nil {
		t.Fatalf("SetState err = %v", err)
	}
	if got := a.GetState(); got != "cruising" {
		t.Errorf("state = %q, want cruising", got)
	}
}

// Edge-contract test: TurnOn reports the ignition event to the telemetry
// service with the VIN and the "on" flag set. This IS a call-verification
// test, and it's appropriate here because TelemetryReporter is a genuine
// external-system boundary — the interaction itself is the contract.
func TestTurnOn_ReportsIgnitionToTelemetry(t *testing.T) {
	reporter := &recordingReporter{}
	a := honda.NewAccordWithTelemetry("VIN123", reporter)

	if err := a.TurnOn(); err != nil {
		t.Fatalf("TurnOn err = %v", err)
	}

	if len(reporter.events) != 1 {
		t.Fatalf("emitted %d events, want 1", len(reporter.events))
	}
	if got := reporter.events[0]; got.vin != "VIN123" || !got.on {
		t.Errorf("event = %+v, want {vin: VIN123, on: true}", got)
	}
}

func TestTurnOff_ReportsIgnitionOff(t *testing.T) {
	reporter := &recordingReporter{}
	a := honda.NewAccordWithTelemetry("VIN123", reporter)
	_ = a.TurnOn()

	if err := a.TurnOff(); err != nil {
		t.Fatalf("TurnOff err = %v", err)
	}

	if len(reporter.events) != 2 {
		t.Fatalf("emitted %d events, want 2", len(reporter.events))
	}
	if got := reporter.events[1]; got.vin != "VIN123" || got.on {
		t.Errorf("event = %+v, want {vin: VIN123, on: false}", got)
	}
}

func TestTurnOn_ReturnsTelemetryError(t *testing.T) {
	upstreamErr := errors.New("telemetry offline")
	reporter := &recordingReporter{err: upstreamErr}
	a := honda.NewAccordWithTelemetry("VIN123", reporter)

	err := a.TurnOn()

	if !errors.Is(err, upstreamErr) {
		t.Errorf("err = %v, want telemetry error to surface", err)
	}
	// The state STILL changed — this repo's design chooses "state change
	// commits before telemetry emits" and returns the telemetry error to
	// the caller. Assert the observable state to lock that design in.
	if got := a.GetState(); got != "on" {
		t.Errorf("state = %q, want on (state should commit before telemetry)", got)
	}
}

// Compile-time check: *honda.Accord satisfies mocking.Car. This is a
// structural assertion (no test-time cost); if the interface grows a method
// the compile breaks here rather than at some caller far away.
var _ mocking.Car = (*honda.Accord)(nil)
