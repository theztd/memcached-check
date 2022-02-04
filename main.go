package main

import (
	"flag"
	"log" // https://golang.org/pkg/net/
	"time"
)

// Connection address
var address string

// Service TAG
var service string

// Output file for metrics
var metrics_file string

// Have to be true to proceed wite check
var write_check bool

func main() {
	flag.StringVar(&address, "addr", "127.0.0.1:11211", "Memcached server address.")
	flag.StringVar(&service, "service", "example", "Service tag variable (usefull to identify which memcached we measure).")
	flag.StringVar(&metrics_file, "metrics", "./TMP_metrics.prom", "Define path where the metrics have to be saved")
	flag.BoolVar(&write_check, "write", false, "Do a write and read test")

	flag.Parse()

	log.Println("Check memcached", service, "on address", address, "write_check", write_check, "...")

	if metrics_file != "" {
		PrometheusMetrics(address, service, metrics_file)
	}

	if write_check {
		FunctionCheck(address, string(time.Now().Format(time.RFC3339)))
	}

	// fmt.Printf("%#v\n", stats)
	log.Println("Chech has been done")
}
