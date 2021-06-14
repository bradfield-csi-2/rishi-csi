package main

import (
	"fmt"
	"os"
	"syscall"
)

const BACKLOG = 20
const PORT = 8000

var localhost = [4]byte{127, 0, 0, 1}

func main() {
	fmt.Printf("This will be a reverse proxy.\n")

	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Printf("proxy: could not open socket: %s\n", err)
		os.Exit(1)
	}

	// Open a socket to localhost:8000 and listen
	sockAddr := &syscall.SockaddrInet4{Addr: localhost, Port: PORT}
	err = syscall.Bind(sock, sockAddr)
	if err != nil {
		fmt.Printf("proxy: could not bind to port %d: %s\n", PORT, err)
		os.Exit(1)
	}
	err = syscall.Listen(sock, BACKLOG)
	if err != nil {
		fmt.Printf("proxy: could not listen on port %d: %s\n", PORT, err)
		os.Exit(1)
	}

	nfd, sa, err := syscall.Accept(sock)
	if err != nil {
		fmt.Printf("proxy: could not accept on port %d: %s\n", PORT, err)
		os.Exit(1)
	}

	fmt.Printf("New FD: %d\tSockAddr: %+v\n", nfd, sa)
}
