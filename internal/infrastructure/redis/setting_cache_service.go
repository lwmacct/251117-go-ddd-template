package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

const (
	settingCacheTTL    = 10 * time.Minute
	warmupLockTTL      = 30 * time.Second
	warmupPollInterval = 100 * time.Millisecond
	settingKeyPrefix   = "setting:"
	warmupLockKey      = "setting:_warmup_lock"
	warmedUpKey        = "setting:_warmed_up"
)

// settingCacheDTO 用于 Redis 缓存的 Setting 数据结构。
//
// 将领域实体序列化为 JSON 存储，避免在领域层添加 JSON tags。
// Version 字段用于回写竞争检测。
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
	Version        int64  `json:"v"` // 版本号，使用 UpdatedAt.UnixNano()
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
		Version:        s.UpdatedAt.UnixNano(),
	}
}

// toSettingCacheDTOWithVersion 将领域实体转换为缓存 DTO，使用指定版本号。
func toSettingCacheDTOWithVersion(s *setting.Setting, version int64) settingCacheDTO {
	dto := toSettingCacheDTO(s)
	dto.Version = version
	return dto
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

// =========================================================================
// 版本化写入
// =========================================================================

// setWithVersionScript Lua 脚本：版本化写入
// KEYS[1]: 缓存 key
// ARGV[1]: JSON 数据
// ARGV[2]: 新版本号
// ARGV[3]: TTL 秒数
// 返回：1 表示写入成功，0 表示版本过旧被跳过
//
//nolint:dupword
var setWithVersionScript = redis.NewScript(`
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
func (s *settingCacheService) SetWithVersion(ctx context.Context, st *setting.Setting, version int64) (bool, error) {
	if st == nil {
		return false, nil
	}

	dto := toSettingCacheDTOWithVersion(st, version)
	data, err := json.Marshal(dto)
	if err != nil {
		return false, fmt.Errorf("failed to marshal setting: %w", err)
	}

	ttlSeconds := int(settingCacheTTL.Seconds())
	result, err := setWithVersionScript.Run(ctx, s.client, []string{s.buildKey(st.Key)}, string(data), version, ttlSeconds).Int()
	if err != nil {
		return false, fmt.Errorf("failed to run version script: %w", err)
	}

	return result == 1, nil
}

// =========================================================================
// 批量操作
// =========================================================================

// GetAll 获取所有缓存的 Setting。
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

		data, err := s.client.Get(ctx, redisKey).Bytes()
		if err != nil {
			continue
		}

		var dto settingCacheDTO
		if json.Unmarshal(data, &dto) == nil {
			result[dto.Key] = dto.toEntity()
		}
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan keys: %w", err)
	}

	return result, nil
}

// SetAll 批量设置 Setting 缓存（用于预热）。
func (s *settingCacheService) SetAll(ctx context.Context, settings []*setting.Setting) error {
	if len(settings) == 0 {
		return nil
	}

	pipe := s.client.Pipeline()
	for _, st := range settings {
		dto := toSettingCacheDTO(st)
		data, err := json.Marshal(dto)
		if err != nil {
			slog.Warn("failed to marshal setting for batch set", "key", st.Key, "err", err)
			continue
		}
		pipe.Set(ctx, s.buildKey(st.Key), data, settingCacheTTL)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}

	return nil
}

// GetByKeys 批量获取指定 key 的缓存。
func (s *settingCacheService) GetByKeys(ctx context.Context, keys []string) (map[string]*setting.Setting, error) {
	if len(keys) == 0 {
		return map[string]*setting.Setting{}, nil
	}

	redisKeys := make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = s.buildKey(k)
	}

	values, err := s.client.MGet(ctx, redisKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to mget: %w", err)
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

		var dto settingCacheDTO
		if json.Unmarshal([]byte(data), &dto) == nil {
			result[keys[i]] = dto.toEntity()
		}
	}

	return result, nil
}

// =========================================================================
// 预热控制
// =========================================================================

// IsWarmedUp 检查缓存是否已预热完成。
func (s *settingCacheService) IsWarmedUp(ctx context.Context) bool {
	exists, err := s.client.Exists(ctx, s.keyPrefix+warmedUpKey).Result()
	if err != nil {
		slog.Warn("failed to check warmup status", "err", err)
		return false
	}
	return exists > 0
}

// SetWarmedUp 标记缓存已预热完成。
func (s *settingCacheService) SetWarmedUp(ctx context.Context) error {
	// 使用 0 表示永不过期
	return s.client.Set(ctx, s.keyPrefix+warmedUpKey, "1", 0).Err()
}

// TryAcquireWarmupLock 尝试获取预热分布式锁。
func (s *settingCacheService) TryAcquireWarmupLock(ctx context.Context) (bool, func()) {
	lockKey := s.keyPrefix + warmupLockKey

	// 尝试获取锁（SETNX + TTL）
	ok, err := s.client.SetNX(ctx, lockKey, "1", warmupLockTTL).Result()
	if err != nil {
		slog.Warn("failed to acquire warmup lock", "err", err)
		return false, nil
	}

	if !ok {
		return false, nil
	}

	// 返回释放锁的函数
	// 使用独立 context 确保释放操作不受调用方 context 取消影响
	return true, func() { //nolint:contextcheck // 故意使用独立 context
		if err := s.client.Del(context.Background(), lockKey).Err(); err != nil {
			slog.Warn("failed to release warmup lock", "err", err)
		}
	}
}

// WaitForWarmup 等待缓存预热完成。
func (s *settingCacheService) WaitForWarmup(ctx context.Context) error {
	ticker := time.NewTicker(warmupPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if s.IsWarmedUp(ctx) {
				return nil
			}
		}
	}
}

// =========================================================================
// 内部辅助方法
// =========================================================================

func (s *settingCacheService) buildKey(settingKey string) string {
	return s.keyPrefix + settingKeyPrefix + settingKey
}

var _ cache.SettingCacheService = (*settingCacheService)(nil)
