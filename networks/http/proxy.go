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
	fmt.Printf("Listening on port %d...\n", PORT)

	// Open a socket and connect to a "remote" server for forwarding
	srvSock, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	srvSockAddr := &syscall.SockaddrInet4{Addr: localhost, Port: 9000}
	err = syscall.Connect(srvSock, srvSockAddr)
	if err != nil {
		fmt.Printf("proxy: could not connect to remote server: %s\n", err)
		os.Exit(1)
	}

	for {
		nfd, _, err := syscall.Accept(sock)
		if err != nil {
			fmt.Printf("proxy: could not accept on port %d: %s\n", PORT, err)
			os.Exit(1)
		}

		// Receive from the client and send to the remote
		buf := make([]byte, 1024)
		n, _, _, _, err := syscall.Recvmsg(nfd, buf, nil, 0)
		if err != nil {
			fmt.Printf("proxy: error receiving message: %s\n", err)
			os.Exit(1)
		}
		syscall.Sendmsg(srvSock, buf[:n], nil, srvSockAddr, 0)

		// Receive from the remote and send back to the client
		buf = make([]byte, 1024)
		n, _, _, from, err := syscall.Recvmsg(srvSock, buf, nil, 0)
		if err != nil {
			fmt.Printf("proxy: error receiving message: %s\n", err)
			os.Exit(1)
		}
		syscall.Sendmsg(nfd, buf[:n], nil, from, 0)

		syscall.Close(nfd)
	}
}
