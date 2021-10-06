package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
)

type row map[string]string

type Node interface {
	Next() row
}

type SortNode struct {
	sortCols []string
	data     []row
	nRows    int
	child    Node
	cursor   int
}

type ProjectionNode struct {
	cols  []string
	child Node
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

type SeqScanNode struct {
	data   []row
	nRows  int
	cursor int
	child  Node
}

type PredFn func(row) bool

func newSortNode(sortCols []string, child Node) *SortNode {
	return &SortNode{sortCols: sortCols, data: nil, child: child}
}

func newProjectionNode(cols []string, child Node) *ProjectionNode {
	return &ProjectionNode{cols: cols, child: child}
}

func newLimitNode(limit int, child Node) *LimitNode {
	return &LimitNode{limit: limit, cursor: 0, child: child}
}

func newSelectionNode(pred PredFn, child Node) *SelectionNode {
	return &SelectionNode{pred: pred, child: child}
}

func newSeqScanNode() *SeqScanNode {
	data := make([]row, 0)
	f, err := os.Open("data/movies.csv")
	if err != nil {
		fmt.Printf("Could not open movies file.")
		return nil
	}
	r := csv.NewReader(f)
	r.Read() // Skip header
	records, err := r.ReadAll()
	if err != nil {
		fmt.Printf("Could not read movies file.")
		return nil
	}
	for _, rec := range records {
		movie := map[string]string{"id": rec[0], "title": rec[1], "genres": rec[2]}
		data = append(data, movie)
	}

	return &SeqScanNode{
		data:   data,
		nRows:  len(data),
		cursor: 0,
		child:  nil,
	}
}

func (s *SortNode) Next() row {
	// Gather up all the records from the node below if we haven't sorted yet
	if s.data == nil {
		rows := make([]row, 0)
		for r := s.child.Next(); r != nil; r = s.child.Next() {
			rows = append(rows, r)
		}
		sort.Slice(rows, func(i, j int) bool {
			var val bool
			for _, sortCol := range s.sortCols {
				parts := strings.Split(sortCol, ":")
				col := parts[0]
				sortAsc := true
				if len(parts) > 1 && parts[1] == "desc" {
					sortAsc = false
				}
				if sortAsc {
					if rows[i][col] < rows[j][col] {
						return true
					}
					if rows[i][col] > rows[j][col] {
						return false
					}
				} else {
					if rows[i][col] > rows[j][col] {
						return true
					}
					if rows[i][col] < rows[j][col] {
						return false
					}
				}
			}
			return val
		})
		s.data = rows
		s.nRows = len(rows)
	}

	if s.cursor >= s.nRows {
		return nil
	}

	// Now cursor through the sorted results
	row := s.data[s.cursor]
	s.cursor++
	return row
}

func (p *ProjectionNode) Next() row {
	row := p.child.Next()
	if row == nil {
		return nil
	}
	proj := make(map[string]string)
	for _, col := range p.cols {
		proj[col] = row[col]
	}
	return proj
}

func (l *LimitNode) Next() row {
	if l.cursor < l.limit {
		l.cursor++
		return l.child.Next()
	}
	return nil
}

func (s *SelectionNode) Next() row {
	for m := s.child.Next(); m != nil; m = s.child.Next() {
		if s.pred(m) {
			return m
		}
	}
	return nil
}

func (s *SeqScanNode) Next() row {
	if s.cursor >= s.nRows {
		return nil
	}

	row := s.data[s.cursor]
	s.cursor++
	return row
}

func Execute(root Node) []row {
	rows := make([]row, 0)
	for rec := root.Next(); rec != nil; rec = root.Next() {
		rows = append(rows, rec)
	}
	return rows
}

func main() {
	fmt.Println("Sample Query 1\n==============\nExecuting: SELECT title FROM movies WHERE id = 5000")

	id := "5000"
	pred := func(r row) bool {
		return r["id"] == id
	}
	cols := []string{"title"}
	s := newSeqScanNode()
	sel := newSelectionNode(pred, s)
	root := newProjectionNode(cols, sel)

	rows := Execute(root)
	for _, row := range rows {
		fmt.Printf("%+v\n", row)
	}

	fmt.Println()
	fmt.Println("Sample Query 2\n==============\nExecuting: SELECT title, genres FROM movies ORDER BY genres, title DESC LIMIT 3")

	sortCols := []string{"genres:asc", "title:desc"}
	projCols := []string{"title", "genres"}
	s = newSeqScanNode()
	sort := newSortNode(sortCols, s)
	l := newLimitNode(3, sort)
	root = newProjectionNode(projCols, l)

	rows = Execute(root)
	for _, row := range rows {
		fmt.Printf("%+v\n", row)
	}
}
