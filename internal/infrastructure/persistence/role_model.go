package persistence

import (
	"time"

	domainRole "github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// RoleModel 定义角色的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type RoleModel struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt    `gorm:"index"`
	Name        string            `gorm:"size:50;uniqueIndex;not null"`
	DisplayName string            `gorm:"size:100;not null"`
	Description string            `gorm:"size:255"`
	IsSystem    bool              `gorm:"default:false;not null"`
	Permissions []PermissionModel `gorm:"many2many:role_permissions;"`
}

// TableName 指定角色表名
func (RoleModel) TableName() string {
	return "roles"
}

func newRoleModelFromEntity(entity *domainRole.Role) *RoleModel {
	if entity == nil {
		return nil
	}

	model := &RoleModel{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Name:        entity.Name,
		DisplayName: entity.DisplayName,
		Description: entity.Description,
		IsSystem:    entity.IsSystem,
		Permissions: mapPermissionEntitiesToModels(entity.Permissions),
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity（实现 Model[E] 接口）
func (m *RoleModel) ToEntity() *domainRole.Role {
	if m == nil {
		return nil
	}

	entity := &domainRole.Role{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Description: m.Description,
		IsSystem:    m.IsSystem,
		Permissions: mapPermissionModelsToEntities(m.Permissions),
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

func mapRoleModelsToEntities(models []RoleModel) []domainRole.Role {
	if len(models) == 0 {
		return nil
	}

	roles := make([]domainRole.Role, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			roles = append(roles, *entity)
		}
	}
	return roles
}

func mapRoleEntitiesToModels(roles []domainRole.Role) []RoleModel {
	if len(roles) == 0 {
		return nil
	}

	models := make([]RoleModel, 0, len(roles))
	for i := range roles {
		r := roles[i]
		if model := newRoleModelFromEntity(&r); model != nil {
			models = append(models, *model)
		}
	}
	return models
}
