package main

import (
	"encoding/binary"
	"hash/fnv"
	"math/bits"
)

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type BloomFilter struct {
	data       []byte
	sizeInBits uint64
	k          uint64
}

// Size of /usr/share/dict/words = 235886
// Loading every other word, N ~= 117943
// For m = 1600000 bits, optimal k = 1600000/117943 * ln(2) = 9
func newBloomFilter() *BloomFilter {
	size := uint64(200000)
	return &BloomFilter{
		data:       make([]byte, size),
		sizeInBits: size * 8,
		k:          9,
	}
}

func (b *BloomFilter) add(item string) {
	h1 := fnv.New64()
	h2 := fnv.New64a()
	h1.Write([]byte(item))
	h2.Write([]byte(item))
	for k := uint64(0); k < b.k; k++ {
		// Use hash mixing from Kirsch 2006
		// https://www.eecs.harvard.edu/~michaelm/postscripts/rsa2008.pdf
		sum := h1.Sum64() + k*h2.Sum64()
		// Get overall bit position in array of bits
		// Then split into array index and remainder
		quo, rem := bits.Div64(0, sum%(b.sizeInBits), 8)
		// Remainder is the bit offset into that byte
		// Turn on that bit in the byte
		b.data[quo] |= (1 << rem)
	}
}

func (b *BloomFilter) maybeContains(item string) bool {
	h1 := fnv.New64()
	h2 := fnv.New64a()
	h1.Write([]byte(item))
	h2.Write([]byte(item))
	for k := uint64(0); k < b.k; k++ {
		sum := h1.Sum64() + k*h2.Sum64()
		quo, rem := bits.Div64(0, sum%(b.sizeInBits), 8)
		if b.data[quo]&(1<<rem) == 0 {
			return false
		}
	}
	return true
}

func (b *BloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}
