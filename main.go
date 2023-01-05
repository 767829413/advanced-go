package main

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
)

var m = map[[12]byte]int{}
// var m = map[string]int{}

func init() {
	for i := 0; i < 1000000; i++ {
		var key [12]byte
		copy(key[:],fmt.Sprint(i))
		m[key] = i
		// m[fmt.Sprint(i)] = i
	}
}

func sayHeelo(wr http.ResponseWriter, r *http.Request) {
	io.WriteString(wr, "hello")
}

func main() {
	http.HandleFunc("/", sayHeelo)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		fmt.Println(err)
	}
}
