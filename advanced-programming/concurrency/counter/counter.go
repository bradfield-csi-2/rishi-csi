package main

import "fmt"

type counterService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently without any additional synchronization.
	getNext() uint64
}

// Part 1: Implement with no synchronization
type NoSyncCounter struct {
	counter uint64
}

func (c *NoSyncCounter) getNext() uint64 {
	c.counter++
	return c.counter
}

// Function to exercise counters
func GetCounters(cs counterService, n int) {
	for i := 0; i < n; i++ {
		go func() {
			first := cs.getNext()
			second := cs.getNext()
			third := cs.getNext()

			if first > second || second > third || first > third {
				fmt.Printf(
					"Not monotonic: first %d, second %d, third %d\n",
					first,
					second,
					third,
				)
			}
		}()
	}
}
