package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var rw sync.RWMutex

func main() {
	rand.Seed(time.Now().UnixNano())

	ch := make(chan int)

	for i := 0; i < 5; i++ {
		go readGo(ch, i)
	}

	for i := 0; i < 5; i++ {
		go writeGo(ch, i)
	}

	<-time.After(20 * time.Second)
	close(ch)
}

func readGo(in <-chan int, idx int) {
	for {
		// rw.RLock()
		num := <-in
		fmt.Printf("-------第%d read goroutine 读到数字 %d\n", idx, num)
		// rw.RUnlock()
	}
}

func writeGo(out chan<- int, idx int) {
	for {
		num := rand.Intn(1000)
		rw.Lock()
		out <- num
		fmt.Printf("-------第%d write goroutine 读到数字 %d\n", idx, num)
		time.Sleep(time.Millisecond * 300)
		rw.Unlock()
	}
}
