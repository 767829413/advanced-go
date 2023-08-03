package main

import (
	"fmt"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "localhost:9898")
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer l.Close()
	fmt.Println("start wait data")
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("l.Accept() err: ", err)
	}
	defer conn.Close()
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err: ", err)
		return
	}
	fmt.Println("Read data: ", string(buf[:n]))
	_, err = conn.Write([]byte("ok ok ok!!!!"))
	if err != nil {
		fmt.Println("conn.Write failed: ", err)
		return
	}
}
