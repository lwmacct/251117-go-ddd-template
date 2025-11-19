package persistence

import (
	"time"

	domainTwoFA "github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	"gorm.io/gorm"
)

// TwoFAModel 2FA 的 GORM 实体
type TwoFAModel struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID uint `gorm:"uniqueIndex;not null"`
	Enabled bool `gorm:"default:false;not null"`
	Secret  string `gorm:"size:255;not null"`
	RecoveryCodes domainTwoFA.RecoveryCodes `gorm:"type:text"`
	SetupCompletedAt *time.Time
	LastUsedAt       *time.Time
}

// TableName 指定 2FA 表名
func (TwoFAModel) TableName() string {
	return "user_2fas"
}

func newTwoFAModelFromEntity(entity *domainTwoFA.TwoFA) *TwoFAModel {
	if entity == nil {
		return nil
	}

	model := &TwoFAModel{
		ID:              entity.ID,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
		UserID:          entity.UserID,
		Enabled:         entity.Enabled,
		Secret:          entity.Secret,
		RecoveryCodes:   entity.RecoveryCodes,
		SetupCompletedAt: entity.SetupCompletedAt,
		LastUsedAt:       entity.LastUsedAt,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

func (m *TwoFAModel) toEntity() *domainTwoFA.TwoFA {
	if m == nil {
		return nil
	}

	entity := &domainTwoFA.TwoFA{
		ID:              m.ID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		UserID:          m.UserID,
		Enabled:         m.Enabled,
		Secret:          m.Secret,
		RecoveryCodes:   m.RecoveryCodes,
		SetupCompletedAt: m.SetupCompletedAt,
		LastUsedAt:       m.LastUsedAt,
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}
