package persistence

import (
	"encoding/json"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"gorm.io/datatypes"
)

// SettingModel 配置定义的 GORM 实体
//
// 索引设计：
//   - idx_settings_category_sort: 复合索引 (category_id, group, order, key) 覆盖分类查询和排序
//   - idx_settings_scope: 单列索引用于 scope 过滤
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type SettingModel struct {
	ID           uint           `gorm:"primaryKey"`
	Key          string         `gorm:"uniqueIndex;size:100;not null"`
	DefaultValue datatypes.JSON `gorm:"type:jsonb;not null;default:'null'"` // JSONB 原生值
	Scope        string         `gorm:"size:20;not null;default:'user';index:idx_settings_scope"`

	// 复合索引：覆盖 FindByCategoryID 的 WHERE + ORDER BY
	CategoryID uint   `gorm:"not null;index:idx_settings_category_sort,priority:1"`
	Group      string `gorm:"size:50;default:'';index:idx_settings_category_sort,priority:2"`
	Order      int    `gorm:"default:0;index:idx_settings_category_sort,priority:3"`

	ValueType string `gorm:"size:20;default:'string'"`
	Label     string `gorm:"size:200"`

	// UI 配置
	InputType  string         `gorm:"column:input_type;size:32;not null;default:'text'"` // 控件类型
	Validation string         `gorm:"column:validation;type:text"`                       // JSON Logic 规则
	UIConfig   datatypes.JSON `gorm:"type:jsonb;default:'{}'"`                           // hint/options/depends_on

	// 权限控制
	ViewPermission string `gorm:"size:100;not null;default:'*:settings:read'"`
	EditPermission string `gorm:"size:100;not null;default:'admin:settings:update'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定配置定义表名
func (SettingModel) TableName() string {
	return "settings"
}

func newSettingModelFromEntity(entity *setting.Setting) *SettingModel {
	if entity == nil {
		return nil
	}

	// 将 any 类型的 DefaultValue 序列化为 JSON
	defaultValueJSON, _ := json.Marshal(entity.DefaultValue) //nolint:errchkjson // DefaultValue 是任意 JSONB 值

	return &SettingModel{
		ID:             entity.ID,
		Key:            entity.Key,
		DefaultValue:   datatypes.JSON(defaultValueJSON),
		Scope:          entity.Scope,
		CategoryID:     entity.CategoryID,
		Group:          entity.Group,
		ValueType:      entity.ValueType,
		Label:          entity.Label,
		Order:          entity.Order,
		InputType:      entity.InputType,
		Validation:     entity.Validation,
		UIConfig:       datatypes.JSON(entity.UIConfig),
		ViewPermission: entity.ViewPermission,
		EditPermission: entity.EditPermission,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *SettingModel) ToEntity() *setting.Setting {
	if m == nil {
		return nil
	}

	// 将 JSON 反序列化为 any 类型
	var defaultValue any
	_ = json.Unmarshal(m.DefaultValue, &defaultValue)

	return &setting.Setting{
		ID:             m.ID,
		Key:            m.Key,
		DefaultValue:   defaultValue,
		Scope:          m.Scope,
		CategoryID:     m.CategoryID,
		Group:          m.Group,
		ValueType:      m.ValueType,
		Label:          m.Label,
		Order:          m.Order,
		InputType:      m.InputType,
		Validation:     m.Validation,
		UIConfig:       string(m.UIConfig),
		ViewPermission: m.ViewPermission,
		EditPermission: m.EditPermission,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func mapSettingModelsToEntities(models []SettingModel) []*setting.Setting {
	if len(models) == 0 {
		return nil
	}

	defs := make([]*setting.Setting, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			defs = append(defs, entity)
		}
	}

	return defs
}
