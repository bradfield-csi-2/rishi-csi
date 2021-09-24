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
	id := "5000"
	pred := func(r row) bool {
		return r["id"] == id
	}
	s := newSeqScanNode()
	root := newSelectionNode(pred, s)

	rows := Execute(root)
	movie := rows[0]

	if movie["id"] != id {
		t.Fatalf("SELECT * FROM movies WHERE id = 5000\nReturned id %s, wanted movieId %s",
			movie["id"],
			id,
		)
	}
}

func TestProjection(t *testing.T) {
	id := "5000"
	pred := func(r row) bool {
		return r["id"] == id
	}
	cols := []string{"title"}
	s := newSeqScanNode()
	sel := newSelectionNode(pred, s)
	root := newProjectionNode(cols, sel)

	rows := Execute(root)
	movie := rows[0]

	for key := range movie {
		found := false
		for _, col := range cols {
			if key == col {
				found = true
			}
		}

		if !found {
			t.Fatalf("SELECT title FROM movies WHERE id = 5000\nReturned map %v, wanted %v",
				movie,
				map[string]string{"title": "Medium Cool (1969)"},
			)
		}
	}
}
