package main

import (
	"fmt"
)

const numTasks = 3

func main_orig() {
	var done chan struct{}
	for i := 0; i < numTasks; i++ {
		go func() {
			fmt.Println("running task...")

			// Signal that task is done
			done <- struct{}{}
		}()
	}

	// Wait for tasks to complete
	for i := 0; i < numTasks; i++ {
		<-done
	}
	fmt.Printf("all %d tasks done!\n", numTasks)
}

// This is causing a deadlock because we're sending in the main goroutine and
// so we will block before we can move forward in the loop
func main() {
	var done chan struct{}
	for i := 0; i < numTasks; i++ {
		go func() {
			fmt.Println("running task...")

			// Signal that task is done
			done <- struct{}{}
		}()
	}

	// Wait for tasks to complete
	for i := 0; i < numTasks; i++ {
		go func() {
			<-done
		}
	}
	fmt.Printf("all %d tasks done!\n", numTasks)
}
