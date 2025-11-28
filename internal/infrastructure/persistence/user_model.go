// Package persistence 提供领域模型的持久化实现。
//
// 本包是 DDD 架构中基础设施层的核心组件，职责：
//   - 定义 GORM Model：包含数据库映射的持久化模型（带 GORM Tag）
//   - 实现 Repository：实现领域层定义的仓储接口
//   - 模型映射：提供 Domain Entity ↔ GORM Model 的双向转换
//
// CQRS 模式：
//   - CommandRepository：处理写操作（Create, Update, Delete）
//   - QueryRepository：处理读操作（Get, List, Search, Count）
//
// 文件命名规范：
//   - {module}_model.go: GORM 模型定义和映射函数
//   - {module}_command_repository.go: 写仓储实现
//   - {module}_query_repository.go: 读仓储实现
//   - {module}_repositories.go: 仓储聚合（可选）
package persistence

import (
	"time"

	domainUser "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"gorm.io/gorm"
)

// UserModel 定义用户的 GORM 持久化模型
type UserModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Username string `gorm:"uniqueIndex;size:50;not null"`
	Email    string `gorm:"uniqueIndex;size:100;not null"`
	Password string `gorm:"size:255;not null"`
	FullName string `gorm:"size:100"`
	Avatar   string `gorm:"size:255"`
	Bio      string `gorm:"type:text"`
	Status   string `gorm:"size:20;default:'active'"`

	Roles []RoleModel `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID;foreignKey:ID;references:ID"`
}

// TableName 指定用户表名
func (UserModel) TableName() string {
	return "users"
}

func newUserModelFromEntity(entity *domainUser.User) *UserModel {
	if entity == nil {
		return nil
	}

	model := &UserModel{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Username:  entity.Username,
		Email:     entity.Email,
		Password:  entity.Password,
		FullName:  entity.FullName,
		Avatar:    entity.Avatar,
		Bio:       entity.Bio,
		Status:    entity.Status,
		Roles:     mapRoleEntitiesToModels(entity.Roles),
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

func (m *UserModel) toEntity() *domainUser.User {
	if m == nil {
		return nil
	}

	entity := &domainUser.User{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		FullName:  m.FullName,
		Avatar:    m.Avatar,
		Bio:       m.Bio,
		Status:    m.Status,
		Roles:     mapRoleModelsToEntities(m.Roles),
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

func mapUserModelsToEntities(models []UserModel) []*domainUser.User {
	if len(models) == 0 {
		return nil
	}

	users := make([]*domainUser.User, 0, len(models))
	for i := range models {
		if entity := models[i].toEntity(); entity != nil {
			users = append(users, entity)
		}
	}
	return users
}
