package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

const (
	settingCacheTTL  = 7 * 24 * time.Hour // 7 天（数据变更不频繁，且有主动失效机制）
	settingKeyPrefix = "setting:"
)

// settingCacheDTO 用于 Redis 缓存的 Setting 数据结构。
//
// 将领域实体序列化为 JSON 存储，避免在领域层添加 JSON tags。
type settingCacheDTO struct {
	ID             uint   `json:"id"`
	Key            string `json:"key"`
	DefaultValue   any    `json:"default_value"`
	Scope          string `json:"scope"`
	CategoryID     uint   `json:"category_id"`
	Group          string `json:"group"`
	ValueType      string `json:"value_type"`
	Label          string `json:"label"`
	InputType      string `json:"input_type"`
	Validation     string `json:"validation"`
	UIConfig       string `json:"ui_config"`
	Order          int    `json:"order"`
	ViewPermission string `json:"view_permission"`
	EditPermission string `json:"edit_permission"`
}

// toDTO 将领域实体转换为缓存 DTO。
func toSettingCacheDTO(s *setting.Setting) settingCacheDTO {
	return settingCacheDTO{
		ID:             s.ID,
		Key:            s.Key,
		DefaultValue:   s.DefaultValue,
		Scope:          s.Scope,
		CategoryID:     s.CategoryID,
		Group:          s.Group,
		ValueType:      s.ValueType,
		Label:          s.Label,
		InputType:      s.InputType,
		Validation:     s.Validation,
		UIConfig:       s.UIConfig,
		Order:          s.Order,
		ViewPermission: s.ViewPermission,
		EditPermission: s.EditPermission,
	}
}

// toEntity 将缓存 DTO 转换为领域实体。
func (d settingCacheDTO) toEntity() *setting.Setting {
	return &setting.Setting{
		ID:             d.ID,
		Key:            d.Key,
		DefaultValue:   d.DefaultValue,
		Scope:          d.Scope,
		CategoryID:     d.CategoryID,
		Group:          d.Group,
		ValueType:      d.ValueType,
		Label:          d.Label,
		InputType:      d.InputType,
		Validation:     d.Validation,
		UIConfig:       d.UIConfig,
		Order:          d.Order,
		ViewPermission: d.ViewPermission,
		EditPermission: d.EditPermission,
	}
}

// settingCacheService Setting 缓存服务的 Redis 实现。
//
// 提供按 key 粒度的缓存操作：
//   - Key 格式：{prefix}setting:{key}
//   - TTL：7 天
//   - JSON 序列化存储（RedisJSON）
type settingCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewSettingCacheService 创建 Setting 缓存服务。
func NewSettingCacheService(client *redis.Client, keyPrefix string) cache.SettingCacheService {
	return &settingCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// Get 获取缓存的 Setting（使用 RedisJSON）。
// 缓存未命中返回 nil, nil。
func (s *settingCacheService) Get(ctx context.Context, key string) (*setting.Setting, error) {
	data, err := s.client.JSONGet(ctx, s.buildKey(key), "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss returns nil, nil
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	var wrapper []settingCacheDTO
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, s.buildKey(key))
		return nil, nil //nolint:nilerr,nilnil // corrupted cache treated as miss
	}

	if len(wrapper) == 0 {
		return nil, nil //nolint:nilnil // empty wrapper treated as miss
	}

	return wrapper[0].toEntity(), nil
}

// Set 设置 Setting 缓存（使用 RedisJSON）。
func (s *settingCacheService) Set(ctx context.Context, st *setting.Setting) error {
	if st == nil {
		return nil
	}

	dto := toSettingCacheDTO(st)

	// 使用 Pipeline 执行 JSON.SET + EXPIRE
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, s.buildKey(st.Key), "$", dto)
	pipe.Expire(ctx, s.buildKey(st.Key), settingCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set setting cache: %w", err)
	}

	return nil
}

// Delete 删除指定 key 的缓存。
func (s *settingCacheService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, s.buildKey(key)).Err()
}

// DeleteByKeys 批量删除缓存。
func (s *settingCacheService) DeleteByKeys(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	redisKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		redisKeys = append(redisKeys, s.buildKey(k))
	}

	return s.client.Del(ctx, redisKeys...).Err()
}

// DeleteAll 删除所有 Setting 缓存。
func (s *settingCacheService) DeleteAll(ctx context.Context) error {
	pattern := s.keyPrefix + "setting:*"

	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
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
// 批量操作
// =========================================================================

// GetAll 获取所有缓存的 Setting（使用 RedisJSON）。
func (s *settingCacheService) GetAll(ctx context.Context) (map[string]*setting.Setting, error) {
	pattern := s.keyPrefix + settingKeyPrefix + "*"
	result := make(map[string]*setting.Setting)

	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		redisKey := iter.Val()

		// 跳过预热相关 key
		if strings.HasSuffix(redisKey, "_warmed_up") || strings.HasSuffix(redisKey, "_warmup_lock") {
			continue
		}

		data, err := s.client.JSONGet(ctx, redisKey, "$").Result()
		if err != nil {
			continue
		}

		// JSON.GET $ 返回数组包装
		var wrapper []settingCacheDTO
		if json.Unmarshal([]byte(data), &wrapper) == nil && len(wrapper) > 0 {
			result[wrapper[0].Key] = wrapper[0].toEntity()
		}
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan keys: %w", err)
	}

	return result, nil
}

// SetAll 批量设置 Setting 缓存（使用 RedisJSON，用于预热）。
func (s *settingCacheService) SetAll(ctx context.Context, settings []*setting.Setting) error {
	if len(settings) == 0 {
		return nil
	}

	pipe := s.client.Pipeline()
	for _, st := range settings {
		dto := toSettingCacheDTO(st)
		key := s.buildKey(st.Key)
		pipe.JSONSet(ctx, key, "$", dto)
		pipe.Expire(ctx, key, settingCacheTTL)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}

	return nil
}

// GetByKeys 批量获取指定 key 的缓存（使用 RedisJSON）。
func (s *settingCacheService) GetByKeys(ctx context.Context, keys []string) (map[string]*setting.Setting, error) {
	if len(keys) == 0 {
		return map[string]*setting.Setting{}, nil
	}

	redisKeys := make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = s.buildKey(k)
	}

	// JSON.MGET 返回每个 key 的 JSON 字符串
	values, err := s.client.JSONMGet(ctx, "$", redisKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to json mget: %w", err)
	}

	result := make(map[string]*setting.Setting)
	for i, v := range values {
		if v == nil {
			continue
		}

		data, ok := v.(string)
		if !ok {
			continue
		}

		// JSON.MGET 返回数组包装
		var wrapper []settingCacheDTO
		if json.Unmarshal([]byte(data), &wrapper) == nil && len(wrapper) > 0 {
			result[keys[i]] = wrapper[0].toEntity()
		}
	}

	return result, nil
}

// =========================================================================
// 内部辅助方法
// =========================================================================

func (s *settingCacheService) buildKey(settingKey string) string {
	return s.keyPrefix + settingKeyPrefix + settingKey
}

var _ cache.SettingCacheService = (*settingCacheService)(nil)
