package main

import "sync"

var slicePool = sync.Pool{}
var mapPool = sync.Pool{}

// slice can be easily reused
func processUserRequest1() {
	sl := slicePool.Get().([]any)
	defer func() {
		sl := sl[:0]
		slicePool.Put(sl)
	}()
	// processs user logic
}

// what about map?
func processUserRequest2() {
	m := mapPool.Get()
	defer func() {
		// how to reset a map
		mapPool.Put(m)
	}()
}
