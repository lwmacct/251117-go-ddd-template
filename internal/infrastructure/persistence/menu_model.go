package persistence

import (
	"time"

	domainMenu "github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
	"gorm.io/gorm"
)

// MenuModel 菜单的 GORM 实体
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type MenuModel struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"type:varchar(100);not null"`
	Path      string `gorm:"type:varchar(255);not null"`
	Icon      string `gorm:"type:varchar(100)"`
	ParentID  *uint  `gorm:"index"`
	Order     int    `gorm:"default:0"`
	Visible   bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 指定菜单表名
func (MenuModel) TableName() string {
	return "menus"
}

func newMenuModelFromEntity(entity *domainMenu.Menu) *MenuModel {
	if entity == nil {
		return nil
	}

	model := &MenuModel{
		ID:        entity.ID,
		Title:     entity.Title,
		Path:      entity.Path,
		Icon:      entity.Icon,
		ParentID:  entity.ParentID,
		Order:     entity.Order,
		Visible:   entity.Visible,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity（实现 Model[E] 接口）
func (m *MenuModel) ToEntity() *domainMenu.Menu {
	if m == nil {
		return nil
	}

	entity := &domainMenu.Menu{
		ID:        m.ID,
		Title:     m.Title,
		Path:      m.Path,
		Icon:      m.Icon,
		ParentID:  m.ParentID,
		Order:     m.Order,
		Visible:   m.Visible,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}
