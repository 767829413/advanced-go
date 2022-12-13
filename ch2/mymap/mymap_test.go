package mymap

import (
	"sync"
	"testing"
)

func BenchmarkMyMap(b *testing.B) {
	m := NewMyMap()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Load("test")
			m.LoadOrStore("fff", 333)
			m.LoadAndDelete("kkk")
		}
	})
}

func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Store("kkk", 123)
			m.Load("test")
			m.LoadOrStore("fff", 333)
			m.LoadAndDelete("kkk")
		}
	})
}
