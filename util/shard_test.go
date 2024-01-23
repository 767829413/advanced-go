package util

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkShardCounter(b *testing.B) {
	counter := NewShard[atomic.Int64]()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Get().Add(1)
		}
	})
}

func BenchmarkMutexCounter(b *testing.B) {
	var counter int64
	var mu sync.Mutex
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter += 1
			mu.Unlock()
		}
	})
}

func BenchmarkAtomicCounter(b *testing.B) {
	var counter atomic.Int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Add(1)
		}
	})
}
