package main

// #include <stdlib.h>
//
// int square(int x) {
//   return x * x;
// }
import "C"
import "fmt"

func main() {
	fmt.Printf("square(4) = %d\n", C.square(4))
}
