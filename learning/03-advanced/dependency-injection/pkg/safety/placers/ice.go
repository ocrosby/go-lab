package placers

import "fmt"

type IceSafetyPlacer struct{}

func (sp IceSafetyPlacer) PlaceSafeties() {
	fmt.Println("placing my Ice safeties")
}
