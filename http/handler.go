package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type cacheHandler struct {
	*Server
}

type statusHandler struct {
	*Server
}

type clusterHandler struct {
	*Server
}

func (c *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	escapedURL := r.URL.EscapedPath()
	key := strings.Split(escapedURL, "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		b, err := c.Get(key)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(b)
		return
	case http.MethodPut:
		b, err := ioutil.ReadAll(r.Body)
		log.Println("read from body: ", string(b))
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = c.Set(key, b); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	case http.MethodDelete:
		if err := c.Del(key); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(s.GetStat())
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (c *clusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	m := c.Members()
	b, e := json.Marshal(m)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
