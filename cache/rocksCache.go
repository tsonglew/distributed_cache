package cache

// #include "rocksdb/c.h"
// #include <stdlib.h>
// #cgo CFLAGS: -I${SRCDIR}/../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"

import (
	"errors"
	"regexp"
	"runtime"
	"strconv"
	"unsafe"
)

type rocksdbCache struct {
	db *C.rocksdb_t              `rocksDB type`
	ro *C.rocksdb_readoptions_t  `rocksDB read options`
	wo *C.rocksdb_writeoptions_t `rocksDB write options`
	e  *C.char                   `error string from rocksDB C API`
}

func newRocksdbCache() *rocksdbCache {
	options := C.rocksdb_options_create()                                    // brand new options pointer
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU())) // parallelism threads num
	C.rocksdb_options_set_create_if_missing(options, 1)
	var e *C.char
	db := C.rocksdb_open(options, C.CString("/mnt/rocksdb"), &e)
	if e != nil {
		panic(C.GoString(e))
	}
	C.rocksdb_options_destroy(options)
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), C.rocksdb_writeoptions_create(), e}
}

func (c *rocksdbCache) Get(key string) ([]byte, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	var valueLen C.size_t
	value := C.rocksdb_get(c.db, c.ro, k, C.size_t(len(key)), &valueLen, &c.e)
	defer C.free(unsafe.Pointer(value))
	if c.e != nil {
		return nil, errors.New(C.GoString(c.e))
	}
	return C.GoBytes(unsafe.Pointer(value), C.int(valueLen)), nil
}

func (c *rocksdbCache) Set(key string, value []byte) error {
	k := C.CString(key)
	v := C.CBytes(value)
	defer C.free(unsafe.Pointer(k))
	defer C.free(unsafe.Pointer(v))

	C.rocksdb_put(c.db, c.wo, k, C.size_t(len(key)), (*C.char)(v), C.size_t(len(value)), &c.e)
	if c.e != nil {
		return errors.New(C.GoString(c.e))
	}
	return nil
}

func (c *rocksdbCache) Del(key string) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	C.rocksdb_delete(c.db, c.wo, k, C.size_t(len(key)), &c.e)
	if c.e != nil {
		return errors.New(C.GoString(c.e))
	}
	return nil
}

func (c *rocksdbCache) GetStat() Stat {
	k := C.CString("rocksdb.aggregated-table-properties")
	defer C.free(unsafe.Pointer(k))
	v := C.rocksdb_property_value(c.db, k)
	defer C.free(unsafe.Pointer(v))
	p := C.GoString(v)
	r := regexp.MustCompile(`([^;]+)=([^;]+);`)
	s := Stat{}
	for _, submatches := range r.FindAllStringSubmatch(p, -1) {
		switch submatches[1] {
		case " # entries":
			s.Count, _ = strconv.ParseInt(submatches[2], 10, 64)
		case " raw key size":
			s.KeySize, _ = strconv.ParseInt(submatches[2], 10, 64)
		case " raw value size":
			s.ValueSize, _ = strconv.ParseInt(submatches[2], 10, 64)
		}
	}
	return s
}
