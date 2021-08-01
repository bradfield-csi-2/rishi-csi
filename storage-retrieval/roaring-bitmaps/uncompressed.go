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
	data := make([]uint64, len(b.data))
	for i, _ := range data {
		data[i] = b.data[i] | other.data[i]
	}
	return &uncompressedBitmap{
		data: data,
	}
}

func (b *uncompressedBitmap) Intersect(other *uncompressedBitmap) *uncompressedBitmap {
	data := make([]uint64, len(b.data))
	for i, _ := range data {
		data[i] = b.data[i] & other.data[i]
	}
	return &uncompressedBitmap{
		data: data,
	}
}
