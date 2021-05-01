package main

import (
	"fmt"
	"unsafe"
)

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type itab struct {
	inter *interfacetype
	_type *_type
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}

type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}

type name struct {
	bytes *byte
}

type nameOff int32
type typeOff int32
type tflag uint8

type imethod struct {
	name nameOff
	ityp typeOff
}

func ExtractInt(i interface{}) int {
	return *(*int)((*iface)(unsafe.Pointer(&i)).data)
}

func MethodInfo(i interface{}) {
	iface_ptr := (*iface)(unsafe.Pointer(&i))
	fmt.Printf("#v\n", iface_ptr)
}

func main() {
	x := 10101
	var i interface{} = x
	fmt.Printf("%d\n", ExtractInt(i))

	MethodInfo(2)
}
