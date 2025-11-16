package main

import "fmt"

func main() {
	charChannel := make(chan string, 3) // a buffered channel

	chars := []string{"a", "b", "c"}

	for _, s := range chars {
		select {
		case charChannel <- s:

		}
	}

	close(charChannel)

	for result := range charChannel {
		fmt.Println(result)
	}
}
