package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const SockAddr = "/tmp/dkvs.sock"

func HandleInt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nExiting. Bye!")
		os.Exit(0)
	}()
}

func main() {
	HandleInt()

	conn, err := net.Dial("unix", SockAddr)
	if err != nil {
		fmt.Println("Error connecting to server: %s", err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("dkvs> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input")
			continue
		}
		fmt.Fprintf(conn, line)
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error getting response: %s", err)
			continue
		}
		fmt.Println(response)
	}
}
