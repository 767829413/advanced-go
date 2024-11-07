//go:generate go run github.com/cilium/ebpf/cmd/bpf2go --go-package minimal -output-dir minimal Minimal hello.c -- -I/usr/include/bpf -I/usr/include/linux
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"C"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"

	minimal "github.com/767829413/advanced-go/eBPF/helloworld/minimal"
)

func main() {
	// 允许当前进程锁定内存以用于 eBPF 资源.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("failed to remove memlock: %v", err)
	}

	// 加载 eBPF 程序.
	objs := minimal.MinimalObjects{}
	if err := minimal.LoadMinimalObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// 将 eBPF 程序附加到 sys_enter_write tracepoint.
	tp, err := link.Tracepoint("syscalls", "sys_enter_write", objs.HandleTp, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %v", err)
	}
	defer tp.Close()

	log.Println("eBPF program loaded and attached. Press Ctrl+C to exit.")

	// 创建一个通道，用于接收信号
	sigs := make(chan os.Signal, 1)
	// 注册信号处理器，监听 SIGINT 和 SIGTERM 信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
