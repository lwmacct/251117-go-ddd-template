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
	schemaCacheTTL     = 30 * time.Minute
	schemaKeyPrefix    = "schema:"
	schemaUserPrefix   = "user:"
	schemaAdminPrefix  = "admin:"
	schemaCategoryAll  = "_all"
	scanBatchSize      = 100
	deletePipelineSize = 100
)

// =========================================================================
// 内部缓存 DTO（序列化用）
// =========================================================================

// schemaCategoryDTO Schema 分类缓存 DTO。
type schemaCategoryDTO struct {
	Category string           `json:"category"`
	Label    string           `json:"label"`
	Icon     string           `json:"icon"`
	Groups   []schemaGroupDTO `json:"groups"`
}

// schemaGroupDTO Schema 分组缓存 DTO。
type schemaGroupDTO struct {
	Group    string             `json:"group"`
	Label    string             `json:"label"`
	Settings []schemaSettingDTO `json:"settings"`
}

// schemaSettingDTO Schema 设置项缓存 DTO。
type schemaSettingDTO struct {
	Key            string         `json:"key"`
	Value          any            `json:"value"`
	DefaultValue   any            `json:"default_value"`
	IsCustomized   bool           `json:"is_customized"`
	Scope          string         `json:"scope,omitempty"`
	ValueType      string         `json:"value_type"`
	Label          string         `json:"label"`
	UIConfig       schemaUIConfig `json:"ui_config"`
	Order          int            `json:"order"`
	ViewPermission string         `json:"view_permission,omitempty"`
	EditPermission string         `json:"edit_permission,omitempty"`
}

// schemaUIConfig UI 配置缓存 DTO。
type schemaUIConfig struct {
	InputType  string               `json:"input_type"`
	Hint       string               `json:"hint,omitempty"`
	Options    []schemaSelectOption `json:"options,omitempty"`
	Validation any                  `json:"validation,omitempty"`
	DependsOn  *schemaDependsOn     `json:"depends_on,omitempty"`
}

// schemaSelectOption 下拉选项缓存 DTO。
type schemaSelectOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// schemaDependsOn 依赖关系缓存 DTO。
type schemaDependsOn struct {
	Key      string `json:"key"`
	Value    any    `json:"value,omitempty"`
	Operator string `json:"operator,omitempty"`
}

// =========================================================================
// DTO 转换函数
// =========================================================================

// toSchemaCategoryDTOs 将 Application DTO 转换为缓存 DTO。
func toSchemaCategoryDTOs(dtos []setting.SchemaCategoryDTO) []schemaCategoryDTO {
	result := make([]schemaCategoryDTO, 0, len(dtos))
	for _, dto := range dtos {
		result = append(result, toSchemaCategoryDTO(dto))
	}
	return result
}

func toSchemaCategoryDTO(dto setting.SchemaCategoryDTO) schemaCategoryDTO {
	groups := make([]schemaGroupDTO, 0, len(dto.Groups))
	for _, g := range dto.Groups {
		groups = append(groups, toSchemaGroupDTO(g))
	}
	return schemaCategoryDTO{
		Category: dto.Category,
		Label:    dto.Label,
		Icon:     dto.Icon,
		Groups:   groups,
	}
}

func toSchemaGroupDTO(dto setting.SchemaGroupDTO) schemaGroupDTO {
	settings := make([]schemaSettingDTO, 0, len(dto.Settings))
	for _, s := range dto.Settings {
		settings = append(settings, toSchemaSettingDTO(s))
	}
	return schemaGroupDTO{
		Group:    dto.Group,
		Label:    dto.Label,
		Settings: settings,
	}
}

func toSchemaSettingDTO(dto setting.SchemaSettingDTO) schemaSettingDTO {
	return schemaSettingDTO{
		Key:            dto.Key,
		Value:          dto.Value,
		DefaultValue:   dto.DefaultValue,
		IsCustomized:   dto.IsCustomized,
		Scope:          dto.Scope,
		ValueType:      dto.ValueType,
		Label:          dto.Label,
		UIConfig:       toSchemaUIConfig(dto.UIConfig),
		Order:          dto.Order,
		ViewPermission: dto.ViewPermission,
		EditPermission: dto.EditPermission,
	}
}

