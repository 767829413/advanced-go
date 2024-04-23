package main

import "fmt"

const (
	PIC_SET_DESK = 1 << 0 //1
	PIC_SET_HEAD = 1 << 1 //2
	PIC_SET_DOC  = 1 << 2 //4
)
func main() {
	fmt.Println(PIC_SET_DESK + PIC_SET_HEAD)
}