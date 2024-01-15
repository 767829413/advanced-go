package main

import (
	"fmt"
	"github.com/kortschak/goroutine"
	"sort"
	"time"
)

func main() {
	links := goroutine.All()
	sort.Slice(links, func(i, j int) bool {
		return links[i].Child < links[j].Child
	})

	for _, link := range links {
		fmt.Printf("%d -> %d\n", link.Parent, link.Child)
	}
}

func foo(depth int) {
	fmt.Printf("goroutine #%d: depth=%d, parent: #%d\n", goroutine.ID(), depth, goroutine.ParentID())

	if depth == 0 {
		return
	}

	go foo(depth - 1)
}

func mainxxx() {
	depth := 5
	go foo(depth)

	time.Sleep(1 * time.Second)
}
