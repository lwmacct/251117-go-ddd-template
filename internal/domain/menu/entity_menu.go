package menu

import (
	"time"
)

// Menu 菜单实体
type Menu struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	Title     string     `json:"title" gorm:"type:varchar(100);not null"`
	Path      string     `json:"path" gorm:"type:varchar(255);not null"`
	Icon      string     `json:"icon" gorm:"type:varchar(100)"`
	ParentID  *uint      `json:"parent_id" gorm:"index"`
	Order     int        `json:"order" gorm:"default:0"`
	Visible   bool       `json:"visible" gorm:"default:true"`
	Children  []*Menu    `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}
