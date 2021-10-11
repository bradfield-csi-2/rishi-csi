package main

import (
	"bytes"
	"io"
	"testing"
)

func TestFileScanner(t *testing.T) {
	var (
		columns = []string{"first_name", "last_name"}

		tuples = []Tuple{
			newTuple(
				"first_name", "john",
				"last_name", "doe"),
			newTuple(
				"first_name", "jane",
				"last_name", "doe"),
			newTuple(
				"first_name", "erica",
				"last_name", "doe"),
		}
	)

	var (
		fileReader  = NewMemoryTestFile(t, columns, tuples)
		fileScanner = NewFileScanner(fileReader)
	)
	for _, tuple := range tuples {
		assertEq(t, true, fileScanner.Next())
		assertEq(t, tuple, fileScanner.Execute())
	}
}

// NewMemoryTestFile creates a new Bradfield file format file with the provided tuples
// that is stored in memory.
func NewMemoryTestFile(t *testing.T, columns []string, tuples []Tuple) io.Reader {
	var (
		buf    = bytes.NewBuffer(nil)
		writer = NewFileWriter(columns, len(tuples), buf)
	)

	for _, tuple := range tuples {
		if err := writer.Append(tuple); err != nil {
			t.Fatalf("error appending tuple: %v, err: %v", tuple, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("error closing writer: %v", err)
	}

	return buf
}
