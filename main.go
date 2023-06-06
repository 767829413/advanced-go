package main

import "time"

func main() {
	var ch = make(chan int)

	go func() {
		i := 0
		for {
			ch <- i
			i++
		}
	}()

	go func() {
		for {
			println(<-ch)
		}
	}()

	af := time.After(1000 * time.Microsecond)

	<-af
}
