package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i <= 3; i++ {
		wg.Add(1)

		// Below code prints
		// 4 4 4 4
		// because when gorountine get a chance to execute the code at that i value got change to 4
		// go func() {
		// 	defer wg.Done()
		// 	fmt.Println(i)
		// }()

		// To Avoid above issue we need to
		// Passing a argument to the goroutine closure to get the current value
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
		// So Goroutines will operate on the current update value at time of execution
		// if we need to operate goroutine to work specific value than we need pass
		// those values as an argument to the goroutine

	}
	// Hey main goroutine wait until other goroutine lets finish
	// their job
	wg.Wait()
}
