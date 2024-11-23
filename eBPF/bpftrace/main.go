//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -Wall -Werror" -target amd64 --go-package tool -output-dir tool Tool mysql_trace/mysql_trace.c -- -I/usr/include/bpf -I/usr/include/linux

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/767829413/advanced-go/eBPF/bpftrace/tool"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

// Event 对应 eBPF 程序中定义的 struct event
type Event struct {
	Pid     uint32    // 进程 ID
	Tid     uint32    // 线程 ID
	DeltaNs uint64    // 执行时间（纳秒）
	Query   [256]byte // SQL 查询字符串，最大长度为 256 字节
}

// GetQuery 返回查询字符串，去除末尾的空字节
func (e *Event) GetQuery() string {
	return string(bytes.TrimRight(e.Query[:], "\x00"))
}

func main() {
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

	// Open a uretprobe for the runtime.casgstatus function.
	ex, err := link.OpenExecutable(
		"/var/lib/docker/overlay2/270506fa70024c98a81543a6a274e9f03258f39471655b69fa22fce5fada6d13/merged/sbin/mysqld",
	)
	if err != nil {
		log.Fatalf("opening executable: %v", err)
	}

	up, err := ex.Uprobe(
		"_Z16dispatch_commandP3THDPK8COM_DATA19enum_server_command",
		objs.UprobeMysqlDispatchCommand,
		nil,
	)
	if err != nil {
		log.Fatalf("creating uretprobe: %v", err)
	}
	defer up.Close()

	upret, err := ex.Uretprobe(
		"_Z16dispatch_commandP3THDPK8COM_DATA19enum_server_command",
		objs.UretprobeMysqlDispatchCommand,
		nil,
	)
	if err != nil {
		log.Fatalf("creating uretprobe: %v", err)
	}
	defer upret.Close()

	// Open a ring buffer to receive events from the eBPF program.
	rb, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		log.Fatalf("opening ringbuf reader: %v", err)
	}
	defer rb.Close()

	// 捕获 SIGINT 和 SIGTERM 信号以正确退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("listening for events...")

	// 读取 ring buffer 中的数据
	go func() {
		for {
			// 读取事件
			// 读取到的 `Event` 数据，打印出来，这样就可以看到慢查询 `SQL` 了
			record, err := rb.Read()
			if err != nil {
				log.Printf("reading from ringbuf: %v", err)
				return
			}

			// 将事件解析为 Event 结构
			var data Event
			err = binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &data)
			if err != nil {
				log.Printf("parsing event: %v", err)
				return
			}

			// 输出 goroutine 信息
			fmt.Printf("PID: %d, TID: %d, Latency: %v, SQL: %s\n",
				data.Pid, data.Tid, time.Duration(data.DeltaNs), data.GetQuery())
		}
	}()

	// 等待退出信号
	<-sig
	fmt.Println("Exiting...")
}
