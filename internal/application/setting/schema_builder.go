package setting

import (
	"sort"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// SchemaBuilder 构建 Category → Group → Settings 层级结构
// 用于 ListSchemaHandler 和 UserListSchemaHandler 共享构建逻辑
type SchemaBuilder struct {
	categoryByID map[uint]*setting.SettingCategory
}

// NewSchemaBuilder 创建 Schema 构建器
func NewSchemaBuilder(categories []*setting.SettingCategory) *SchemaBuilder {
	categoryByID := make(map[uint]*setting.SettingCategory, len(categories))
	for _, cat := range categories {
		categoryByID[cat.ID] = cat
	}
	return &SchemaBuilder{
		categoryByID: categoryByID,
	}
}

// SettingMapper 将 Setting 转换为 SchemaSettingDTO 的函数类型
// Admin 场景使用 ToSchemaSettingDTO，User 场景使用 ToUserSchemaSettingDTO
type SettingMapper func(s *setting.Setting, us *setting.UserSetting) *SchemaSettingDTO

// Build 构建 Schema 层级结构
// settings: 配置定义列表
// userSettingMap: 用户配置映射（Admin 场景传 nil）
// mapper: 转换函数
func (b *SchemaBuilder) Build(
	settings []*setting.Setting,
	userSettingMap map[string]*setting.UserSetting,
	mapper SettingMapper,
) []SchemaCategoryDTO {
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

		var us *setting.UserSetting
		if userSettingMap != nil {
			us = userSettingMap[s.Key]
		}
		dto := mapper(s, us)
		if dto != nil {
			categoryMap[categoryID][group] = append(categoryMap[categoryID][group], *dto)
		}
	}

	// 构建响应
	result := make([]SchemaCategoryDTO, 0, len(categoryMap))
	for categoryID, groupMap := range categoryMap {
		cat, ok := b.categoryByID[categoryID]
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
				Label:    GetGroupLabel(group),
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
		catI := b.categoryByID[b.getCategoryIDByKey(result[i].Category)]
		catJ := b.categoryByID[b.getCategoryIDByKey(result[j].Category)]
		if catI == nil || catJ == nil {
			return result[i].Category < result[j].Category
		}
		return catI.Order < catJ.Order
	})

	return result
}

// getCategoryIDByKey 根据 key 查找 CategoryID（用于排序）
func (b *SchemaBuilder) getCategoryIDByKey(key string) uint {
	for id, cat := range b.categoryByID {
		if cat.Key == key {
			return id
		}
	}
	return 0
}

// GetGroupLabel 获取分组显示名称
func GetGroupLabel(group string) string {
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

// AdminSettingMapper Admin 场景的 Setting 转换器（包含全部字段）
func AdminSettingMapper(s *setting.Setting, _ *setting.UserSetting) *SchemaSettingDTO {
	return ToSchemaSettingDTO(s)
}

// UserSettingMapper User 场景的 Setting 转换器（省略权限字段，合并用户值）
func UserSettingMapper(s *setting.Setting, us *setting.UserSetting) *SchemaSettingDTO {
	return ToUserSchemaSettingDTO(s, us)
}
