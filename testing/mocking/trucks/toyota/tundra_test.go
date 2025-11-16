package toyota

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tundra", func() {
	Describe("GetWheelCount", func() {
		It("should return 4", func() {
			// Arrange
			tundra := NewTundra()

			// Act
			wheelCount := tundra.GetWheelCount()

			// Assert
			Expect(wheelCount).To(Equal(4))
		})
	})

	Describe("GetMake", func() {
		It("should return Toyota", func() {
			// Arrange
			tundra := NewTundra()

			// Act
			make := tundra.GetMake()

			// Assert
			Expect(make).To(Equal("Toyota"))
		})
	})

	Describe("GetModel", func() {
		It("should return Tundra", func() {
			// Arrange
			tundra := NewTundra()

			// Act
			model := tundra.GetModel()

			// Assert
			Expect(model).To(Equal("Tundra"))
		})
	})

	Describe("GetYear", func() {
		It("should return 0 by default", func() {
			// Arrange
			tundra := NewTundra()

			// Act
			year := tundra.GetYear()

			// Assert
			Expect(year).To(Equal(0))
		})

		It("should return 2021 when set", func() {
			// Arrange
			tundra := NewTundra()
			tundra.year = 2021

			// Act
			year := tundra.GetYear()

			// Assert
			Expect(year).To(Equal(2021))
		})
	})

	Describe("SetYear", func() {
		It("should set the year", func() {
			// Arrange
			tundra := NewTundra()

			// Act
			tundra.SetYear(2021)

			// Assert
			Expect(tundra.year).To(Equal(2021))
		})
	})
})
