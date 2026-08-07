package main

import (
	"github.com/ocrosby/go-lab/lessons/13-dependency-injection/pkg"
	"github.com/ocrosby/go-lab/lessons/13-dependency-injection/pkg/safety/placers"
)

func main() {
	pkg.NewRockClimber(placers.ConcreteSafetyPlacer{}).ClumbRock()
	pkg.NewRockClimber(placers.NOPSafetyPlacer{}).ClumbRock()
	pkg.NewRockClimber(placers.RockSafetyPlacer{}).ClumbRock()
	pkg.NewRockClimber(placers.IceSafetyPlacer{}).ClumbRock()
}
