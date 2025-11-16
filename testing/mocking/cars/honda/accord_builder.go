package honda

//go:generate mockgen -destination=./mock_accord_builder.go -package=honda github.com/ocrosby/golab/mocking/cars/honda IAccordBuilder

import "fmt"

// IAccordBuilder is a builder interface for the Accord
type IAccordBuilder interface {
	GetInstance() *Accord
	Build()
	BuildState(state string) error
	BuildYear(year int) error
}

// AccordBuilder is a builder for the Accord
type AccordBuilder struct {
	instance *Accord
}

// NewAccordBuilder creates a new AccordBuilder
func NewAccordBuilder() *AccordBuilder {
	return &AccordBuilder{instance: nil}
}

// GetInstance returns the instance of the Accord
func (builder *AccordBuilder) GetInstance() *Accord {
	return builder.instance
}

// Build creates a new instance of the Accord
func (builder *AccordBuilder) Build() {
	builder.instance = NewAccord()
}

// BuildState sets the state of the Accord
func (builder *AccordBuilder) BuildState(state string) error {
	if builder.instance == nil {
		return fmt.Errorf("instance is nil")
	}

	builder.instance.state = state

	return nil
}

// BuildYear sets the year of the Accord
func (builder *AccordBuilder) BuildYear(year int) error {
	if builder.instance == nil {
		return fmt.Errorf("instance is nil")
	}

	builder.instance.year = year

	return nil
}
