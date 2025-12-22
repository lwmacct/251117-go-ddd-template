package persistence

import (
	"time"

	domainSetting "github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// SettingModel 配置的 GORM 实体
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type SettingModel struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"uniqueIndex;size:100;not null"`
	Value     string `gorm:"type:text"`
	Category  string `gorm:"size:50;index;not null"`
	ValueType string `gorm:"size:20;default:'string'"`
	Label     string `gorm:"size:200"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定配置表名
func (SettingModel) TableName() string {
	return "settings"
}

func newSettingModelFromEntity(entity *domainSetting.Setting) *SettingModel {
	if entity == nil {
		return nil
	}

	return &SettingModel{
		ID:        entity.ID,
		Key:       entity.Key,
		Value:     entity.Value,
		Category:  entity.Category,
		ValueType: entity.ValueType,
		Label:     entity.Label,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity（实现 Model[E] 接口）
func (m *SettingModel) ToEntity() *domainSetting.Setting {
	if m == nil {
		return nil
	}

	return &domainSetting.Setting{
		ID:        m.ID,
		Key:       m.Key,
		Value:     m.Value,
		Category:  m.Category,
		ValueType: m.ValueType,
		Label:     m.Label,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func mapSettingModelsToEntities(models []SettingModel) []*domainSetting.Setting {
	if len(models) == 0 {
		return nil
	}

	settings := make([]*domainSetting.Setting, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			settings = append(settings, entity)
		}
	}

	return settings
}
