package main

import (
	"github.com/kasheemlew/distribute_cache/cache"
	"github.com/kasheemlew/distribute_cache/http"
	"github.com/kasheemlew/distribute_cache/tcp"
)

func main() {
	c := cache.New("inmemory")
	go tcp.New(c).Listen(":1236")
	http.New(c).Listen(":8888")
}
