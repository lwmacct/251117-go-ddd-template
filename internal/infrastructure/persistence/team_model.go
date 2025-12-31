package persistence

import (
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
	"gorm.io/gorm"
)

// TeamModel 定义团队的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type TeamModel struct {
	ID             uint   `gorm:"primaryKey"`
	OrganizationID uint   `gorm:"index;not null"`
	Name           string `gorm:"size:50;not null;index:idx_team_org_name,unique,priority:2"`
	DisplayName    string `gorm:"size:100;not null"`
	Description    string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 聚合关系
	Members []TeamMemberModel `gorm:"foreignKey:TeamID"`
}

// TableName 指定团队表名
func (TeamModel) TableName() string {
	return "teams"
}

func init() {
	// 确保组织ID+名称唯一的复合索引
	// GORM 会通过 tag 自动创建: idx_team_org_name
}

func newTeamModelFromEntity(entity *organization.Team) *TeamModel {
	if entity == nil {
		return nil
	}

	model := &TeamModel{
		ID:             entity.ID,
		OrganizationID: entity.OrganizationID,
		Name:           entity.Name,
		DisplayName:    entity.DisplayName,
		Description:    entity.Description,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *TeamModel) ToEntity() *organization.Team {
	if m == nil {
		return nil
	}

	entity := &organization.Team{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
		DisplayName:    m.DisplayName,
		Description:    m.Description,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		Members:        mapTeamMemberModelsToEntities(m.Members),
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

func mapTeamModelsToEntities(models []TeamModel) []organization.Team {
	if len(models) == 0 {
		return nil
	}

	teams := make([]organization.Team, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			teams = append(teams, *entity)
		}
	}
	return teams
}

func mapTeamModelPtrsToEntities(models []*TeamModel) []*organization.Team {
	if len(models) == 0 {
		return nil
	}

	teams := make([]*organization.Team, 0, len(models))
	for _, m := range models {
		if entity := m.ToEntity(); entity != nil {
			teams = append(teams, entity)
		}
	}
	return teams
}
