package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unicode/utf8"
)

var rw sync.RWMutex
var w sync.WaitGroup

func main(){
	fmt.Println(Len("张"))
	fmt.Println(Len("z"))
}

func Len(s string) int {
	return utf8.RuneCountInString(s)
}

func main111() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(22222)
			panic(r)
		}
	}()

	w.Add(2)

	go func() {
		time.Sleep(10 * time.Second)
		w.Done()
	}()

	go func() {
		time.Sleep(2 * time.Second)
		panic(232323)
		w.Done()
	}()
	w.Wait()
}

func test() {

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
