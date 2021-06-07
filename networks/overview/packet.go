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

type Packet struct {
	EthernetFrameHeader
	Payload       [64]byte
	FrameCheckSeq [4]byte
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

	for offset < capturedBytes {
		packetHeaderReader := io.NewSectionReader(f, offset, snapshotLen)
		packetHeader := parsePacketHeader(packetHeaderReader)
		offset += int64(unsafe.Sizeof(*packetHeader))
		fmt.Printf("Pcap Packet Header:\n%+v\n", packetHeader)

		packetReader := io.NewSectionReader(f, offset, int64(packetHeader.Length))
		packet := parsePacket(packetReader)
		offset += int64(packetHeader.Length)
		fmt.Printf("Packet:\n%+v\n", packet)
	}
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

func parsePacketHeader(f io.Reader) *PcapPacketHeader {
	header := new(PcapPacketHeader)
	binary.Read(f, binary.LittleEndian, header)
	return header
}
