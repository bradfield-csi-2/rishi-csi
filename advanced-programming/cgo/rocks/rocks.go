package rocks

/*
#cgo CFLAGS: -I/opt/homebrew/Cellar/rocksdb/6.20.3/include
#cgo LDFLAGS: -L/opt/homebrew/Cellar/rocksdb/6.20.3/lib -lrocksdb -lz -lbz2 -lpthread
#include <stdlib.h>
#include <strings.h>
#include <rocksdb/c.h>
*/
import "C"
import "unsafe"

type database struct {
	db      *C.rocksdb_t
	options *C.rocksdb_options_t
	err     *C.char
}

func (db *database) Get(key string) string {
	k := C.CString(key)
	length := C.size_t(0)
	readOptions := C.rocksdb_readoptions_create()
	v := C.rocksdb_get(db.db, readOptions, k, C.strlen(k), &length, &db.err)
	val := C.GoString(v)

	defer func() {
		C.free(unsafe.Pointer(k))
		C.free(unsafe.Pointer(v))
		C.free(unsafe.Pointer(readOptions))
	}()

	return val
}

func (db *database) Put(key, val string) {
	k := C.CString(key)
	v := C.CString(val)
	writeOptions := C.rocksdb_writeoptions_create()
	defer func() {
		C.free(unsafe.Pointer(k))
		C.free(unsafe.Pointer(v))
		C.free(unsafe.Pointer(writeOptions))
	}()

	C.rocksdb_put(db.db, writeOptions, k, C.strlen(k), v, C.strlen(v)+1, &db.err)
}

func DestroyDB(db *database) {
	C.free(unsafe.Pointer(db.options))
	C.free(unsafe.Pointer(db.err))
	C.rocksdb_close(db.db)
}

func CreateDB() *database {
	options := C.rocksdb_options_create()
	C.rocksdb_options_set_create_if_missing(options, 1)

	dbpath := C.CString("/tmp/rocksdb_test")
	err := C.CString("")
	rocksdb := C.rocksdb_open(options, dbpath, &err)

	return &database{db: rocksdb, options: options, err: err}
}
