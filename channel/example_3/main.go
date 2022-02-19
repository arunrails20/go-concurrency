package main

import "fmt"

// Creating Ownership of the channels Avoids
// Deadlocking by writing to nil channel
// Avoid Panic
func main() {
	// Create channel owner, which creates channel
	// return receive only channel to consumer
	// spine a goroutines, which writes data to channel and
	// close the channel when its done

	// Owner of the channel is goroutines, that initailze, create and write
	owner := func() <-chan int {
		// Owner create the channel
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := 0; i < 5; i++ {
				ch <- i
			}
		}()
		return ch
	}

	// Consumer of the channel only read the channel
	consumer := func(ch <-chan int) {
		for v := range ch {
			fmt.Printf("Received: %d\n", v)
		}
		fmt.Println("Done!")
	}

	ch := owner()
	consumer(ch)
}
