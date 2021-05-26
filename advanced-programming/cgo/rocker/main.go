package main

import (
	"fmt"
	"log"
	"rocks"
)

func main() {
	key := "carrot"
	val := "orange"

	db, err := rocks.CreateDB("/tmp/rocksdb_test")
	if err != nil {
		log.Fatalf("rocker: error creating database: %s", err)
	}

	err = db.Put(key, val)
	if err != nil {
		log.Fatalf("rocker: error putting (key, val) (%s, %s): %s", key, val, err)
	}

	newVal, err := db.Get(key)
	if err != nil {
		log.Fatalf("rocker: error getting key %s: %s", key, err)
	}

	fmt.Printf("key %s has value: %s\n", key, newVal)
	rocks.DestroyDB(db)
}
