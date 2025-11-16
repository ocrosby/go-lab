package honda

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("AccordFactory", Ordered, func() {
	var mockCtrl *gomock.Controller
	var mockedBuilder *MockIAccordBuilder
	var factory *AccordFactory

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockedBuilder = NewMockIAccordBuilder(mockCtrl)
		factory = NewAccordFactory(mockedBuilder)
	})

	AfterEach(func() {
		factory = nil
		mockedBuilder = nil
		mockCtrl.Finish()
		mockCtrl = nil
	})

	Describe("NewAccordFactory", func() {
		It("should have a reference to the mocked builder", func() {
			Expect(factory.builder).To(Equal(mockedBuilder))
		})
	})

	Describe("Create", func() {
		It("should create a new Accord", func() {
			// Arrange
			accord := NewAccord()

			mockedBuilder.EXPECT().Build()
			mockedBuilder.EXPECT().GetInstance().Return(accord)

			// Act
			result := factory.Create()

			// Assert
			Expect(result).NotTo(BeNil())
			Expect(result).To(Equal(accord))
		})
	})

	Describe("CreateWithYear", func() {
		It("should create a new Accord with the year", func() {
			// Arrange
			accord := NewAccord()
			accord.year = 2019

			mockedBuilder.EXPECT().Build()
			mockedBuilder.EXPECT().BuildYear(2019).Return(nil)
			mockedBuilder.EXPECT().GetInstance().Return(accord)

			// Act
			accord, err := factory.CreateWithYear(2019)

			// Assert
			Expect(err).To(BeNil())
			Expect(accord).NotTo(BeNil())
			Expect(accord.year).To(Equal(2019))
		})

		It("should return an error when the builder returns an error", func() {
			// Arrange
			accord := NewAccord()

			mockedBuilder.EXPECT().Build()
			mockedBuilder.EXPECT().BuildYear(2019).Return(fmt.Errorf("something blew up"))

			// Act
			accord, err := factory.CreateWithYear(2019)

			// Assert
			Expect(accord).To(BeNil())

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("something blew up"))
		})
	})

	Describe("CreateWithState", func() {
		It("should create a new Accord with the state", func() {
			// Arrange
			accord := NewAccord()
			accord.state = "on"

			mockedBuilder.EXPECT().Build()
			mockedBuilder.EXPECT().BuildState("on").Return(nil)
			mockedBuilder.EXPECT().GetInstance().Return(accord)

			// Act
			accord, err := factory.CreateWithState("on")

			// Assert
			Expect(err).To(BeNil())
			Expect(accord).NotTo(BeNil())
			Expect(accord.state).To(Equal("on"))
		})

		It("should return an error when the builder returns an error", func() {
			// Arrange
			accord := NewAccord()

			mockedBuilder.EXPECT().Build()
			mockedBuilder.EXPECT().BuildState("on").Return(fmt.Errorf("something blew up"))

			// Act
			accord, err := factory.CreateWithState("on")

			// Assert
			Expect(accord).To(BeNil())

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("something blew up"))
		})
	})

	Describe("CreateWithStateAndYear", func() {
		It("should create a new Accord with the state and year", func() {
			// Arrange
			accord := NewAccord()
			accord.state = "on"
			accord.year = 2019

			mockedBuilder.EXPECT().Build()
			mockedBuilder.EXPECT().BuildState("on").Return(nil)
			mockedBuilder.EXPECT().BuildYear(2019).Return(nil)
			mockedBuilder.EXPECT().GetInstance().Return(accord)

			// Act
			accord, err := factory.CreateWithStateAndYear("on", 2019)

			// Assert
			Expect(err).To(BeNil())
			Expect(accord).NotTo(BeNil())
			Expect(accord.state).To(Equal("on"))
			Expect(accord.year).To(Equal(2019))
		})

		It("should return an error when the builder returns an error", func() {
			// Arrange
			accord := NewAccord()

			mockedBuilder.EXPECT().Build()
			mockedBuilder.EXPECT().BuildState("on").Return(fmt.Errorf("something blew up"))

			// Act
			accord, err := factory.CreateWithStateAndYear("on", 2019)

			// Assert
			Expect(accord).To(BeNil())

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("something blew up"))
		})
	})
})
