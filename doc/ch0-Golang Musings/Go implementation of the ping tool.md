# Golang 实现ping工具方法

## 介绍ping

基本介绍请看维基: [ping](https://zh.wikipedia.org/wiki/Ping)

一般ping的工具是基于[rfc792](https://datatracker.ietf.org/doc/html/rfc792)来实现的,主要是通过ICMP协议，该协议是TCP/IP网络协议套件中的一个重要组成部分.

至于该协议是如何作用的以及相关概念原理请看维基: [ICMP](https://zh.wikipedia.org/wiki/%E4%BA%92%E8%81%94%E7%BD%91%E6%8E%A7%E5%88%B6%E6%B6%88%E6%81%AF%E5%8D%8F%E8%AE%AE)

## 使用Go实现ping的几种方式

### 使用系统调用

```go
package main

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/exec"
)

func main() {
	host := os.Args[1]
	output, err := exec.Command("ping", host).CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	// 处理命令行中文转码的问题
	newByte, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(output)
	fmt.Println(string(newByte))
}
```

### 使用 ip4:icmp 实现

Go的使用golang.org/x/net/icmp包实现ping.

network需要是ip4:icmp,能够发送ICMP包

发送额度内容是ICMP Echo消息，地址不是UDP地址，是IP 地址

```go
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
		log.Fatal(err)
	}

	// 创建 ICMP 连接
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// 发送 ICMP 报文
	start := time.Now()
	_, err = conn.WriteTo(msgBytes, dst)
	if err != nil {
		log.Fatal(err)
	}

	// 接收 ICMP 报文
	reply := make([]byte, 1500)
	for i := 0; i < 3; i++ {
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Fatal(err)
		}
		n, peer, err := conn.ReadFrom(reply)
		if err != nil {
			log.Fatal(err)
		}
		duration := time.Since(start)

		// 解析 ICMP 报文
		msg, err = icmp.ParseMessage(protocolICMP, reply[:n])
		if err != nil {
			log.Fatal(err)
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
```

### 使用pro-bing

Go net扩展库提供了icmp包，方便实现ping能力，这里使用一个第三方包: [github.com/prometheus-community/pro-bing](https://github.com/prometheus-community/pro-bing)

下面代码就是一个ping的基本功能，没什么好说的，ping3次得到结果:

```go
package proBing

import (
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"log"
	"os"
)

func ProBing() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s host\n", os.Args[0])
		os.Exit(1)
	}
	host := os.Args[1]

	pinger, err := probing.NewPinger(host)
	if err != nil {
		log.Fatal("probing.NewPinger error: ", err)
	}
	// Windows 
	pinger.SetPrivileged(true)
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Fatal("pinger.Run error: ", err)
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats

	fmt.Println(stats)
}
```

其他示例可以看看这个库的其他例子

## 总结

介绍了使用Go语言实现ping工具的几种方式。

1. 通过系统调用的方式，使用os/exec包执行系统命令来实现ping功能
2. 使用 golang.org/x/net/icmp 包以及 golang.org/x/net/ipv4 包实现了基于ICMP协议的ping功能
3. 使用第三方库 github.com/prometheus-community/pro-bing 来实现ping功能
4. 每种方式都有其特点和适用场景，可以根据自己的需求选择合适的方式来实现ping功能