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

const (
	userSettingQueryCacheTTL  = 30 * time.Minute
	userSettingQueryKeyPrefix = "usersetting:user:"
)

// userSettingQueryCacheDTO 用于 Redis 缓存的 UserSetting 数据结构。
type userSettingQueryCacheDTO struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	SettingKey string `json:"setting_key"`
	Value      any    `json:"value"`
	Version    int64  `json:"v"` // 版本号，使用 UpdatedAt.UnixNano()
}

func toUserSettingQueryCacheDTO(s *setting.UserSetting) userSettingQueryCacheDTO {
	return userSettingQueryCacheDTO{
		ID:         s.ID,
		UserID:     s.UserID,
		SettingKey: s.SettingKey,
		Value:      s.Value,
		Version:    s.UpdatedAt.UnixNano(),
	}
}

func (d userSettingQueryCacheDTO) toEntity() *setting.UserSetting {
	return &setting.UserSetting{
		ID:         d.ID,
		UserID:     d.UserID,
		SettingKey: d.SettingKey,
		Value:      d.Value,
	}
}

// userSettingQueryCacheService 用户设置查询缓存的 Redis 实现。
//
// 存储原始 UserSetting 记录，用于 Repository 层减少数据库查询。
// 采用用户维度全量缓存策略：一次查询缓存用户所有自定义配置。
type userSettingQueryCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewUserSettingQueryCacheService 创建用户设置查询缓存服务。
func NewUserSettingQueryCacheService(client *redis.Client, keyPrefix string) cache.UserSettingQueryCacheService {
	return &userSettingQueryCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// GetByUser 获取用户的所有自定义配置缓存。
func (s *userSettingQueryCacheService) GetByUser(ctx context.Context, userID uint) (map[string]*setting.UserSetting, error) {
	key := s.buildKey(userID)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var dtos []userSettingQueryCacheDTO
	if err := json.Unmarshal(data, &dtos); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, key)
		return nil, nil //nolint:nilerr,nilnil // corrupted cache treated as miss
	}

	result := make(map[string]*setting.UserSetting, len(dtos))
	for _, dto := range dtos {
		result[dto.SettingKey] = dto.toEntity()
	}
	return result, nil
}

// SetByUser 设置用户的所有自定义配置缓存。
func (s *userSettingQueryCacheService) SetByUser(ctx context.Context, userID uint, settings []*setting.UserSetting) error {
	dtos := make([]userSettingQueryCacheDTO, 0, len(settings))
	for _, us := range settings {
		dtos = append(dtos, toUserSettingQueryCacheDTO(us))
	}

	data, err := json.Marshal(dtos)
	if err != nil {
		return fmt.Errorf("failed to marshal user settings: %w", err)
	}

	key := s.buildKey(userID)
	return s.client.Set(ctx, key, data, userSettingQueryCacheTTL).Err()
}

// DeleteByUser 删除用户的所有配置缓存。
func (s *userSettingQueryCacheService) DeleteByUser(ctx context.Context, userID uint) error {
	key := s.buildKey(userID)
	return s.client.Del(ctx, key).Err()
}

// buildKey 构建缓存 key。
// 格式：{prefix}usersetting:user:{userID}
func (s *userSettingQueryCacheService) buildKey(userID uint) string {
	return fmt.Sprintf("%s%s%d", s.keyPrefix, userSettingQueryKeyPrefix, userID)
}

var _ cache.UserSettingQueryCacheService = (*userSettingQueryCacheService)(nil)
