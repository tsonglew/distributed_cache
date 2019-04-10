package clients

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type tcpClient struct {
	net.Conn
	r *bufio.Reader
}

func (c *tcpClient) sendGet(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("G%d %s", klen, key)))
}

func (c *tcpClient) sendSet(key, value string) {
	klen := len(key)
	vlen := len(value)
	c.Write([]byte(fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)))
}

func (c *tcpClient) sendDel(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
}

func readLen(r *bufio.Reader) int {
	respLenStr, err := r.ReadString(' ')
	if err != nil {
		log.Println(err)
		return 0
	}
	respLen, err := strconv.Atoi(strings.TrimSpace(respLenStr))
	if err != nil {
		log.Println(err)
		return 0
	}
	return respLen
}

func (c *tcpClient) recvResponse() (string, error) {
	vlen := readLen(c.r)
	if vlen == 0 {
		return "", nil
	}
	if vlen < 0 {
		respErr := make([]byte, -vlen)
		_, err := io.ReadFull(c.r, respErr)
		if err != nil {
			return "", err
		}
		return "", errors.New(string(respErr))
	}
	value := make([]byte, vlen)
	_, err := io.ReadFull(c.r, value)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (c *tcpClient) Run(cmd *Cmd) {
	switch cmd.Name {
	case "get":
		c.sendGet(cmd.Key)
		cmd.Value, cmd.Error = c.recvResponse()
	case "set":
		c.sendSet(cmd.Key, cmd.Value)
		_, cmd.Error = c.recvResponse()
	case "del":
		c.sendDel(cmd.Key)
		_, cmd.Error = c.recvResponse()
	default:
		panic("unkown cmd name:" + cmd.Name)
	}
}

func (c *tcpClient) PipelineRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	for _, cmd := range cmds {
		switch cmd.Name {
		case "get":
			c.sendGet(cmd.Key)
		case "set":
			c.sendSet(cmd.Key, cmd.Value)
		case "del":
			c.sendDel(cmd.Key)
		default:
			panic("unkown cmd name:" + cmd.Name)
		}
	}
	for _, cmd := range cmds {
		cmd.Value, cmd.Error = c.recvResponse()
	}
}

func newTCPClient(server string) *tcpClient {
	c, err := net.Dial("tcp", server+":1235")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(c)
	return &tcpClient{c, r}
}
