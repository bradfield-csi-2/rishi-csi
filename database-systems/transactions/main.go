package main

import "fmt"

type Lock struct {
	kind string
}

type LockManager struct {
	granted      map[string]Lock
	requestQueue []string
}

func main() {
	fmt.Println("Hello")
	/*
		API:
		  acquire
			release
			detectDeadlock
	*/
}

func (lm *LockManager) acquire(objName string, op string) {
	/*
		  If the request is a read request,
			  -- check the operation. If read, can acquire if there is a shared lock
				-- if write, must acquire an exclusive lock
	*/
	lock, ok := lm.granted[objName]
	if !ok {
		fmt.Println("Lock not held")
		lm.granted[objName] = append(lm.granted[objName], new(Lock))
		return
	}
	fmt.Println("Lock held")

	// TODO: make this an enum
	if lock.kind == "SHARED" && op == "READ" {
		fmt.Println("Granting shared lock to reader")
		granted[objName] = append(granted[objName], new(Lock))
	} else {
		fmt.Println("Adding to request queue")
		lm.requestQueue = append(lm.requestQueue, objName)
	}

	fmt.Printf("%v\n", lock)
}

func (lm *LockManager) release(objName string) {
	/*
	  Remove from granted
	*/

}

func (lm *LockManager) detectDeadlock() {
}
