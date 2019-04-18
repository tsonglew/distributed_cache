package cache

import (
	"unsafe"
)

type rocksdbScanner struct {
	i           *C.rocksdb_iterator_t
	initialized bool
}

func (s *rocksdbScanner) Close() {
	C.rocksdb_iter_destroy(s.i)
}

func (s *rocksdbScanner) Scan() bool {
	if !s.initialized {
		C.rocksdb_iter_seek_to_first(s.i)
		s.initialized = true
	} else {
		C.rocksdb_iter_next(s.i)
	}
	return C.rocksdb_iter_valid(s.i) != 0
}

func (s *rocksdbScanner) Key() string {
	var length C.size_t
	k := C.rocksdb_iter_key(s.i, &length)
	return C.GoString(k)
}

func (s *rocksdbScanner) Value() string {
	var length C.size_t
	v := C.rocksdb_iter_value(s.i, &length)
	return C.GoBytes(unsafe.Pointer(v), C.int(length))
}
