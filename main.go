package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "ewqewqewq:2323232323"
	fmt.Println(strings.HasPrefix(s, "1ewqewqewq:"))
}
