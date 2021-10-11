package main

import (
	"encoding/csv"
	"log"
	"os"
)

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
}
