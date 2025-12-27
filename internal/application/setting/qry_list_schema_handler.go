package setting

import (
	"context"
	"fmt"
	"sort"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// ListSchemaHandler 获取设置 Schema 查询处理器
type ListSchemaHandler struct {
	settingQueryRepo setting.QueryRepository
}

// NewListSchemaHandler 创建 ListSchemaHandler 实例
func NewListSchemaHandler(settingQueryRepo setting.QueryRepository) *ListSchemaHandler {
	return &ListSchemaHandler{
		settingQueryRepo: settingQueryRepo,
	}
}

// Handle 处理获取设置 Schema 查询
// 返回按 Category → Group → Settings 层级组织的精简数据
func (h *ListSchemaHandler) Handle(ctx context.Context, _ ListSchemaQuery) ([]SchemaCategoryDTO, error) {
	settings, err := h.settingQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch settings: %w", err)
	}

	return h.buildSchema(settings), nil
}

// buildSchema 构建 Category → Group → Settings 层级结构
func (h *ListSchemaHandler) buildSchema(settings []*setting.Setting) []SchemaCategoryDTO {
	// 按 Category 分组
	categoryMap := make(map[string]map[string][]SchemaSettingDTO)

	for _, s := range settings {
		category := s.Category
		group := s.Group
		if group == "" {
			group = "default"
		}

		if _, ok := categoryMap[category]; !ok {
			categoryMap[category] = make(map[string][]SchemaSettingDTO)
		}

		dto := ToSchemaSettingDTO(s)
		if dto != nil {
			categoryMap[category][group] = append(categoryMap[category][group], *dto)
		}
	}

	// 获取 Category 元数据
	categoryMetas := setting.DefaultCategoryMetas()

	// 构建响应
	categories := make([]SchemaCategoryDTO, 0, len(categoryMap))
	for category, groupMap := range categoryMap {
		meta, ok := categoryMetas[category]
		if !ok {
			meta = setting.CategoryMeta{
				Label: category,
				Icon:  "mdi-cog",
				Order: 99,
			}
			// 将未知 category 的元数据添加到 map 中，用于后续排序
			categoryMetas[category] = meta
		}

		groups := make([]SchemaGroupDTO, 0, len(groupMap))
		for group, settingDTOs := range groupMap {
			// 按 Order 排序设置项
			sort.Slice(settingDTOs, func(i, j int) bool {
				return settingDTOs[i].Order < settingDTOs[j].Order
			})

			groups = append(groups, SchemaGroupDTO{
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

		categories = append(categories, SchemaCategoryDTO{
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
func (h *ListSchemaHandler) getGroupLabel(group string) string {
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
