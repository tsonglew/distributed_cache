package cache

import (
	"log"
)

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
}

func New(typ string) Cache {
	var c Cache
	switch typ {
	case "inmemory":
		c = newInMemoryCache()
	case "rocksdb":
		c = newRocksdbCache()
	default:
		panic("unkown cache type: " + typ)
	}
	if c == nil {
		panic("unkown cache type: " + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}
