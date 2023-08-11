package main

import (
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	http.HandleFunc("/hello", HelloHandler)

	http.ListenAndServe("wsl:9999", nil)
}
