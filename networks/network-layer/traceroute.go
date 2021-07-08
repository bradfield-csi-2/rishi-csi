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
	"time"
)

const (
	IP_HEADER_LENGTH = 20
	MAX_HOPS         = 64
	PING_REPLY_CODE  = 0
	PING_REPLY_TYPE  = 0
	PING_REQ_CODE    = 0
	PING_REQ_TYPE    = 8
	RESP_BUFFER_SIZE = 1024
	TIMEOUT          = time.Duration(5) * time.Second
)

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
	host := "google.com"
	dest := &syscall.SockaddrInet4{Addr: getIPFromHost(host)}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	check(err)
	err = syscall.SetNonblock(fd, true)
	check(err)
	err = syscall.Bind(fd, &syscall.SockaddrInet4{})
	check(err)

	var hop int = 1
	for hop < MAX_HOPS {
		probeStart := time.Now()
		req := newICMPRequest(hop)
		syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TTL, hop)
		syscall.Sendto(fd, req, 0, dest)
		resp := make([]byte, RESP_BUFFER_SIZE)
		for {
			n, _, err := syscall.Recvfrom(fd, resp, 0)
			if err == nil {
				resp = resp[:n]
				//fmt.Printf("Length: %d\n% x\n", n, resp)
				rtt := time.Since(probeStart)
				fmt.Printf("%d\t%s\t%v\n", hop, formatIp(resp[12:16]), rtt)
				break
			}
			if time.Since(probeStart) > TIMEOUT {
				fmt.Printf("%d\t%s\n", hop, "*")
				break
			}
		}
		if isLastHop(resp[IP_HEADER_LENGTH:]) {
			break
		}
		hop++
	}
}

func isLastHop(resp []byte) bool {
	header := new(ICMPHeader)
	buf := bytes.NewReader(resp)
	binary.Read(buf, binary.BigEndian, header)
	return header.Type == PING_REPLY_TYPE && header.Code == PING_REPLY_CODE
}

func formatIp(rawIp []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", rawIp[0], rawIp[1], rawIp[2], rawIp[3])
}

func newICMPRequest(seqNum int) []byte {
	h := &ICMPHeader{
		Type:       PING_REQ_TYPE,
		Code:       PING_REQ_CODE,
		Identifier: uint16(rand.Uint32()),
		SeqNumber:  uint16(seqNum),
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
