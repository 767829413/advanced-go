package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func worker(id int, wg *sync.WaitGroup, done <-chan struct{}) {
	defer wg.Done()

	for i := 0; ; i++ {
		select {
		case <-done:
			fmt.Printf("Goroutine %d: Exiting\n", id)
			return
		default:
			fmt.Printf("Goroutine %d: Iteration %d\n", id, i)
			time.Sleep(time.Second)
		}
	}
}

func main() {
	const numGoroutines = 5
	var wg sync.WaitGroup
	done := make(chan struct{})

	fmt.Printf("Main goroutine ID: %d\n", getGoroutineID())

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go worker(i, &wg, done)
	}

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 主 goroutine 打印，同时监听信号
	go func() {
		for i := 0; ; i++ {
			select {
			case <-done:
				return
			default:
				fmt.Printf("Main goroutine: Iteration %d\n", i)
				time.Sleep(time.Second)
			}
		}
	}()

	// 等待中断信号
	<-sigChan
	fmt.Println("\nReceived interrupt, shutting down...")

	// 关闭 done channel，通知所有 goroutine 退出
	close(done)

	// 等待所有 worker goroutine 完成
	wg.Wait()
	fmt.Println("All goroutines have finished. Exiting.")
}

// getGoroutineID 返回当前 goroutine 的 ID
// 注意：这不是官方支持的方法，仅用于演示目的
func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.ParseInt(idField, 10, 64)
	return id
}
