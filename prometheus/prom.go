package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hello_world_req_count",
		Help: "hello world requested. ",
	})
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	requests.Inc()
	w.Write([]byte("hello world"))
}

func main() {
	http.HandleFunc("/hello", sayHello)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8888", nil))
}
