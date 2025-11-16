package pkg

import "fmt"

type NinjaInterface interface {
	HumanInterface
	Attack()
}

type Ninja struct {
	human Human
}

func newNinja(human Human) Ninja {
	return Ninja{human: human}
}

func (n Ninja) Chill() {
	n.human.Chill()
}

func (n Ninja) Attack() {
	fmt.Println("Throwing ninja stars")
}
