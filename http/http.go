package http

import (
	"github.com/kasheemlew/distribute_cache/cache"
	"github.com/kasheemlew/distribute_cache/server"
)

func New(c cache.Cache) *server.Server {
	return &server.Server{c}
}