func toSchemaUIConfig(dto setting.UIConfigDTO) schemaUIConfig {
	options := make([]schemaSelectOption, 0, len(dto.Options))
	for _, o := range dto.Options {
		options = append(options, schemaSelectOption{Value: o.Value, Label: o.Label})
	}
	var dependsOn *schemaDependsOn
	if dto.DependsOn != nil {
		dependsOn = &schemaDependsOn{
			Key:      dto.DependsOn.Key,
			Value:    dto.DependsOn.Value,
			Operator: dto.DependsOn.Operator,
		}
	}
	return schemaUIConfig{
		InputType:  dto.InputType,
		Hint:       dto.Hint,
		Options:    options,
		Validation: dto.Validation,
		DependsOn:  dependsOn,
	}
}

// fromSchemaCategoryDTOs 将缓存 DTO 转换为 Application DTO。
func fromSchemaCategoryDTOs(dtos []schemaCategoryDTO) []setting.SchemaCategoryDTO {
	result := make([]setting.SchemaCategoryDTO, 0, len(dtos))
	for _, dto := range dtos {
		result = append(result, fromSchemaCategoryDTO(dto))
	}
	return result
}

func fromSchemaCategoryDTO(dto schemaCategoryDTO) setting.SchemaCategoryDTO {
	groups := make([]setting.SchemaGroupDTO, 0, len(dto.Groups))
	for _, g := range dto.Groups {
		groups = append(groups, fromSchemaGroupDTO(g))
	}
	return setting.SchemaCategoryDTO{
		Category: dto.Category,
		Label:    dto.Label,
		Icon:     dto.Icon,
		Groups:   groups,
	}
}

func fromSchemaGroupDTO(dto schemaGroupDTO) setting.SchemaGroupDTO {
	settings := make([]setting.SchemaSettingDTO, 0, len(dto.Settings))
	for _, s := range dto.Settings {
		settings = append(settings, fromSchemaSettingDTO(s))
	}
	return setting.SchemaGroupDTO{
		Group:    dto.Group,
		Label:    dto.Label,
		Settings: settings,
	}
}

func fromSchemaSettingDTO(dto schemaSettingDTO) setting.SchemaSettingDTO {
	return setting.SchemaSettingDTO{
		Key:            dto.Key,
		Value:          dto.Value,
		DefaultValue:   dto.DefaultValue,
		IsCustomized:   dto.IsCustomized,
		Scope:          dto.Scope,
		ValueType:      dto.ValueType,
		Label:          dto.Label,
		UIConfig:       fromSchemaUIConfig(dto.UIConfig),
		Order:          dto.Order,
		ViewPermission: dto.ViewPermission,
		EditPermission: dto.EditPermission,
	}
}

func fromSchemaUIConfig(dto schemaUIConfig) setting.UIConfigDTO {
	options := make([]setting.SelectOptionDTO, 0, len(dto.Options))
	for _, o := range dto.Options {
		options = append(options, setting.SelectOptionDTO{Value: o.Value, Label: o.Label})
	}
	var dependsOn *setting.DependsOnConfigDTO
	if dto.DependsOn != nil {
		dependsOn = &setting.DependsOnConfigDTO{
			Key:      dto.DependsOn.Key,
			Value:    dto.DependsOn.Value,
			Operator: dto.DependsOn.Operator,
		}
	}
	return setting.UIConfigDTO{
		InputType:  dto.InputType,
		Hint:       dto.Hint,
		Options:    options,
		Validation: dto.Validation,
		DependsOn:  dependsOn,
	}
}

// =========================================================================
// Service 实现
// =========================================================================

// schemaCacheService Schema 缓存服务的 Redis 实现。
//
// 缓存 Setting Schema API 的最终响应：
//   - Key 格式：{prefix}schema:user:{userID}:{categoryKey}
//   - Key 格式：{prefix}schema:admin:{categoryKey}
//   - TTL：30 分钟
//   - JSON 序列化存储
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

// get 通用获取缓存方法。
func (s *schemaCacheService) get(ctx context.Context, key string) ([]setting.SchemaCategoryDTO, error) {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var result []schemaCategoryDTO
	if err := json.Unmarshal(data, &result); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, key)
		slog.Warn("corrupted schema cache, deleted", "key", key, "err", err)
		return nil, nil
	}

	return fromSchemaCategoryDTOs(result), nil
}

// set 通用设置缓存方法。
func (s *schemaCacheService) set(ctx context.Context, key string, schema []setting.SchemaCategoryDTO) error {
	dtos := toSchemaCategoryDTOs(schema)
	data, err := json.Marshal(dtos)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	return s.client.Set(ctx, key, data, schemaCacheTTL).Err()
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
