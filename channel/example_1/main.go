package main

import "fmt"

func main() {
	// make will allocate memory for channel
	// Below we are initailze integer type channel
	ch := make(chan int)

	go func(a, b int) {
		c := a + b
		// Value send to channel
		ch <- c
	}(3, 4)
	// data Receive to the variable
	// get the value computed from goroutines
	r := <-ch
	fmt.Printf("Results: %v\n", r)
}
