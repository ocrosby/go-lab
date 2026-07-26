package mocking

//go:generate mockgen -destination=./mocks/mock_car.go -package=mocks github.com/ocrosby/go-lab/testing/mocking Car

type Car interface {
	Vehicle
}
