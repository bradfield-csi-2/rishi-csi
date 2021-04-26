package main

import (
	"reflect"
	"unsafe"
)

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
