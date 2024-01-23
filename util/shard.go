package util

import (
	"runtime"

	"github.com/767829413/advanced-go/goroutine"
	"golang.org/x/sys/cpu"
)

// Shard类型
type Shard[T any] struct {
	values []value[T]
}

// NewShard创建一个新的Shard.
func NewShard[T any]() *Shard[T] {
	n := runtime.GOMAXPROCS(0)

	return &Shard[T]{
		values: make([]value[T], n),
	}
}

// 避免伪共享
type value[T any] struct {
	_ cpu.CacheLinePad
	v T
	_ cpu.CacheLinePad
}

// 得到当前P的值
func (s *Shard[T]) Get() *T {
	if len(s.values) == 0 {
		panic("sync: Sharded is empty and has not been initialized")
	}

	return &s.values[int(goroutine.PID())%len(s.values)].v
}

// 遍历所有的值
func (s *Shard[T]) Range(f func(*T)) {
	if len(s.values) == 0 {
		panic("sync: Sharded is empty and has not been initialized")
	}

	for i := range s.values {
		f(&s.values[i].v)
	}
}
