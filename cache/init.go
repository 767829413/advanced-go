package cache

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var cacheManagerIns *cacheManager

func InitByRedisMysql(redisClient *redis.Client, db *gorm.DB) {
	// 初始化Redis和MySQL缓存实例
	redisCache := newRedisCache(redisClient)
	mysqlCache := newMySQLCache(db)

	// 创建CacheManager实例，以Redis为主缓存，MySQL为备份缓存
	cacheManagerIns = newCacheManager(redisCache, mysqlCache)
}

func GetCacheManager() *cacheManager {
	return cacheManagerIns
}
