package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	conn, err := net.Dial("tcp", "localhost:1236")
	if err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Print("input: ")
		sWithEnd, err := reader.ReadString('\n')
		s := strings.TrimRight(sWithEnd, "\n")
		if err != nil {
			log.Fatal(err)
		}
		infos := strings.Split(s, " ")
		op := infos[0]
		words := make([]string, len(infos)-1)
		lens := make([]string, len(infos)-1)
		for i := 1; i < len(infos); i++ {
			words[i-1] = infos[i]
			lens[i-1] = strconv.Itoa(len(infos[i]))
		}
		msg := fmt.Sprintf("%s%s %s", op, strings.Join(lens, " "), strings.Join(words, ""))
		log.Println(msg)
		fmt.Fprint(conn, msg)
		connReader := bufio.NewReader(conn)
		resLenStr, err := connReader.ReadString(' ')
		if err != nil {
			log.Fatal(err)
		}
		resLen, err := strconv.Atoi(strings.TrimSpace(resLenStr))
		if err != nil {
			log.Fatal(err)
		}
		var resp []byte
		if resLen < 0 {
			resp = make([]byte, -resLen)
		} else {
			resp = make([]byte, resLen)
		}
		_, err = io.ReadFull(connReader, resp)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(resp))
	}
}
