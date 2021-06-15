package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

const BACKLOG = 20
const PORT = 8000

var localhost = [4]byte{127, 0, 0, 1}

func main() {
	cache := make(map[string][]byte)
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

		req := make([]byte, 1024)
		resp := make([]byte, 1024)

		// Receive from the client
		n, _, _, _ := recv(nfd, req, nil, 0)
		req = req[:n]
		path := getPath(req)
		if cachedResp, ok := cache[path]; !ok {
			fmt.Printf("Cache miss, requesting %s from server\n", path)
			send(srvSock, req, nil, srvSockAddr, 0)
			n, _, _, _ = recv(srvSock, resp, nil, 0)
			resp = resp[:n]
			cache[path] = resp
		} else {
			resp = cachedResp
		}
		// Send back to the client
		send(nfd, resp, nil, nil, 0)

		syscall.Close(nfd)
	}
}

func getPath(req []byte) string {
	reqline := strings.Split(string(req), "\n")[0]
	return strings.Split(reqline, " ")[1]
}

func recv(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from syscall.Sockaddr) {
	n, oobn, recvflags, from, err := syscall.Recvmsg(fd, p, oob, flags)
	if err != nil {
		fmt.Printf("proxy: error receiving message: %s\n", err)
		os.Exit(1)
	}
	return
}

func send(fd int, p, oob []byte, to syscall.Sockaddr, flags int) {
	err := syscall.Sendmsg(fd, p, oob, to, flags)
	if err != nil {
		fmt.Printf("proxy: error sending message: %s\n", err)
		os.Exit(1)
	}
}
