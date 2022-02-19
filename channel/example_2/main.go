package main

import (
	"fmt"
)

// Send msg to ch1
// Direction of the arrow indicates
// sending or receive message
func genMsg(ch1 chan<- string) {
	ch1 <- "Happy coding"
}

func receiveMsg(ch1 <-chan string, ch2 chan<- string) {
	// receive messsage from ch1
	msg := <-ch1
	// send that message to ch2
	ch2 <- msg
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go genMsg(ch1)
	go receiveMsg(ch1, ch2)

	results := <-ch2
	fmt.Println(results)
}
