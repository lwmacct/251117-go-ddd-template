package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
	"gorm.io/gorm"
)

// teamCommandRepository 团队命令仓储的 GORM 实现
type teamCommandRepository struct {
	*GenericCommandRepository[organization.Team, *TeamModel]
}

// NewTeamCommandRepository 创建团队命令仓储实例
func NewTeamCommandRepository(db *gorm.DB) organization.TeamCommandRepository {
	return &teamCommandRepository{
		GenericCommandRepository: NewGenericCommandRepository(
			db, newTeamModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供
