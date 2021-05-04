package main

import "testing"

func TestNoSyncCounter(t *testing.T) {
	ns := new(NoSyncCounter)
	got := GetCounters(ns, 10)
	want := uint64(31) // Expect ns to have value 30 (10 * 3 per goroutine) + 1
	if got != want {
		t.Errorf("Expected getNext() to have value %d, got %d", want, got)
	}
}

func TestSyncAtomicCounter(t *testing.T) {
	sa := new(SyncAtomicCounter)
	got := GetCounters(sa, 10)
	want := uint64(31) // Expect sa to have value 30 (10 * 3 per goroutine) + 1
	if got != want {
		t.Errorf("Expected getNext() to have value %d, got %d", want, got)
	}
}

func TestMutexCounter(t *testing.T) {
	mc := new(MutexCounter)
	got := GetCounters(mc, 10)
	want := uint64(31) // Expect mc to have value 30 (10 * 3 per goroutine) + 1
	if got != want {
		t.Errorf("Expected getNext() to have value %d, got %d", want, got)
	}
}
