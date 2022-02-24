package main

import (
	"fmt"
	"sync"
)

var shareRsc = make(map[string]interface{})

func main() {
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	c := sync.NewCond(&mu)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// acquire the lock
		c.L.Lock()
		for len(shareRsc) < 1 {
			c.Wait()
		}
		fmt.Println(shareRsc["var1"])
		c.L.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// acquire the lock
		c.L.Lock()
		for len(shareRsc) < 2 {
			c.Wait()
		}
		fmt.Println(shareRsc["var2"])
		c.L.Unlock()
	}()

	c.L.Lock()
	// Writes values to shareRsc variables
	shareRsc["var1"] = "var1"
	shareRsc["var2"] = "var2"
	// Broadcast method will send signal to all waiting goroutines
	// so that GR will execute
	c.Broadcast()
	c.L.Unlock()
	wg.Wait()

}
