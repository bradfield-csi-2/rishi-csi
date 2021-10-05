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

func TestSortAndLimit(t *testing.T) {
	lim := 3
	sortCols := []string{"title"}
	s := newSeqScanNode()
	sort := newSortNode(sortCols, s)
	root := newLimitNode(lim, sort)

	rows := Execute(root)

	firstTitle := "\"Great Performances\" Cats (1998)"
	if rows[0]["title"] != firstTitle {
		t.Fatalf("SELECT * FROM movies ORDER BY title LIMIT 3\nReturned title %s as first movie, wanted %s",
			rows[0]["title"],
			firstTitle,
		)
	}

	if len(rows) != lim {
		t.Fatalf("SELECT * FROM movies ORDER BY title LIMIT 3\nReturned %d rows, want %d rows",
			len(rows),
			lim,
		)
	}
}

func TestSortDescending(t *testing.T) {
	lim := 3
	sortCols := []string{"id:desc"}
	s := newSeqScanNode()
	sort := newSortNode(sortCols, s)
	root := newLimitNode(lim, sort)

	rows := Execute(root)

	// This is the last ID by string but not numerically
	// TODO: Add types to the schema
	lastId := "99999"
	if rows[0]["id"] != lastId {
		t.Fatalf("SELECT * FROM movies ORDER BY id DESC\nReturned title %s as last id, wanted %s",
			rows[0]["id"],
			lastId,
		)
	}
}

func TestSortMultipleCols(t *testing.T) {
	sortCols := []string{"genres:asc", "title:desc"}
	projCols := []string{"title", "genres"}
	s := newSeqScanNode()
	sort := newSortNode(sortCols, s)
	root := newProjectionNode(projCols, sort)

	rows := Execute(root)

	firstGenre := "(no genres listed)"
	lastTitle := "Zanjeer (1973)"
	if rows[0]["genres"] != firstGenre || rows[0]["title"] != lastTitle {
		t.Fatalf("SELECT title, genres FROM movies ORDER BY genres, title DESC\nReturned first movie %v, wanted %v",
			rows[0],
			map[string]string{"title": lastTitle, "genre": firstGenre},
		)
	}
}

func TestReadWrite(t *testing.T) {
	wr := newWriter("movies_db")
	r := row{"title": "Toy Story", "year": "1995"}
	wr.Write(r)
}
