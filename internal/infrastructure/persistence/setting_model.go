package persistence

import (
	"encoding/json"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"gorm.io/datatypes"
)

// SettingModel 配置定义的 GORM 实体
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type SettingModel struct {
	ID           uint           `gorm:"primaryKey"`
	Key          string         `gorm:"uniqueIndex;size:100;not null"`
	DefaultValue datatypes.JSON `gorm:"type:jsonb;not null;default:'null'"`    // JSONB 原生值
	Scope        string         `gorm:"size:20;not null;default:'user';index"` // system | user
	CategoryID   uint           `gorm:"index;not null"`                        // 逻辑关联 setting_categories.id（无物理 FK）
	Group        string         `gorm:"size:50;index;default:''"`
	ValueType    string         `gorm:"size:20;default:'string'"`
	Label        string         `gorm:"size:200"`
	UIConfig     datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	Order        int            `gorm:"default:0;index"`

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
		UIConfig:       datatypes.JSON(entity.UIConfig),
		Order:          entity.Order,
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
		UIConfig:       string(m.UIConfig),
		Order:          m.Order,
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
