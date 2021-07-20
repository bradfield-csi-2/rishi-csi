package table

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
)

type Item struct {
	Key, Value string
}

// Given a sorted list of key/value pairs, write them out according to the format you designed.
func Build(path string, sortedItems []Item) error {
	data := new(bytes.Buffer)
	index := new(bytes.Buffer)
	var dataOffset uint16 = 0
	var indexLen uint16 = 0
	for _, item := range sortedItems {
		keyLen := uint16(len(item.Key))
		valLen := uint16(len(item.Value))

		// Write an index entry, containing the value's offset into the data
		// portion of the file, the length of the value, length of the key, and the
		// key itself
		entry := &IndexEntry{dataOffset, valLen, keyLen}
		binary.Write(index, binary.BigEndian, entry)
		index.WriteString(item.Key)

		// Write the value in data portion of file
		data.WriteString(item.Value)

		indexLen += (keyLen + 6)
		dataOffset += valLen
	}
	file := new(bytes.Buffer)
	// Place the length of the index as the first two bytes in the file, then
	// write the index and data
	binary.Write(file, binary.BigEndian, indexLen)
	file.Write(index.Bytes())
	file.Write(data.Bytes())
	ioutil.WriteFile(path, file.Bytes(), 0644)
	return nil
}

// A Table provides efficient access into sorted key/value data that's organized according
// to the format you designed.
//
// Although a Table shouldn't keep all the key/value data in memory, it should contain
// some metadata to help with efficient access (e.g. size, index, optional Bloom filter).
type Table struct {
	indexLen uint16
	index    map[string]IndexEntry
	file     *os.File
}

type IndexEntry struct {
	Offset uint16
	ValLen uint16
	KeyLen uint16
}

// Prepares a Table for efficient access. This will likely involve reading some metadata
// in order to populate the fields of the Table struct.
func LoadTable(path string) (*Table, error) {
	f, _ := os.Open(path)
	d := make([]byte, 2)
	f.Read(d)
	indexLen := binary.BigEndian.Uint16(d)

	indexBytes := make([]byte, indexLen)
	f.ReadAt(indexBytes, 2)
	index := make(map[string]IndexEntry)

	var i uint16 = 0
	for i < indexLen {
		entry := new(IndexEntry)
		r := bytes.NewReader(indexBytes[i : i+6])
		binary.Read(r, binary.BigEndian, entry)
		i += 6

		key := string(indexBytes[i : i+entry.KeyLen])
		index[key] = *entry
		i += entry.KeyLen
	}

	t := &Table{
		file:     f,
		indexLen: indexLen,
		index:    index,
	}

	return t, nil
}

func (t *Table) Get(key string) (string, bool, error) {
	entry, ok := t.index[key]
	if !ok {
		return "", false, nil
	}
	value := make([]byte, entry.ValLen)
	t.file.ReadAt(value, int64(entry.Offset+t.indexLen+2))
	return string(value), true, nil
}

func (t *Table) RangeScan(startKey, endKey string) (Iterator, error) {
	return nil, nil
}

type Iterator interface {
	// Advances to the next item in the range. Assumes Valid() == true.
	Next()

	// Indicates whether the iterator is currently pointing to a valid item.
	Valid() bool

	// Returns the Item the iterator is currently pointing to. Assumes Valid() == true.
	Item() Item
}
