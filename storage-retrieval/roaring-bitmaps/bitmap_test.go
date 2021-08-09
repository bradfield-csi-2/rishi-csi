package bitmap

import (
	"math/rand"
	"testing"
)

const (
	start = 1000 * 1000
	limit = 100 * 1000
	items = 10 * 1000
)

func TestBitmap(t *testing.T) {
	b1 := newUncompressedBitmap()
	m1 := make(map[uint32]struct{})

	// Call Set a bunch
	for i := 0; i < items; i++ {
		x := start + uint32(rand.Intn(limit))
		b1.Set(x)
		m1[x] = struct{}{}
	}

	// Make sure subsequent Get works as expected
	for x := uint32(0); x < start+limit+wordSize; x++ {
		_, ok := m1[x]
		if ok != b1.Get(x) {
			t.Fatalf("Get should've returned %t for %d\n", ok, x)
		}
	}

	b2 := newUncompressedBitmap()
	m2 := make(map[uint32]struct{})

	// Call Set a bunch
	for i := 0; i < items; i++ {
		x := uint32(rand.Intn(limit))
		b2.Set(x)
		m2[x] = struct{}{}
	}

	union := b1.Union(b2)
	intersect := b1.Intersect(b2)
	for x := uint32(0); x < start+limit+wordSize; x++ {
		_, ok1 := m1[x]
		_, ok2 := m2[x]
		if (ok1 || ok2) != union.Get(x) {
			t.Fatalf("Union: Get should've returned %t for %d\n", ok1 || ok2, x)
		}
		if (ok1 && ok2) != intersect.Get(x) {
			t.Fatalf("Intersect: Get should've returned %t for %d\n", ok1 && ok2, x)
		}
	}

	compressed := compress(b1)
	t.Logf("Uncompressed size: %d words, compressed size: %d words\n", len(b1.data), len(compressed))
	b := decompress(compressed)
	for x := uint32(0); x < start+limit+wordSize; x++ {
		if b1.Get(x) != b.Get(x) {
			t.Fatalf("Compression then decompression produced inconsistent result for %d\n", x)
		}
	}
}

func TestCompress(t *testing.T) {
	b := newUncompressedBitmap()
	b.Set(64)
	// Should see compressed value 100...01 0...01
	var elem1 uint64 = 0x8000000000000001
	var elem2 uint64 = 0b10
	c := compress(b)
	if c[0] != elem1 {
		t.Fatalf("First element is not correct, got %#064b, want %#064b\n", c[0], elem1)
	}
	if c[1] != elem2 {
		t.Fatalf("Second element is not correct, got %#064b, want %#064b\n", c[1], elem2)
	}
}
