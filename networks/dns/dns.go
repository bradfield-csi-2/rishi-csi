package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fmt.Printf("Looking up DNS record for %s\n", os.Args[1])

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		fmt.Errorf("dns: could not open socket %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Socket File Descriptor: %d\n", socket)
}
