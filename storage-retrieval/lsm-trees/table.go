package table

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"sort"
)

const BLOCK_CAPACITY = 4096

type Item struct {
	Key, Value string
}

// Given a sorted list of key/value pairs, write them out according to the format you designed.
func Build(path string, sortedItems []Item) error {
	data := new(bytes.Buffer)
	index := new(bytes.Buffer)
	blockSize := 0
	var blockOffset uint64 = 0
	var lastKey string
	for _, item := range sortedItems {
		buf := new(bytes.Buffer)
		// Write the keys and values to a temporary buffer
		keyLen, err1 := buf.WriteString(item.Key)
		valLen, err2 := buf.WriteString(item.Value)
		elemSize := (4 + keyLen + valLen)

		// If the current element will not fit in block, pad out the remaining
		// bytes with 0s, start a new block and write the previous key to the index
		if (blockSize + elemSize) > BLOCK_CAPACITY {
			data.Write(make([]byte, BLOCK_CAPACITY-blockSize))
			// binary.Write(index, binary.BigEndian, uint64(blockOffset))
			binary.Write(index, binary.BigEndian, uint16(len(lastKey)))
			index.WriteString(lastKey)
			blockSize = 0
			blockOffset++
		}
		blockSize += elemSize

		// Write the elements to the data:
		// key length, value length, key, value
		err3 := binary.Write(data, binary.BigEndian, uint16(keyLen))
		err4 := binary.Write(data, binary.BigEndian, uint16(valLen))
		_, err5 := data.Write(buf.Bytes())

		lastKey = item.Key
		// Something went wrong writing to the buffers, bail and return the first
		// error that occured
		if err := check(err1, err2, err3, err4, err5); err != nil {
			return err
		}
	}
	// Pad out the final block and write it to the index
	data.Write(make([]byte, BLOCK_CAPACITY-blockSize))
	binary.Write(index, binary.BigEndian, uint16(len(lastKey)))
	index.WriteString(lastKey)

	// File looks like:
	// <-- data section   -->
	// <-- index section  -->
	// <-- footer section -->
	file := new(bytes.Buffer)
	_, err1 := file.Write(data.Bytes())
	_, err2 := file.Write(index.Bytes())
	err3 := binary.Write(file, binary.BigEndian, blockOffset+1)
	err4 := ioutil.WriteFile(path, file.Bytes(), 0644)

	// Will return whatever error happened first when writing the file
	return check(err1, err2, err3, err4)
}

func check(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// A Table provides efficient access into sorted key/value data that's organized according
// to the format you designed.
//
// Although a Table shouldn't keep all the key/value data in memory, it should contain
// some metadata to help with efficient access (e.g. size, index, optional Bloom filter).
type Table struct {
	index     []string
	numBlocks uint64
	keys      []string
	file      *os.File
}

// Prepares a Table for efficient access. This will likely involve reading some metadata
// in order to populate the fields of the Table struct.
func LoadTable(path string) (*Table, error) {
	f, err := os.Open(path)
	fi, err := f.Stat()
	fileSize := uint64(fi.Size())
	d := make([]byte, 8)
	_, err = f.ReadAt(d, int64(fileSize-8))
	numBlocks := binary.BigEndian.Uint64(d)
	dataSize := numBlocks * BLOCK_CAPACITY

	indexBytes := make([]byte, fileSize-dataSize-8)
	_, err = f.ReadAt(indexBytes, int64(dataSize))
	if err != nil {
		return nil, err
	}
	index := make([]string, numBlocks)
	keys := make([]string, 1)

	i := uint64(0)
	block := 0
	for i < uint64(len(indexBytes)) {
		keyLen := uint64(binary.BigEndian.Uint16(indexBytes[i : i+2]))
		key := string(indexBytes[i+2 : i+2+keyLen])
		index[block] = key
		if err != nil {
			return nil, err
		}
		i += (2 + keyLen)
		block++
	}

	t := &Table{
		file:      f,
		index:     index,
		numBlocks: numBlocks,
		keys:      keys,
	}

	return t, nil
}

func (t *Table) Get(key string) (string, bool, error) {
	// Binary Search the blocks to find the block where the key might reside
	blockOffset := uint64(sort.SearchStrings(t.index, key))
	// If the blockOffset is equal to the number of blocks, it means the key is
	// greater than all the elements in the table, and hence is not in the table
	if blockOffset == t.numBlocks {
		return "", false, nil
	}

	// Read the block into memory and linear scan to find the key. This is OK
	// because each block is fixed and small-ish (4kb). A binary search would
	// improve the constant factor, so we can consider this a TODO optimization
	block := make([]byte, BLOCK_CAPACITY)
	t.file.ReadAt(block, int64(blockOffset*BLOCK_CAPACITY))
	i := uint64(0)
	for i < BLOCK_CAPACITY {
		keyLen := uint64(binary.BigEndian.Uint16(block[i : i+2]))
		valLen := uint64(binary.BigEndian.Uint16(block[i+2 : i+4]))
		k := string(block[i+4 : i+4+keyLen])
		if k == key {
			val := string(block[i+4+keyLen : i+4+keyLen+valLen])
			return val, true, nil
		}
		i += (4 + keyLen + valLen)
	}
	return "", false, nil
}

func (t *Table) RangeScan(startKey, endKey string) (Iterator, error) {
	index := sort.SearchStrings(t.keys, startKey)
	return &tableIterator{t, index, endKey}, nil
}

type tableIterator struct {
	table  *Table
	index  int
	endKey string
}

func (t *tableIterator) Valid() bool {
	return t.index < len(t.table.keys) && t.table.keys[t.index] <= t.endKey
}

func (t *tableIterator) Item() Item {
	key := t.table.keys[t.index]
	val, _, _ := t.table.Get(key)
	return Item{key, val}
}

func (t *tableIterator) Next() {
	t.index++
}

type Iterator interface {
	// Advances to the next item in the range. Assumes Valid() == true.
	Next()

	// Indicates whether the iterator is currently pointing to a valid item.
	Valid() bool

	// Returns the Item the iterator is currently pointing to. Assumes Valid() == true.
	Item() Item
}
