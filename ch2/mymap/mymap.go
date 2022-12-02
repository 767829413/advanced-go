package mymap

import (
	"sync"
)

// 只支持 int 即可。
var SHARD_COUNT = 32

type MyMap []*ConcurrentMapShared

type ConcurrentMapShared struct {
	m map[string]interface{}
	sync.RWMutex
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func NewMyMap() MyMap {
	m := make(MyMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &ConcurrentMapShared{m: make(map[string]interface{})}
	}
	return m
}

func (m MyMap) GetShard(key string) *ConcurrentMapShared {
	return m[uint(fnv32(key))%uint(SHARD_COUNT)]
}

func (m *MyMap) Load(key string) (value interface{}, ok bool) {
	shard := m.GetShard(key)
	shard.RLock()
	value, ok = shard.m[key]
	shard.RUnlock()
	return value, ok
}

func (m *MyMap) Store(key string, value interface{}) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.m[key] = value
	shard.Unlock()
}

func (m *MyMap) Delete(key string) {
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.m, key)
	shard.Unlock()
}

func (m *MyMap) LoadOrStore(key string, value interface{}) (actual interface{}, loaded bool) {
	shard := m.GetShard(key)
	shard.Lock()
	actual, ok := shard.m[key]
	if !ok {
		shard.m[key] = value
	}
	shard.Unlock()
	return actual, ok
}

func (m *MyMap) LoadAndDelete(key string) (value interface{}, loaded bool) {
	shard := m.GetShard(key)
	shard.Lock()
	value, loaded = shard.m[key]
	delete(shard.m, key)
	shard.Unlock()
	return value, loaded
}
