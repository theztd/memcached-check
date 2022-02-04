package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net" // https://golang.org/pkg/net/
	"os"
	"regexp"
	"strings"
	"time"
)

type Connection struct {
	conn net.Conn
	buf  bufio.ReadWriter
}

func Connect(addr string) (conn *Connection, err error) {
	nc, err := net.Dial("tcp", addr)

	// return no data and error
	if err != nil {
		return nil, err
	}

	// return data and nil error
	return &Connection{
		conn: nc,
		buf: bufio.ReadWriter{
			Reader: bufio.NewReader(nc),
			Writer: bufio.NewWriter(nc),
		},
	}, nil
}

func WriteCheck(data string) bool {
	log.Println("Writing data to memcache", data)
	time.Sleep(time.Duration(500))
	log.Println("Reading data from memcache")
	return true
}

// Connection address
var address string
var service string
var metrics_file string
var write_check bool

func main() {
	flag.StringVar(&address, "addr", "127.0.0.1:11211", "Memcached server address.")
	flag.StringVar(&service, "service", "example", "Service tag variable (usefull to identify which memcached we measure).")
	flag.StringVar(&metrics_file, "metrics", "./metrics_test.prom", "Define path where the metrics have to be saved")
	flag.BoolVar(&write_check, "write", false, "Do a write and read test")

	flag.Parse()

	log.Println("Check memcached", service, "on address", address, "write_check:", write_check)

	// Open file for writing metrics
	f, err := os.Create(metrics_file)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	// connect to a Memcached server
	mc, err := Connect(address)
	if err != nil {
		// memcached_up 0
		f.WriteString(fmt.Sprintf("memcached_up{service=\"%s\"} %d\n", service, 0))

		log.Fatal(err)
		os.Exit(1)
	}
	defer mc.conn.Close()

	//
	if write_check {
		WriteCheck(time.Now().Format(time.RFC3339))
	}

	// send "stats" command
	mc.buf.WriteString("stats\r\n")
	mc.buf.Flush()

	re := regexp.MustCompile("^STAT (.*) ([0-9]+)$")
	stats := map[string]string{}

	for {
		bl, _, _ := mc.buf.ReadLine()
		l := string(bl)

		if strings.HasPrefix(l, "END") {
			break
		}

		if strings.Contains(l, "ERROR") {
			panic("ERROR")
		}

		// get value of each STAT
		stat := re.FindStringSubmatch(l)
		if len(stat) > 2 {
			stats[stat[1]] = stat[2]
		}
	}

	// memcached_up
	f.WriteString(fmt.Sprintf("memcached_up{service=\"%s\"} %d\n", service, 1))

	for k, v := range stats {
		//fmt.Printf("memcached_stats{type=\"%s\", service=\"%s\"} %s\n", k, service, v)
		f.WriteString(fmt.Sprintf("memcached_stats{type=\"%s\", service=\"%s\"} %s\n", k, service, v))
	}
	// fmt.Printf("%#v\n", stats)
	log.Println("Chech has been done")
}
