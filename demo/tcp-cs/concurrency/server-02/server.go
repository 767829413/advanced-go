package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "localhost:9898")
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer l.Close()
	fmt.Println("start wait data")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("l.Accept() err: ", err)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	fmt.Println("Remote address: ", addr)
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			fmt.Println("end of connection")
			return
		}
		if err != nil {
			fmt.Println("conn.Read err: ", err)
			return
		}
		if strings.TrimSpace(string(buf[:n])) == "exit" {
			fmt.Println("client request close")
			return
		}
		fmt.Println("read data: ", string(buf[:n]))
		conn.Write([]byte(strings.ToUpper(string(buf[:n]))))
		if err != nil {
			fmt.Println("conn.Write err: ", err)
			return
		}
	}
}
