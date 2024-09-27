package main

import ()

func main() {
	var a []string
	a = nil
	for _, v := range a {
		_ = v
	}
}
