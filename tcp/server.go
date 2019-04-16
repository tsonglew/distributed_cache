package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/kasheemlew/distribute_cache/cache"
	"github.com/kasheemlew/distribute_cache/cluster"
)

type Server struct {
	cache.Cache
	cluster.Node
}

type result struct {
	v []byte
	e error
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}

func (s *Server) Listen(port string) {
	l, err := net.Listen("tcp", s.Addr()+port)
	if err != nil {
		log.Fatal(err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		go s.process(conn)
	}
}

func (s *Server) readLen(r *bufio.Reader) (int, error) {
	lenStr, err := r.ReadString(' ')
	if err != nil {
		return -1, err
	}
	lenNum, err := strconv.Atoi(strings.TrimSpace(lenStr))
	if err != nil {
		return -1, err
	}
	return lenNum, nil
}

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	keyLen, err := s.readLen(r)
	if err != nil {
		return "", err
	}
	keyBuf := make([]byte, keyLen)
	_, err = io.ReadFull(r, keyBuf)
	if err != nil {
		return "", err
	}
	key := string(keyBuf)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", errors.New("redirect " + addr)
	}
	return key, nil
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	keyLen, err := s.readLen(r)
	if err != nil {
		return "", []byte{}, err
	}
	valueLen, err := s.readLen(r)
	if err != nil {
		return "", []byte{}, err
	}
	keyBuf := make([]byte, keyLen)
	_, err = io.ReadFull(r, keyBuf)
	if err != nil {
		return "", []byte{}, err
	}
	key := string(keyBuf)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", nil, errors.New("redirect " + addr)
	}
	valueBuf := make([]byte, valueLen)
	_, err = io.ReadFull(r, valueBuf)
	if err != nil {
		return "", []byte{}, err
	}
	return key, valueBuf, nil
}

func sendResponse(value []byte, err error, conn net.Conn) error {
	var respStr string
	if err != nil {
		respStr = fmt.Sprintf("-%d %s", len(err.Error()), err.Error())
	} else {
		respStr = fmt.Sprintf("%d %s", len(value), string(value))
	}
	defer log.Printf("sent response: %s", respStr)
	_, writeErr := conn.Write([]byte(respStr))
	return writeErr
}

func (s *Server) get(conn net.Conn, r *bufio.Reader, chChan chan chan *result) {
	ch := make(chan *result)
	chChan <- ch
	key, err := s.readKey(r)
	if err != nil {
		ch <- &result{nil, err}
		return
	}
	go func() {
		v, err := s.Get(key)
		ch <- &result{v, err}
	}()
}

func (s *Server) set(conn net.Conn, r *bufio.Reader, chChan chan chan *result) {
	ch := make(chan *result)
	chChan <- ch
	k, v, err := s.readKeyAndValue(r)
	if err != nil {
		ch <- &result{nil, err}
		return
	}
	go func() {
		ch <- &result{nil, s.Set(k, v)}
	}()
}

func (s *Server) del(conn net.Conn, r *bufio.Reader, chChan chan chan *result) {
	ch := make(chan *result)
	key, err := s.readKey(r)
	if err != nil {
		ch <- &result{nil, err}
		return
	}
	go func() {
		ch <- &result{nil, s.Del(key)}
	}()
}

func (s *Server) process(conn net.Conn) {
	chChan := make(chan chan *result, 5000)
	r := bufio.NewReader(conn)
	go reply(conn, chChan)
	for {
		op, err := r.ReadByte()
		if err != nil {
			log.Printf("connection closed: %s\n", err)
		}
		switch op {
		case 'S':
			s.set(conn, r, chChan)
		case 'G':
			s.get(conn, r, chChan)
		case 'D':
			s.del(conn, r, chChan)
		default:
			log.Printf("connection closed for wrong op: %b\n", op)
			return
		}
		if err != nil {
			log.Printf("connection closed: %s", err)
			return
		}
	}
}

func reply(conn net.Conn, chChan chan chan *result) {
	defer conn.Close()
	for {
		ch, open := <-chChan
		if !open {
			return
		}
		r := <-ch
		err := sendResponse(r.v, r.e, conn)
		if err != nil {
			log.Println("send response error:", err)
			return
		}
	}
}
