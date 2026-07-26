package mocking

//go:generate mockgen -destination=./mocks/mock_truck.go -package=mocks github.com/ocrosby/go-lab/testing/mocking Truck

// Truck is a vehicle with a bed
type Truck interface {
	Vehicle
}
