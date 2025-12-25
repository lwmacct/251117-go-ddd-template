package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// permissionCommandRepository 权限命令仓储的 GORM 实现
// 嵌入 GenericCommandRepository 以复用基础 CRUD 操作
type permissionCommandRepository struct {
	*GenericCommandRepository[role.Permission, *PermissionModel]
}

// NewPermissionCommandRepository 创建权限命令仓储实例
func NewPermissionCommandRepository(db *gorm.DB) role.PermissionCommandRepository {
	return &permissionCommandRepository{
		GenericCommandRepository: NewGenericCommandRepository(
			db, newPermissionModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供
