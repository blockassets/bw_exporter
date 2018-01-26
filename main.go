package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
	"os"
)

var (
	// Makefile build
	version = ""
)

func main() {
	port := flag.String("port", ":4030", "The address to listen on for /metrics HTTP requests.")
	cgHost := flag.String("cghost", "127.0.0.1", "The address of the worker.")
	cgPort := flag.Int64("cgport", 4028, "The port cgminer runs on.")
	cgTimeout := flag.Duration("cgtimeout", 5 * time.Second, "The amount of time to wait for cgminer to return.")
	flag.Parse()

	prometheus.MustRegister(NewExporter(*cgHost, *cgPort, *cgTimeout))

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("%s %s", os.Args[0], version)
	log.Fatal(http.ListenAndServe(*port, nil))
}
