package bitmap

var chunkSize uint64 = wordSize - 1

func compress(b *uncompressedBitmap) []uint64 {
	compressed := make([]uint64, 0)
	var consumed, lastChunk, runLength uint64
	totalBits := uint64(len(b.data) * wordSize)
	for consumed < totalBits {
		// A chunk is a w-1 bit value, though the last chunk may be smaller
		// than this so we also return the size of the chunk in bits
		chunk, bitsInChunk := nextChunk(b.data, consumed)
		consumed += bitsInChunk

		// Here we have a bunch of trailing zeros so we can ignore these bits
		// instead of appending a useless word
		if chunk == 0 && bitsInChunk < chunkSize {
			continue
		}

		// When examining a chunk, there are two possibilites. Either we're:
		// 1) extending a run, so update the last element with the new length
		// 2) writing a literal or starting a new run, so append an element
		if chunk == lastChunk && runLength > 0 {
			runLength++
			compressed[len(compressed)-1] = fillWord(chunk, runLength)
		} else {
			if chunk == 0 || chunk == 1<<chunkSize-1 {
				// If the chunk was all 0 or 1, then start a run by appending
				// a new fill word
				runLength = 1
				compressed = append(compressed, fillWord(chunk, runLength))
			} else {
				// Otherwise it's not a fill word and we can just append the
				// literal chunk
				runLength = 0
				compressed = append(compressed, chunk)
			}
			lastChunk = chunk
		}
	}
	return compressed
}

func decompress(compressed []uint64) *uncompressedBitmap {
	data := make([]uint64, 1)
	var written uint64
	for _, w := range compressed {
		firstBit, chunk := w>>chunkSize, (1<<chunkSize-1)&w
		if firstBit == 1 {
			// This is a fill word
			// Write a chunk of the fill bit runLength times
			fillBit := chunk >> (chunkSize - 1)
			runLength := (1<<(chunkSize-1) - 1) & w
			for i := uint64(0); i < runLength; i++ {
				data = writeChunk(data, fillBit*(1<<chunkSize-1), written)
				written += chunkSize
			}
		} else {
			// This is a literal word, so just need to write
			data = writeChunk(data, chunk, written)
			written += chunkSize
		}
	}

	return &uncompressedBitmap{data: data}
}

func writeChunk(data []uint64, chunk, written uint64) []uint64 {
	i, o := written/wordSize, written%wordSize
	if o+chunkSize > wordSize || i >= uint64(len(data)) {
		// We may need to make space for another chunk
		data = append(data, 0)
	}
	wordBits := min(chunkSize, wordSize-o)
	nextWordBits := chunkSize - wordBits

	wordMask := uint64(1<<wordBits - 1)
	nextWordMask := uint64((1<<nextWordBits - 1) << wordBits)

	// Need to set MSBs of data[i] and LSBs of data[i+1]
	data[i] |= ((chunk & wordMask) << o)
	if i+1 < uint64(len(data)) {
		data[i+1] |= ((chunk & nextWordMask) >> wordBits)
	}
	return data
}

func fillWord(chunk, runLength uint64) uint64 {
	return (1 << chunkSize) | chunk&1<<(chunkSize-1) | runLength
}

func nextChunk(bitmap []uint64, consumed uint64) (uint64, uint64) {
	var word, nextWord uint64
	// A chunk may straddle two consecutive elements in the bitmap
	// i: the element to start reading bits from
	// o: the bit within element i to start reading from
	i, o := consumed/wordSize, consumed%wordSize
	word = bitmap[i]

	bitsInChunk := chunkSize
	if i == uint64(len(bitmap)-1) {
		// If i is the last element, then we might not be writing a full chunk
		bitsInChunk = min(chunkSize, wordSize-o)
	} else if o+chunkSize > wordSize {
		// Otherwise, if o + chunkSize > wordSize, then some of the bits in the
		// chunk lie in the next element in the bitmap
		nextWord = bitmap[i+1]
	}

	/*
		A chunk has up to two pieces, those bits that come from word and maybe
		some "spill over" bits that need to come from the next word.
		Consider the following bitmap with wordSize=4 and chunkSize=3:

		0011 0101

		Suppose we've read the first chunk and we're trying to get the next
		one. That is i=0, o=3. This corresponds to the chunk 010: the most
		significant bit of word and the two least significant bits of nextWord.

		The number of bits to come from word is the number of bits from o,
		i.e. wordSize - o. We take the minimum of this and chunkSize.
		chunkSize being the minimum corresponds to the case where there
		are no spill over bits. Otherwise, the remaining bits come from nextWord.
	*/
	wordBits := min(chunkSize, wordSize-o)
	nextWordBits := chunkSize - wordBits

	/*
		The steps to build a chunk are therefore:
		1) Left shift o bits into word and then mask off all remaining bits to
		grab the first portion of the chunk. These are the MSBs of word.
		2) Then we mask off the remaining bits from nextWord. These are the
		LSBs of nextWord.
		3) Finally, concatenate the masked bits from word and nextWord to
		create the chunk. We need to shift the nextWord bits to the left
		before concatenating it with the bits from word, otherwise they
		will clobber each other.

		Note here that the MSBs from word become the LSBs in chunk. Similarly,
		the LSBs from nextWord become the MSBs in the chunk, hence the left shift.
	*/
	wordMask := uint64((1<<wordBits - 1) << o)
	nextWordMask := uint64(1<<nextWordBits - 1)
	chunk := (word&wordMask)>>o | (nextWord&nextWordMask)<<wordBits

	return chunk, bitsInChunk
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
