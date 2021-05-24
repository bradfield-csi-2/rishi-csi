package main

import (
	"fmt"
	"rocks"
)

func main() {
	key := "carrot"
	val := "orange"
	rocks.Put(key, val)
	newVal := rocks.Get(key)
	fmt.Printf("key %s has value = %s\n", key, newVal)
}
