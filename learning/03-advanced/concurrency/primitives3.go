package main

import (
	"fmt"
	"time"
)

func main() {
	// Infinite
	go func() {
		for {
			select {
			default:
				fmt.Println("DOING WORK")
			}
		}
	}()

	time.Sleep(time.Second * 10)
}
