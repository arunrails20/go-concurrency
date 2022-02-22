// Using Select default to achive the non blocking commuication
package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)

	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			ch1 <- "Happy Coding"
		}
	}()

	for i := 0; i <= 2; i++ {
		select {
		case m := <-ch1:
			fmt.Println(m)
		default:
			fmt.Println("No Message Received")
		}
	}

	fmt.Println("Processing other tasks")
	time.Sleep(1500 * time.Millisecond)
}
