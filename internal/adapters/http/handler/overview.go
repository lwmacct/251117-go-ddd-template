package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"gorm.io/gorm"
)

type OverviewHandler struct {
	db *gorm.DB
}

func NewOverviewHandler(db *gorm.DB) *OverviewHandler {
	return &OverviewHandler{db: db}
}

type SystemStatsResponse struct {
	TotalUsers       int64 `json:"total_users"`
	ActiveUsers      int64 `json:"active_users"`
	InactiveUsers    int64 `json:"inactive_users"`
	BannedUsers      int64 `json:"banned_users"`
	TotalRoles       int64 `json:"total_roles"`
	TotalPermissions int64 `json:"total_permissions"`
	TotalMenus       int64 `json:"total_menus"`
}

func (h *OverviewHandler) GetStats(c *gin.Context) {
	stats := SystemStatsResponse{}

	// 统计用户数（排除软删除）
	h.db.Table("users").Where("deleted_at IS NULL").Count(&stats.TotalUsers)
	h.db.Table("users").Where("deleted_at IS NULL AND status = ?", "active").Count(&stats.ActiveUsers)
	h.db.Table("users").Where("deleted_at IS NULL AND status = ?", "inactive").Count(&stats.InactiveUsers)
	h.db.Table("users").Where("deleted_at IS NULL AND status = ?", "banned").Count(&stats.BannedUsers)

	// 统计角色和权限
	h.db.Table("roles").Where("deleted_at IS NULL").Count(&stats.TotalRoles)
	h.db.Table("permissions").Where("deleted_at IS NULL").Count(&stats.TotalPermissions)

	// 统计菜单（使用表名）
	h.db.Table("menus").Count(&stats.TotalMenus)

	response.OK(c, stats)
}
