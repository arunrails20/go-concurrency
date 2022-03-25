package main

import (
	"context"
	"fmt"
)

func main() {
	// generator generates integers in a separate GR and sends them
	// to the returned channel.
	// the callers of gen need to cancel the GR once
	// they consume 5 intergers
	// so that internal GR
	// started by gen is not leaked
	generator := func(ctx context.Context) <-chan int {
		ch := make(chan int)
		n := 1
		go func() {
			defer close(ch)
			for {
				select {
				case ch <- n:
				case <-ctx.Done():
					return
				}
				n++
			}
		}()
		return ch

	}

	// Create a context that is cancellable.
	ctx, cancel := context.WithCancel(context.Background())

	ch := generator(ctx)

	for n := range ch {
		fmt.Println(n)
		if n == 5 {
			cancel()
		}
	}
}
