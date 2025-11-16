package honda

import (
	"fmt"
	"github.com/ocrosby/golab/mocking"
)

// Accord is a car
type Accord struct {
	mocking.Car
	state string
	year  int
}

// NewAccord creates a new Accord
func NewAccord() *Accord {
	return &Accord{
		year:  0,
		state: "parked",
	}
}

// TurnOn turns the car on
func (a *Accord) TurnOn() error {
	if a.state == "on" {
		return fmt.Errorf("car already on")
	}

	a.state = "on"

	return nil
}

// TurnOff turns the car off
func (a *Accord) TurnOff() error {
	if a.state == "off" {
		return fmt.Errorf("car already off")
	}

	a.state = "off"

	return nil
}

// GetState returns the state of the car
func (t *Accord) GetState() string {
	return t.state
}

// SetState sets the state of the car
func (t *Accord) SetState(state string) error {
	t.state = state
	return nil
}

// GetWheelCount returns the number of wheels on the car
func (t *Accord) GetWheelCount() int {
	return 4
}

// GetMake returns the make of the car
func (t *Accord) GetMake() string {
	return "Honda"
}

// GetModel returns the model of the car
func (t *Accord) GetModel() string {
	return "Accord"
}

// GetYear returns the year of the car
func (t *Accord) GetYear() int {
	return t.year
}

// SetYear sets the year of the car
func (t *Accord) SetYear(year int) {
	t.year = year
}
