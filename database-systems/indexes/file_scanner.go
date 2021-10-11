package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

// FileScanner is a scan operator (implements the Operator interface) that reads
// Bradfield file format files one tuple at a time.
type FileScanner struct {
	header  *Header
	r       *byteReader
	numRead int
	next    Tuple
}

// NewFileScanner creates a new FileScanner that scans the provided file.
func NewFileScanner(r io.Reader) *FileScanner {
	return &FileScanner{
		r: newByteReader(r),
	}
}

// Next returns a boolean indicating whether the file contains an additional tuple.
func (f *FileScanner) Next() bool {
	if f.header == nil {
		f.readHeader()
	}

	if f.numRead < f.header.NumRows {
		f.readTuple()
		return true
	}

	return false
}

// Execute returns the next tuple, should only be called if a previous call to Next()
// returned true.
func (f *FileScanner) Execute() Tuple {
	return f.next
}

func (f *FileScanner) readHeader() {
	headerLength, err := binary.ReadUvarint(f.r)
	if err != nil {
		panic(fmt.Sprintf("FileScanner: error reading header length: %v", err))
	}

	headerBytes := make([]byte, headerLength)
	if _, err := io.ReadFull(f.r, headerBytes); err != nil {
		panic(fmt.Sprintf("FileScanner: error reading header bytes: %v", err))
	}

	header := &Header{}
	if err := json.Unmarshal(headerBytes, header); err != nil {
		panic(fmt.Sprintf("FileScanner: error unmarshaling header: %v", err))
	}

	f.header = header
}

func (f *FileScanner) readTuple() {
	tuple := Tuple{}

	for _, col := range f.header.ColumnNames {
		valLen, err := binary.ReadUvarint(f.r)
		if err != nil {
			panic(fmt.Sprintf("readTuple: error reading next value length: %v", err))
		}

		// In practice we wouldn't want to allocate a byte slice each time and then
		// immediately wrap a string to wrap it for performance reasons, but this is
		// fine for educational purposes.
		valBytes := make([]byte, valLen)
		if _, err := io.ReadFull(f.r, valBytes); err != nil {
			panic(fmt.Sprintf("readTuple: error reading value bytes: %v", err))
		}

		tuple.Values = append(tuple.Values, Value{
			Name:        col,
			StringValue: string(valBytes),
		})
	}

	f.next = tuple
}

// byteReader wraps an io.Reader so that it implements io.ByteReader.
type byteReader struct {
	io.Reader
	byteBuf []byte
}

// newByteReaders create a byteReader from an io.Reader.
func newByteReader(r io.Reader) *byteReader {
	return &byteReader{
		Reader:  r,
		byteBuf: make([]byte, 1),
	}
}

func (b *byteReader) ReadByte() (byte, error) {
	n, err := b.Reader.Read(b.byteBuf)
	if err != nil {
		return 0, fmt.Errorf("byteReader: ReadByte: error reading byte: %v", err)
	}
	if n != 1 {
		return 0, fmt.Errorf("byteReader: ReadByte: expected to read one byte, but read: %d", n)
	}

	return b.byteBuf[0], nil
}
