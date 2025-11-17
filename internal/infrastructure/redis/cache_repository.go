// Package redis 提供 Redis 缓存仓储实现
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheRepository Redis 缓存仓储接口
type CacheRepository interface {
	// Get 获取缓存值
	Get(ctx context.Context, key string, dest interface{}) error
	// Set 设置缓存值
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Delete 删除缓存值
	Delete(ctx context.Context, key string) error
	// Exists 检查键是否存在
	Exists(ctx context.Context, key string) (bool, error)
	// SetNX 仅当键不存在时设置值（分布式锁）
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
}

// cacheRepository Redis 缓存仓储实现
type cacheRepository struct {
	client *redis.Client
}

// NewCacheRepository 创建缓存仓储实例
func NewCacheRepository(client *redis.Client) CacheRepository {
	return &cacheRepository{
		client: client,
	}
}

// Get 获取缓存值并反序列化到目标对象
func (r *cacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}

	return nil
}

// Set 序列化并设置缓存值
func (r *cacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// Delete 删除缓存值
func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists 检查键是否存在
func (r *cacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return n > 0, nil
}

// SetNX 仅当键不存在时设置值（用于分布式锁）
func (r *cacheRepository) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	ok, err := r.client.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to setnx key %s: %w", key, err)
	}

	return ok, nil
}
