package setting

// SelectOption 下拉/单选/多选选项。
type SelectOption struct {
	Value string `json:"value"`          // 选项值
	Label string `json:"label"`          // 显示文本
	Icon  string `json:"icon,omitempty"` // 可选图标（mdi-xxx）
}

// OptionsConfig 选项配置。
// 用于 select, radio, checkbox 等控件的选项列表。
//
// JSON 示例:
//
//	{
//	  "items": [
//	    {"value": "Asia/Shanghai", "label": "上海 (UTC+8)"},
//	    {"value": "UTC", "label": "UTC"}
//	  ]
//	}
type OptionsConfig struct {
	Items []SelectOption `json:"items"`
}

// DependsOnConfig 依赖关系配置。
// 用于控制设置项的可用状态，当依赖的设置项满足条件时才可编辑。
//
// JSON 示例:
//
//	{"key": "notification.enable_notifications", "value": true}
//
// 语义：当 notification.enable_notifications == true 时，本设置项才可编辑。
type DependsOnConfig struct {
	Key      string `json:"key"`                // 依赖的设置项 key
	Value    any    `json:"value,omitempty"`    // 期望的值
	Operator string `json:"operator,omitempty"` // 比较操作符: eq, ne, gt, lt（默认 eq）
}

// CategoryMeta Category 元数据。
// 用于定义 Category 的显示属性。
type CategoryMeta struct {
	Label string `json:"label"` // 显示名称
	Icon  string `json:"icon"`  // Tab 图标（mdi-xxx）
	Order int    `json:"order"` // 排序权重
}

// DefaultCategoryMetas 返回默认的 Category 元数据映射。
func DefaultCategoryMetas() map[string]CategoryMeta {
	return map[string]CategoryMeta{
		CategoryGeneral: {
			Label: "常规设置",
			Icon:  "mdi-cog",
			Order: 1,
		},
		CategorySecurity: {
			Label: "安全设置",
			Icon:  "mdi-shield-lock",
			Order: 2,
		},
		CategoryNotification: {
			Label: "通知设置",
			Icon:  "mdi-bell",
			Order: 3,
		},
		CategoryBackup: {
			Label: "备份设置",
			Icon:  "mdi-backup-restore",
			Order: 4,
		},
	}
}
