package main

import (
	"encoding/binary"
	"hash"
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
	k          int
	hashFns    []hash.Hash64
}

// |/usr/share/dict/words| = 235886
func newBloomFilter() *BloomFilter {
	size := uint64(100000)
	return &BloomFilter{
		data:       make([]byte, size),
		sizeInBits: size * 8,
		k:          2,
		hashFns:    []hash.Hash64{fnv.New64(), fnv.New64a()},
	}
}

func (b *BloomFilter) add(item string) {
	for _, h := range b.hashFns {
		h.Write([]byte(item))
		sum := h.Sum64()
		// Get overall bit position in array of bits
		// Then split into array index and remainder
		quo, rem := bits.Div64(0, sum%(b.sizeInBits), 8)
		// Remainder is the bit offset into that byte
		// Turn on that bit in the byte
		b.data[quo] |= (1 << rem)

		h.Reset()
	}
}

func (b *BloomFilter) maybeContains(item string) bool {
	for _, h := range b.hashFns {
		h.Write([]byte(item))
		sum := h.Sum64()
		h.Reset()
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
