package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"syscall"
)

var googleDNSAddr [4]byte = [4]byte{0x8, 0x8, 0x8, 0x8}
var googleDNSPort int = 53

var TYPE_A = [2]byte{0x00, 0x1}
var CLASS_IN = [2]byte{0x0, 0x1}

type DNSHeader struct {
	ID                 uint16
	QR_Opcode_AA_TC_RD byte
	RA_Z_Rcode         byte
	QDCount            uint16
	ANCount            uint16
	NSCount            uint16
	ARCount            uint16
}

type DNSQuestion struct {
	QType  [2]byte
	QClass [2]byte
}

type DNSResourceRecord struct {
	Name     []byte
	Type     [2]byte
	Class    [2]byte
	TTL      uint32
	RDLength uint16
	RData    []byte
}

func main() {
	domain := os.Args[1]
	fmt.Printf("Looking up DNS record for %s\n", domain)

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		fmt.Printf("dns: could not open socket: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Socket File Descriptor: %d\n", socket)

	sa := &syscall.SockaddrInet4{}
	err = syscall.Bind(socket, sa)
	if err != nil {
		fmt.Printf("dns: could not bind to port: %s\n", err)
		os.Exit(1)
	}

	dnsQuery := buildDNSQuery(domain)
	fmt.Printf("DNS Query: % x\n", dnsQuery)
}

func buildDNSQuery(domain string) []byte {
	buf := new(bytes.Buffer)
	header := DNSHeader{
		ID:                 1234,
		QR_Opcode_AA_TC_RD: 0b00000001,
		RA_Z_Rcode:         0b00100000,
		QDCount:            1,
		ANCount:            0,
		NSCount:            0,
		ARCount:            0,
	}
	question := DNSQuestion{
		QType:  TYPE_A,
		QClass: CLASS_IN,
	}

	binary.Write(buf, binary.BigEndian, header)
	domainToQName(buf, domain)
	binary.Write(buf, binary.BigEndian, question)
	return buf.Bytes()
}

func domainToQName(buf *bytes.Buffer, domain string) {
	parts := strings.Split(domain, ".")
	for _, label := range parts {
		binary.Write(buf, binary.BigEndian, uint8(len(label)))
		binary.Write(buf, binary.BigEndian, []byte(label))
	}
	binary.Write(buf, binary.BigEndian, uint8(0))
}
