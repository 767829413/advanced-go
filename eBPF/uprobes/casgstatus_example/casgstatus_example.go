package main

import (
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		runtime.Gosched() // 让出 CPU，触发 goroutine 状态变化
	}()

	wg.Wait()
}
