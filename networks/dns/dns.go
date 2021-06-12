package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var googleDNSAddr [4]byte = [4]byte{0x8, 0x8, 0x8, 0x8}
var googleDNSPort int = 53

var TYPE_A uint16 = 0x0001
var CLASS_IN uint16 = 0x0001

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
	Name   []byte
	QType  uint16
	QClass uint16
}

type DNSResourceRecord struct {
	Name     []byte
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

func main() {
	domain := os.Args[1]
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		fmt.Printf("dns: could not open socket: %s\n", err)
		os.Exit(1)
	}

	sa := &syscall.SockaddrInet4{}
	err = syscall.Bind(socket, sa)
	if err != nil {
		fmt.Printf("dns: could not bind to port: %s\n", err)
		os.Exit(1)
	}

	dnsQuery := buildDNSQuery(domain)
	googleSockAddr := &syscall.SockaddrInet4{Addr: googleDNSAddr, Port: googleDNSPort}
	err = syscall.Sendto(socket, dnsQuery, 0, googleSockAddr)
	if err != nil {
		fmt.Printf("dns: could not send query: %s\n", err)
		os.Exit(1)
	}

	resp := make([]byte, 100)
	n, _, err := syscall.Recvfrom(socket, resp, 0)
	if err != nil {
		fmt.Printf("dns: could not receive response: %s\n", err)
		os.Exit(1)
	}

	parseDNSResponse(resp[:n])
}

func buildDNSQuery(domain string) []byte {
	buf := new(bytes.Buffer)
	header := DNSHeader{
		ID:                 uint16(rand.Uint32()),
		QR_Opcode_AA_TC_RD: 0b00000001,
		RA_Z_Rcode:         0b00100000,
		QDCount:            1,
		ANCount:            0,
		NSCount:            0,
		ARCount:            0,
	}
	question := DNSQuestion{
		Name:   domainToQName(domain),
		QType:  TYPE_A,
		QClass: CLASS_IN,
	}

	binary.Write(buf, binary.BigEndian, header)
	binary.Write(buf, binary.BigEndian, question.Name)
	binary.Write(buf, binary.BigEndian, question.QType)
	binary.Write(buf, binary.BigEndian, question.QClass)
	return buf.Bytes()
}

func domainToQName(domain string) []byte {
	buf := new(bytes.Buffer)
	parts := strings.Split(domain, ".")
	for _, label := range parts {
		binary.Write(buf, binary.BigEndian, uint8(len(label)))
		binary.Write(buf, binary.BigEndian, []byte(label))
	}
	binary.Write(buf, binary.BigEndian, uint8(0))
	return buf.Bytes()
}

func parseDNSResponse(resp []byte) {
	header := new(DNSHeader)
	binary.Read(bytes.NewReader(resp), binary.BigEndian, header)
	header.display()
	n := int(unsafe.Sizeof(*header))
	resp = resp[n:]

	var i uint16
	q := new(DNSQuestion)
	for i < header.QDCount {
		n = bytes.Index(resp, []byte{0x00})
		q.Name = qNametoDomain(resp[:n])
		resp = resp[n+1:]
		q.QType = binary.BigEndian.Uint16(resp[:2])
		q.QClass = binary.BigEndian.Uint16(resp[2:4])
		resp = resp[4:]
		i++
	}
	q.display()

	i = 0
	rr := new(DNSResourceRecord)
	for i < header.ANCount {
		n = bytes.Index(resp, []byte{0x00})
		rr.Name = qNametoDomain(resp[:n])
		resp = resp[n:]
		rr.Type = binary.BigEndian.Uint16(resp[:2])
		rr.Class = binary.BigEndian.Uint16(resp[2:4])
		rr.TTL = binary.BigEndian.Uint32(resp[4:8])
		rdlen := binary.BigEndian.Uint16(resp[8:10])
		resp = resp[10:]
		rr.RData = resp[:rdlen]
		resp = resp[rdlen:]
		i++
	}
	rr.display()
}

func (h DNSHeader) display() {
	var opcode string
	if h.QR_Opcode_AA_TC_RD&0x1e == 0 {
		opcode = "QUERY"
	}
	var flags []string
	if h.QR_Opcode_AA_TC_RD&0x80 == 0x80 {
		flags = append(flags, "qr")
	}
	if h.QR_Opcode_AA_TC_RD&0x40 == 0x40 {
		flags = append(flags, "aa")
	}
	if h.QR_Opcode_AA_TC_RD&0x20 == 0x20 {
		flags = append(flags, "tc")
	}
	if h.QR_Opcode_AA_TC_RD&0x01 == 0x01 {
		flags = append(flags, "rd")
	}
	if h.RA_Z_Rcode&0x80 == 0x80 {
		flags = append(flags, "ra")
	}
	var status string
	if h.RA_Z_Rcode&0x0f == 0 {
		status = "NO ERROR"
	}
	fmt.Printf(";; ->>HEADER<<- opcode: %s, status: %s, id: %d\n", opcode, status, h.ID)
	fmt.Printf(
		";; flags: %s; QUERY: %d, ANSWER: %d, AUTHORITY: %d, ADDITIONAL: %d\n\n",
		strings.Join(flags, " "),
		h.QDCount,
		h.ANCount,
		h.NSCount,
		h.ARCount,
	)
}

func (q DNSQuestion) display() {
	fmt.Println(";; QUESTION SECTION:")
	var qclass, qtype string
	if q.QClass == 1 {
		qclass = "IN"
	}
	if q.QType == 1 {
		qtype = "A"
	}
	fmt.Printf("%s\t%s\t%v\t%v\n", q.Name, "", qclass, qtype)
}

func (rr DNSResourceRecord) display() {
	fmt.Println("\n;; ANSWER SECTION:")
	var rrclass, rrtype string
	if rr.Class == 1 {
		rrclass = "IN"
	}
	if rr.Type == 1 {
		rrtype = "A"
	}
	ip := fmt.Sprintf("%d.%d.%d.%d", rr.RData[0], rr.RData[1], rr.RData[2], rr.RData[3])
	fmt.Printf("%s\t%d\t%v\t%v\t%s\n", rr.Name, rr.TTL, rrclass, rrtype, ip)
}

func qNametoDomain(qname []byte) []byte {
	buf := new(bytes.Buffer)
	var i int
	for i < len(qname) {
		// Skip compressed labels for now -- they start with two bits 11
		if (qname[i] & 0xc0) == 0xc0 {
			return []byte("Compressed")
		}
		labelLen := int(qname[i])
		buf.Write(qname[i+1 : i+1+labelLen])
		buf.Write([]byte("."))
		i = i + 1 + labelLen
	}

	return buf.Bytes()
}
