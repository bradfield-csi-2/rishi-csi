package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const PING_REQ_TYPE = 8

type ICMPHeader struct {
	Type       byte
	Code       byte
	Checksum   uint16
	Identifier uint16
	SeqNumber  uint16
}

func (h *ICMPHeader) calculateChecksum() {
	sum := binary.BigEndian.Uint16([]byte{h.Type, h.Code}) + h.Identifier + h.SeqNumber
	h.Checksum = ^sum
}

func main() {
	fmt.Printf("Traceroute\n")
	host := "lichess.org"
	dest := &syscall.SockaddrInet4{Addr: getIPFromHost(host)}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	check(err)
	err = syscall.Bind(fd, &syscall.SockaddrInet4{})
	check(err)

	req := newICMPRequest(0)
	syscall.Sendto(fd, req, 0, dest)
}

func newICMPRequest(seqNum uint16) []byte {
	h := &ICMPHeader{
		Type:       PING_REQ_TYPE,
		Code:       0,
		Identifier: uint16(rand.Uint32()),
		SeqNumber:  seqNum,
	}
	h.calculateChecksum()
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, h)
	return buf.Bytes()
}

func getIPFromHost(host string) [4]byte {
	addrs, err := net.LookupHost(host)
	check(err)

	if len(addrs) == 0 {
		fmt.Printf("traceroute: %s", "Could not find IP for host")
		os.Exit(1)
	}
	parts := strings.Split(addrs[0], ".")
	addr := [4]byte{0, 0, 0, 0}

	for i, part := range parts {
		p, err := strconv.Atoi(part)
		check(err)
		addr[i] = byte(p)
	}

	return addr
}

func check(err error) {
	if err != nil {
		fmt.Printf("traceroute: %s", err)
		os.Exit(1)
	}
}
