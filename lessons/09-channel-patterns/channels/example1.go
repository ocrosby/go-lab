package channels

type Order struct {
	Sum int
}

// processData generates a running total of values passed into the given channel.
// Because we will invoke this function on its own goroutine, we also pass a channel
// to this function for it to use to tell us that it is finished processing.
func (o *Order) ProcessData(values <-chan int, done chan<- struct{}) {
	for v := range values {
		o.Sum += v
	}
	done <- struct{}{}
	close(done)
}
