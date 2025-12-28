package setting

import (
	"context"
	"fmt"
	"sort"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// UserListSchemaHandler 获取用户配置 Schema 查询处理器
type UserListSchemaHandler struct {
	settingQueryRepo  setting.QueryRepository
	queryRepo         setting.UserSettingQueryRepository
	categoryQueryRepo setting.SettingCategoryQueryRepository
}

// NewUserListSchemaHandler 创建 UserListSchemaHandler 实例
func NewUserListSchemaHandler(
	settingQueryRepo setting.QueryRepository,
	queryRepo setting.UserSettingQueryRepository,
	categoryQueryRepo setting.SettingCategoryQueryRepository,
) *UserListSchemaHandler {
	return &UserListSchemaHandler{
		settingQueryRepo:  settingQueryRepo,
		queryRepo:         queryRepo,
		categoryQueryRepo: categoryQueryRepo,
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
func (h *UserListSchemaHandler) Handle(ctx context.Context, query UserListSchemaQuery) ([]UserSchemaCategoryDTO, error) {
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

	// 5. 构建 Schema
	return h.buildSchema(defs, userSettingMap, categories), nil
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

// buildSchema 构建 Category → Group → Settings 层级结构
//
//nolint:dupl // 与 ListSchemaHandler.buildSchema 结构相似但使用不同 DTO 类型
func (h *UserListSchemaHandler) buildSchema(
	defs []*setting.Setting,
	userSettingMap map[string]*setting.UserSetting,
	categoryEntities []*setting.SettingCategory,
) []UserSchemaCategoryDTO {
	// 构建 CategoryID → Category 实体映射
	categoryByID := make(map[uint]*setting.SettingCategory, len(categoryEntities))
	for _, cat := range categoryEntities {
		categoryByID[cat.ID] = cat
	}

	// 按 CategoryID 分组
	categoryMap := make(map[uint]map[string][]UserSchemaSettingDTO)

	for _, def := range defs {
		categoryID := def.CategoryID
		group := def.Group
		if group == "" {
			group = "default"
		}

		if _, ok := categoryMap[categoryID]; !ok {
			categoryMap[categoryID] = make(map[string][]UserSchemaSettingDTO)
		}

		dto := ToUserSchemaSettingDTO(def, userSettingMap[def.Key])
		if dto != nil {
			categoryMap[categoryID][group] = append(categoryMap[categoryID][group], *dto)
		}
	}

	// 构建响应
	result := make([]UserSchemaCategoryDTO, 0, len(categoryMap))
	for categoryID, groupMap := range categoryMap {
		cat, ok := categoryByID[categoryID]
		if !ok {
			// 跳过未知 category
			continue
		}

		groups := make([]UserSchemaGroupDTO, 0, len(groupMap))
		for group, settingDTOs := range groupMap {
			// 按 Order 排序设置项
			sort.Slice(settingDTOs, func(i, j int) bool {
				return settingDTOs[i].Order < settingDTOs[j].Order
			})

			groups = append(groups, UserSchemaGroupDTO{
				Group:    group,
				Label:    h.getGroupLabel(group),
				Settings: settingDTOs,
			})
		}

		// 按 Group 名排序（default 在前）
		sort.Slice(groups, func(i, j int) bool {
			if groups[i].Group == "default" {
				return true
			}
			if groups[j].Group == "default" {
				return false
			}
			return groups[i].Group < groups[j].Group
		})

		result = append(result, UserSchemaCategoryDTO{
			Category: cat.Key,
			Label:    cat.Label,
			Icon:     cat.Icon,
			Groups:   groups,
		})
	}

	// 按 Category Order 排序
	sort.Slice(result, func(i, j int) bool {
		catI := categoryByID[h.getCategoryIDByKey(categoryByID, result[i].Category)]
		catJ := categoryByID[h.getCategoryIDByKey(categoryByID, result[j].Category)]
		if catI == nil || catJ == nil {
			return result[i].Category < result[j].Category
		}
		return catI.Order < catJ.Order
	})

	return result
}

// getCategoryIDByKey 根据 key 查找 CategoryID（用于排序）
func (h *UserListSchemaHandler) getCategoryIDByKey(categoryByID map[uint]*setting.SettingCategory, key string) uint {
	for id, cat := range categoryByID {
		if cat.Key == key {
			return id
		}
	}
	return 0
}

// getGroupLabel 获取分组显示名称
func (h *UserListSchemaHandler) getGroupLabel(group string) string {
	labels := map[string]string{
		"default":    "",
		"basic":      "基本设置",
		"locale":     "本地化",
		"appearance": "外观",
		"password":   "密码策略",
		"session":    "会话管理",
		"advanced":   "高级设置",
		"general":    "基本设置",
		"email":      "邮件通知",
		"sms":        "短信通知",
		"schedule":   "备份计划",
	}

	if label, ok := labels[group]; ok {
		return label
	}
	return group
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
