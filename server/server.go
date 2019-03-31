package server

import (
	"net/http"

	"github.com/kasheemlew/distribute_cache/cache"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen(port string) {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.ListenAndServe(port, nil)
}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}

func (s *Server) statusHandler() http.Handler {
	return &statusHandler{s}
}
