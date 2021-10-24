package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

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
	reader := bufio.NewReader(os.Stdin)
	HandleInt()

	for {
		fmt.Printf("dkvs> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input")
			continue
		}
		line = strings.TrimSpace(line)
		fmt.Println(line)
	}
}
