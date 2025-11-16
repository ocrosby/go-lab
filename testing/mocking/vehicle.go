package mocking

// Vehicle is an interface that defines the methods that all vehicles must implement.
type Vehicle interface {
	TurnOn() error
	TurnOff() error
	GetState() string
	SetState(string) error
	GetWheelCount() int
	GetMake() string
	GetModel() string
	GetYear() int
	SetYear(int)
}
