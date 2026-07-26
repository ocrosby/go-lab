package honda

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Accord", func() {
	Describe("NewAccord", func() {
		It("should return a new Accord", func() {
			// Act
			accord := NewAccord()

			// Assert
			Expect(accord).NotTo(BeNil())
			Expect(accord.state).To(Equal("parked"))
			Expect(accord.year).To(Equal(0))
		})
	})

	Describe("GetWheelCount", func() {
		It("should return 4", func() {
			// Arrange
			accord := NewAccord()

			// Act
			wheelCount := accord.GetWheelCount()

			// Assert
			Expect(wheelCount).To(Equal(4))
		})
	})

	Describe("TurnOff", func() {
		It("should turn the car off", func() {
			// Arrange
			accord := NewAccord()
			accord.state = "on"

			// Act
			err := accord.TurnOff()

			// Assert
			Expect(err).To(BeNil())
			Expect(accord.state).To(Equal("off"))
		})

		Context("when the car is already off", func() {
			It("should return an error", func() {
				// Arrange
				accord := NewAccord()
				accord.state = "off"

				// Act
				err := accord.TurnOff()

				// Assert
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("car already off"))
			})
		})
	})

	Describe("TurnOn", func() {
		It("should turn the car on", func() {
			// Arrange
			accord := NewAccord()

			// Act
			err := accord.TurnOn()

			// Assert
			Expect(err).To(BeNil())
			Expect(accord.state).To(Equal("on"))
		})

		Context("when the car is already on", func() {
			It("should return an error", func() {
				// Arrange
				accord := NewAccord()
				accord.state = "on"

				// Act
				err := accord.TurnOn()

				// Assert
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("car already on"))
			})
		})
	})

	Describe("GetMake", func() {
		It("should return Honda", func() {
			// Arrange
			accord := NewAccord()

			// Act
			make := accord.GetMake()

			// Assert
			Expect(make).To(Equal("Honda"))
		})
	})

	Describe("GetModel", func() {
		It("should return Accord", func() {
			// Arrange
			accord := NewAccord()

			// Act
			model := accord.GetModel()

			// Assert
			Expect(model).To(Equal("Accord"))
		})
	})

	Describe("GetYear", func() {
		It("should return 0", func() {
			// Arrange
			accord := NewAccord()

			// Act
			year := accord.GetYear()

			// Assert
			Expect(year).To(Equal(0))
		})

		It("should return 2020", func() {
			// Arrange
			accord := NewAccord()
			accord.SetYear(2020)

			// Act
			year := accord.GetYear()

			// Assert
			Expect(year).To(Equal(2020))
		})
	})

	Describe("GetState", func() {
		It("should return parked by default", func() {
			// Arrange
			accord := NewAccord()

			// Act
			state := accord.GetState()

			// Assert
			Expect(state).To(Equal("parked"))
		})

		It("should return on when the state is set to on", func() {
			// Arrange
			accord := NewAccord()
			accord.state = "on"

			// Act
			state := accord.GetState()

			// Assert
			Expect(state).To(Equal("on"))
		})

		It("should return off when the state is set to off", func() {
			// Arrange
			accord := NewAccord()
			accord.state = "off"

			// Act
			state := accord.GetState()

			// Assert
			Expect(state).To(Equal("off"))
		})
	})

	Describe("SetState", func() {
		It("should set the state", func() {
			// Arrange
			accord := NewAccord()

			// Act
			err := accord.SetState("on")

			// Assert
			Expect(err).To(BeNil())
			Expect(accord.state).To(Equal("on"))
		})
	})
})
