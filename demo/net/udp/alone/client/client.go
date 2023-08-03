package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("udp", "wsl:9899")
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return
	}
	defer conn.Close()
	fmt.Println("start request server")
	_, err = conn.Write([]byte("hello world"))
	if err != nil {
		fmt.Println("conn.Write failed: ", err)
		return
	}
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read failed: ", err)
	}
	fmt.Println("conn.Read: ", string(buf[:n]))
}