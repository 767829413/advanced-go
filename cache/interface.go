package cache

import (
	"context"
	"time"
)

// Cache 接口定义了缓存操作的标准方法
type Cache interface {
	// 设置缓存过期时间的单位是秒
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	// 获取缓存
	Get(ctx context.Context, key string) (string, error)
	// 删除缓存
	Delete(ctx context.Context, key string) error
}
