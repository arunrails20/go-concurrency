package main

import (
	"fmt"
	"sync"
)

// Goroutines will execute and value also avail
// even after the enclose function return
func main() {
	var wg sync.WaitGroup

	incr := func(wg *sync.WaitGroup) {
		var i int
		wg.Add(1)
		go func() {
			defer wg.Done()
			i++
			fmt.Printf("value of i: %v\n", i)
		}()
		fmt.Println("return from function")
		return
	}
	incr(&wg)
	wg.Wait()
	fmt.Println("done..")
}
