package main

import (
	"github.com/kasheemlew/distribute_cache/cache"
	"github.com/kasheemlew/distribute_cache/http"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen(":12345")
}
