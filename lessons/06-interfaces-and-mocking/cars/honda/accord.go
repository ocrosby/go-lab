package honda

import (
	"fmt"

	mocking "github.com/ocrosby/go-lab/lessons/06-interfaces-and-mocking"
)

// Accord is a Honda Accord that satisfies mocking.Car.
//
// The optional Telemetry field is the seam this lesson uses to demonstrate
// mocking at a genuine system edge: when set, ignition events are reported
// to an external fleet-tracking service via the TelemetryReporter contract.
// A nil Telemetry means the car is off-network — a valid state, no error.
type Accord struct {
	mocking.Car
	Telemetry mocking.TelemetryReporter

	vin   string
	state string
	year  int
}

// NewAccord creates a new Accord in the "parked" state with no telemetry
// wired up. Use NewAccordWithTelemetry to inject a reporter — nearly always
// what production code wants; tests that don't care about telemetry can
// use NewAccord and skip it.
func NewAccord() *Accord {
	return &Accord{
		vin:   "unset",
		state: "parked",
		year:  0,
	}
}

// NewAccordWithTelemetry creates a new Accord wired to a TelemetryReporter.
// The vin is stamped once at construction time and included in every
// telemetry event.
func NewAccordWithTelemetry(vin string, reporter mocking.TelemetryReporter) *Accord {
	return &Accord{
		Telemetry: reporter,
		vin:       vin,
		state:     "parked",
		year:      0,
	}
}

// TurnOn transitions the car to the "on" state and, if telemetry is wired,
// reports the ignition event to the external service. Telemetry failures
// are returned to the caller so the caller can decide whether to consider
// the ignition itself failed — this repo's tests take the "return the
// telemetry error but the state has already changed" position.
func (a *Accord) TurnOn() error {
	if a.state == "on" {
		return fmt.Errorf("car already on")
	}
	a.state = "on"
	if a.Telemetry != nil {
		return a.Telemetry.ReportIgnition(a.vin, true)
	}
	return nil
}

// TurnOff transitions the car to the "off" state and reports the ignition
// event, symmetrically with TurnOn.
func (a *Accord) TurnOff() error {
	if a.state == "off" {
		return fmt.Errorf("car already off")
	}
	a.state = "off"
	if a.Telemetry != nil {
		return a.Telemetry.ReportIgnition(a.vin, false)
	}
	return nil
}

// VIN returns the vehicle identification number the Accord was constructed
// with. Present so tests can assert on it without reaching into the struct.
func (a *Accord) VIN() string { return a.vin }

// GetState returns the current state ("parked", "on", "off").
func (a *Accord) GetState() string { return a.state }

// SetState sets the state directly. Provided for callers that need to
// bypass the on/off state machine (e.g. loading from persistence).
func (a *Accord) SetState(state string) error {
	a.state = state
	return nil
}

// GetWheelCount returns the number of wheels on the car.
func (a *Accord) GetWheelCount() int { return 4 }

// GetMake returns the make of the car.
func (a *Accord) GetMake() string { return "Honda" }

// GetModel returns the model of the car.
func (a *Accord) GetModel() string { return "Accord" }

// GetYear returns the year of the car.
func (a *Accord) GetYear() int { return a.year }

// SetYear sets the year of the car.
func (a *Accord) SetYear(year int) { a.year = year }
