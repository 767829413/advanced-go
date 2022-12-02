package mymap

import (
	"sync"
	"testing"
)

func BenchmarkMyMap(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := NewMyMap()
			m.Store("kkk", 123)
			m.Load("test")
			m.LoadOrStore("fff", 333)
			m.LoadAndDelete("kkk")
		}
	})
}

func BenchmarkSyncMap(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var m sync.Map
			m.Store("kkk", 123)
			m.Load("test")
			m.LoadOrStore("fff", 333)
			m.LoadAndDelete("kkk")
		}
	})
}
