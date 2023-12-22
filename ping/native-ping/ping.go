package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

const (
	protocolICMP = 1
)

func NativePing() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host\n", os.Args[0])
		os.Exit(1)
	}
	host := os.Args[1]
	// 使用icmp得到一个*packetconn,注意这里的network我们设置的`udp4`
	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	// 生成一个Echo消息
	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("Hello, are you there!"),
		},
	}
	wb, err := msg.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	// 发送，注意这里必须是一个UDP地址
	start := time.Now()
	if _, err := c.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP(host)}); err != nil {
		log.Fatal(err)
	}
	// 读取回包
	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Fatal(err)
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		log.Fatal(err)
	}
	duration := time.Since(start)
	// 得到的回包是一个ICMP消息，先解析出来
	msg, err = icmp.ParseMessage(protocolICMP, reply[:n])
	if err != nil {
		log.Fatal(err)
	}
	// 打印结果
	switch msg.Type {
	case ipv4.ICMPTypeEchoReply: // 如果是Echo Reply消息
		echoReply, ok := msg.Body.(*icmp.Echo) // 消息体是Echo类型
		if !ok {
			log.Fatal("invalid ICMP Echo Reply message")
			return
		}
		// 这里可以通过ID, Seq、远程地址来进行判断，下面这个只使用了两个判断条件，是有风险的
		// 如果此时有其他程序也发送了ICMP Echo,序列号一样，那么就可能是别的程序的回包，只不过这个几率比较小而已
		// 如果再加上ID的判断，就精确了
		if peer.(*net.UDPAddr).IP.String() == host && echoReply.Seq == 1 {
			fmt.Printf("Reply from %s: seq=%d time=%v\n", host, msg.Body.(*icmp.Echo).Seq, duration)
			return
		}
	default:
		fmt.Printf("Unexpected ICMP message type: %v\n", msg.Type)
	}
}
