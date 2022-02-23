package main

// condition variable are waiting for a certion condition to met
// once it met GR will execute
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

	// Suspend the GR unitl sharedRsc is populated
	go func() {
		defer wg.Done()
		c.L.Lock()
		for len(shareRsc) == 0 {
			// will suspend the GR
			c.Wait()
		}
		fmt.Println(shareRsc["rsc1"])
		c.L.Unlock()
	}()

	// Acquire the Lock
	c.L.Lock()
	shareRsc["rsc1"] = "foo"
	// Indicate the GR to condition is met
	c.Signal()
	c.L.Unlock()
	wg.Wait()
}
