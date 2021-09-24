package main

import "fmt"

type Node interface {
	Next() []string
}

type SeqScanNode struct {
	data   [][]string
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
	data := [][]string{
		{"a1", "a2", "a3"},
		{"b1", "b2", "b3"},
		{"c1", "c2", "c3"},
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

func (l *LimitNode) Next() []string {
	if l.cursor < l.limit {
		//fmt.Printf("c: %+v\n", c)
		l.cursor++
		return l.child.Next()
	}
	return nil
}

func (s *SeqScanNode) Next() []string {
	if s.cursor >= s.nRows {
		return nil
	}

	row := s.data[s.cursor]
	s.cursor++
	return row
}

func main() {
	fmt.Println("Executing: SELECT * FROM test LIMIT 2")
	s := newSeqScanNode()
	l := newLimitNode(2, s)
	for row := l.Next(); row != nil; row = l.Next() {
		fmt.Printf("%+v\n", row)
	}
}
