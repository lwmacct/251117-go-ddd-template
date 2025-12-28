package setting

import (
	"context"
	"fmt"
	"sort"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// ListSchemaHandler 获取设置 Schema 查询处理器
type ListSchemaHandler struct {
	settingQueryRepo  setting.QueryRepository
	categoryQueryRepo setting.SettingCategoryQueryRepository
}

// NewListSchemaHandler 创建 ListSchemaHandler 实例
func NewListSchemaHandler(
	settingQueryRepo setting.QueryRepository,
	categoryQueryRepo setting.SettingCategoryQueryRepository,
) *ListSchemaHandler {
	return &ListSchemaHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
	}
}

// Handle 处理获取设置 Schema 查询
// 返回按 Category → Group → Settings 层级组织的精简数据
//
// 支持 CategoryKey 过滤：
//   - 为空时返回全量系统设置（用于总配置页）
//   - 指定 Key 时只返回该分类（用于分散页面的懒加载）
func (h *ListSchemaHandler) Handle(ctx context.Context, query ListSchemaQuery) ([]SchemaCategoryDTO, error) {
	// 1. 根据 CategoryKey 决定查询范围
	settings, err := h.fetchSettings(ctx, query.CategoryKey)
	if err != nil {
		return nil, err
	}

	// 2. 查询所有分类元数据
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch setting categories: %w", err)
	}

	return h.buildSchema(settings, categories), nil
}

// fetchSettings 根据 CategoryKey 获取设置列表
func (h *ListSchemaHandler) fetchSettings(ctx context.Context, categoryKey string) ([]*setting.Setting, error) {
	// 全量查询
	if categoryKey == "" {
		settings, err := h.settingQueryRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch settings: %w", err)
		}
		return settings, nil
	}

	// 按 Category Key 过滤
	category, err := h.categoryQueryRepo.FindByKey(ctx, categoryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find category by key %q: %w", categoryKey, err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found: %s", categoryKey)
	}

	settings, err := h.settingQueryRepo.FindByCategoryID(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch settings: %w", err)
	}
	return settings, nil
}

// buildSchema 构建 Category → Group → Settings 层级结构
//
//nolint:dupl // 与 UserListSchemaHandler.buildSchema 结构相似但使用不同 DTO 类型
func (h *ListSchemaHandler) buildSchema(
	settings []*setting.Setting,
	categoryEntities []*setting.SettingCategory,
) []SchemaCategoryDTO {
	// 构建 CategoryID → Category 实体映射
	categoryByID := make(map[uint]*setting.SettingCategory, len(categoryEntities))
	for _, cat := range categoryEntities {
		categoryByID[cat.ID] = cat
	}

	// 按 CategoryID 分组
	categoryMap := make(map[uint]map[string][]SchemaSettingDTO)

	for _, s := range settings {
		categoryID := s.CategoryID
		group := s.Group
		if group == "" {
			group = "default"
		}

		if _, ok := categoryMap[categoryID]; !ok {
			categoryMap[categoryID] = make(map[string][]SchemaSettingDTO)
		}

		dto := ToSchemaSettingDTO(s)
		if dto != nil {
			categoryMap[categoryID][group] = append(categoryMap[categoryID][group], *dto)
		}
	}

	// 构建响应
	result := make([]SchemaCategoryDTO, 0, len(categoryMap))
	for categoryID, groupMap := range categoryMap {
		cat, ok := categoryByID[categoryID]
		if !ok {
			// 跳过未知 category
			continue
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

		result = append(result, SchemaCategoryDTO{
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
func (h *ListSchemaHandler) getCategoryIDByKey(categoryByID map[uint]*setting.SettingCategory, key string) uint {
	for id, cat := range categoryByID {
		if cat.Key == key {
			return id
		}
	}
	return 0
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
