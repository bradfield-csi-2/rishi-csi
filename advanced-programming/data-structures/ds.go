package main

import (
	"unsafe"
)

func Float64ToUint64(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}
