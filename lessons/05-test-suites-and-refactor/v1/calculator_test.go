package calculator_test

import (
	"testing"

	calculator "github.com/ocrosby/go-lab/lessons/05-test-suites-and-refactor/v1"
)

var numberSet = []struct {
	x      int
	y      int
	result int
}{
	{1, 2, 3},
	{2, 2, 4},
	{3, 3, 6},
}

func TestAdd(t *testing.T) {
	result := calculator.Add(1, 3)
	if result != 4 {
		// t.Fail()
		t.Errorf("Expect 1 + 3 == 4, got %d instead", result)
	}
}

func TestAdds(t *testing.T) {
	for _, set := range numberSet {
		result := calculator.Add(set.x, set.y)

		if result != set.result {
			t.Errorf("Expected %d + %d == %d, got %d instead.", set.x, set.y, set.result, result)
		}
	}
}
