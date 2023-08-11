package main

import (
	"fmt"
	"log"
	"net"
	"testing"
)


func TestMain(t *testing.T) {
	conn, err := net.Dial("tcp", "wsl:9999")
	if err != nil {
		log.Println("net.Dial error: ", err)
		return
	}
	defer conn.Close()

	httpRequset := "GET /hello1 HTTP/1.1\r\nHost:wsl:9999\r\n\r\n"
	conn.Write([]byte(httpRequset))

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	if err != nil {
		log.Println("conn.Read error: ", err)
		return
	}
	if n == 0 {
		return
	}
	fmt.Println(string(buf[:n]))
}
