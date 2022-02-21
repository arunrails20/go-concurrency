package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Happing Coding"
	}()

	// It should Print the timeout Message, because the above GR delayed for 2 seconds
	// But Timeout print for 1 second
	select {
	case m1 := <-ch1:
		fmt.Println(m1)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
	}

}
