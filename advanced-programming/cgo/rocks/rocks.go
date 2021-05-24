package rocks

import "C"
import "fmt"

func Get(key string) string {
	val := "value"
	return val
}

func Put(key, val string) {
	fmt.Printf("Putting %s: %s\n", key, val)
}
