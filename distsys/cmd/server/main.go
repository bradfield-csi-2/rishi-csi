package main

import (
	"fmt"
	"os"

	"dkvs/server"
)

func main() {
	s := server.NewServer()
	defer s.Store.File.Close()

	s.Start()
	defer s.Ln.Close()
	fmt.Println("Listening for requests...")

	for {
		conn, err := s.Ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting client connection: %s", err)
			os.Exit(1)
		}
		go s.Handle(conn)
	}
}
