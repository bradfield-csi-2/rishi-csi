package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"
)

type PcapFileHeader struct {
	MagicNum            [4]byte
	MajVer              uint16
	MinVer              uint16
	TzOffset            uint32 // Always 0
	TzAcc               uint32 // Always 0
	SnapshotLen         uint32
	LinkLayerHeaderType uint32 // 1 == Ethernet (https://www.tcpdump.org/linktypes.html)
}

type PcapPacketHeader struct {
	TimestampSec       uint32
	TimestampMicroNano uint32
	Length             uint32
	UntruncatedLength  uint32
}

type EthernetFrameHeader struct {
	MACDst    [6]byte
	MACSrc    [6]byte
	Ethertype uint16
}

type IPDatagramHeader struct {
	Version_IHL          byte
	DSCP_ECN             byte
	Length               uint16
	ID                   uint16
	Flags_FragmentOffset [2]byte
	TTL                  uint8
	Protocol             uint8
	HeaderChecksum       uint16
	SrcIP                [4]byte
	DestIP               [4]byte
}

type TCPSegmentHeader struct {
	SrcPort                   uint16
	DstPort                   uint16
	SeqNum                    uint32
	AckNum                    uint32
	DataOffset_Reserved_Flags [2]byte
	WindowSize                uint16
	Checksum                  uint16
	UrgentPtr                 uint16
}

func main() {
	fmt.Printf("Parsing the packet capture...\n\n")
	f, err := os.Open("net.cap")
	defer f.Close()
	if err != nil {
		fmt.Errorf("Could not read packet capture file: %v", err)
	}
	stat, _ := f.Stat()
	capturedBytes := stat.Size()
	var offset int64

	fileHeader := parseFileHeader(f)
	fmt.Printf("Pcap File Header:\n%+v\n", fileHeader)

	offset += int64(unsafe.Sizeof(*fileHeader))
	snapshotLen := int64(fileHeader.SnapshotLen)

	var packetCount int
	for offset < capturedBytes {
		pcapHeaderReader := io.NewSectionReader(f, offset, snapshotLen)
		pcapHeader := parsePcapHeader(pcapHeaderReader)
		offset += int64(unsafe.Sizeof(*pcapHeader))
		fmt.Printf("Pcap Packet Header:\n%+v\n", pcapHeader)

		ethReader := io.NewSectionReader(f, offset, int64(pcapHeader.Length))
		ethHeader := parsePacket(ethReader)
		ethHeaderLen := int64(unsafe.Sizeof(*ethHeader))
		offset += ethHeaderLen
		fmt.Printf("Ethernet Header:\n%+v\n", ethHeader)

		ipHeaderReader := io.NewSectionReader(f, offset, int64(pcapHeader.Length))
		ipHeader := parseIPHeader(ipHeaderReader)
		ipHeaderLen := int64((ipHeader.Version_IHL<<4)>>4) * 4
		offset += ipHeaderLen
		fmt.Printf("IP Datagram Header:\n%+v\nLength: %d\n", ipHeader, ipHeaderLen)

		tcpSegmentReader := io.NewSectionReader(f, offset, int64(pcapHeader.Length))
		tcpSegmentHeader := parseTCPHeader(tcpSegmentReader)
		dataOffset := int64(((tcpSegmentHeader.DataOffset_Reserved_Flags[0]) >> 4) * 4)
		fmt.Printf("TCP Segment Header:\n%+v\nData Offset: %d\n", tcpSegmentHeader, dataOffset)

		tcpSegmentDataReader := io.NewSectionReader(f, offset+dataOffset, int64(pcapHeader.Length))
		tcpData := make([]byte, int64(ipHeader.Length)-ipHeaderLen-dataOffset)
		tcpSegmentDataReader.Read(tcpData)
		fmt.Printf("%s\n", tcpData)

		offset += int64(pcapHeader.Length) - ethHeaderLen - ipHeaderLen

		packetCount++
		fmt.Println("======================================")
	}

	fmt.Printf("Total Packets: %d\n", packetCount)
}

func parseTCPHeader(f io.Reader) *TCPSegmentHeader {
	tcpHeader := new(TCPSegmentHeader)
	binary.Read(f, binary.BigEndian, tcpHeader)
	return tcpHeader
}

func parseIPHeader(f io.Reader) *IPDatagramHeader {
	ipHeader := new(IPDatagramHeader)
	binary.Read(f, binary.BigEndian, ipHeader)
	return ipHeader
}

func parsePacket(f io.Reader) *EthernetFrameHeader {
	packet := new(EthernetFrameHeader)
	binary.Read(f, binary.BigEndian, packet)
	return packet
}

func parseFileHeader(f io.Reader) *PcapFileHeader {
	header := new(PcapFileHeader)
	binary.Read(f, binary.LittleEndian, header)
	return header
}

func parsePcapHeader(f io.Reader) *PcapPacketHeader {
	header := new(PcapPacketHeader)
	binary.Read(f, binary.LittleEndian, header)
	return header
}
