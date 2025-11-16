package toyota

import "github.com/ocrosby/golab/mocking"

// Tundra is a truck
type Tundra struct {
	mocking.Truck
	state string
	year  int
}

// NewTundra creates a new Tundra
func NewTundra() *Tundra {
	return &Tundra{
		year:  0,
		state: "parked",
	}
}

// TurnOn turns the truck on
func (t *Tundra) TurnOn() error {
	return nil
}

// TurnOff turns the truck off
func (t *Tundra) TurnOff() error {
	return nil
}

// GetState returns the state of the truck
func (t *Tundra) GetState() string {
	return t.state
}

// SetState sets the state of the truck
func (t *Tundra) SetState(state string) error {
	t.state = state
	return nil
}

// GetWheelCount returns the number of wheels on the truck
func (t *Tundra) GetWheelCount() int {
	return 4
}

// GetMake returns the make of the truck
func (t *Tundra) GetMake() string {
	return "Toyota"
}

// GetModel returns the model of the truck
func (t *Tundra) GetModel() string {
	return "Tundra"
}

// GetYear returns the year of the truck
func (t *Tundra) GetYear() int {
	return t.year
}

// SetYear sets the year of the truck
func (t *Tundra) SetYear(year int) {
	t.year = year
}
