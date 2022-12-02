package main

import "time"

func main() {
	ch := make(chan int)

	go func() {
		time.Sleep(3 * time.Second)
		ch <- 123354
	}()

	println(<-ch)
}
