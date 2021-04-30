package main

import (
	"fmt"
	"unsafe"
)

type iface struct {
	tab  unsafe.Pointer
	data unsafe.Pointer
}

func ExtractInt(i interface{}) int {
	return *(*int)((*iface)(unsafe.Pointer(&i)).data)
}

func main() {
	x := 10101
	var i interface{} = x
	fmt.Printf("%d\n", ExtractInt(i))
}
