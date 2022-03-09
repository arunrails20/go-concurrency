package main

import "fmt"

// Build the Pipeline System
// generator() -> square() -> print
// Above each functions are separation of concerns

// convert the list of integer to channel
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, i := range nums {
			out <- i
		}
		close(out)
	}()
	return out
}

// Receive an inbound channel, square the number, output on outbound the channel
func square(nums <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for i := range nums {
			out <- i * i
		}
		close(out)
	}()
	return out
}

func main() {
	// Compose the generator method, since the input type same
	//
	for n := range square(generator(2, 3)) {
		fmt.Println(n)
	}
}
