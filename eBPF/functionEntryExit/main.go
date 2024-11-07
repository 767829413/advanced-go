//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -target bpf -D__TARGET_ARCH_x86" --go-package tool -output-dir tool Tool func_entry_exit.c -- -I/usr/include/bpf -I/usr/include/linux

package main

import (
	"C"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"

	"github.com/767829413/advanced-go/eBPF/functionEntryExit/tool"
)

func main() {
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// 加载编译的 eBPF 程序
	objs := tool.ToolObjects{}
	if err := tool.LoadToolObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// 将 fentry 程序附加到 do_unlinkat 函数的入口点
	linkFentry, err := link.AttachTracing(link.TracingOptions{
		Program:    objs.DoUnlinkat,
		AttachType: ebpf.AttachTraceFEntry,
	})
	if err != nil {
		log.Fatalf("Attaching fentry program failed: %v", err)
	}
	defer linkFentry.Close()

	// 将 fexit 程序附加到 do_unlinkat 函数的退出点
	linkFexit, err := link.AttachTracing(link.TracingOptions{
		Program:    objs.DoUnlinkatExit,
		AttachType: ebpf.AttachTraceFExit,
	})
	if err != nil {
		log.Fatalf("Attaching fexit program failed: %v", err)
	}
	defer linkFexit.Close()

	fmt.Println("eBPF programs successfully attached.")

	// 捕获退出信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	fmt.Println("Exiting program...")
}
