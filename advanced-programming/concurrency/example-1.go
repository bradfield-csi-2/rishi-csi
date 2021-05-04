package main

import (
	"fmt"
	"time"
)

func main_orig() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Printf("launched goroutine %d\n", i)
		}()
	}
	// Wait for goroutines to finish
	time.Sleep(time.Second)
}

// The bug here is not passing the loop variable to the goroutine. Each one
// will just see the last value, which is 10 in this case.
func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("launched goroutine %d\n", i)
		}(i)
	}
	// Wait for goroutines to finish
	time.Sleep(time.Second)
}
