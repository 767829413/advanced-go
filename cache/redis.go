package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// redisCache 实现了Cache接口
type redisCache struct {
	client *redis.Client
}

func newRedisCache(client *redis.Client) *redisCache {
	return &redisCache{client: client}
}

func (c *redisCache) Set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
