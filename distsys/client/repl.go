package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pb "dkvs/dkvspb"

	"google.golang.org/protobuf/proto"
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

func ParseRequest(cmd string) (*pb.Request, error) {
	cmd = strings.TrimSpace(cmd)
	i := strings.Index(cmd, " ")
	if i < 0 {
		return nil, fmt.Errorf("Invalid command: not enough arguments")
	}

	op := strings.ToLower(cmd[0:i])
	args := cmd[i+1:]
	req := &pb.Request{}
	pair := &pb.Pair{}
	if op == "get" {
		req.Op = pb.Request_GET
		pair.Key = args
	} else if op == "set" {
		req.Op = pb.Request_SET
		parts := strings.Split(args, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid command: must provide 'key=value' for set")
		}
		pair.Key = parts[0]
		pair.Value = parts[1]
	} else {
		return nil, fmt.Errorf("Invalid command: unknown operation")
	}
	req.Pair = pair
	return req, nil
}

func main() {
	HandleInt()

	conn, err := net.Dial("unix", SockAddr)
	if err != nil {
		fmt.Printf("Error connecting to server: %s\n", err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("dkvs> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input")
			continue
		}
		req, err := ParseRequest(line)
		if err != nil {
			fmt.Println("Error parsing request")
			continue
		}
		out, err := proto.Marshal(req)
		conn.Write(out)
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		buf = buf[:n]
		resp := &pb.Response{}
		err = proto.Unmarshal(buf, resp)
		if err != nil {
			fmt.Printf("Error getting response: %s", err)
			continue
		}
		fmt.Printf("%+v\n", resp)
	}
}
