package main

import (
	"log"
	"net"

	"github.com/767829413/advanced-go/util"
)

// 用户结构体类型
type user struct {
	C    chan string
	Name string
	Addr string
}

// 创建全局map来存储在线用户
var clients = make(map[uint64]*user)

// 创建全局 channel 来传递用户消息
var msg = make(chan string)

func main() {
	// 创建socket
	l, err := net.Listen("tcp", "localhost:9990")
	if err != nil {
		log.Println("net.Listen error: ", err)
		return
	}
	defer l.Close()
	log.Println("server listening on ", l.Addr().String())

	// 启动一个 goroutine 来管理全局map和全局channel
	go manager()
	log.Println("manager started successfully")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("l.Accept error: ", err)
			return
		}

		go handlerConn(conn)
	}
}

func manager() {
	// 监听全局channel获取信息
	for {
		message := <-msg
		// 获取消息后循环发送消息
		for _, user := range clients {
			user.C <- message
		}
	}
}

func handlerConn(conn net.Conn) {
	defer conn.Close()
	// 获取用户远程地址(ip + port),作为默认名称
	addr := conn.RemoteAddr().String()
	// 创建连接用户
	user := user{
		C:    make(chan string),
		Name: addr,
		Addr: addr,
	}
	// 加入全局map中
	clients[util.Str2HashInt(addr)] = &user

	// 发送上线信息到全局通道
	msg <- user.Name + " login success"

	// 启动一个 goroutine 来读取消息
	go readUserMsg(conn, &user)

	// 获取用户输入的信息
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("conn.Read error: ", err)
			return
		}
		user.C <- user.Name + ": " + string(buf[:n])
	}
}

func readUserMsg(conn net.Conn, user *user) {
	// 监听用户自带 channel 消息
	for message := range user.C {
		// 添加 \n 来标记消息结束
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Println("conn.Write error: ", err)
			return
		}
	}
}
