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

	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
)

const (
	schemaCacheTTL         = 30 * time.Minute
	schemaKeyPrefix        = "schema:"
	schemaUserPrefix       = "user:"
	schemaAdminPrefix      = "admin:"
	schemaCategoryAll      = "_all"
	schemaCategoriesPrefix = "categories:"
	scanBatchSize          = 100
	deletePipelineSize     = 100
)

// schemaCacheService Schema 缓存服务的 Redis 实现。
//
// 缓存 Setting Schema API 的最终响应：
//   - Key 格式：{prefix}schema:user:{userID}:{categoryKey}
//   - Key 格式：{prefix}schema:admin:{categoryKey}
//   - TTL：30 分钟
//   - RedisJSON 原生 JSON 类型存储
//   - 直接序列化 Application DTO（无独立缓存 DTO）
type schemaCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewSchemaCacheService 创建 Schema 缓存服务。
func NewSchemaCacheService(client *redis.Client, keyPrefix string) setting.SchemaCacheService {
	return &schemaCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// =========================================================================
// 用户 Schema 操作
// =========================================================================

// GetUserSchema 获取用户 Schema 缓存。
func (s *schemaCacheService) GetUserSchema(ctx context.Context, userID uint, categoryKey string) ([]setting.SchemaCategoryDTO, error) {
	key := s.buildUserKey(userID, categoryKey)
	return s.get(ctx, key)
}

// SetUserSchema 设置用户 Schema 缓存。
func (s *schemaCacheService) SetUserSchema(ctx context.Context, userID uint, categoryKey string, schema []setting.SchemaCategoryDTO) error {
	key := s.buildUserKey(userID, categoryKey)
	return s.set(ctx, key, schema)
}

// DeleteUserSchema 删除用户的指定 category Schema 缓存。
func (s *schemaCacheService) DeleteUserSchema(ctx context.Context, userID uint, categoryKey string) error {
	key := s.buildUserKey(userID, categoryKey)
	return s.client.Del(ctx, key).Err()
}

// DeleteUserSchemaAll 删除用户的所有 Schema 缓存。
func (s *schemaCacheService) DeleteUserSchemaAll(ctx context.Context, userID uint) error {
	pattern := s.keyPrefix + schemaKeyPrefix + schemaUserPrefix + strconv.FormatUint(uint64(userID), 10) + ":*"
	return s.deleteByPattern(ctx, pattern)
}

// =========================================================================
// 管理员 Schema 操作
// =========================================================================

// GetAdminSchema 获取管理员 Schema 缓存。
func (s *schemaCacheService) GetAdminSchema(ctx context.Context, categoryKey string) ([]setting.SchemaCategoryDTO, error) {
	key := s.buildAdminKey(categoryKey)
	return s.get(ctx, key)
}

// SetAdminSchema 设置管理员 Schema 缓存。
func (s *schemaCacheService) SetAdminSchema(ctx context.Context, categoryKey string, schema []setting.SchemaCategoryDTO) error {
	key := s.buildAdminKey(categoryKey)
	return s.set(ctx, key, schema)
}

// DeleteAdminSchema 删除管理员的指定 category Schema 缓存。
func (s *schemaCacheService) DeleteAdminSchema(ctx context.Context, categoryKey string) error {
	key := s.buildAdminKey(categoryKey)
	return s.client.Del(ctx, key).Err()
}

// DeleteAdminSchemaAll 删除管理员的所有 Schema 缓存。
func (s *schemaCacheService) DeleteAdminSchemaAll(ctx context.Context) error {
	pattern := s.keyPrefix + schemaKeyPrefix + schemaAdminPrefix + "*"
	return s.deleteByPattern(ctx, pattern)
}

// =========================================================================
// 批量失效操作
// =========================================================================

// DeleteByCategoryKey 删除所有用户和管理员的指定 category Schema 缓存。
func (s *schemaCacheService) DeleteByCategoryKey(ctx context.Context, categoryKey string) error {
	if categoryKey == "" {
		categoryKey = schemaCategoryAll
	}

	// 1. 删除管理员的指定 category 缓存
	adminKey := s.buildAdminKey(categoryKey)
	if err := s.client.Del(ctx, adminKey).Err(); err != nil {
		slog.Warn("failed to delete admin schema cache", "categoryKey", categoryKey, "err", err)
	}

	// 2. 删除所有用户的指定 category 缓存
	pattern := s.keyPrefix + schemaKeyPrefix + schemaUserPrefix + "*:" + categoryKey
	if err := s.deleteByPattern(ctx, pattern); err != nil {
		return fmt.Errorf("failed to delete user schema caches: %w", err)
	}

	// 3. 同时删除所有用户的 _all 缓存（因为 _all 包含该 category）
	if categoryKey != schemaCategoryAll {
		allPattern := s.keyPrefix + schemaKeyPrefix + schemaUserPrefix + "*:" + schemaCategoryAll
		if err := s.deleteByPattern(ctx, allPattern); err != nil {
			slog.Warn("failed to delete user _all schema caches", "err", err)
		}

		// 同时删除管理员的 _all 缓存
		adminAllKey := s.buildAdminKey(schemaCategoryAll)
		if err := s.client.Del(ctx, adminAllKey).Err(); err != nil {
			slog.Warn("failed to delete admin _all schema cache", "err", err)
		}
	}

	return nil
}

// DeleteAll 删除所有 Schema 缓存。
func (s *schemaCacheService) DeleteAll(ctx context.Context) error {
	pattern := s.keyPrefix + schemaKeyPrefix + "*"
	return s.deleteByPattern(ctx, pattern)
}

// =========================================================================
// 分类列表缓存操作
// =========================================================================

// GetUserCategories 获取用户分类列表缓存。
func (s *schemaCacheService) GetUserCategories(ctx context.Context) ([]setting.CategoryMetaDTO, error) {
	key := s.buildCategoriesKey("user")
	return s.getCategories(ctx, key)
}

// SetUserCategories 设置用户分类列表缓存。
func (s *schemaCacheService) SetUserCategories(ctx context.Context, categories []setting.CategoryMetaDTO) error {
	key := s.buildCategoriesKey("user")
	return s.setCategories(ctx, key, categories)
}

// DeleteUserCategories 删除用户分类列表缓存。
func (s *schemaCacheService) DeleteUserCategories(ctx context.Context) error {
	key := s.buildCategoriesKey("user")
	return s.client.Del(ctx, key).Err()
}

// =========================================================================
// 内部辅助方法
// =========================================================================

// buildUserKey 构建用户 Schema 缓存 key。
func (s *schemaCacheService) buildUserKey(userID uint, categoryKey string) string {
	if categoryKey == "" {
		categoryKey = schemaCategoryAll
	}
	return s.keyPrefix + schemaKeyPrefix + schemaUserPrefix + strconv.FormatUint(uint64(userID), 10) + ":" + categoryKey
}

// buildAdminKey 构建管理员 Schema 缓存 key。
func (s *schemaCacheService) buildAdminKey(categoryKey string) string {
	if categoryKey == "" {
		categoryKey = schemaCategoryAll
	}
	return s.keyPrefix + schemaKeyPrefix + schemaAdminPrefix + categoryKey
}

// buildCategoriesKey 构建分类列表缓存 key。
func (s *schemaCacheService) buildCategoriesKey(scope string) string {
	return s.keyPrefix + schemaKeyPrefix + schemaCategoriesPrefix + scope
}

// get 通用获取缓存方法（使用 RedisJSON）。
func (s *schemaCacheService) get(ctx context.Context, key string) ([]setting.SchemaCategoryDTO, error) {
	// 使用 JSON.GET 命令读取
	data, err := s.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	// 需要解包外层数组
	var wrapper [][]setting.SchemaCategoryDTO
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, key)
		slog.Warn("corrupted schema cache, deleted", "key", key, "err", err)
		return nil, nil // corrupted cache treated as miss
	}

	if len(wrapper) == 0 {
		return nil, nil // empty wrapper treated as miss
	}

	return wrapper[0], nil
}

