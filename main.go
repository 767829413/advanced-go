package main

import (
	"fmt"
	"strconv"
)

func main() {
	var i int64 = 0
	fmt.Println(23232, strconv.FormatInt(i, 10))
	// m := mystruct{0}
	// test(m)  //错误
	// test(*m) //错误
}

type myinterface interface {
	print()
}

func test(value *myinterface) { // 这里不应该使用 myinterface 的指针
	//someting to do ...
}

type mystruct struct {
	i int
}

// 实现接口
func (this *mystruct) print() {
	fmt.Println(this.i)
	this.i = 1
}
