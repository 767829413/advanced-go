//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -Wall -Werror" -target bpf --go-package tool -output-dir tool Tool uprobes_goroutine.c -- -I/usr/include/bpf -I/usr/include/linux

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"

	"github.com/767829413/advanced-go/eBPF/uprobes/goroutineState"
	"github.com/767829413/advanced-go/eBPF/uprobes/tool"
)

var pidTarget int

func main() {
	// 解析命令行参数
	flag.IntVar(&pidTarget, "pid", 0, "Target PID for filtering")
	flag.Parse()

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// 加载编译的 eBPF 程序
	spec, err := tool.LoadTool()
	if err != nil {
		log.Fatalf("loading BPF program: %v", err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := tool.ToolObjects{}
	if err := spec.LoadAndAssign(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// Open a uprobe for the runtime.casgstatus function.
	ex, err := link.OpenExecutable(
		"/home/fangyuan/code/go/src/github.com/767829413/advanced-go/eBPF/uprobes/demo/main",
	)
	if err != nil {
		log.Fatalf("opening executable: %v", err)
	}

	err = spec.RewriteConstants(map[string]interface{}{
		"pid_target": uint32(pidTarget),
	})
	if err != nil {
		log.Fatalf("rewriting constants: %v", err)
	}

	// up, err := ex.Uretprobe("runtime.casgstatus", objs.UprobeRuntimeCasgstatus, nil)
	up, err := ex.Uprobe("runtime.casgstatus", objs.UprobeRuntimeCasgstatus, nil)
	if err != nil {
		log.Fatalf("creating uretprobe: %v", err)
	}
	defer up.Close()

	// Open a ring buffer to receive events from the eBPF program.
	rb, err := ringbuf.NewReader(objs.Rb)
	if err != nil {
		log.Fatalf("opening ringbuf reader: %v", err)
	}
	defer rb.Close()

	// 捕获 SIGINT 和 SIGTERM 信号以正确退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for events...")

	// 读取 ring buffer 中的数据
	go func() {
		for {
			// 读取事件
			record, err := rb.Read()
			if err != nil {
				log.Printf("reading from ringbuf: %v", err)
				return
			}

			// 将事件解析为 GoroutineExecuteData 结构

			var data goroutineState.GoroutineExecuteData
			err = binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &data)
			if err != nil {
				log.Printf("parsing event: %v", err)
				return
			}

			// 输出 goroutine 信息
			fmt.Printf("TGID: %d, PID: %d, GoID: %d, OldState: %s, NewState： %s\n",
				data.Tgid, data.Pid, data.Goid, data.OldState.String(), data.NewState.String())
		}
	}()

	// 等待退出信号
	<-sig
	fmt.Println("Exiting...")
}
