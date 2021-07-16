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
	data    []byte
	size    uint64
	k       int
	hashFns []hash.Hash64
}

// |/usr/share/dict/words| = 235886
func newBloomFilter() *BloomFilter {
	size := uint64(100000)
	return &BloomFilter{
		data:    make([]byte, size),
		size:    size,
		k:       2,
		hashFns: []hash.Hash64{fnv.New64(), fnv.New64a()},
	}
}

func (b *BloomFilter) add(item string) {
	for _, h := range b.hashFns {
		h.Write([]byte(item))
		sum := h.Sum64()
		// Get overall bit position in array of size*8 bits
		i := sum % (b.size * 8)
		// Then split into array index and remainder
		quo, rem := bits.Div64(0, i, 8)
		// Remainder is the bit offset into that byte
		var offset byte = 1 << rem
		// Turn on that bit in the byte
		b.data[quo] |= offset
		h.Reset()
	}
}

func (b *BloomFilter) maybeContains(item string) bool {
	for _, h := range b.hashFns {
		h.Write([]byte(item))
		sum := h.Sum64()
		h.Reset()
		i := sum % (b.size * 8)
		quo, rem := bits.Div64(0, i, 8)
		var offset byte = 1 << rem
		if b.data[quo]&(offset) == 0 {
			return false
		}
	}
	return true
}

func (b *BloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}
