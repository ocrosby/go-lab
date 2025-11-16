package channels_test

import (
	"github.com/ocrosby/golab/test/channels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type testCase struct {
	name        string
	inputValues []int
	wantSum     int
}

var _ = Describe("Example1", Ordered, func() {
	var tests []testCase

	BeforeAll(func() {
		tests = append(tests, testCase{
			name:        "base case",
			inputValues: []int{0, 1, 2, 3},
			wantSum:     6,
		})
	})

	AfterAll(func() {
		tests = nil
	})

	It("should process order data", func() {
		// Arrange
		order := channels.Order{}
		testChan := make(chan int)
		doneChan := make(chan struct{})
		inputValues := []int{0, 1, 2, 3}

		// Act
		// execute my function on a separate goroutine, so that I can move on
		// and send my data to the input channel
		go order.ProcessData(testChan, doneChan)

		// loop to send data to the input channel
		// if we created this loop before executing the goroutine, we would
		// be stuck trying to write the second value to the channel
		for _, v := range inputValues {
			testChan <- v
		}

		// close channel
		// we must do this here because this is the condition that will cause the
		// for loop inside of order.ProcessData to exit
		close(testChan)

		// wait for the done signal
		<-doneChan

		// Assert
		Expect(order.Sum).To(Equal(6))
	})
})
