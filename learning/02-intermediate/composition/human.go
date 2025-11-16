package pkg

import "fmt"

type HumanInterface interface {
	Chill()
}

type Human struct{}

func (h Human) Chill() {
	fmt.Println("Chilling")
}
