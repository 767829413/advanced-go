# Go 实现ping工具

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

### 使用 golang.org/x/net/icmp

Go的net扩展库专门实现了icmp协议,我们可以使用它来实现ping.

PS: 

1. `如果使用SOCK_RAW实现ping，是需要cap_net_raw权限的,你可以通过下面的命令设置:`

    ```bash
    setcap cap_net_raw=+ep /path/to/your/compiled/binary
    ```

2. `在Linux 3.0新实现了一种Socket方式，可以实现普通用户也能执行ping命令:`

    ```bash
    socket(PF_INET, SOCK_DGRAM, IPPROTO_ICMP)
    sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
    ```

实现非特权(non-privileged) 方式的ping, icmp包为我们做了封装，所以我们不必使用底层的socket,而是直接使用icmp.ListenPacket("udp4", "0.0.0.0")来实现.

```go

```

## References

* <https://wudaijun.com/2018/02/go-sync-map-implement/>
* <https://github.com/kat-co/concurrency-in-go-src>
* <https://speakerdeck.com/kavya719/understanding-channels>
* <https://www.zenlife.tk/concurrency-with-keep-order.md?hmsr=joyk.com&utm_source=joyk.com&utm_medium=referral>
* <https://golang.org/ref/mem>
* <https://www.hardwaretimes.com/difference-between-l1-l2-and-l3-cache-what-is-cpu-cache/>
* <https://github.com/lotusirous/go-concurrency-patterns>
* <https://songlh.github.io/paper/go-study.pdf>
* <https://github.com/cch123/golang-notes/blob/master/memory_barrier.md>

## 未涉及

* 内置并发结构：sync.Cond
* 进阶话题：如 acquire、release、sequential consistency、Lock-Free，Wait-free 等等
* 扩展并发原语：SingleFlight，ErrGroup 等