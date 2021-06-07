package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
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

func main() {
	fmt.Printf("Parsing the packet capture...\n\n")
	f, err := os.Open("net.cap")
	if err != nil {
		fmt.Errorf("Could not read packet capture file: %v", err)
	}

	fileHeader := parseFileHeader(f)
	fmt.Printf("Pcap File Header:\n%+v\n", fileHeader)
}

func parseFileHeader(f io.Reader) *PcapFileHeader {
	header := new(PcapFileHeader)
	binary.Read(f, binary.LittleEndian, header)
	return header
}
