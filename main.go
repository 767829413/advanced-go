package main

var c = make(chan int)
var a string

func f() {
	a = "hi"
	<-c
}

func main() {
	go f()
	c <- 0
	println(a)
}
