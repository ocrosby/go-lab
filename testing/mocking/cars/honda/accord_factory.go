package honda

//go:generate mockgen -destination=./mock_accord_factory.go -package=honda github.com/ocrosby/golab/mocking/cars/honda IAccordFactory

// IAccordFactory is a factory interface for the Accord
type IAccordFactory interface {
	Create() *Accord
	CreateWithState(state string) (*Accord, error)
	CreateWithYear(year int) (*Accord, error)
	CreateWithStateAndYear(state string, year int) (*Accord, error)
}

// AccordFactory is a factory for the Accord
type AccordFactory struct {
	builder IAccordBuilder
}

// NewAccordFactory creates a new AccordFactory
func NewAccordFactory(builder IAccordBuilder) *AccordFactory {
	return &AccordFactory{builder: builder}
}

// Create creates a new instance of the Accord
func (factory *AccordFactory) Create() *Accord {
	factory.builder.Build()

	return factory.builder.GetInstance()
}

// CreateWithState creates a new instance of the Accord with a state
func (factory *AccordFactory) CreateWithState(state string) (*Accord, error) {
	factory.builder.Build()
	err := factory.builder.BuildState(state)
	if err != nil {
		return nil, err
	}

	accord := factory.builder.GetInstance()

	return accord, nil
}

// CreateWithYear creates a new instance of the Accord with a year
func (factory *AccordFactory) CreateWithYear(year int) (*Accord, error) {
	factory.builder.Build()
	err := factory.builder.BuildYear(year)
	if err != nil {
		return nil, err
	}

	accord := factory.builder.GetInstance()

	return accord, nil
}

// CreateWithStateAndYear creates a new instance of the Accord with a state and year
func (factory *AccordFactory) CreateWithStateAndYear(state string, year int) (*Accord, error) {
	factory.builder.Build()
	err := factory.builder.BuildState(state)
	if err != nil {
		return nil, err
	}
	err = factory.builder.BuildYear(year)
	if err != nil {
		return nil, err
	}

	accord := factory.builder.GetInstance()

	return accord, nil
}
