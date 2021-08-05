package bitmap

const wordSize = 64

type uncompressedBitmap struct {
	data []uint64
}

func newUncompressedBitmap() *uncompressedBitmap {
	return &uncompressedBitmap{data: make([]uint64, 1)}
}

func (b *uncompressedBitmap) Get(x uint32) bool {
	index, offset := x/wordSize, x%wordSize
	// Bitmap is too small to contain this value
	if index >= uint32(len(b.data)) {
		return false
	}
	return b.data[index]&(1<<offset) != 0
}

func (b *uncompressedBitmap) Set(x uint32) {
	index, offset := x/wordSize, x%wordSize

	// If our bitmap isn't big enough, grow it
	overshoot := int(index) - int(len(b.data)) + 1
	if overshoot > 0 {
		extra := make([]uint64, overshoot)
		b.data = append(b.data, extra...)
	}

	b.data[index] |= 1 << offset
}

func (b *uncompressedBitmap) Union(other *uncompressedBitmap) *uncompressedBitmap {
	// Set the new bitmap length to be the longest of the two
	var longest int
	if len(b.data) > len(other.data) {
		longest = len(b.data)
	} else {
		longest = len(other.data)
	}
	data := make([]uint64, longest)
	for i, _ := range data {
		if i >= len(b.data) {
			data[i] = other.data[i]
		} else if i >= len(other.data) {
			data[i] = b.data[i]
		} else {
			data[i] = b.data[i] | other.data[i]
		}
	}
	return &uncompressedBitmap{
		data: data,
	}
}

func (b *uncompressedBitmap) Intersect(other *uncompressedBitmap) *uncompressedBitmap {
	// Set the new bitmap length to be the shortest of the two
	var shortest int
	if len(b.data) < len(other.data) {
		shortest = len(b.data)
	} else {
		shortest = len(other.data)
	}
	data := make([]uint64, shortest)
	for i, _ := range data {
		data[i] = b.data[i] & other.data[i]
	}
	return &uncompressedBitmap{
		data: data,
	}
}
