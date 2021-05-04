package main

import "testing"

func TestNoSyncCounter(t *testing.T) {
	ns := new(NoSyncCounter)
	GetCounters(ns, 10)
	// Expect ns to have value 30 (10 * 3 per goroutine)
	want := uint64(30)
	got := ns.getNext()
	if got != want {
		t.Errorf("Expected getNext() to have value %d, got %d", want, got)
	}
}
