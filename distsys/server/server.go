package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	pb "dkvs/dkvspb"

	"google.golang.org/protobuf/proto"
)

const ClientSockAddr = "/tmp/dkvs.sock"
const ServerSockAddr = "/tmp/dkvs-server.sock"

var store map[string]string
var storeFile *os.File

func SetupStore() *os.File {
	store = make(map[string]string)
	f, err := os.OpenFile("store.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening data store")
		os.Exit(1)
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading data store")
		os.Exit(1)
	}
	if len(bytes) == 0 {
		return f
	}

	var out map[string]interface{}
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		fmt.Println("Error unmarshaling JSON data from store")
		os.Exit(1)
	}

	for k, v := range out {
		store[k] = v.(string)
	}
	return f
}

func ExecuteCommand(req *pb.Request, secondaryConn net.Conn) *pb.Response {
	resp := &pb.Response{}

	switch op := req.Op; op {
	case pb.Request_GET:
		if val, ok := store[req.Pair.Key]; ok {
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
		store[req.Pair.Key] = req.Pair.Value
		jsonData, err := json.Marshal(store)
		if err != nil {
			msg := "Error marshaling JSON"
			resp.Status = pb.Response_ERROR
			resp.Message = msg
			return resp
		}
		storeFile.WriteAt(jsonData, 0)
		// Synchronous Replication here
		out, err := proto.Marshal(req)
		if err != nil {
			msg := "Error marshaling replication message"
			resp.Status = pb.Response_ERROR
			resp.Message = msg
			return resp
		}
		secondaryConn.Write(out)
		buf := make([]byte, 1024)
		n, err := secondaryConn.Read(buf)
		buf = buf[:n]
		secondaryResp := &pb.Response{}
		err = proto.Unmarshal(buf, secondaryResp)
		if err != nil || secondaryResp.Status == pb.Response_ERROR {
			msg := fmt.Sprintf("Error from seconday", err)
			resp.Status = pb.Response_ERROR
			resp.Message = msg
			return resp
		}

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

func handleConnection(conn, secondaryConn net.Conn) {
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
		response := ExecuteCommand(req, secondaryConn)
		out, err := proto.Marshal(response)
		if err != nil {
			fmt.Println("Error marshaling response proto")
			continue
		}
		conn.Write(out)
	}
	conn.Close()
}

func main() {
	storeFile = SetupStore()
	defer storeFile.Close()

	if err := os.RemoveAll(ClientSockAddr); err != nil {
		fmt.Printf("Error cleaning up sockets")
		os.Exit(1)
	}

	l, err := net.Listen("unix", ClientSockAddr)
	if err != nil {
		fmt.Printf("Error listening on socket")
		os.Exit(1)
	}
	defer l.Close()
	secondaryConn, err := net.Dial("unix", ServerSockAddr)
	if err != nil {
		fmt.Printf("Error listening on server socket")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s", err)
			os.Exit(1)
		}
		go handleConnection(conn, secondaryConn)
	}
}
