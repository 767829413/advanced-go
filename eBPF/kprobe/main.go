//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -D__TARGET_ARCH_x86" --go-package minimal -output-dir minimal Minimal kprobe_unlink.c -- -I/usr/include/bpf -I/usr/include/linux

package main

import (
	"os"
	"os/signal"
	"syscall"

	"C"
	"log"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"

	minimal "github.com/767829413/advanced-go/eBPF/kprobe/minimal"
)

func main() {
	// 允许锁定内存
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("Removing memlock limit: %v", err)
	}

	// 加载编译的 eBPF 程序
	objs := minimal.MinimalObjects{}
	if err := minimal.LoadMinimalObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// 附加 kprobe 和 kretprobe
	kp, err := link.Kprobe("do_unlinkat", objs.DoUnlinkat, nil)
	if err != nil {
		log.Fatalf("attaching kprobe: %v", err)
	}
	defer kp.Close()

	krp, err := link.Kretprobe("do_unlinkat", objs.DoUnlinkatExit, nil)
	if err != nil {
		log.Fatalf("attaching kretprobe: %v", err)
	}
	defer krp.Close()

	log.Println("eBPF program successfully attached, press Ctrl+C to exit...")

	// 创建一个通道，用于接收信号
	sigs := make(chan os.Signal, 1)
	// 注册信号处理器，监听 SIGINT 和 SIGTERM 信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
