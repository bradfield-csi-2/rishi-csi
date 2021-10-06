package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

type Writer struct {
	heapFile *HeapFile
	relation string
}

type Reader struct {
	heapFile    *HeapFile
	relation    string
	currPage    *Page
	currPageNum int64
	currSlot    int
}

type HeapFile struct {
	file     *os.File
	pages    []*Page
	numPages int64
}

// 8 kb pages
type Page struct {
	numEntries     uint16
	startFreeSpace uint16
	endFreeSpace   uint16
	slotArray      []Slot
}

type Slot struct {
	size   uint16
	offset uint16
}

func newWriter(relName string) *Writer {
	// TODO: Check if file exists before creating
	f, err := os.Create(relName)
	if err != nil {
		return nil
	}

	h := newHeapFile(f)
	return &Writer{heapFile: h, relation: relName}
}

func newReader(relName string) *Reader {
	f, err := os.Open(relName)
	if err != nil {
		return nil
	}

	h := newHeapFile(f)
	return &Reader{
		heapFile:    h,
		relation:    relName,
		currPage:    nil,
		currPageNum: 0,
		currSlot:    0,
	}
}

func newHeapFile(f *os.File) *HeapFile {
	return &HeapFile{file: f, pages: []*Page{newPage()}, numPages: 1}
}

func newPage() *Page {
	return &Page{
		numEntries:     0,
		startFreeSpace: 6,
		endFreeSpace:   0xffff,
		slotArray:      make([]Slot, 0),
	}
}

var PAGE_SIZE int64 = 0xffff

func (w *Writer) Write(r row) {
	p := w.heapFile.pages[w.heapFile.numPages-1]
	freeSpace := p.endFreeSpace - p.startFreeSpace

	rowBytes := make([]byte, 0)
	for _, v := range r {
		rowBytes = append(rowBytes, []byte(v)...)
	}
	rowSize := uint16(len(rowBytes))

	if freeSpace < rowSize {
		// Add a new page
		p = newPage()
		w.heapFile.numPages += 1
		w.heapFile.pages = append(w.heapFile.pages, p)
	}

	p.numEntries++
	slot := Slot{size: rowSize, offset: p.endFreeSpace - rowSize}

	p.startFreeSpace += 4 // New slot array entry is four bytes
	p.endFreeSpace -= rowSize

	// Will append to the slot array -- not handling deletes right now
	p.slotArray = append(p.slotArray, slot)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, p.numEntries)
	binary.Write(buf, binary.LittleEndian, p.startFreeSpace)
	binary.Write(buf, binary.LittleEndian, p.endFreeSpace)
	for _, s := range p.slotArray {
		binary.Write(buf, binary.LittleEndian, s.size)
		binary.Write(buf, binary.LittleEndian, s.offset)
	}
	headerBytes := buf.Bytes()

	// Now flush this page back to disk
	pageStartOffset := (w.heapFile.numPages - 1) * PAGE_SIZE
	byteOffsetInFile := pageStartOffset + int64(p.endFreeSpace)

	w.heapFile.file.WriteAt(headerBytes, pageStartOffset)
	w.heapFile.file.WriteAt(rowBytes, byteOffsetInFile)
}

func (r *Reader) Read() []byte {
	p := r.currPage
	pageStartOffset := r.currPageNum * PAGE_SIZE

	// If page is empty, read it into memory
	if p == nil {
		b := make([]byte, PAGE_SIZE)
		_, err := r.heapFile.file.ReadAt(b, pageStartOffset)
		if err != nil {
			return nil
		}
		numEntries := binary.LittleEndian.Uint16(b[0:2])
		startFreeSpace := binary.LittleEndian.Uint16(b[2:4])
		endFreeSpace := binary.LittleEndian.Uint16(b[4:6])
		slotArray := make([]Slot, 0)
		for i := 6; i < len(b); i += 4 {
			size := binary.LittleEndian.Uint16(b[i : i+2])
			if size == 0 {
				break
			}
			offset := binary.LittleEndian.Uint16(b[i+2 : i+4])
			slotArray = append(slotArray, Slot{offset, size})
		}

		p = &Page{numEntries, startFreeSpace, endFreeSpace, slotArray}
		// Need to load tuples into memory
		// The tuple array is indexed the same as the slot array,
		// So build back to front
	}

	slot := p.slotArray[r.currSlot]
	b := make([]byte, slot.size)
	byteOffsetInFile := pageStartOffset + int64(slot.offset)
	_, err := r.heapFile.file.ReadAt(b, byteOffsetInFile)
	if err != nil {
		return nil
	}
	return b
}
