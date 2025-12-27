package setting

import (
	"context"
	"fmt"
	"sort"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// UserListSchemaHandler 获取用户配置 Schema 查询处理器
type UserListSchemaHandler struct {
	settingQueryRepo setting.QueryRepository
	queryRepo        setting.UserSettingQueryRepository
}

// NewUserListSchemaHandler 创建 UserListSchemaHandler 实例
func NewUserListSchemaHandler(
	settingQueryRepo setting.QueryRepository,
	queryRepo setting.UserSettingQueryRepository,
) *UserListSchemaHandler {
	return &UserListSchemaHandler{
		settingQueryRepo: settingQueryRepo,
		queryRepo:        queryRepo,
	}
}

// Handle 处理获取用户配置 Schema 查询
// 返回按 Category → Group → Settings 层级组织的数据，包含用户自定义值
//

func (h *UserListSchemaHandler) Handle(ctx context.Context, query UserListSchemaQuery) ([]UserSchemaCategoryDTO, error) {
	// 1. 查找所有配置定义
	defs, err := h.settingQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definitions: %w", err)
	}

	// 2. 查找用户的所有自定义配置
	userSettings, err := h.queryRepo.FindByUser(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user settings: %w", err)
	}

	// 3. 构建用户配置映射
	userSettingMap := make(map[string]*setting.UserSetting, len(userSettings))
	for _, us := range userSettings {
		userSettingMap[us.SettingKey] = us
	}

	// 4. 构建 Schema
	return h.buildSchema(defs, userSettingMap), nil
}

// buildSchema 构建 Category → Group → Settings 层级结构
//
//nolint:dupl // 与 ListSchemaHandler.buildSchema 结构相似但使用不同 DTO 类型
func (h *UserListSchemaHandler) buildSchema(
	defs []*setting.Setting,
	userSettingMap map[string]*setting.UserSetting,
) []UserSchemaCategoryDTO {
	// 按 Category 分组
	categoryMap := make(map[string]map[string][]UserSchemaSettingDTO)

	for _, def := range defs {
		category := def.Category
		group := def.Group
		if group == "" {
			group = "default"
		}

		if _, ok := categoryMap[category]; !ok {
			categoryMap[category] = make(map[string][]UserSchemaSettingDTO)
		}

		dto := ToUserSchemaSettingDTO(def, userSettingMap[def.Key])
		if dto != nil {
			categoryMap[category][group] = append(categoryMap[category][group], *dto)
		}
	}

	// 获取 Category 元数据
	categoryMetas := setting.DefaultCategoryMetas()

	// 构建响应
	categories := make([]UserSchemaCategoryDTO, 0, len(categoryMap))
	for category, groupMap := range categoryMap {
		meta, ok := categoryMetas[category]
		if !ok {
			meta = setting.CategoryMeta{
				Label: category,
				Icon:  "mdi-cog",
				Order: 99,
			}
			categoryMetas[category] = meta
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

		categories = append(categories, UserSchemaCategoryDTO{
			Category: category,
			Label:    meta.Label,
			Icon:     meta.Icon,
			Groups:   groups,
		})
	}

	// 按 Category Meta Order 排序
	sort.Slice(categories, func(i, j int) bool {
		metaI := categoryMetas[categories[i].Category]
		metaJ := categoryMetas[categories[j].Category]
		return metaI.Order < metaJ.Order
	})

	return categories
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
