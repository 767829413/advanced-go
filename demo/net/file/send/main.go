package main

/*
1. 获取文件名
2. 建立连接
3. 发送文件名
4. 接收服务器的回执
5. 判断是否准备好接收文件
6. 发送文件内容
*/

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var path string

func init() {
	flag.StringVar(&path, "path", "", "file path")
	flag.Parse()
}

func main() {
	// 判断文件名是否设置
	if path == "" {
		fmt.Println("-path or --path value is required")
		return
	}
	// 获取文件属性
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("os.Stat error:", err)
		return
	}
	// 提取文件名
	fileName := fileInfo.Name()

	// 主动发起连接
	conn, err := net.Dial("tcp", "127.0.0.1:9899")
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return
	}
	defer conn.Close()

	// 发送文件名给接收端
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return
	}

	// 接收服务器的回执 OK
	buff := make([]byte, 16)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Println("conn.Read error:", err)
		return
	}
	if string(buff[:n]) == "OK" {
		// 发送文件内容
		send(conn, path)
	}
}

func send(conn net.Conn, path string) {
	// 只读打开文件
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("os.Open error:", err)
		return
	}
	defer f.Close()

	// 从本地文件读内容,读多少写多少
	buf := make([]byte, 4096)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("file read ended successfully")
			} else {
				fmt.Println("f.Read error:", err)
			}
			return
		}
		// 写入网络中
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("conn.Write error:", err)
			return
		}
	}
}
