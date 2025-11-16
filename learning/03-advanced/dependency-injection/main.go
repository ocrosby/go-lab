package main

import (
	"github.com/ocrosby/golab/dependency-injection/pkg"
	"github.com/ocrosby/golab/dependency-injection/pkg/safety/placers"
)

func main() {
	pkg.NewRockClimber(placers.ConcreteSafetyPlacer{}).ClumbRock()
	pkg.NewRockClimber(placers.NOPSafetyPlacer{}).ClumbRock()
	pkg.NewRockClimber(placers.RockSafetyPlacer{}).ClumbRock()
	pkg.NewRockClimber(placers.IceSafetyPlacer{}).ClumbRock()
}
