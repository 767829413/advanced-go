package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "wsl:9898")
	exit := make(chan struct{})
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return
	}
	defer conn.Close()
	// 获取用户键盘输入 stdin,将输入发给服务器
	go func() {
		wBuf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(wBuf)
			if err != nil {
				fmt.Println("os.Stdin.Read err: ", err)
				continue
			}
			conn.Write(wBuf[:n])
			if strings.TrimSpace(string(wBuf[:n])) == "exit" {
				exit <- struct{}{}
			}
		}
	}()
	// 回显服务器回发的大写数据
	rBuf := make([]byte, 4096)
	for {
		select {
		case <-exit:
			return
		default:
			n, err := conn.Read(rBuf)
			if n == 0 {
				fmt.Println("read end")
				continue
			}
			if err != nil {
				fmt.Println("conn.Read err: ", err)
				continue
			}
			fmt.Println(string(rBuf[:n]))
		}
	}
}
