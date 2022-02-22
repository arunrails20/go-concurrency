package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	balance := 0

	runtime.GOMAXPROCS(4)

	var wg sync.WaitGroup

	// Using mutex to guard the access to the sharing variable
	// Sharing variable is balance
	// prevent from data race
	var mu sync.Mutex

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			deposit(&balance, 1, &mu)
		}()
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			withdraw(&balance, 1, &mu)
		}()
	}
	wg.Wait()
	fmt.Println(balance)
}

func deposit(balance *int, amount int, mu *sync.Mutex) {
	mu.Lock()
	*balance += amount
	defer mu.Unlock()
}

func withdraw(balance *int, amount int, mu *sync.Mutex) {
	mu.Lock()
	*balance -= amount
	defer mu.Unlock()
}
