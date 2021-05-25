package rocks

/*
#cgo CFLAGS: -I/opt/homebrew/Cellar/rocksdb/6.20.3/include
#cgo LDFLAGS: -L/opt/homebrew/Cellar/rocksdb/6.20.3/lib -lrocksdb -lz -lbz2 -lpthread
#include <strings.h>
#include <rocksdb/c.h>
*/
import "C"
import "fmt"

type database struct {
	db *C.rocksdb_t
}

func (db *database) Get(key string) string {
	k := C.CString(key)
	err := C.CString("")
	length := C.size_t(0)

	readOptions := C.rocksdb_readoptions_create()
	val := C.rocksdb_get(db.db, readOptions, k, C.strlen(k), &length, &err)
	return C.GoString(val)
}

func (db *database) Put(key, val string) {
	fmt.Printf("Putting %s: %s\n", key, val)

	// Convert to C Strings
	k := C.CString(key)
	v := C.CString(val)
	err := C.CString("")

	writeOptions := C.rocksdb_writeoptions_create()
	C.rocksdb_put(db.db, writeOptions, k, C.strlen(k), v, C.strlen(v)+1, &err)
}

func CreateDB() *database {
	// Set up options
	options := C.rocksdb_options_create()
	C.rocksdb_options_set_create_if_missing(options, 1)

	dbpath := C.CString("/tmp/rocksdb_test")
	err := C.CString("")

	return &database{db: C.rocksdb_open(options, dbpath, &err)}
}
