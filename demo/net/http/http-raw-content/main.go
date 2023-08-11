package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "wsl:9999")
	if err != nil {
		log.Println("net.Listen error: ", err)
		return
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Println("l.Accept error: ", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("conn.Read error: ", err)
		return
	}
	fmt.Println(string(buf[:n]))
	conn.Write([]byte(``))
}
