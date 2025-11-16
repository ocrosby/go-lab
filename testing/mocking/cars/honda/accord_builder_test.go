package honda

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccordBuilder", func() {
	var builder *AccordBuilder

	BeforeEach(func() {
		builder = NewAccordBuilder()
	})

	AfterEach(func() {
		builder = nil
	})

	Describe("GetInstance", func() {
		It("should return the instance", func() {
			// Arrange
			builder.instance = NewAccord()

			// Act
			instance := builder.GetInstance()

			// Assert
			Expect(instance).To(Equal(builder.instance))
		})

		Context("when the instance is nil", func() {
			It("should return nil", func() {
				// Arrange
				builder.instance = nil

				// Act
				instance := builder.GetInstance()

				// Assert
				Expect(instance).To(BeNil())
			})
		})
	})

	Describe("Build", func() {
		It("should create an instance", func() {
			// Arrange
			builder.instance = nil

			// Act
			builder.Build()

			// Assert
			Expect(builder.instance).NotTo(BeNil())
		})
	})

	Describe("BuildState", func() {
		Context("when the instance is nil", func() {
			It("should return an error", func() {
				// Act
				err := builder.BuildState("on")

				// Assert
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("instance is nil"))
			})
		})

		Context("when the instance is not nil", func() {
			It("should build the state when given 'on'", func() {
				// Arrange
				builder.instance = NewAccord()

				// Act
				err := builder.BuildState("on")

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(builder.instance.state).To(Equal("on"))
			})
		})
	})

	Describe("BuildYear", func() {
		Context("when the instance is nil", func() {
			It("should return an error", func() {
				// Arrange
				builder := NewAccordBuilder()
				builder.instance = nil

				// Act
				err := builder.BuildYear(2020)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("instance is nil"))
			})
		})

		Context("when the instance is not nil", func() {
			It("should build the year", func() {
				// Arrange
				builder := NewAccordBuilder()
				builder.instance = NewAccord()

				// Act
				err := builder.BuildYear(2020)

				// Assert
				Expect(err).To(BeNil())
				Expect(builder.instance.year).To(Equal(2020))
			})
		})

		Context("when the instance is not nil", func() {
			It("should build the year", func() {
				// Arrange
				builder := NewAccordBuilder()
				builder.Build()

				// Act
				err := builder.BuildYear(2020)

				// Assert
				Expect(err).To(BeNil())
				Expect(builder.instance.year).To(Equal(2020))
			})
		})
	})

	It("should build an Accord", func() {
		// Arrange
		builder := NewAccordBuilder()

		// Act
		builder.Build()
		builder.BuildState("on")
		builder.BuildYear(2020)

		accord := builder.GetInstance()

		// Assert
		Expect(accord).NotTo(BeNil())
		Expect(accord.GetState()).To(Equal("on"))
		Expect(accord.GetYear()).To(Equal(2020))
	})
})
