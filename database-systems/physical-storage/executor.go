package main

import (
	"bytes"
	"encoding/binary"
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

// Each relation is stored in one file (more or less), as an array of pages
// that are of the size block_size (usually 8kb)

// Build up an array of Pages and write out the file
// Pages contain records
// Need a separate schema file
// When writing,

type Writer struct {
	heapFile *HeapFile
	relation string
}

type Reader struct {
	heapFile    *HeapFile
	relation    string
	currPage    *Page
	currPageNum int64
	currSlot    int
}

type HeapFile struct {
	file     *os.File
	pages    []*Page
	numPages int64
}

// 8 kb pages
type Page struct {
	numEntries     uint16
	startFreeSpace uint16
	endFreeSpace   uint16
	slotArray      []Slot
}

type Slot struct {
	size   uint16
	offset uint16
}

func newWriter(relName string) *Writer {
	// TODO: Check if file exists before creating
	f, err := os.Create(relName)
	if err != nil {
		return nil
	}

	h := newHeapFile(f)
	return &Writer{heapFile: h, relation: relName}
}

func newReader(relName string) *Reader {
	f, err := os.Open(relName)
	if err != nil {
		return nil
	}

	h := newHeapFile(f)
	return &Reader{
		heapFile:    h,
		relation:    relName,
		currPage:    nil,
		currPageNum: 0,
		currSlot:    0,
	}
}

func newHeapFile(f *os.File) *HeapFile {
	return &HeapFile{file: f, pages: []*Page{newPage()}, numPages: 1}
}

func newPage() *Page {
	return &Page{
		numEntries:     0,
		startFreeSpace: 6,
		endFreeSpace:   0xffff,
		slotArray:      make([]Slot, 0),
	}
}

var PAGE_SIZE int64 = 0xffff

func (w *Writer) Write(r row) {
	p := w.heapFile.pages[w.heapFile.numPages-1]
	freeSpace := p.endFreeSpace - p.startFreeSpace

	rowBytes := make([]byte, 0)
	for _, v := range r {
		rowBytes = append(rowBytes, []byte(v)...)
	}
	rowSize := uint16(len(rowBytes))

	if freeSpace < rowSize {
		// Add a new page
		p = newPage()
		w.heapFile.numPages += 1
		w.heapFile.pages = append(w.heapFile.pages, p)
	}

	p.numEntries++
	slot := Slot{size: rowSize, offset: p.endFreeSpace - rowSize}

	p.startFreeSpace += 4 // New slot array entry is four bytes
	p.endFreeSpace -= rowSize

	// Will append to the slot array -- not handling deletes right now
	p.slotArray = append(p.slotArray, slot)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, p.numEntries)
	binary.Write(buf, binary.LittleEndian, p.startFreeSpace)
	binary.Write(buf, binary.LittleEndian, p.endFreeSpace)
	for _, s := range p.slotArray {
		binary.Write(buf, binary.LittleEndian, s.size)
		binary.Write(buf, binary.LittleEndian, s.offset)
	}
	headerBytes := buf.Bytes()

	// Now flush this page back to disk
	pageStartOffset := (w.heapFile.numPages - 1) * PAGE_SIZE
	byteOffsetInFile := pageStartOffset + int64(p.endFreeSpace)

	w.heapFile.file.WriteAt(headerBytes, pageStartOffset)
	w.heapFile.file.WriteAt(rowBytes, byteOffsetInFile)
}

func (r *Reader) Read() []byte {
	p := r.currPage
	pageStartOffset := r.currPageNum * PAGE_SIZE

	// If page is empty, read it into memory
	if p == nil {
		b := make([]byte, PAGE_SIZE)
		_, err := r.heapFile.file.ReadAt(b, pageStartOffset)
		if err != nil {
			return nil
		}
		numEntries := binary.LittleEndian.Uint16(b[0:2])
		startFreeSpace := binary.LittleEndian.Uint16(b[2:4])
		endFreeSpace := binary.LittleEndian.Uint16(b[4:6])
		slotArray := make([]Slot, 0)
		for i := 6; i < len(b); i += 4 {
			size := binary.LittleEndian.Uint16(b[i : i+2])
			if size == 0 {
				break
			}
			offset := binary.LittleEndian.Uint16(b[i+2 : i+4])
			slotArray = append(slotArray, Slot{offset, size})
		}

		p = &Page{numEntries, startFreeSpace, endFreeSpace, slotArray}
	}

	slot := p.slotArray[r.currSlot]
	b := make([]byte, slot.size)
	byteOffsetInFile := pageStartOffset + int64(slot.offset)
	_, err := r.heapFile.file.ReadAt(b, byteOffsetInFile)
	if err != nil {
		return nil
	}
	return b
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
