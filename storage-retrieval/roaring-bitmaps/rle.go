package bitmap

import (
	"fmt"
	"math/bits"
)

// c  i s l l+1
// 0  0 0 3 0
// 3  0 3 1 2
// 6  1 2 2 1
// 9  2 1 3 0
// 12 3 0 3 0
// 15 3 3 1 2
// 18 4 2 2 1

// consumed = 0
// get w-1 bits from index 0, 0 bits from index 1
//	 which index are we starting at?
//     consumed / w = 0 / 4 = 0
//   which bit are we starting at?
//     consumed % w = 0 % 4 = 0
//   how many bits are getting from this index?
//     leftover bits = w - (consumed % w) = 4
//     a = max(w-1, leftover)
//   how many bits from index+1?
//		 (w-1) - a = 3 - 3 = 0
// consumed = 3
//	 which index are we starting at?
//     consumed / w = 3 / 4 = 0
//   which bit are we starting at?
//     consumed % w = 3 % 4 = 3
//   how many bits are getting from this index?
//     leftover bits = w - (consumed % w) = 1
//     a = (w-1)%leftover = 3 % 1 = 1
//   how many bits from index+1?
//		 (w-1) - a = 3 - 1 = 2
// consumed = 6
//	 which index are we starting at?
//     consumed / w = 6 / 4 = 1
//   which bit are we starting at?
//     consumed % w = 6 % 4 = 2
//   how many bits are getting from this index?
//     leftover bits = w - (consumed % w) = 2
//     a = leftover %(w-1) = 2 % 3 = 2
//   how many bits from index+1?
//		 (w-1) - a = 3 - 2 = 1
// consumed = 9
//	 which index are we starting at?
//     consumed / w = 9 / 4 = 2
//   which bit are we starting at?
//     consumed % w = 9 % 4 = 1
//   how many bits are getting from this index?
//     leftover bits = w - (consumed % w) = 3
//     a = leftover %(w) = 3 % 4 = 3
//   how many bits from index+1?
//		 (w-1) - a = 3 - 3 = 0
// consumed = 12
//	 which index are we starting at?
//     consumed / w = 12 / 4 = 3
//   which bit are we starting at?
//     consumed % w = 12 % 4 = 0
//   how many bits are getting from this index?
//     leftover bits = w - (consumed % w) = 4
//     a = leftover %(w-1) = 3 % 4 = 3
//   how many bits from index+1?
//		 (w-1) - a = 3 - 3 = 0

const LITERAL = 2

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func compress(b *uncompressedBitmap) []uint64 {
	c := make([]uint64, 0)
	uncompressedLen := len(b.data)
	totalBits := uint64(uncompressedLen * wordSize)
	var consumed, currRunBit, runBit, runLength uint64
	runBit = LITERAL
	//fmt.Printf("Uncompressed: %#0b\n", b.data)
	for consumed < totalBits {
		//for _, w := range b.data {

		// Suppose wordSize is 4.
		// That means we need to examine 3 bits at a time.
		// At any point we could be straddling two words, call them left and right,
		// where right is the current word at index i in the loop. E.g.
		// *xxx_  ____  ____  ____  ____ (left = 0, right = 3)
		//  ___x *xx__  ____  ____  ____ (left = 1, right = 2)
		//  ____  __xx *x___  ____  ____ (left = 2, right = 1)
		//  ____  ____  _xxx *____  ____ (left = 3, right = 0)
		//  ____  ____  ____ *xxx_  ____ (left = 0, right = 3)
		//  ____  ____  ____  ___x  *xx__ *(left = 0, right = 3)
		// leftOverlap and rightOverlap are the lengths of the L and R portions in
		// the diagram

		// Now that we have the overlaps, we want the most significant bits of the
		// left sided element and the least significant bits of the right sided
		// element. In the example above leftOverlap = 1 and rightOverlap = 2
		// We achieve this by bit shifting and masking.
		// leftMask  will look like 0b1000
		// rightMask will look like 0b0011

		// Finally, we can just OR left and right together to get our value

		remaining := totalBits - consumed
		bitsToCompress := min(remaining, uint64(wordSize-1))
		var val uint64
		for bitsToCompress > 0 {
			i, offset := consumed/wordSize, consumed%wordSize
			//fmt.Printf("bitsToCompress: %d, consumed: %d, i: %d, offset: %d\n", bitsToCompress, consumed, i, offset)
			// if i >= len(b.data) {
			// 	w := 0
			// }
			w := b.data[i]
			// How many bits are we able to write to the current index ?
			bitsAvailable := min(bitsToCompress, wordSize-offset)
			mask := uint64(((1 << bitsAvailable) - 1) << offset)
			// fmt.Printf("Word: %#0b\n", w)
			// fmt.Printf("Mask: %#0b\n", mask)
			val |= w & mask
			// fmt.Printf("Writing %d bitsAvailable to val %#0b\n", bitsAvailable, val)
			consumed += bitsAvailable
			bitsToCompress -= bitsAvailable
		}
		// fmt.Printf("Val: %#0b\n", val)
		if remaining < wordSize-1 {
			currRunBit = LITERAL
		} else {
			currRunBit = getRunBit(val)
		}
		// fmt.Printf("Run Bit: %#0b\n", currRunBit)

		// Three cases:
		// 1) Literal word, just append the literal word
		// 2) Extending a run, update the fill word with new length
		// 3) Starting a new run, append the fill word with length 1
		// fmt.Printf("=========\n")
		if currRunBit == LITERAL {
			runBit = LITERAL
			runLength = 0
			fmt.Printf("Appending literal word: %#064b\n", val)
			c = append(c, val)
		} else if currRunBit == runBit {
			runLength += (wordSize - 1)
			fmt.Printf("Incrementing run, length: %d\n", runLength)
			c[len(c)-1] = encodeFillWord(runBit, runLength)
		} else {
			// fmt.Printf("Starting new run\n")
			runBit = currRunBit
			runLength = wordSize - 1
			c = append(c, encodeFillWord(runBit, runLength))
		}
	}
	//fmt.Printf("Compressed: %#064b\n", c)
	return c
}

