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
	DefaultValue datatypes.JSON `gorm:"type:jsonb;not null;default:'null'"` // JSONB 原生值
	Category     string         `gorm:"size:50;index;not null"`
	Group        string         `gorm:"size:50;index;default:''"`
	ValueType    string         `gorm:"size:20;default:'string'"`
	Label        string         `gorm:"size:200"`
	UIConfig     datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	Order        int            `gorm:"default:0;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定配置定义表名
func (SettingModel) TableName() string {
	return "setting_definitions"
}

func newSettingModelFromEntity(entity *setting.Setting) *SettingModel {
	if entity == nil {
		return nil
	}

	// 将 any 类型的 DefaultValue 序列化为 JSON
	defaultValueJSON, _ := json.Marshal(entity.DefaultValue) //nolint:errchkjson // DefaultValue 是任意 JSONB 值

	return &SettingModel{
		ID:           entity.ID,
		Key:          entity.Key,
		DefaultValue: datatypes.JSON(defaultValueJSON),
		Category:     entity.Category,
		Group:        entity.Group,
		ValueType:    entity.ValueType,
		Label:        entity.Label,
		UIConfig:     datatypes.JSON(entity.UIConfig),
		Order:        entity.Order,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
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
		ID:           m.ID,
		Key:          m.Key,
		DefaultValue: defaultValue,
		Category:     m.Category,
		Group:        m.Group,
		ValueType:    m.ValueType,
		Label:        m.Label,
		UIConfig:     string(m.UIConfig),
		Order:        m.Order,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
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
