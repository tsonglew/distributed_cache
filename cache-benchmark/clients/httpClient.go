package clients

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type httpClient struct {
	*http.Client
	server string
}

func (c *httpClient) get(key string) string {
	resp, err := c.Get(c.server + key)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return ""
	case resp.StatusCode != http.StatusOK:
		panic(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (c *httpClient) set(key, value string) {
	req, err := http.NewRequest(http.MethodPut, c.server+key, strings.NewReader(value))
	if err != nil {
		log.Println(err)
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}

func (c *httpClient) Run(cmd *Cmd) {
	switch cmd.Name {
	case "get":
		cmd.Value = c.get(cmd.Key)
	case "set":
		c.set(cmd.Key, cmd.Value)
	default:
		panic("unkown cmd name: " + cmd.Name)
	}
}

func newHTTPClient(server string) *httpClient {
	// Transport is a struct used by clients to manage the underlying TCP connection
	// MaxIdleConnsPerHost restricts the no. of connections which clients has not closed
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, "http://" + server + ":1234/cache/"}
}

func (c *httpClient) PipelineRun([]*Cmd) {
	panic("httpClient PipelineRun not implemented")
}
