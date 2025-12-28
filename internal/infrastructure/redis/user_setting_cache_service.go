package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
)

const userSettingCacheTTL = 30 * time.Minute

// userSettingCacheService 用户设置缓存服务的 Redis 实现。
//
// Key 格式：{prefix}user:{userID}:setting:{key}
// TTL：30 分钟
type userSettingCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewUserSettingCacheService 创建用户设置缓存服务。
func NewUserSettingCacheService(client *redis.Client, keyPrefix string) cache.UserSettingCacheService {
	return &userSettingCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// =========================================================================
// 单条操作
// =========================================================================

// Get 获取用户的有效设置值。
func (s *userSettingCacheService) Get(ctx context.Context, userID uint, key string) (*cache.EffectiveUserSetting, error) {
	data, err := s.client.Get(ctx, s.buildKey(userID, key)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var result cache.EffectiveUserSetting
	if err := json.Unmarshal(data, &result); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, s.buildKey(userID, key))
		return nil, nil //nolint:nilerr,nilnil // corrupted cache
	}

	return &result, nil
}

// Set 缓存用户的有效设置值。
func (s *userSettingCacheService) Set(ctx context.Context, userID uint, value *cache.EffectiveUserSetting) error {
	if value == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal effective user setting: %w", err)
	}

	return s.client.Set(ctx, s.buildKey(userID, value.Key), data, userSettingCacheTTL).Err()
}

// userSetWithVersionScript Lua 脚本：版本化写入用户设置
// KEYS[1]: 缓存 key
// ARGV[1]: JSON 数据
// ARGV[2]: 新版本号
// ARGV[3]: TTL 秒数
// 返回：1 表示写入成功，0 表示版本过旧被跳过
//
//nolint:dupword
var userSetWithVersionScript = redis.NewScript(`
local current = redis.call('GET', KEYS[1])
if current then
    local data = cjson.decode(current)
    if data.v and data.v >= tonumber(ARGV[2]) then
        return 0
    end
end
redis.call('SETEX', KEYS[1], ARGV[3], ARGV[1])
return 1
`)

// SetWithVersion 版本化设置缓存。
func (s *userSettingCacheService) SetWithVersion(ctx context.Context, userID uint, value *cache.EffectiveUserSetting, version int64) (bool, error) {
	if value == nil {
		return false, nil
	}

	// 确保版本号写入
	value.Version = version

	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal effective user setting: %w", err)
	}

	ttlSeconds := int(userSettingCacheTTL.Seconds())
	result, err := userSetWithVersionScript.Run(ctx, s.client, []string{s.buildKey(userID, value.Key)}, string(data), version, ttlSeconds).Int()
	if err != nil {
		return false, fmt.Errorf("failed to run version script: %w", err)
	}

	return result == 1, nil
}

// =========================================================================
// 批量操作
// =========================================================================

// GetByKeys 批量获取用户的有效设置值。
func (s *userSettingCacheService) GetByKeys(ctx context.Context, userID uint, keys []string) (map[string]*cache.EffectiveUserSetting, error) {
	if len(keys) == 0 {
		return map[string]*cache.EffectiveUserSetting{}, nil
	}

	redisKeys := make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = s.buildKey(userID, k)
	}

	values, err := s.client.MGet(ctx, redisKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to mget: %w", err)
	}

	result := make(map[string]*cache.EffectiveUserSetting)
	for i, v := range values {
		if v == nil {
			continue
		}

		data, ok := v.(string)
		if !ok {
			continue
		}

		var setting cache.EffectiveUserSetting
		if json.Unmarshal([]byte(data), &setting) == nil {
			result[keys[i]] = &setting
		}
	}

	return result, nil
}

// SetBatch 批量设置用户的有效设置值。
func (s *userSettingCacheService) SetBatch(ctx context.Context, userID uint, values []*cache.EffectiveUserSetting) error {
	if len(values) == 0 {
		return nil
	}

	pipe := s.client.Pipeline()
	for _, v := range values {
		if v == nil {
			continue
		}

		data, err := json.Marshal(v)
		if err != nil {
			slog.Warn("failed to marshal effective user setting for batch set", "key", v.Key, "err", err)
			continue
		}
		pipe.Set(ctx, s.buildKey(userID, v.Key), data, userSettingCacheTTL)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}

	return nil
}

// =========================================================================
// 删除操作
// =========================================================================

// Delete 删除用户的指定设置缓存。
func (s *userSettingCacheService) Delete(ctx context.Context, userID uint, key string) error {
	return s.client.Del(ctx, s.buildKey(userID, key)).Err()
}

// DeleteByKeys 批量删除用户的指定设置缓存。
func (s *userSettingCacheService) DeleteByKeys(ctx context.Context, userID uint, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	redisKeys := make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = s.buildKey(userID, k)
	}

	return s.client.Del(ctx, redisKeys...).Err()
}

// DeleteByUser 删除用户的所有设置缓存。
func (s *userSettingCacheService) DeleteByUser(ctx context.Context, userID uint) error {
	pattern := s.buildUserPattern(userID)
	return s.deleteByPattern(ctx, pattern)
}

// DeleteBySettingKey 删除所有用户的某个设置缓存。
func (s *userSettingCacheService) DeleteBySettingKey(ctx context.Context, key string) error {
	pattern := s.buildSettingKeyPattern(key)
	return s.deleteByPattern(ctx, pattern)
}

// DeleteBySettingKeys 批量删除所有用户的多个设置缓存。
func (s *userSettingCacheService) DeleteBySettingKeys(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := s.DeleteBySettingKey(ctx, key); err != nil {
			slog.Warn("failed to delete user setting cache by key", "key", key, "err", err)
		}
	}
	return nil
}

// deleteByPattern 使用 SCAN 删除匹配模式的所有 key。
func (s *userSettingCacheService) deleteByPattern(ctx context.Context, pattern string) error {
	iter := s.client.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}

	if len(keys) > 0 {
		return s.client.Del(ctx, keys...).Err()
	}

	return nil
}

// =========================================================================
// 内部辅助方法
// =========================================================================

func (s *userSettingCacheService) buildKey(userID uint, settingKey string) string {
	return fmt.Sprintf("%suser:%d:setting:%s", s.keyPrefix, userID, settingKey)
}

func (s *userSettingCacheService) buildUserPattern(userID uint) string {
	return fmt.Sprintf("%suser:%d:setting:*", s.keyPrefix, userID)
}

func (s *userSettingCacheService) buildSettingKeyPattern(settingKey string) string {
	return fmt.Sprintf("%suser:*:setting:%s", s.keyPrefix, settingKey)
}

var _ cache.UserSettingCacheService = (*userSettingCacheService)(nil)
