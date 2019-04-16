package main

import (
	"flag"
	"log"

	"github.com/kasheemlew/distributed_cache/cache"
	"github.com/kasheemlew/distributed_cache/cluster"
	"github.com/kasheemlew/distributed_cache/http"
	"github.com/kasheemlew/distributed_cache/tcp"
)

func main() {
	typ := flag.String("type", "inmemory", "memory type")
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.Parse()
	log.Printf("using memory type: %s\n", *typ)
	log.Println("node address: ", *node)
	log.Println("cluster address: ", *clus)
	c := cache.New(*typ)
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	go tcp.New(c, n).Listen(":1235")
	http.New(c, n).Listen(":1234")
}
