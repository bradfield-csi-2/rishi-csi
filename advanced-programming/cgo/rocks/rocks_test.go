package rocks

import (
	"testing"
)

func TestRocksDB(t *testing.T) {
	db := CreateDB()

	var tests = []struct {
		key string
		val string
	}{
		{"red", "#ff0000"},
		{"green", "#00ff00"},
		{"blue", "#0000ff"},
	}

	for _, test := range tests {
		db.Put(test.key, test.val)
		if got := db.Get(test.key); got != test.val {
			t.Errorf("db.Get(%s) returned %s, want %s", test.key, got, test.val)
		}
	}
}
