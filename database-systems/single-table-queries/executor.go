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
	Next() *movie
}

type SeqScanNode struct {
	data   []*movie
	nRows  int
	cursor int
	child  Node
}

type LimitNode struct {
	limit  int
	cursor int
	child  Node
}

type SelectionNode struct {
	pred  PredFn
	child Node
}

type PredFn func(*movie) bool

func newLimitNode(limit int, child Node) *LimitNode {
	return &LimitNode{limit: limit, cursor: 0, child: child}
}

func newSelectionNode(pred PredFn, child Node) *SelectionNode {
	return &SelectionNode{pred: pred, child: child}
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

func (l *LimitNode) Next() *movie {
	if l.cursor < l.limit {
		l.cursor++
		return l.child.Next()
	}
	return nil
}

func (s *SelectionNode) Next() *movie {
	for m := s.child.Next(); m != nil; m = s.child.Next() {
		if s.pred(m) {
			return m
		}
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

func Execute(root Node) []*movie {
	results := make([]*movie, 0)
	for row := root.Next(); row != nil; row = root.Next() {
		results = append(results, row)
	}
	return results
}

func main() {
	// Sample Query
	fmt.Println("Sample Query\n\nExecuting: SELECT * FROM movies WHERE id = 5000")
	pred := func(m *movie) bool {
		return m.movieId == "5000"
	}
	s := newSeqScanNode()
	rows := Execute(newSelectionNode(pred, s))
	for _, row := range rows {
		fmt.Printf("%+v\n", row)
	}
}
