package cache

import (
	"context"
	"log"
	"time"
)

var ctx = context.Background()

// cacheManager 管理不同的缓存实现，并提供自动切换功能
type cacheManager struct {
	primaryCache   Cache
	secondaryCache Cache
}

// newCacheManager 创建一个新的 cacheManager 实例
func newCacheManager(primaryCache Cache, secondaryCache Cache) *cacheManager {
	return &cacheManager{
		primaryCache:   primaryCache,
		secondaryCache: secondaryCache,
	}
}

// Set 尝试在主缓存中设置值，如果失败则使用备份缓存
func (m *cacheManager) Set(
	key string,
	value any,
	expiration time.Duration,
) error {
	if m.primaryCache == nil {
		return m.secondaryCache.Set(ctx, key, value, expiration)
	}
	err := m.primaryCache.Set(ctx, key, value, expiration)
	// 如果主缓存设置失败，则使用备份缓存
	if err != nil {
		// 记录错误日志
		log.Println("cacheManager use redis Set error: %v", err)
	}
	return m.secondaryCache.Set(ctx, key, value, expiration)
}

// Get 尝试从主缓存获取值，如果失败则从备份缓存获取
func (m *cacheManager) Get(key string) (string, bool) {
	if m.primaryCache == nil {
		res, isExist, err := m.secondaryCache.Get(ctx, key)
		// 记录错误日志
		log.Println("cacheManager use mysql Get error: %v", err)
		return res, isExist
	}
	value, isExist, err := m.primaryCache.Get(ctx, key)
	if err != nil {
		// 记录错误日志
		log.Println("cacheManager use redis Get error: %v", err)
		res, isExist, err := m.secondaryCache.Get(ctx, key)
		log.Println("cacheManager use after mysql Get error: %v", err)
		return res, isExist
	}
	return value, isExist
}

// Delete 尝试在主缓存中删除值，如果失败则在备份缓存中删除
func (m *cacheManager) Del(key string) error {
	if m.primaryCache == nil {
		return m.secondaryCache.Delete(ctx, key)
	}
	err := m.primaryCache.Delete(ctx, key)
	if err != nil {
		// 记录错误日志
		log.Println("cacheManager use redis Delete error: %v", err)
	}
	return m.secondaryCache.Delete(ctx, key)
}
