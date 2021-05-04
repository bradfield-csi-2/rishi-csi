package main

import (
	"fmt"
	"sync"
)

type coordinator struct {
	lock   sync.RWMutex
	leader string
}

func newCoordinator(leader string) *coordinator {
	return &coordinator{
		lock:   sync.RWMutex{},
		leader: leader,
	}
}

func (c *coordinator) logState() {
	c.lock.RLock()
	defer c.lock.RUnlock()

	fmt.Printf("leader = %q\n", c.leader)
}

func (c *coordinator) setLeader_orig(leader string, shouldLog bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.leader = leader

	if shouldLog {
		c.logState()
	}
}

// We were locking the entire method, including a call to logState which also
// used a lock. If we only protect the critical section then we avoid deadlock
func (c *coordinator) setLeader(leader string, shouldLog bool) {
	c.lock.Lock()
	c.leader = leader
	c.lock.Unlock()

	if shouldLog {
		c.logState()
	}
}

func main() {
	c := newCoordinator("us-east")
	c.logState()
	c.setLeader("us-west", true)
}
