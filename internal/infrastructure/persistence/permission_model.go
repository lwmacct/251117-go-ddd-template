package persistence

import (
	"time"

	domainRole "github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// PermissionModel 定义权限的 GORM 持久化模型
type PermissionModel struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Domain      string         `gorm:"size:50;not null;index"`
	Resource    string         `gorm:"size:50;not null;index"`
	Action      string         `gorm:"size:50;not null;index"`
	Description string         `gorm:"size:255"`
	Code        string         `gorm:"size:150;uniqueIndex;not null"`
}

// TableName 指定权限表名
func (PermissionModel) TableName() string {
	return "permissions"
}

func newPermissionModelFromEntity(entity *domainRole.Permission) *PermissionModel {
	if entity == nil {
		return nil
	}

	model := &PermissionModel{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Domain:      entity.Domain,
		Resource:    entity.Resource,
		Action:      entity.Action,
		Description: entity.Description,
		Code:        entity.Code,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

func (m *PermissionModel) toEntity() *domainRole.Permission {
	if m == nil {
		return nil
	}

	entity := &domainRole.Permission{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Domain:      m.Domain,
		Resource:    m.Resource,
		Action:      m.Action,
		Description: m.Description,
		Code:        m.Code,
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

func mapPermissionModelsToEntities(models []PermissionModel) []domainRole.Permission {
	if len(models) == 0 {
		return nil
	}

	permissions := make([]domainRole.Permission, 0, len(models))
	for i := range models {
		if entity := models[i].toEntity(); entity != nil {
			permissions = append(permissions, *entity)
		}
	}
	return permissions
}

func mapPermissionEntitiesToModels(perms []domainRole.Permission) []PermissionModel {
	if len(perms) == 0 {
		return nil
	}

	models := make([]PermissionModel, 0, len(perms))
	for i := range perms {
		p := perms[i]
		if model := newPermissionModelFromEntity(&p); model != nil {
			models = append(models, *model)
		}
	}
	return models
}
