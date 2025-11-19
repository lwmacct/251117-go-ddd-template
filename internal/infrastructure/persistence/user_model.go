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

	Roles []RoleModel `gorm:"many2many:user_roles;"`
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
