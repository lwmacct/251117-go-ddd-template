package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

const (
	settingCategoryCacheTTL       = 10 * time.Minute
	settingCategoryWarmupFlagTTL  = 24 * time.Hour // 预热标记独立 TTL，避免与数据同时过期
	settingCategoryWarmupLockTTL  = 30 * time.Second
	settingCategoryWarmupPollItvl = 100 * time.Millisecond
	settingCategoryKeyPrefix      = "setting_category:"
	settingCategoryIDPrefix       = "setting_category:id:"
	settingCategoryKeyByKeyPrefix = "setting_category:key:"
	settingCategoryAllKey         = "setting_category:all"
	settingCategoryWarmupLockKey  = "setting_category:_warmup_lock"
	settingCategoryWarmedUpKey    = "setting_category:_warmed_up"
)

// settingCategoryCacheDTO 用于 Redis 缓存的 SettingCategory 数据结构。
//
// 将领域实体序列化为 JSON 存储，避免在领域层添加 JSON tags。
type settingCategoryCacheDTO struct {
	ID    uint   `json:"id"`
	Key   string `json:"key"`
	Label string `json:"label"`
	Icon  string `json:"icon"`
	Order int    `json:"order"`
}

// toSettingCategoryCacheDTO 将领域实体转换为缓存 DTO。
func toSettingCategoryCacheDTO(c *setting.SettingCategory) settingCategoryCacheDTO {
	return settingCategoryCacheDTO{
		ID:    c.ID,
		Key:   c.Key,
		Label: c.Label,
		Icon:  c.Icon,
		Order: c.Order,
	}
}

// toEntity 将缓存 DTO 转换为领域实体。
func (d settingCategoryCacheDTO) toEntity() *setting.SettingCategory {
	return &setting.SettingCategory{
		ID:    d.ID,
		Key:   d.Key,
		Label: d.Label,
		Icon:  d.Icon,
		Order: d.Order,
	}
}

