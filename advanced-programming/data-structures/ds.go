package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

const (
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits

	dataOffset = unsafe.Offsetof(struct {
		b bmap
		v int64
	}{}.v)
)

type hmap struct {
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}

type mapextra struct {
	// If both key and elem do not contain pointers and are inline, then we mark bucket
	// type as containing no pointers. This avoids scanning such maps.
	// However, bmap.overflow is a pointer. In order to keep overflow buckets
	// alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.
	// overflow and oldoverflow are only used if key and elem do not contain pointers.
	// overflow contains overflow buckets for hmap.buckets.
	// oldoverflow contains overflow buckets for hmap.oldbuckets.
	// The indirection allows to store a pointer to the slice in hiter.
	overflow    *[]*bmap
	oldoverflow *[]*bmap

	// nextOverflow holds a pointer to a free overflow bucket.
	nextOverflow *bmap
}

// A bucket for a Go map.
type bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	tophash [bucketCnt]uint8
	// Followed by bucketCnt keys and then bucketCnt elems.
	// NOTE: packing all the keys together and then all the elems together makes the
	// code a bit more complicated than alternating key/elem/key/elem/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}

func Float64ToUint64(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

func StringsPointToSame(s, t string) bool {
	sptr := *(*uint64)(unsafe.Pointer(&s))
	slen := *(*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Sizeof(&s)))
	tptr := *(*uint64)(unsafe.Pointer(&t))
	tlen := *(*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(&t)) + unsafe.Sizeof(&t)))
	send := sptr + slen
	tend := tptr + tlen

	// Check if one wholly contains another
	// Either s starts on or after t and ends before or on s
	// or the other way around
	return (sptr >= tptr && send <= tend) || (tptr >= sptr && tend <= send)
}

func SliceSum(s []int) int {
	hdptr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	arr := (*hdptr).Data
	length := (*hdptr).Len
	offset := unsafe.Sizeof(int(0))

	acc := 0
	for i := 0; i < length; i++ {
		acc += *(*int)(unsafe.Pointer(arr + uintptr(i)*offset))
	}
	return acc
}

func HashSum(m map[int]int) (ksum, vsum int) {
	// Need this type to pull out the actual pointer to the hmap
	type mapinterface struct {
		maptype unsafe.Pointer
		data    unsafe.Pointer
	}
	hmapptr := (*hmap)((*mapinterface)(unsafe.Pointer(&m)).data)
	buckets := hmapptr.buckets
	numbkts := uintptr(1 << hmapptr.B)

	thoffset := unsafe.Sizeof(uint8(0)) // Offset for tophash elements
	kvoffset := unsafe.Sizeof(int(0))   // Offset for key and value elements
	// Full size of a bucket: tophash, keys + values, and overflow pointer
	bucketoffset := dataOffset + (2 * bucketCnt * kvoffset) + unsafe.Sizeof(hmapptr)

	// Iterate over all the buckets
	var i uintptr = 0
	var overflow, bucket *bmap
	for b := uintptr(0); b < numbkts; b++ {
		// If this value is not nil, then as we go up to the next iteration, undo
		// the increment to the next bucket and instead look at this overflow
		// bucket instead
		if overflow != nil && *(*uint8)(unsafe.Pointer(&overflow.tophash)) != 0 {
			bucket = overflow
			b--
		} else {
			bucket = (*bmap)(unsafe.Pointer(uintptr(buckets) + b*bucketoffset))
		}

		// Iterate through tophash array of the current bucket
		tophashptr := unsafe.Pointer(&bucket.tophash)
		for i = uintptr(0); i < bucketCnt; i++ {
			th := *(*uint8)(unsafe.Pointer(uintptr(tophashptr) + i*thoffset))
			if th == 0 {
				break
			}

			// The keys are stored consecutively after the tophash (dataoffset)
			// The values are stored just below the keys, so we have to jump
			// bucketCnt keys
			kptr := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(bucket)) + dataOffset + i*kvoffset))
			vptr := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(kptr)) + bucketCnt*kvoffset))
			ksum += *kptr
			vsum += *vptr
		}

		// The overflow bucket is after the keys and values
		overflow = (*bmap)(unsafe.Pointer(uintptr(unsafe.Pointer(bucket)) + dataOffset + 2*bucketCnt*kvoffset))
	}

	return ksum, vsum
}

func main() {
	m1 := map[int]int{0: 1, 1: 2, 2: 3, 3: 4, 4: 5, 5: 6, 6: 7, 7: 8, 8: 9, 9: 10}
	k, v := HashSum(m1)
	fmt.Printf("HashSum(%#v) = %d, %d\n", m1, k, v)

	m2 := make(map[int]int)
	k, v = HashSum(m2)
	fmt.Printf("HashSum(%#v) = %d, %d\n", m2, k, v)

	m3 := map[int]int{100: 100, 2: 2, -100: -100}
	k, v = HashSum(m3)
	fmt.Printf("HashSum(%#v) = %d, %d\n", m3, k, v)
}
