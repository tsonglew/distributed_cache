package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kasheemlew/distribute_cache/cache-benchmark/clients"
)

var typ, server, operation string
var total, valueSize, threads, keyspacelen, pipelen int

func init() {
	flag.StringVar(&typ, "type", "redis", "cache server type")
	flag.StringVar(&server, "h", "localhost", "cache server address")
	flag.IntVar(&total, "n", 1000, "total number of requests")
	flag.IntVar(&valueSize, "d", 1000, "data size of SET/GET value in bytes")
	flag.IntVar(&threads, "c", 1, "number of parallel connections")
	flag.StringVar(&operation, "t", "set", "test set, could be get/set/mixed")
	flag.IntVar(&keyspacelen, "r", 0, "keyspacelen, use random keys from 0 to keyspacelen-1")
	flag.IntVar(&pipelen, "P", 1, "pipeline length")
	flag.Parse()
	fmt.Println("type is", typ)
	fmt.Println("server is ", server)
	fmt.Println("requests num is", total)
	fmt.Println("data size is", valueSize)
	fmt.Println("threads num is", threads)
	fmt.Println("operation is", operation)
	fmt.Println("keyspacelen is", keyspacelen)
	fmt.Println("pipeline length is", pipelen)

	rand.Seed(time.Now().UnixNano())
}

func main() {
	ch := make(chan *result, threads)
	res := &result{0, 0, 0, make([]statistic, 0)}
	start := time.Now()
	for i := 0; i < threads; i++ {
		res.addResult(<-ch)
	}
	d := time.Now().Sub(start)
	totalCount := res.getCount + res.missCount + res.setCount
	fmt.Println(res.getCount, "records get")
	fmt.Println(res.missCount, "records miss")
	fmt.Println(res.setCount, "records set")
	fmt.Println(d.Seconds(), "seconds total")
	statCountSum := 0
	statTimeSum := time.Duration(0)
	for b, s := range res.statBuckets {
		if s.count == 0 {
			continue
		}
		statCountSum += s.count
		statTimeSum += s.time
		fmt.Printf("%d%% requests < %d ms\n", statCountSum*100/totalCount, b+1)
	}
	fmt.Printf("%d usec average for each request\n", int64(statTimeSum/time.Microsecond)/int64(statCountSum))
	fmt.Printf("throughput is %f MB/s\n", float64((res.getCount+res.setCount)*valueSize)/1e6/d.Seconds())
	fmt.Printf("rps is %f\n", float64(totalCount)/float64(d.Seconds()))
}

func operate(id, count int, ch chan *result) {
	client := clients.New(typ, server)
	cmds := make([]*clients.Cmd, 0)
	valuePrefix := strings.Repeat("a", valueSize)
	r := &result{0, 0, 0, make([]statistic, 0)}
	for i := 0; i < count; i++ {
		var keyLen int
		if keyspacelen > 0 {
			keyLen = rand.Intn(keyspacelen)
		} else {
			keyLen = id*count + i
		}
		key := fmt.Sprintf("%d", keyLen)
		value := fmt.Sprintf("%s%d", valuePrefix, keyLen)
		name := operation
		if operation == "mixed" {
			if rand.Intn(2) == 1 {
				name = "set"
			} else {
				name = "get"
			}
		}
		c := &clients.Cmd{
			Name:  name,
			Key:   key,
			Value: value,
			Error: nil,
		}
		if pipelen > 1 {
			cmds = append(cmds, c)
			if len(cmds) == pipelen {
				pipeline(client, cmds, r)
				cmds = make([]*clients.Cmd, 0)
			} else {
				run(client, c, r)
			}
		}
		if len(cmds) != 0 {
			pipeline(client, cmds, r)
		}
		ch <- r
	}
}

func run(client clients.Client, c *clients.Cmd, r *result) {
	expect := c.Value
	start := time.Now()
	client.Run(c)
	d := time.Now().Sub(start)
	resultType := c.Name
	if resultType == "get" {
		if c.Value == "" {
			resultType = "miss"
		} else if c.Value != expect {
			panic(c)
		}
	}
	r.addDuration(d, resultType)
}

func pipeline(client clients.Client, cmds []*clients.Cmd, r *result) {
	expect := make([]string, len(cmds))
	for i, c := range cmds {
		if c.Name == "get" {
			expect[i] = c.Value
		}
	}

	start := time.Now()
	client.PipelineRun(cmds)
	d := time.Now().Sub(start)
	for i, c := range cmds {
		resultType := c.Name
		if resultType == "get" {
			if c.Value == "" {
				resultType = "miss"
			} else if c.Value != expect[i] {
				fmt.Println(expect[i])
				panic(c.Value)
			}
		}
		r.addDuration(d, resultType)
	}
}
