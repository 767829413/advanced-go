package icmpPing

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	protocolICMP = 1
)

func IcmpPing() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s host\n", os.Args[0])
		os.Exit(1)
	}
	host := os.Args[1]

	// 解析目标主机的 IP 地址
	dst, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		log.Fatal("net.ResolveIPAddr error: ", err)
	}

	// 创建 ICMP 连接
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal("icmp.ListenPacket error: ", err)
	}
	defer conn.Close()

	// 构造 ICMP 报文
	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("Hello, are you there!"),
		},
	}
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		log.Fatal("msg.Marshal error: ", err)
	}

	// 发送 ICMP 报文
	start := time.Now()
	_, err = conn.WriteTo(msgBytes, dst)
	if err != nil {
		log.Fatal("conn.WriteTo error: ", err)
	}

	// 接收 ICMP 报文
	reply := make([]byte, 1500)
	for i := 0; i < 3; i++ {
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Fatal("conn.SetReadDeadline error: ", err)
		}
		n, peer, err := conn.ReadFrom(reply)
		if err != nil {
			log.Fatal("conn.ReadFrom error:", err)
		}
		duration := time.Since(start)

		// 解析 ICMP 报文
		msg, err = icmp.ParseMessage(protocolICMP, reply[:n])
		if err != nil {
			log.Fatal("icmp.ParseMessage error: ", err)
		}

		// 打印结果
		switch msg.Type {
		case ipv4.ICMPTypeEchoReply:
			echoReply, ok := msg.Body.(*icmp.Echo)
			if !ok {
				log.Fatal("invalid ICMP Echo Reply message")
				return
			}
			if peer.String() == host && echoReply.ID == os.Getpid()&0xffff && echoReply.Seq == 1 {
				fmt.Printf("reply from %s: seq=%d time=%v\n", dst.String(), msg.Body.(*icmp.Echo).Seq, duration)
				return
			}
		default:
			fmt.Printf("unexpected ICMP message type: %v\n", msg.Type)
		}
	}
}
