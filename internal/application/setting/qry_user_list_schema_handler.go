package setting

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// UserListSchemaHandler 获取用户配置 Schema 查询处理器
type UserListSchemaHandler struct {
	settingQueryRepo  setting.QueryRepository
	queryRepo         setting.UserSettingQueryRepository
	categoryQueryRepo setting.SettingCategoryQueryRepository
	schemaCache       SchemaCacheService
}

// NewUserListSchemaHandler 创建 UserListSchemaHandler 实例
func NewUserListSchemaHandler(
	settingQueryRepo setting.QueryRepository,
	queryRepo setting.UserSettingQueryRepository,
	categoryQueryRepo setting.SettingCategoryQueryRepository,
	schemaCache SchemaCacheService,
) *UserListSchemaHandler {
	return &UserListSchemaHandler{
		settingQueryRepo:  settingQueryRepo,
		queryRepo:         queryRepo,
		categoryQueryRepo: categoryQueryRepo,
		schemaCache:       schemaCache,
	}
}

// Handle 处理获取用户配置 Schema 查询
// 返回按 Category → Group → Settings 层级组织的数据，包含用户自定义值
//
// 仅返回 scope="user" 的设置项（用户可配置项），排除系统级设置
//
// 支持 CategoryKey 过滤：
//   - 为空时返回全量用户设置（用于总配置页）
//   - 指定 Key 时只返回该分类（用于分散页面的懒加载）
//
// 缓存策略：
//   - 先查缓存，命中直接返回
//   - 未命中时查数据库，同步回写缓存
func (h *UserListSchemaHandler) Handle(ctx context.Context, query UserListSchemaQuery) ([]SchemaCategoryDTO, error) {
	// 1. 查缓存
	if cached, err := h.schemaCache.GetUserSchema(ctx, query.UserID, query.CategoryKey); err == nil && cached != nil {
		return cached, nil
	}

	// 2. 缓存未命中，执行原有逻辑
	result, err := h.fetchAndBuildSchema(ctx, query)
	if err != nil {
		return nil, err
	}

	// 3. 同步回写缓存
	if err := h.schemaCache.SetUserSchema(ctx, query.UserID, query.CategoryKey, result); err != nil {
		slog.Warn("failed to cache user schema", "userID", query.UserID, "categoryKey", query.CategoryKey, "err", err)
	}

	return result, nil
}

// fetchAndBuildSchema 从数据库获取数据并构建 Schema
func (h *UserListSchemaHandler) fetchAndBuildSchema(ctx context.Context, query UserListSchemaQuery) ([]SchemaCategoryDTO, error) {
	// 1. 根据 CategoryKey 决定查询范围（只查询 user scope）
	defs, err := h.fetchUserSettings(ctx, query.CategoryKey)
	if err != nil {
		return nil, err
	}

	// 2. 查找用户的所有自定义配置
	userSettings, err := h.queryRepo.FindByUser(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user settings: %w", err)
	}

	// 3. 查询所有分类元数据
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch setting categories: %w", err)
	}

	// 4. 构建用户配置映射
	userSettingMap := make(map[string]*setting.UserSetting, len(userSettings))
	for _, us := range userSettings {
		userSettingMap[us.SettingKey] = us
	}

	// 5. 使用共享构建器
	builder := NewSchemaBuilder(categories)
	return builder.Build(defs, userSettingMap, UserSettingMapper), nil
}

// fetchUserSettings 根据 CategoryKey 获取用户可配置的设置列表
func (h *UserListSchemaHandler) fetchUserSettings(ctx context.Context, categoryKey string) ([]*setting.Setting, error) {
	// 全量查询用户可配置的设置
	if categoryKey == "" {
		defs, err := h.settingQueryRepo.FindByScope(ctx, "user")
		if err != nil {
			return nil, fmt.Errorf("failed to find setting definitions: %w", err)
		}
		return defs, nil
	}

	// 按 Category Key 过滤
	category, err := h.categoryQueryRepo.FindByKey(ctx, categoryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find category by key %q: %w", categoryKey, err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found: %s", categoryKey)
	}

	allDefs, err := h.settingQueryRepo.FindByCategoryID(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definitions: %w", err)
	}

	// 过滤只保留 user scope
	return filterUserScopeSettings(allDefs), nil
}

// filterUserScopeSettings 过滤只保留 user scope 的设置
func filterUserScopeSettings(settings []*setting.Setting) []*setting.Setting {
	result := make([]*setting.Setting, 0, len(settings))
	for _, s := range settings {
		if s.IsUserScope() {
			result = append(result, s)
		}
	}
	return result
}
