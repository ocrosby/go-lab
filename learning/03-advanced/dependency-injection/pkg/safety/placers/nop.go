package placers

import "fmt"

type NOPSafetyPlacer struct{}

func (sp NOPSafetyPlacer) PlaceSafeties() {
	fmt.Println("placing no safeties")
}
