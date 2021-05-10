package main

import (
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
	var tests = []struct {
		iter int
		want int
	}{
		{iter: 1000, want: 1000},
		{iter: 1, want: 1},
		{iter: 0, want: 0},
		{iter: 5000, want: 5000},
		{iter: 50000, want: 50000},
	}

	for _, test := range tests {
		got := Count(new(mutex), test.iter)
		if got != test.want {
			t.Errorf("Count(%d) = %d, want %d", test.iter, got, test.want)
		}
	}
}

func BenchmarkMyMutex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Count(new(mutex), 10000)
	}
}

func BenchmarkSyncMutex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Count(new(sync.Mutex), 10000)
	}
}
