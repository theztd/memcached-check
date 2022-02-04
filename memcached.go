package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

type Connection struct {
	conn net.Conn
	buf  bufio.ReadWriter
}

// Response endings
var breaks = map[string]bool{
	"ERROR":  true,
	"END":    true,
	"STORED": true,
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

func Talk(address string, indata []string) (ret []string, err error) {
	/*
		Establish connection to memcached service
	*/
	mc, err := Connect(address)
	if err != nil {
		return nil, err
	}
	defer mc.conn.Close()

	for _, l := range indata {
		// log.Println(l)
		mc.buf.WriteString(l)
	}

	mc.buf.Flush()

	var read_data []string
	for {
		bl, _, _ := mc.buf.ReadLine()
		l := string(bl)
		// log.Println(l)

		read_data = append(read_data, l)

		/*
			If received line does exist in breaking strings
			return result and exit
		*/
		if breaks[l] {
			return read_data, nil
		}
	}
}

func FunctionCheck(address string, data string) (status bool, err error) {
	/*
		Test writing data to memcached
	*/
	data_size := len(data)
	write_test, err := Talk(address, []string{
		fmt.Sprintf("set CHECK 0 100 %d\r\n", data_size),
		fmt.Sprintf("%s\r\n", data),
	})
	if err == nil && write_test[0] == "STORED" {
		log.Println("Data write test: OK")
	} else {
		return false, err
	}

	/*
		Test reading written data from memcached
	*/
	read_test, err := Talk(address, []string{"get CHECK\r\n"})
	if err == nil && read_test[1] == data {
		log.Println("Data read test: OK")
	} else {
		return false, err
	}

	return true, nil
}

func PrometheusMetrics(address string, service string, metrics_file string) bool {
	// Open file for writing metrics
	f, err := os.Create(metrics_file)
	if err != nil {
		log.Panic(err)
		return false
	}
	defer f.Close()

	data, err := Talk(address, []string{"stats\r\n"})
	if err != nil {
		f.WriteString(fmt.Sprintf("memcached_up{service=\"%s\"} %d\n", service, 0))
		log.Fatal("Unable to generate stats", err)
		return false
	}

	re := regexp.MustCompile("^STAT (.*) ([0-9]+)$")
	for _, l := range data {

		if strings.Contains(l, "ERROR") {
			log.Panic("ERROR")
			return false
		}

		// get value of each STAT
		stat := re.FindStringSubmatch(l)
		if len(stat) > 2 {
			// stats[stat[1]] = stat[2]
			f.WriteString(fmt.Sprintf("memcached_stats{type=\"%s\", service=\"%s\"} %s\n", stat[1], service, stat[2]))
		}
	}

	// memcached_up
	f.WriteString(fmt.Sprintf("memcached_up{service=\"%s\"} %d\n", service, 1))

	return true
}
