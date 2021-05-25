package rocks

/*
#cgo CFLAGS: -I/opt/homebrew/Cellar/rocksdb/6.20.3/include
#cgo LDFLAGS: -L/opt/homebrew/Cellar/rocksdb/6.20.3/lib -lrocksdb -lz -lbz2 -lpthread
#include <stdlib.h>
#include <strings.h>
#include <rocksdb/c.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type database struct {
	db           *C.rocksdb_t
	options      *C.rocksdb_options_t
	writeOptions *C.rocksdb_writeoptions_t
	readOptions  *C.rocksdb_readoptions_t
	err          *C.char
}

func (db *database) Get(key string) (string, error) {
	length := C.size_t(0)
	k := C.CString(key)
	v := C.rocksdb_get(db.db, db.readOptions, k, C.strlen(k), &length, &db.err)
	defer func() {
		C.free(unsafe.Pointer(k))
		C.free(unsafe.Pointer(v))
	}()

	if db.err != nil {
		return "", fmt.Errorf("rocks: %s", C.GoString(db.err))
	}
	return C.GoString(v), nil
}

func (db *database) Put(key, val string) error {
	k := C.CString(key)
	v := C.CString(val)
	defer func() {
		C.free(unsafe.Pointer(k))
		C.free(unsafe.Pointer(v))
	}()

	C.rocksdb_put(db.db, db.writeOptions, k, C.strlen(k), v, C.strlen(v)+1, &db.err)
	if db.err != nil {
		return fmt.Errorf("rocks: %s", C.GoString(db.err))
	}

	return nil
}

func CreateDB() (*database, error) {
	options := C.rocksdb_options_create()
	C.rocksdb_options_set_create_if_missing(options, 1)
	writeOptions := C.rocksdb_writeoptions_create()
	readOptions := C.rocksdb_readoptions_create()

	var err *C.char
	dbpath := C.CString("/tmp/rocksdb_test")
	rocksdb := C.rocksdb_open(options, dbpath, &err)

	if err != nil {
		return nil, fmt.Errorf("rocks: %s", C.GoString(err))
	}

	return &database{db: rocksdb,
		options:      options,
		writeOptions: writeOptions,
		readOptions:  readOptions,
		err:          err,
	}, nil
}

func DestroyDB(db *database) {
	C.free(unsafe.Pointer(db.err))
	C.rocksdb_readoptions_destroy(db.readOptions)
	C.rocksdb_writeoptions_destroy(db.writeOptions)
	C.rocksdb_options_destroy(db.options)
	C.rocksdb_close(db.db)
}
