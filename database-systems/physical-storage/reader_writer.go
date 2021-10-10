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
	currSlot    uint16
}

type HeapFile struct {
	file     *os.File
	pages    []*Page
	numPages int64
	fields   []string
}

// 8 kb pages
type Page struct {
	numEntries     uint16
	startFreeSpace uint16
	endFreeSpace   uint16
	slotArray      []Slot
	tuples         []map[string]string
}

type Slot struct {
	size   uint16
	offset uint16
}

func newWriter(relName string, fields []string) *Writer {
	// TODO: Check if file exists before creating
	f, err := os.Create(relName)
	if err != nil {
		return nil
	}

	h := newHeapFile(f, fields)
	return &Writer{heapFile: h, relation: relName}
}

func newReader(relName string, fields []string) *Reader {
	f, err := os.Open(relName)
	if err != nil {
		return nil
	}

	h := newHeapFile(f, fields)
	return &Reader{
		heapFile:    h,
		relation:    relName,
		currPage:    nil,
		currPageNum: 0,
		currSlot:    0,
	}
}

func newHeapFile(f *os.File, fields []string) *HeapFile {
	return &HeapFile{
		file:     f,
		pages:    []*Page{newPage()},
		numPages: 1,
		fields:   fields,
	}
}

var PAGE_SIZE int64 = 0xffff
var NULL_BYTE uint32 = 0x00

func newPage() *Page {
	return &Page{
		numEntries:     0,
		startFreeSpace: 6,
		endFreeSpace:   uint16(PAGE_SIZE),
		slotArray:      make([]Slot, 0),
		tuples:         make([]map[string]string, 0),
	}
}

func (w *Writer) Write(r row) {
	p := w.heapFile.pages[w.heapFile.numPages-1]
	freeSpace := p.endFreeSpace - p.startFreeSpace

	tupleHeaderLen := len(w.heapFile.fields) * 4
	buf := new(bytes.Buffer)
	offset := uint16(tupleHeaderLen)
	tupleData := make([]byte, 0)
	for _, field := range w.heapFile.fields {
		if val, ok := r[field]; ok {
			size := uint16(len(val))
			binary.Write(buf, binary.LittleEndian, offset)
			binary.Write(buf, binary.LittleEndian, size)
			tupleData = append(tupleData, []byte(val)...)
			offset += size
		} else {
			// Write a zero four bytes for null fields
			binary.Write(buf, binary.LittleEndian, NULL_BYTE)
		}
	}
	buf.Write(tupleData)
	tuple := buf.Bytes()
	rowSize := uint16(len(tuple))

	if (freeSpace - 4) <= rowSize {
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

	buf = new(bytes.Buffer)
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
	w.heapFile.file.WriteAt(tuple, byteOffsetInFile)
}

func (r *Reader) Read() map[string]string {
	p := r.currPage

	// If the page is empty, read it into memory
	if p == nil {
		p = r.LoadPage()
		if p == nil {
			return nil
		}
		r.currPage = p
	}

	// Read the tuple from the page in memory and increment the slot for the next
	// call to Read()
	b := p.tuples[r.currSlot]
	r.currSlot++

	// If we have read the entire page, tick over to the next one. The next call
	// to Read() will load the page into memory
	if r.currSlot >= p.numEntries {
		r.currPage = nil
		r.currPageNum++
		r.currSlot = 0
	}

	return b
}

func (r *Reader) LoadPage() *Page {
	pageStartOffset := r.currPageNum * PAGE_SIZE
	b := make([]byte, PAGE_SIZE)
	f := r.heapFile.file
	_, err := f.ReadAt(b, pageStartOffset)
	if err != nil {
		return nil
	}
	numEntries := binary.LittleEndian.Uint16(b[0:2])
	startFreeSpace := binary.LittleEndian.Uint16(b[2:4])
	endFreeSpace := binary.LittleEndian.Uint16(b[4:6])
	slotArray := make([]Slot, 0)
	tuples := make([]map[string]string, 0)
	for i := 6; i < len(b); i += 4 {
		size := binary.LittleEndian.Uint16(b[i : i+2])
		if size == 0 {
			break
		}
		offset := binary.LittleEndian.Uint16(b[i+2 : i+4])
		slotArray = append(slotArray, Slot{offset, size})
		tuples = append(tuples, r.LoadTuple(b[offset:offset+size]))
	}

	return &Page{numEntries, startFreeSpace, endFreeSpace, slotArray, tuples}
}

func (r *Reader) LoadTuple(bytes []byte) map[string]string {
	tuple := make(map[string]string)
	fields := r.heapFile.fields
	for i := 0; i < len(fields); i++ {
		start := i * 4
		offset := binary.LittleEndian.Uint16(bytes[start : start+2])
		size := binary.LittleEndian.Uint16(bytes[start+2 : start+4])
		fieldName := fields[i]
		if offset > 0 {
			value := bytes[offset : offset+size]
			tuple[fieldName] = string(value)
		}
	}
	return tuple
}
