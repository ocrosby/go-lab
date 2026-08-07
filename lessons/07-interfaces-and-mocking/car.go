package mocking

// Car is any Vehicle that behaves like a car. Since it is a same-team
// interface used inside this deployable, tests should exercise it with real
// implementations (e.g. honda.Accord) rather than a generated mock. Mocking
// same-team collaborators pins the shape of an internal collaboration and
// breaks tests on every refactor — see telemetry.go for the anti-example
// this rule of thumb was written against.
type Car interface {
	Vehicle
}
