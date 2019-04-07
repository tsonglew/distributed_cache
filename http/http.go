package http

import (
	"github.com/kasheemlew/distribute_cache/cache"
)

func New(c cache.Cache) *Server {
	return &Server{c}
}
