package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	f2()
}

func f1() {
	runtime.GOMAXPROCS(1)
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			fmt.Println("A: ", i)
		}()

	}

	var ch = make(chan int)
	<-ch
}

func f2() {
	runtime.GOMAXPROCS(1)
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			fmt.Println("A: ", i)
		}()

	}

	time.Sleep(time.Hour)
}
