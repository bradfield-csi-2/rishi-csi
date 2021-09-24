package main

import "testing"

func TestLimit(t *testing.T) {
	lim := 5
	s := newSeqScanNode()
	root := newLimitNode(lim, s)

	rows := Execute(root)

	if len(rows) != lim {
		t.Fatalf("SELECT * FROM movies LIMIT 5\nReturned %d rows, want %d rows",
			len(rows),
			lim,
		)
	}
}

func TestSelection(t *testing.T) {
	pred := func(m *movie) bool {
		return m.movieId == "5000"
	}
	id := "5000"
	s := newSeqScanNode()
	root := newSelectionNode(pred, s)

	rows := Execute(root)
	movie := rows[0]

	if movie.movieId != id {
		t.Fatalf("SELECT * FROM movies WHERE id = 5000\nReturned id %s, wanted movieId %s",
			movie.movieId,
			id,
		)
	}
}
