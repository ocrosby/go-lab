package calculator_test

import (
	"github.com/ocrosby/golab/calculator/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

var _ = Describe("Calculator", func() {
	Describe("Add", func() {
		It("Adds two numbers together", func() {
			for _, set := range numberSet {
				Expect(calculator.Add(set.x, set.y)).To(Equal(set.result))
			}
		})
	})
})