func getRunBit(val uint64) uint64 {
	ones := bits.OnesCount64(val)
	if ones == wordSize-1 {
		return 1
	} else if ones == 0 {
		return 0
	} else {
		return LITERAL
	}
}

func encodeFillWord(runBit, runLength uint64) uint64 {
	return (1 << (wordSize - 1)) | (runBit << (wordSize - 2)) | runLength
}

func decompress(compressed []uint64) *uncompressedBitmap {
	data := make([]uint64, len(compressed))
	var written uint64

	for _, w := range compressed {
		fmt.Printf("Word: %#0b\n", w)
		i, offset := written/wordSize, written%wordSize
		leftOverlap := min(wordSize-1, wordSize-offset)
		rightOverlap := uint64(wordSize - 1 - leftOverlap)

		// Make room for left and right values
		if i >= uint64(len(data)-1) {
			data = append(data, 0, 0)
		}

		leftMask := uint64(((1 << leftOverlap) - 1) << (wordSize - leftOverlap))
		rightMask := uint64((1 << rightOverlap) - 1)
		if w&(1<<(wordSize-1)) == 0 {
			fmt.Printf("Written: %d, i: %d, offset: %d\n", written, i, offset)
			// Literal word, write the next 63 bits split over the left and right
			// elements
			fmt.Printf("Found literal word: %#0b\n", w)
			data[i] |= (w & leftMask)
			data[i+1] |= (w & rightMask)
			fmt.Printf("Decompressed data[i], data[i+1]: %#0b, %#0b\n", data[i], data[i+1])
			written += (wordSize - 1)
		} else {
			// Fill word, write the appropriate number of 0s or 1s
			runBit := w & (1 << (wordSize - 2))
			runLength := w & ((1 << (wordSize - 2)) - 1)
			fmt.Printf("Found run with length: %d\n", runLength)
			// fmt.Printf("Word: %#0b\n", w)
			// A string of length wordSize 1s or 0s
			run := runBit * ((1 << wordSize) - 1)

			for runLength > 0 {
				//fmt.Printf("writting run with length: %d\n", runLength)
				i, offset := written/wordSize, written%wordSize
				if i >= uint64(len(data)) {
					data = append(data, 0)
				}
				// How many bits are we able to write to the current index
				overlap := min(runLength, wordSize-offset)
				mask := uint64(((1 << overlap) - 1) << offset)
				data[i] |= (run & mask)
				written += overlap
				runLength -= overlap
			}
		}
	}
	//fmt.Printf("Uncompressed: %#0b\n", data)
	return &uncompressedBitmap{
		data: data,
	}
}
