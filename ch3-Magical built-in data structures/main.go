package main

import "time"

func main() {
	var ch = make(chan int, 3)

	go func() {
		var i = 100
		for {
			ch <- i
			i++
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			v := <-ch
			println(v)
		}
	}()
	println(1 << 16)
	select {}
}
