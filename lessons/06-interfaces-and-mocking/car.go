package mocking

//go:generate mockgen -destination=./mocks/mock_car.go -package=mocks github.com/ocrosby/go-lab/lessons/06-interfaces-and-mocking Car

type Car interface {
	Vehicle
}
