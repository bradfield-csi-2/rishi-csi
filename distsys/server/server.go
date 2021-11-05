package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	pb "dkvs/command"

	"google.golang.org/protobuf/proto"
)

type Store struct {
	db   map[string]string
	File *os.File
}

type Server struct {
	Ln          net.Listener
	Store       *Store
	isPrimary   bool
	secondaries []*Server
	primary     *Server
}

func NewStore() *Store {
	return &Store{
		db: make(map[string]string),
	}
}
func (s *Store) Setup() {
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
		return
	}

	var out map[string]interface{}
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		fmt.Println("Error unmarshaling JSON data from store")
		os.Exit(1)
	}

	for k, v := range out {
		s.db[k] = v.(string)
	}
	s.File = f
}

func (s *Server) ExecuteCommand(req *pb.Request) *pb.Response {
	resp := &pb.Response{}

	switch op := req.Op; op {
	case pb.Request_GET:
		if val, ok := s.Store.db[req.Pair.Key]; ok {
			resp.Status = pb.Response_OK
			resp.Message = val
		} else {
			resp.Status = pb.Response_ERROR
			resp.Message = fmt.Sprintf("key '%s' not found", req.Pair.Key)
		}
	case pb.Request_SET:
		s.Store.db[req.Pair.Key] = req.Pair.Value
		jsonData, err := json.Marshal(s.Store.db)
		if err != nil {
			resp.Status = pb.Response_ERROR
			resp.Message = "Error marshaling JSON"
			break
		}
		s.Store.File.WriteAt(jsonData, 0)
		// Synchronous Replication here
		// out, err := proto.Marshal(req)
		// if err != nil {
		// 	msg := "Error marshaling replication message"
		// 	resp.Status = pb.Response_ERROR
		// 	resp.Message = msg
		// 	return resp
		// }
		// secondaryConn.Write(out)
		// buf := make([]byte, 1024)
		// n, err := secondaryConn.Read(buf)
		// buf = buf[:n]
		// secondaryResp := &pb.Response{}
		// err = proto.Unmarshal(buf, secondaryResp)
		// if err != nil || secondaryResp.Status == pb.Response_ERROR {
		// 	msg := fmt.Sprintf("Error from seconday", err)
		// 	resp.Status = pb.Response_ERROR
		// 	resp.Message = msg
		// 	return resp
		// }

		resp.Status = pb.Response_OK
		resp.Message = fmt.Sprintf("Set key '%s' to '%s'", req.Pair.Key, req.Pair.Value)
	default:
		resp.Status = pb.Response_ERROR
		resp.Message = fmt.Sprintf("Invalid Operation: %v", op)
	}
	return resp
}

func (s *Server) Handle(conn net.Conn) {
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
		response := s.ExecuteCommand(req)
		out, err := proto.Marshal(response)
		if err != nil {
			fmt.Println("Error marshaling response proto")
			continue
		}
		conn.Write(out)
	}
	conn.Close()
}

func NewServer() *Server {
	return &Server{
		Store: NewStore(),
	}
}

func (s *Server) Start() {
	s.Store.Setup()
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Printf("Error listening on socket")
		os.Exit(1)
	}
	s.Ln = ln
}
