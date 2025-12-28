package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

const settingCacheTTL = 10 * time.Minute

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
//   - TTL：10 分钟
//   - JSON 序列化存储
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

// Get 获取缓存的 Setting。
// 缓存未命中返回 nil, nil。
func (s *settingCacheService) Get(ctx context.Context, key string) (*setting.Setting, error) {
	data, err := s.client.Get(ctx, s.buildKey(key)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss returns nil, nil
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var dto settingCacheDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, s.buildKey(key))
		return nil, nil //nolint:nilerr,nilnil // corrupted cache treated as miss
	}

	return dto.toEntity(), nil
}

// Set 设置 Setting 缓存。
func (s *settingCacheService) Set(ctx context.Context, st *setting.Setting) error {
	if st == nil {
		return nil
	}

	dto := toSettingCacheDTO(st)
	data, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal setting: %w", err)
	}

	return s.client.Set(ctx, s.buildKey(st.Key), data, settingCacheTTL).Err()
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

func (s *settingCacheService) buildKey(settingKey string) string {
	return s.keyPrefix + "setting:" + settingKey
}

var _ cache.SettingCacheService = (*settingCacheService)(nil)
