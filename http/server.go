package http

import (
	"net/http"

	"github.com/kasheemlew/distribute_cache/cache"
	"github.com/kasheemlew/distribute_cache/cluster"
)

// Server is server
type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen(port string) {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.Handle("/cluster", s.clusterHandler())
	http.ListenAndServe(s.Addr()+port, nil)
}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}

func (s *Server) statusHandler() http.Handler {
	return &statusHandler{s}
}

func (s *Server) clusterHandler() http.Handler {
	return &clusterHandler{s}
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
