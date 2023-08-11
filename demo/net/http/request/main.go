package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		log.Println("http.Get error: ", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Header: ", resp.Header)
	fmt.Println("Status: ", resp.Status)
	fmt.Println("StatusCode: ", resp.StatusCode)
	fmt.Println("Proto: ", resp.Proto)

	buf := make([]byte, 4096)
	res := ""
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 || err == io.EOF {
			log.Println("--------Read finish--------")
			break
		}
		if err != nil {
			log.Println("resp.Body.Read error: ", err)
			return
		}
		res += string(buf[:n])
	}
	fmt.Println("Body: ", res)
}
