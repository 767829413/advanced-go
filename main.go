package main

import (
	"fmt"
	"time"
)

type P struct {
	Name string
}

func main() {
	var p P
	time.Sleep(3 * time.Second)
	p = P{Name: "123"}
	go func() {
		fmt.Println(p.Name)
	}()

}
