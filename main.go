package main

import (
	"flag"
	"log"

	"github.com/kasheemlew/distribute_cache/cache"
	"github.com/kasheemlew/distribute_cache/http"
	"github.com/kasheemlew/distribute_cache/tcp"
)

func main() {
	typ := flag.String("type", "inmemory", "memory type")
	flag.Parse()
	c := cache.New(*typ)
	log.Printf("using memory type: %s\n", *typ)
	go tcp.New(c).Listen(":1235")
	http.New(c).Listen(":1234")
}
