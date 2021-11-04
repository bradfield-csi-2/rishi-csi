package main

import (
	pb "dkvs/dkvspb"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"google.golang.org/protobuf/proto"
)

const ServerSockAddr = "/tmp/dkvs-server.sock"

var secondaryStore map[string]string
var secondaryStoreFile *os.File

func Execute(req *pb.Request) *pb.Response {
	resp := &pb.Response{}

	switch op := req.Op; op {
	case pb.Request_GET:
		if val, ok := secondaryStore[req.Pair.Key]; ok {
			resp.Status = pb.Response_OK
			resp.Message = val
			return resp
		} else {
			msg := fmt.Sprintf("key '%s' not found", req.Pair.Key)
			resp.Status = pb.Response_ERROR
			resp.Message = msg
			return resp
		}
	case pb.Request_SET:
		secondaryStore[req.Pair.Key] = req.Pair.Value
		jsonData, err := json.Marshal(secondaryStore)
		if err != nil {
			msg := "Error marshaling JSON"
			resp.Status = pb.Response_ERROR
			resp.Message = msg
			return resp
		}
		secondaryStoreFile.WriteAt(jsonData, 0)
		msg := fmt.Sprintf("Set key '%s' to '%s'", req.Pair.Key, req.Pair.Value)
		resp.Status = pb.Response_OK
		resp.Message = msg
		return resp
	default:
		msg := fmt.Sprintf("Invalid Operation: %v", op)
		resp.Status = pb.Response_ERROR
		resp.Message = msg
		return resp
	}
}

func handleConn(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		buf = buf[:n]
		req := &pb.Request{}
		err = proto.Unmarshal(buf, req)
		if err != nil {
			fmt.Printf("Error unmarshaling proto: %s", err)
			os.Exit(1)
		}
		response := Execute(req)
		out, err := proto.Marshal(response)
		if err != nil {
			fmt.Println("Error marshaling response proto")
			continue
		}
		conn.Write(out)
	}
	conn.Close()
}

func SetupSecondaryStore() *os.File {
	secondaryStore = make(map[string]string)
	f, err := os.OpenFile("secondary-store.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening secondary data store")
		os.Exit(1)
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading data secondary store")
		os.Exit(1)
	}
	if len(bytes) == 0 {
		return f
	}

	var out map[string]interface{}
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		fmt.Println("Error unmarshaling JSON data from secondary store")
		os.Exit(1)
	}

	for k, v := range out {
		secondaryStore[k] = v.(string)
	}
	return f
}

func main() {
	secondaryStoreFile = SetupSecondaryStore()
	defer secondaryStoreFile.Close()

	if err := os.RemoveAll(ServerSockAddr); err != nil {
		fmt.Printf("Error cleaning up sockets")
		os.Exit(1)
	}

	l, err := net.Listen("unix", ServerSockAddr)
	if err != nil {
		fmt.Printf("Error listening on socket")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s", err)
			os.Exit(1)
		}
		go handleConn(conn)
	}
}
