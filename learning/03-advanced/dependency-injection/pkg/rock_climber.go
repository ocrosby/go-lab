package pkg

import "github.com/ocrosby/golab/dependency-injection/pkg/safety"

type RockClimber struct {
	rocksClimbed int
	sp           safety.SafetyPlacer
}

func NewRockClimber(sp safety.SafetyPlacer) *RockClimber {
	return &RockClimber{
		sp: sp,
	}
}

// here the climbRock method is depending on a behavior by utilizing the SafetyPlacer above
// this is a form of dependency injection
func (rc *RockClimber) ClumbRock() {
	rc.rocksClimbed++
	rc.sp.PlaceSafeties()
}
