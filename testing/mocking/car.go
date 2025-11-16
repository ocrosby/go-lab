package mocking

//go:generate mockgen -destination=./mocks/mock_car.go -package=mocks github.com/ocrosby/golab/mocking Car

type Car interface {
	Vehicle
}
