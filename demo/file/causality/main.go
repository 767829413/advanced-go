package main

import (
	"flag"
	"fmt"
)

func main() {
	var file string
	flag.StringVar(&file, "file", "", "file path")
	flag.Parse()


	fmt.Println("ddddd:   ",file)
}
