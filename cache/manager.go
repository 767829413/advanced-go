package cache

import (
	"context"
	"time"
)

// cacheManager 管理不同的缓存实现，并提供自动切换功能
type cacheManager struct {
	primaryCache   Cache
	secondaryCache Cache
}

// newCacheManager 创建一个新的CacheManager实例
func newCacheManager(primaryCache Cache, secondaryCache Cache) *cacheManager {
	return &cacheManager{
		primaryCache:   primaryCache,
		secondaryCache: secondaryCache,
	}
}

// Set 尝试在主缓存中设置值，如果失败则使用备份缓存
func (m *cacheManager) Set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	err := m.primaryCache.Set(ctx, key, value, expiration)
	// 如果主缓存设置失败，则使用备份缓存
	if err != nil {
		// TODO: 记录错误日志
	}
	return m.secondaryCache.Set(ctx, key, value, expiration)
}

// Get 尝试从主缓存获取值，如果失败则从备份缓存获取
func (m *cacheManager) Get(ctx context.Context, key string) (string, error) {
	value, err := m.primaryCache.Get(ctx, key)
	if err != nil {
		return m.secondaryCache.Get(ctx, key)
	}
	return value, nil
}

// Delete 尝试在主缓存中删除值，如果失败则在备份缓存中删除
func (m *cacheManager) Delete(ctx context.Context, key string) error {
	err := m.primaryCache.Delete(ctx, key)
	if err != nil {
		// TODO: 记录错误日志
	}
	return m.secondaryCache.Delete(ctx, key)
}
