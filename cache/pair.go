package cache

import "time"

type pair struct {
	key   string
	value []byte
}

type value struct {
	v       []byte
	created time.Time
}
