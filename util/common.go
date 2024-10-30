package util

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

func PrintResp(v any) {
	d, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(d))
}

func WriteFile(path string, data any) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	d, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}
	_, err = file.WriteString(string(d))
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
}

// printMemUsage 输出当前的内存使用情况
func PrintMemUsage(file *os.File) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// 构建内存使用情况的字符串
	memUsage := fmt.Sprintf("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v\n",
		bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
	// 将内存使用情况写入文件
	_, err := file.WriteString(memUsage)
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
	}
}

// bToMb 将字节转换为兆字节
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
