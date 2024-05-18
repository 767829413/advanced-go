package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe("localhost:8080", nil)
	}()

	go func() {
		// 每秒打印内存分配情况
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
			fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
			fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
			fmt.Printf("\tNumGC = %v\n", m.NumGC)
			time.Sleep(1 * time.Second)
		}

	}()
	time.Sleep(5 * time.Second)

	fmt.Println("start test")

	// 创建一个 200 MiB 的切片
	var memoryLeaks [][]int32
	for i := 0; i < 10; i++ {
		// leak := make([]int32, 5*1024*1024) // 分配 5*1M*4bytes = 20 MiB
		leak := make([]int32, 60*1024) // 分配 5*1M*4bytes = 20 MiB
		memoryLeaks = append(memoryLeaks, leak)
		time.Sleep(1 * time.Second) // 延迟一秒观察内存分配情况
	}
	// 期望至少分配了 200 MiB 内存
	fmt.Println("end test")
	// 这里强制GC一下
	// runtime.GC()
	// 看到上面的文字后，打开go pprof 工具，查看工具的分析
	// go tool pprof -http :8972 http://127.0.0.1:8080/debug/pprof/heap
	time.Sleep(1 * time.Hour)
	fmt.Println("test", memoryLeaks[9][5*1024*1024-1]) // 避免垃圾回收和优化

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
