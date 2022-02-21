package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		// GR sending message in delay of 1 second
		time.Sleep(1 * time.Second)
		ch1 <- "message from first goroutine"
	}()

	go func() {
		// GR sending message in delay of 2 second
		time.Sleep(2 * time.Second)
		ch2 <- "message from second goroutines"
	}()

	// Using for loop to execute two time the select stmt
	// Select Statement will excecute which goroutines will completed first
	// In below first goroutine will print first, its delay of 1 second only
	for i := 0; i < 2; i++ {
		select {
		case m1 := <-ch1:
			fmt.Println(m1)
		case m2 := <-ch2:
			fmt.Println(m2)
		}
	}
}
