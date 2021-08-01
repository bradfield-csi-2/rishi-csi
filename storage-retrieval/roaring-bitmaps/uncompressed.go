package bitmap

const wordSize = 64

type uncompressedBitmap struct {
	data []uint64
}

func newUncompressedBitmap() *uncompressedBitmap {
	return &uncompressedBitmap{data: make([]uint64, 100000)}
}

func (b *uncompressedBitmap) Get(x uint32) bool {
	index, offset := x/wordSize, x%wordSize
	return b.data[index]&(1<<offset) != 0
}

func (b *uncompressedBitmap) Set(x uint32) {
	index, offset := x/wordSize, x%wordSize
	b.data[index] |= 1 << offset
}

func (b *uncompressedBitmap) Union(other *uncompressedBitmap) *uncompressedBitmap {
	var data []uint64
	return &uncompressedBitmap{
		data: data,
	}
}

func (b *uncompressedBitmap) Intersect(other *uncompressedBitmap) *uncompressedBitmap {
	var data []uint64
	return &uncompressedBitmap{
		data: data,
	}
}