// settingCategoryCacheService SettingCategory 缓存服务的 Redis 实现。
//
// 提供按 ID 和 Key 粒度的缓存操作：
//   - ID 索引：{prefix}setting_category:id:{id}
//   - Key 索引：{prefix}setting_category:key:{key}
//   - 全量缓存：{prefix}setting_category:all
//   - TTL：10 分钟
//   - JSON 序列化存储（RedisJSON）
type settingCategoryCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewSettingCategoryCacheService 创建 SettingCategory 缓存服务。
func NewSettingCategoryCacheService(client *redis.Client, keyPrefix string) cache.SettingCategoryCacheService {
	return &settingCategoryCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// =========================================================================
// 单条操作
// =========================================================================

// GetByID 根据 ID 获取缓存的 SettingCategory（使用 RedisJSON）。
func (s *settingCategoryCacheService) GetByID(ctx context.Context, id uint) (*setting.SettingCategory, error) {
	data, err := s.client.JSONGet(ctx, s.buildIDKey(id), "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss returns nil, nil
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	return s.unmarshalSingle(ctx, data, s.buildIDKey(id))
}

// GetByKey 根据 Key 获取缓存的 SettingCategory（使用 RedisJSON）。
func (s *settingCategoryCacheService) GetByKey(ctx context.Context, key string) (*setting.SettingCategory, error) {
	data, err := s.client.JSONGet(ctx, s.buildKeyKey(key), "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss returns nil, nil
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	return s.unmarshalSingle(ctx, data, s.buildKeyKey(key))
}

// Set 设置 SettingCategory 缓存（使用 RedisJSON）。
// 同时写入 ID 索引和 Key 索引。
func (s *settingCategoryCacheService) Set(ctx context.Context, category *setting.SettingCategory) error {
	if category == nil {
		return nil
	}

	dto := toSettingCategoryCacheDTO(category)

	pipe := s.client.Pipeline()

	// 写入 ID 索引
	idKey := s.buildIDKey(category.ID)
	pipe.JSONSet(ctx, idKey, "$", dto)
	pipe.Expire(ctx, idKey, settingCategoryCacheTTL)

	// 写入 Key 索引
	keyKey := s.buildKeyKey(category.Key)
	pipe.JSONSet(ctx, keyKey, "$", dto)
	pipe.Expire(ctx, keyKey, settingCategoryCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set setting category cache: %w", err)
	}

	return nil
}

// =========================================================================
// 批量操作
// =========================================================================

// GetAll 获取所有缓存的 SettingCategory。
// 从全量缓存 key 读取，返回列表。
func (s *settingCategoryCacheService) GetAll(ctx context.Context) ([]*setting.SettingCategory, error) {
	data, err := s.client.JSONGet(ctx, s.keyPrefix+settingCategoryAllKey, "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[[actual_array]]
	var wrapper [][]settingCategoryCacheDTO
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		_ = s.client.Del(ctx, s.keyPrefix+settingCategoryAllKey)
		return nil, nil //nolint:nilerr // corrupted cache treated as miss
	}

	if len(wrapper) == 0 || len(wrapper[0]) == 0 {
		return nil, nil // empty wrapper treated as miss
	}

	result := make([]*setting.SettingCategory, 0, len(wrapper[0]))
	for _, dto := range wrapper[0] {
		result = append(result, dto.toEntity())
	}

	return result, nil
}

// SetAll 批量设置 SettingCategory 缓存（用于预热）。
// 同时写入全量缓存、ID 索引和 Key 索引。
func (s *settingCategoryCacheService) SetAll(ctx context.Context, categories []*setting.SettingCategory) error {
	if len(categories) == 0 {
		return nil
	}

	pipe := s.client.Pipeline()

	// 构建全量缓存的 DTO 数组
	dtos := make([]settingCategoryCacheDTO, 0, len(categories))

	for _, c := range categories {
		dto := toSettingCategoryCacheDTO(c)
		dtos = append(dtos, dto)

		// 写入 ID 索引
		idKey := s.buildIDKey(c.ID)
		pipe.JSONSet(ctx, idKey, "$", dto)
		pipe.Expire(ctx, idKey, settingCategoryCacheTTL)

		// 写入 Key 索引
		keyKey := s.buildKeyKey(c.Key)
		pipe.JSONSet(ctx, keyKey, "$", dto)
		pipe.Expire(ctx, keyKey, settingCategoryCacheTTL)
	}

	// 写入全量缓存
	allKey := s.keyPrefix + settingCategoryAllKey
	pipe.JSONSet(ctx, allKey, "$", dtos)
	pipe.Expire(ctx, allKey, settingCategoryCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}

	return nil
}

// GetByIDs 批量获取指定 ID 的缓存。
func (s *settingCategoryCacheService) GetByIDs(ctx context.Context, ids []uint) ([]*setting.SettingCategory, error) {
	if len(ids) == 0 {
		return []*setting.SettingCategory{}, nil
	}

	redisKeys := make([]string, len(ids))
	for i, id := range ids {
		redisKeys[i] = s.buildIDKey(id)
	}

	// JSON.MGET 返回每个 key 的 JSON 字符串
	values, err := s.client.JSONMGet(ctx, "$", redisKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to json mget: %w", err)
	}

	result := make([]*setting.SettingCategory, 0, len(ids))
	for _, v := range values {
		if v == nil {
			continue
		}

		data, ok := v.(string)
		if !ok {
			continue
		}

		// JSON.MGET 返回数组包装
		var wrapper []settingCategoryCacheDTO
		if json.Unmarshal([]byte(data), &wrapper) == nil && len(wrapper) > 0 {
			result = append(result, wrapper[0].toEntity())
		}
	}

	return result, nil
}

// =========================================================================
// 删除操作
// =========================================================================

// Delete 删除指定的缓存。
// 同时删除 ID 索引、Key 索引和全量缓存。
func (s *settingCategoryCacheService) Delete(ctx context.Context, id uint, key string) error {
	keys := []string{
		s.buildIDKey(id),
		s.buildKeyKey(key),
		s.keyPrefix + settingCategoryAllKey, // 全量缓存也需要失效
	}
	return s.client.Del(ctx, keys...).Err()
}

// DeleteAll 删除所有 SettingCategory 缓存。
func (s *settingCategoryCacheService) DeleteAll(ctx context.Context) error {
	pattern := s.keyPrefix + settingCategoryKeyPrefix + "*"

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
// 预热控制
// =========================================================================

// IsWarmedUp 检查缓存是否已预热完成。
func (s *settingCategoryCacheService) IsWarmedUp(ctx context.Context) bool {
	exists, err := s.client.Exists(ctx, s.keyPrefix+settingCategoryWarmedUpKey).Result()
	if err != nil {
		slog.Warn("failed to check setting category warmup status", "err", err)
		return false
	}
	return exists > 0
}

// SetWarmedUp 标记缓存已预热完成。
// 使用独立的 24 小时 TTL，避免与数据缓存同时过期。
func (s *settingCategoryCacheService) SetWarmedUp(ctx context.Context) error {
	return s.client.Set(ctx, s.keyPrefix+settingCategoryWarmedUpKey, "1", settingCategoryWarmupFlagTTL).Err()
}

// TryAcquireWarmupLock 尝试获取预热分布式锁。
func (s *settingCategoryCacheService) TryAcquireWarmupLock(ctx context.Context) (bool, func()) {
	lockKey := s.keyPrefix + settingCategoryWarmupLockKey

	// 尝试获取锁（SETNX + TTL）
	ok, err := s.client.SetNX(ctx, lockKey, "1", settingCategoryWarmupLockTTL).Result()
	if err != nil {
		slog.Warn("failed to acquire setting category warmup lock", "err", err)
		return false, nil
	}

	if !ok {
		return false, nil
	}

	// 返回释放锁的函数
	return true, func() { //nolint:contextcheck // 故意使用独立 context
		if err := s.client.Del(context.Background(), lockKey).Err(); err != nil {
			slog.Warn("failed to release setting category warmup lock", "err", err)
		}
	}
}

// WaitForWarmup 等待缓存预热完成。
func (s *settingCategoryCacheService) WaitForWarmup(ctx context.Context) error {
	ticker := time.NewTicker(settingCategoryWarmupPollItvl)
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

func (s *settingCategoryCacheService) buildIDKey(id uint) string {
	return s.keyPrefix + settingCategoryIDPrefix + strconv.FormatUint(uint64(id), 10)
}

func (s *settingCategoryCacheService) buildKeyKey(key string) string {
	return s.keyPrefix + settingCategoryKeyByKeyPrefix + key
}

// unmarshalSingle 解析单条缓存数据。
func (s *settingCategoryCacheService) unmarshalSingle(ctx context.Context, data, redisKey string) (*setting.SettingCategory, error) {
	// JSON.GET $ 返回数组包装：[actual_data]
	var wrapper []settingCategoryCacheDTO
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, redisKey)
		return nil, nil //nolint:nilerr,nilnil // corrupted cache treated as miss
	}

	if len(wrapper) == 0 {
		return nil, nil //nolint:nilnil // empty wrapper treated as miss
	}

	return wrapper[0].toEntity(), nil
}

var _ cache.SettingCategoryCacheService = (*settingCategoryCacheService)(nil)
