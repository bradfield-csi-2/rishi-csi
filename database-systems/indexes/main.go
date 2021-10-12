package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
)

type RecordsByTitle [][]string

func (r RecordsByTitle) Len() int {
	return len(r)
}

func (r RecordsByTitle) Less(i, j int) bool {
	return r[i][1] < r[j][1]
}

func (r RecordsByTitle) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func main() {
	// Read the movies CSV
	columns := []string{"id", "title", "genres"}
	f, err := os.Create("movies_db")
	if err != nil {
		return
	}
	defer f.Close()

	moviesCsv, _ := os.Open("movies.csv")
	records, _ := csv.NewReader(moviesCsv).ReadAll()
	writer := NewFileWriter(columns, len(records)-1, f)
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		tuple := newTuple("id", record[0], "title", record[1], "genres", record[2])
		if err := writer.Append(tuple); err != nil {
			log.Fatalf("error appending tuple: %v, err: %v", tuple, err)
		}
	}
	if err := writer.Close(); err != nil {
		log.Fatalf("error closing writer: %v", err)
	}
	indexScanner := NewIndexScanner(bytes.NewBuffer(nil), buildIndex(records), "")

	for indexScanner.Next() {
		fmt.Printf("%v", indexScanner.Execute())
	}
}

func buildIndex(records [][]string) []Tuple {
	tuples := make([]Tuple, len(records))
	sort.Sort(RecordsByTitle(records))
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		tuples[i] = newTuple("id", record[0], "title", record[1], "genres", record[2])
	}

	return tuples
}
