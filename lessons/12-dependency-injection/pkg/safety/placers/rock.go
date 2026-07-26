package placers

import "fmt"

type RockSafetyPlacer struct{}

func (sp RockSafetyPlacer) PlaceSafeties() {
	fmt.Println("placing my Rock safeties")
}
