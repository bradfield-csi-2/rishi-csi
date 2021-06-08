package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
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

type ImagePacket struct {
	SeqNum uint32
	Data   []byte
}

type BySeqNum []ImagePacket

func (a BySeqNum) Len() int           { return len(a) }
func (a BySeqNum) Less(i, j int) bool { return a[i].SeqNum < a[j].SeqNum }
func (a BySeqNum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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
	offset += int64(unsafe.Sizeof(*fileHeader))
	snapshotLen := int64(fileHeader.SnapshotLen)

	imagePackets := make([]ImagePacket, 100)
	for offset < capturedBytes {
		pcapHeaderReader := io.NewSectionReader(f, offset, snapshotLen)
		pcapHeader := parsePcapHeader(pcapHeaderReader)
		offset += int64(unsafe.Sizeof(*pcapHeader))

		ethReader := io.NewSectionReader(f, offset, int64(pcapHeader.Length))
		ethHeader := parsePacket(ethReader)
		ethHeaderLen := int64(unsafe.Sizeof(*ethHeader))
		offset += ethHeaderLen

		ipHeaderReader := io.NewSectionReader(f, offset, int64(pcapHeader.Length))
		ipHeader := parseIPHeader(ipHeaderReader)
		ipHeaderLen := int64((ipHeader.Version_IHL<<4)>>4) * 4
		offset += ipHeaderLen

		tcpSegmentReader := io.NewSectionReader(f, offset, int64(pcapHeader.Length))
		tcpSegmentHeader := parseTCPHeader(tcpSegmentReader)
		dataOffset := int64(((tcpSegmentHeader.DataOffset_Reserved_Flags[0]) >> 4) * 4)

		tcpSegmentDataReader := io.NewSectionReader(f, offset+dataOffset, int64(pcapHeader.Length))
		tcpData := make([]byte, int64(ipHeader.Length)-ipHeaderLen-dataOffset)
		tcpSegmentDataReader.Read(tcpData)

		// Stick all the data from the server (from port 80) into an array for
		// later combining
		if tcpSegmentHeader.SrcPort == 80 {
			imagePackets = append(
				imagePackets,
				ImagePacket{SeqNum: tcpSegmentHeader.SeqNum, Data: tcpData},
			)
		}

		offset += int64(pcapHeader.Length) - ethHeaderLen - ipHeaderLen
	}

	combined := combineTcpResponse(imagePackets)
	httpHeader, imageData := parseHTTPResponse(combined)
	fmt.Printf("HTTP Header:\n\n%s\n", httpHeader)
	ioutil.WriteFile("img.jpg", imageData, 0644)
	fmt.Println("\nWrote img.jpg")
}

func parseHTTPResponse(resp []byte) (header string, image []byte) {
	respStr := string(resp)
	emptyLine := regexp.MustCompile(`\r\n\r\n`)
	parts := emptyLine.Split(respStr, 2)
	header = parts[0]
	image = []byte(parts[1])
	return
}

func combineTcpResponse(imagePackets []ImagePacket) []byte {
	sort.Sort(BySeqNum(imagePackets))
	resp := make([]byte, 1000)
	for i, pkt := range imagePackets {
		if i > 0 && (imagePackets[i-1].SeqNum == pkt.SeqNum) {
			continue
		}
		resp = append(resp, pkt.Data...)
	}

	return resp
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

func parsePcapHeader(f io.Reader) *PcapPacketHeader {
	header := new(PcapPacketHeader)
	binary.Read(f, binary.LittleEndian, header)
	return header
}

func parseFileHeader(f io.Reader) *PcapFileHeader {
	header := new(PcapFileHeader)
	binary.Read(f, binary.LittleEndian, header)
	return header
}
