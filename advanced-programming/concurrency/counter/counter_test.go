package main

import "testing"

func TestNoSyncCounter(t *testing.T) {
	ns := new(NoSyncCounter)
	got := GetCounters(ns, 10)
	// Expect ns to have value 30 (10 * 3 per goroutine)
	want := uint64(31)
	if got != want {
		t.Errorf("Expected getNext() to have value %d, got %d", want, got)
	}
}

func TestSyncAtomicCounter(t *testing.T) {
	sa := new(SyncAtomicCounter)
	got := GetCounters(sa, 10)
	// Expect ns to have value 30 (10 * 3 per goroutine)

	want := uint64(31)
	if got != want {
		t.Errorf("Expected getNext() to have value %d, got %d", want, got)
	}
}
