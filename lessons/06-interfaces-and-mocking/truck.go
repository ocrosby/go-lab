package mocking

// Truck is any Vehicle with a bed. Same-team; use real implementations
// (e.g. toyota.Tundra) in tests, not a generated mock — same rationale
// as Car above.
type Truck interface {
	Vehicle
}
