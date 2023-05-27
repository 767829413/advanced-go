package main

func main() {
	var a = make(chan int)
	select {
	case <-a:
	default:
	}
}
