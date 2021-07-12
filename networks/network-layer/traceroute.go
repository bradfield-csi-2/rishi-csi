package main

import (
	"bytes"
	"encoding/binary"
	"flag"
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
	host := flag.String("h", "", "host name to trace")
	flag.Parse()
	destIp := getIPFromHost(*host)
	dest := &syscall.SockaddrInet4{Addr: destIp, Port: 33434}

	sendSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	check(err)
	recvSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	check(err)
	err = syscall.SetsockoptTimeval(recvSock, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &syscall.Timeval{Sec: 5})
	check(err)

	ip := ""
	for hop := 1; hop < MAX_HOPS; hop++ {
		probeStart := time.Now()
		req := newICMPRequest(hop)
		syscall.SetsockoptInt(sendSock, syscall.IPPROTO_IP, syscall.IP_TTL, hop)
		syscall.Sendto(sendSock, req, 0, dest)
		resp := make([]byte, RESP_BUFFER_SIZE)
		for {
			n, _, err := syscall.Recvfrom(recvSock, resp, 0)
			if err == nil {
				resp = resp[:n]
				rtt := time.Since(probeStart)
				ip = formatIp(resp[12:16])
				fmt.Printf("%d\t%s\t%v\n", hop, ip, rtt)
				break
			} else {
				fmt.Printf("%d\t%s\n", hop, "*")
				break
			}
		}
		// Since we just get ping replies instead of "Port Unreachable" errors
		// because we're using raw sockets, we have to check that the current IP is
		// the destination IP
		if ip == formatIp(destIp[:]) {
			break
		}
	}
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
