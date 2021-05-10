package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type mutex struct {
	locked uint32
}

func (m *mutex) Lock() {
	for atomic.SwapUint32(&m.locked, 1) == 1 {
		// Spin until we're unlocked
	}
}

func (m *mutex) Unlock() {
	atomic.SwapUint32(&m.locked, 0)
}

func Count(lock sync.Locker, iter int) int {
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < iter; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock.Lock()
			counter++
			lock.Unlock()
		}()
	}

	wg.Wait()
	return counter
}

func main() {
	lock := new(mutex)
	x := Count(lock, 10000)
	fmt.Printf("Counter %d\n", x)
}
