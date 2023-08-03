package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

/*
1. 创建监听的socket
2. 阻塞等待连接建立
3. 接受读取文件名
4. 回发确定准备完毕回执(OK字符串)
5. 接收文件内容,写入文件中
*/

func main() {
	// 创建用于监听的socket
	listen, err := net.Listen("tcp", "localhost:9899")
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	defer listen.Close()

	for {
		// 阻塞监听
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept() error:", err)
			return
		}
		go recev(conn)
	}
}

func recev(conn net.Conn) {
	defer conn.Close()
	// 获取文件名,保存
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read error:", err)
		return
	}
	fileName := string(buf[:n])

	// 回写OK,表示可以执行接收
	conn.Write([]byte("OK"))

	// 获取文件内容
	recvFile(conn, fileName)
}

func recvFile(conn net.Conn, fileName string) {
	// 按照文件名创建文件
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("os.Create error:", err)
		return
	}
	defer f.Close()

	// 从网络中读取数据,写入本地文件
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 || err == io.EOF {
			fmt.Println("file receive successfully")
			return
		}

		if err != nil {
			fmt.Println("conn.Read error:", err)
			return
		}
		//写入本地文件
		f.Write(buf)
	}
}
