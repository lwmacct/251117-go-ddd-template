package persistence

import (
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
)

// OrganizationMemberModel 定义组织成员的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type OrganizationMemberModel struct {
	ID             uint      `gorm:"primaryKey"`
	OrganizationID uint      `gorm:"index:idx_org_member,unique,priority:1;not null"`
	UserID         uint      `gorm:"index:idx_org_member,unique,priority:2;index;not null"`
	Role           string    `gorm:"size:20;default:'member';not null"`
	JoinedAt       time.Time `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定组织成员表名
func (OrganizationMemberModel) TableName() string {
	return "organization_members"
}

func newOrgMemberModelFromEntity(entity *organization.Member) *OrganizationMemberModel {
	if entity == nil {
		return nil
	}

	return &OrganizationMemberModel{
		ID:             entity.ID,
		OrganizationID: entity.OrganizationID,
		UserID:         entity.UserID,
		Role:           string(entity.Role),
		JoinedAt:       entity.JoinedAt,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *OrganizationMemberModel) ToEntity() *organization.Member {
	if m == nil {
		return nil
	}

	return &organization.Member{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		Role:           organization.MemberRole(m.Role),
		JoinedAt:       m.JoinedAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func mapOrgMemberModelsToEntities(models []OrganizationMemberModel) []organization.Member {
	if len(models) == 0 {
		return nil
	}

	members := make([]organization.Member, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			members = append(members, *entity)
		}
	}
	return members
}

func mapOrgMemberModelPtrsToEntities(models []*OrganizationMemberModel) []*organization.Member {
	if len(models) == 0 {
		return nil
	}

	members := make([]*organization.Member, 0, len(models))
	for _, m := range models {
		if entity := m.ToEntity(); entity != nil {
			members = append(members, entity)
		}
	}
	return members
}
