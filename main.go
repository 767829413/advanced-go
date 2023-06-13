package main

func main() {
	allocOnHeap()
}

func allocOnHeap() {
	var m = make([]int, 10240)
	println(m)
}