// set 通用设置缓存方法（使用 RedisJSON）。
func (s *schemaCacheService) set(ctx context.Context, key string, schema []setting.SchemaCategoryDTO) error {
	// 使用 Pipeline 执行 JSON.SET + EXPIRE
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", schema)
	pipe.Expire(ctx, key, schemaCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set schema cache: %w", err)
	}

	return nil
}

// getCategories 获取分类列表缓存（使用 RedisJSON）。
func (s *schemaCacheService) getCategories(ctx context.Context, key string) ([]setting.CategoryMetaDTO, error) {
	data, err := s.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	var wrapper [][]setting.CategoryMetaDTO
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		_ = s.client.Del(ctx, key)
		slog.Warn("corrupted categories cache, deleted", "key", key, "err", err)
		return nil, nil
	}

	if len(wrapper) == 0 {
		return nil, nil
	}

	return wrapper[0], nil
}

// setCategories 设置分类列表缓存（使用 RedisJSON）。
func (s *schemaCacheService) setCategories(ctx context.Context, key string, categories []setting.CategoryMetaDTO) error {
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", categories)
	pipe.Expire(ctx, key, schemaCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set categories cache: %w", err)
	}

	return nil
}

// deleteByPattern 使用 SCAN 批量删除匹配模式的 key。
func (s *schemaCacheService) deleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var allKeys []string

	for {
		keys, nextCursor, err := s.client.Scan(ctx, cursor, pattern, scanBatchSize).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		allKeys = append(allKeys, keys...)
		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	if len(allKeys) == 0 {
		return nil
	}

	// 使用 Pipeline 批量删除
	pipe := s.client.Pipeline()
	for i := 0; i < len(allKeys); i += deletePipelineSize {
		end := min(i+deletePipelineSize, len(allKeys))
		pipe.Del(ctx, allKeys[i:end]...)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute delete pipeline: %w", err)
	}

	slog.Debug("deleted schema caches", "pattern", pattern, "count", len(allKeys))
	return nil
}

var _ setting.SchemaCacheService = (*schemaCacheService)(nil)
