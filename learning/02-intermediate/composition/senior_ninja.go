package pkg

import "fmt"

type SeniorNinjaInterface interface {
	NinjaInterface
}

type SeniorNinja struct {
	ninja Ninja
}

func newSeniorNinja(ninja Ninja) SeniorNinja {
	return SeniorNinja{ninja: ninja}
}

func (sn SeniorNinja) Chill() {
	sn.ninja.Chill()
}

func (sn SeniorNinja) Attack() {
	sn.ninja.Attack()
	fmt.Println("Swinging ninja swords")
	sn.Chill()
}
