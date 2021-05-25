package main

import (
	"fmt"
	"rocks"
)

func main() {
	key := "carrot"
	val := "orange"

	db := rocks.CreateDB()
	db.Put(key, val)
	newVal := db.Get(key)
	fmt.Printf("key %s has value: %s\n", key, newVal)
	rocks.DestroyDB(db)
}
