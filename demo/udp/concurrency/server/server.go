package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// 组织一个 udp 地址结构,指定服务器的 ip:port
	udpServerAddr, err := net.ResolveUDPAddr("udp", "wsl:9899")
	if err != nil {
		fmt.Println("net.ResolveUDPAddr err: ", err)
	}
	fmt.Println("server udp address build successfully")
	// 创建用户通信 socket
	conn, err := net.ListenUDP("udp", udpServerAddr)
	if err != nil {
		fmt.Println("net.ListenUDP err: ", err)
	}
	defer conn.Close()
	fmt.Println("server ListenUDP successfully")
	buf := make([]byte, 4096)
	for {
		// 读取客户端发送的数据
		// 读取的字节数,客户端的udp地址,错误信息
		n, udpClientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("conn.ReadFromUDP err: ", err)
		}
		go func() {
			fmt.Println("server ReadFromUDP successfully")
			fmt.Println("client udp addr: ", udpClientAddr)

			// 数据处理
			fmt.Println("get data: ", string(buf[:n]))
			// 回复客户端
			_, err = conn.WriteToUDP([]byte(time.Now().String()+"\n"), udpClientAddr)
			if err != nil {
				fmt.Println("conn.WriteToUDP err: ", err)
			}
		}()

	}
}
