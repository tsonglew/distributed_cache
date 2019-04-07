package tcp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/kasheemlew/distribute_cache/cache"
)

type Server struct {
	cache.Cache
}

func New(c cache.Cache) *Server {
	return &Server{c}
}

func (s *Server) Listen(port string) {
	l, err := net.Listen("tcp", port)
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
	return string(keyBuf), nil
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
	valueBuf := make([]byte, valueLen)
	_, err = io.ReadFull(r, valueBuf)
	if err != nil {
		return "", []byte{}, err
	}
	return string(keyBuf), valueBuf, nil
}

func (s *Server) sendResponse(value []byte, err error, conn net.Conn) error {
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

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	v, err := s.Get(key)
	return s.sendResponse(v, err, conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	k, v, err := s.readKeyAndValue(r)
	if err != nil {
		return err
	}
	err = s.Set(k, v)
	return s.sendResponse(nil, err, conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	err = s.Del(key)
	return s.sendResponse(nil, err, conn)
}

func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		op, err := r.ReadByte()
		if err != nil {
			log.Printf("connection closed: %s\n", err)
		}
		switch op {
		case 'S':
			err = s.set(conn, r)
		case 'G':
			err = s.get(conn, r)
		case 'D':
			err = s.del(conn, r)
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
