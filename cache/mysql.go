package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// mySQLCache 实现了Cache接口，使用gorm作为数据库操作工具
type mySQLCache struct {
	db *gorm.DB
}

func newMySQLCache(db *gorm.DB) *mySQLCache {
	return &mySQLCache{db: db}
}

func (c *mySQLCache) Set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	valueStr, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// local_cache 表结构为 cache_key VARCHAR, cache_value TEXT, expire_duration BIGINT
	// 使用gorm的clause.OnConflict进行冲突处理，如果键存在则更新记录
	return c.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "cache_key"}}, // 指定冲突发生的列
		DoUpdates: clause.AssignmentColumns(
			[]string{"cache_value", "expire_duration"},
		), // 指定发生冲突时更新的列
	}).Create(&localCache{
		Key:            key,
		Value:          string(valueStr),
		ExpireDuration: time.Now().Add(expiration).UnixMilli(),
	}).Error
}

func (c *mySQLCache) Get(ctx context.Context, key string) (string, error) {
	var cacheEntry localCache
	result := c.db.WithContext(ctx).
		First(&cacheEntry, "cache_key = ? AND expire_duration > ?", key, time.Now().UnixMilli())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", result.Error
	}
	var value string
	err := json.Unmarshal([]byte(cacheEntry.Value), &value)
	return value, err
}

func (c *mySQLCache) Delete(ctx context.Context, key string) error {
	return c.db.WithContext(ctx).Delete(&localCache{}, "cache_key = ?", key).Error
}

/*
CREATE TABLE `local_cache` (

	`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
	`cache_key` varchar(255)  NOT NULL COMMENT '查询的key',
	`cache_value` text NOT NULL COMMENT '存储值',
	`expire_duration` bigint NOT NULL COMMENT '失效时间',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE KEY `unique_redis_key` (`cache_key`) USING BTREE,
	KEY `idx_expire_duration` (`expire_duration`) USING BTREE COMMENT '失效时间索引'

) ENGINE=InnoDB COMMENT='本地缓存数据表';
*/
type localCache struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Key            string `gorm:"column:cache_key"`
	Value          string `gorm:"column:cache_value"`
	ExpireDuration int64  `gorm:"column:expire_duration"`
}

// TableName FileSystem's table name
func (*localCache) TableName() string {
	return "local_cache"
}
