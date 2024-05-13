package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	errNilRedisClient = errors.New("redis client is nil")
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
	if c.client == nil {
		return errNilRedisClient
	}
	return c.client.Set(ctx, key, getValue(value), expiration).Err()
}

func (c *redisCache) Get(ctx context.Context, key string) (string, bool, error) {
	if c.client == nil {
		return "", false, errNilRedisClient
	}
	res := c.client.Get(ctx, key)
	result, err := res.Val(), res.Err()
	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	return result, true, nil
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return errNilRedisClient
	}
	return c.client.Del(ctx, key).Err()
}
