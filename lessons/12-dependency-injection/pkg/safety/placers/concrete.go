package placers

import "fmt"

type ConcreteSafetyPlacer struct{}

func (sp ConcreteSafetyPlacer) PlaceSafeties() {
	fmt.Println("placing my Concrete safeties")
}
