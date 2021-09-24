package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type movie struct {
	movieId string
	title   string
	genres  string
}

type Node interface {
	Next() movie
}

type SeqScanNode struct {
	data   []*movie
	nRows  int
	cursor int
	child  *Node
}

type LimitNode struct {
	limit  int
	cursor int
	child  *SeqScanNode
}

func newSeqScanNode() *SeqScanNode {
	data := make([]*movie, 0)
	f, err := os.Open("data/movies.csv")
	if err != nil {
		fmt.Printf("Could not open movies file.")
		return nil
	}
	r := csv.NewReader(f)
	r.Read() // Skip header
	movies, err := r.ReadAll()
	if err != nil {
		fmt.Printf("Could not read movies file.")
		return nil
	}
	for _, m := range movies {
		data = append(data, &movie{m[0], m[1], m[2]})
	}

	return &SeqScanNode{
		data:   data,
		nRows:  len(data),
		cursor: 0,
		child:  nil,
	}
}

func newLimitNode(limit int, child *SeqScanNode) *LimitNode {
	return &LimitNode{limit: limit, cursor: 0, child: child}
}

func (l *LimitNode) Next() *movie {
	if l.cursor < l.limit {
		l.cursor++
		return l.child.Next()
	}
	return nil
}

func (s *SeqScanNode) Next() *movie {
	if s.cursor >= s.nRows {
		return nil
	}

	row := s.data[s.cursor]
	s.cursor++
	return row
}

func main() {
	fmt.Println("Executing: SELECT * FROM test LIMIT 5")
	s := newSeqScanNode()
	l := newLimitNode(5, s)
	for row := l.Next(); row != nil; row = l.Next() {
		fmt.Printf("%+v\n", row)
	}
}
