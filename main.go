package main

import "fmt"

type A struct {
	Name string
	List []int
}

type B struct {
	A
}

func main() {
	b := &B{}
	fmt.Println(b.List)
}
