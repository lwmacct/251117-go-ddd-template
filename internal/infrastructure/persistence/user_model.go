package persistence

import (
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"gorm.io/gorm"
)

// UserModel 定义用户的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type UserModel struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;size:50;not null"`
	Email    string `gorm:"uniqueIndex;size:100;not null"`
	Password string `gorm:"size:255;not null"`
	FullName string `gorm:"size:100"`
	Avatar   string `gorm:"size:255"`
	Bio      string `gorm:"type:text"`
	Status   string `gorm:"size:20;default:'active'"`

	// Type 用户类型：human（人类用户）或 service（服务账户）
	Type string `gorm:"size:20;default:'human';not null;index"`

	// IsSystem 系统预置用户标记（如 root、admin），不可删除
	IsSystem bool `gorm:"default:false;not null;index"`

	Roles []RoleModel `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID;foreignKey:ID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 指定用户表名
func (UserModel) TableName() string {
	return "users"
}

func newUserModelFromEntity(entity *user.User) *UserModel {
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
		Type:      string(entity.Type),
		IsSystem:  entity.IsSystem,
		// Roles 不在这里映射，通过 user_roles 关联表管理
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity（实现 Model[E] 接口）
func (m *UserModel) ToEntity() *user.User {
	if m == nil {
		return nil
	}

	entity := &user.User{
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
		Type:      user.UserType(m.Type),
		IsSystem:  m.IsSystem,
		Roles:     mapRoleModelsToEntities(m.Roles),
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

func mapUserModelsToEntities(models []UserModel) []*user.User {
	if len(models) == 0 {
		return nil
	}

	users := make([]*user.User, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			users = append(users, entity)
		}
	}
	return users
}
