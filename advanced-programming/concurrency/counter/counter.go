package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

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

// Part 2: Implement with sync/atomic
type SyncAtomicCounter struct {
	counter uint64
}

func (c *SyncAtomicCounter) getNext() uint64 {
	return atomic.AddUint64(&c.counter, 1)
}

// Function to exercise counters
func GetCounters(cs counterService, n int) uint64 {
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
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
			wg.Done()
		}()
	}

	wg.Wait()

	return cs.getNext()
}

func main() {
	ns := new(NoSyncCounter)
	GetCounters(ns, 10)
}
