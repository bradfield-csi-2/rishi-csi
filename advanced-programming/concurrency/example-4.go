package main

import (
	"fmt"
)

func main_orig() {
	done := make(chan struct{}, 1)
	go func() {
		fmt.Println("performing initialization...")
		<-done
	}()

	done <- struct{}{}
	fmt.Println("initialization done, continuing with rest of program")
}

// The initialization can't happen on a buffered channel because the receive
// won't block the main goroutine. Instead, we have to make it an unbuffered
// channel so that we wait for intialization
func main() {
	done := make(chan struct{})
	go func() {
		fmt.Println("performing initialization...")
		<-done
	}()

	done <- struct{}{}
	fmt.Println("initialization done, continuing with rest of program")
}
