package persistence

import (
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// ============================================================================
// Role Model
// ============================================================================

// RoleModel 定义角色的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type RoleModel struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"size:50;uniqueIndex;not null"`
	DisplayName string         `gorm:"size:100;not null"`
	Description string         `gorm:"size:255"`
	IsSystem    bool           `gorm:"default:false;not null"`
}

// TableName 指定角色表名
func (RoleModel) TableName() string {
	return "roles"
}

func newRoleModelFromEntity(entity *role.Role) *RoleModel {
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
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity（不含权限）
func (m *RoleModel) ToEntity() *role.Role {
	if m == nil {
		return nil
	}

	entity := &role.Role{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Description: m.Description,
		IsSystem:    m.IsSystem,
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

// ToEntityWithPermissions 将 GORM Model 转换为 Domain Entity（含权限）
func (m *RoleModel) ToEntityWithPermissions(permissions []RolePermissionModel) *role.Role {
	entity := m.ToEntity()
	if entity == nil {
		return nil
	}

	entity.Permissions = mapRolePermissionModelsToEntities(permissions)
	return entity
}

func mapRoleModelsToEntities(models []RoleModel) []role.Role {
	if len(models) == 0 {
		return nil
	}

	roles := make([]role.Role, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			roles = append(roles, *entity)
		}
	}
	return roles
}

// ============================================================================
// Role Permission Model（关联表）
// ============================================================================

// RolePermissionModel 角色权限关联表
// 存储 Operation + Resource 模式，支持通配符
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type RolePermissionModel struct {
	RoleID           uint   `gorm:"primaryKey"`
	OperationPattern string `gorm:"primaryKey;size:100;not null"`
	ResourcePattern  string `gorm:"primaryKey;size:100;not null;default:'*'"`
	CreatedAt        time.Time
}

// TableName 指定关联表名
func (RolePermissionModel) TableName() string {
	return "role_permissions"
}

func newRolePermissionModelFromEntity(roleID uint, p role.Permission) RolePermissionModel {
	resPattern := p.ResourcePattern
	if resPattern == "" {
		resPattern = "*"
	}
	return RolePermissionModel{
		RoleID:           roleID,
		OperationPattern: p.OperationPattern,
		ResourcePattern:  resPattern,
		CreatedAt:        time.Now(),
	}
}

func (m *RolePermissionModel) toEntity() role.Permission {
	return role.Permission{
		OperationPattern: m.OperationPattern,
		ResourcePattern:  m.ResourcePattern,
	}
}

func mapRolePermissionModelsToEntities(models []RolePermissionModel) []role.Permission {
	if len(models) == 0 {
		return nil
	}

	permissions := make([]role.Permission, len(models))
	for i := range models {
		permissions[i] = models[i].toEntity()
	}
	return permissions
}

func mapPermissionEntitiesToRolePermissionModels(roleID uint, permissions []role.Permission) []RolePermissionModel {
	if len(permissions) == 0 {
		return nil
	}

	models := make([]RolePermissionModel, len(permissions))
	for i, p := range permissions {
		models[i] = newRolePermissionModelFromEntity(roleID, p)
	}
	return models
}
