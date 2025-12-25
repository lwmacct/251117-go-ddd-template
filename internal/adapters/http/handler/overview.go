package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/stats"
)

// OverviewHandler 系统概览处理器
type OverviewHandler struct {
	getStatsHandler *stats.GetStatsHandler
}

// NewOverviewHandler 创建 OverviewHandler 实例
func NewOverviewHandler(getStatsHandler *stats.GetStatsHandler) *OverviewHandler {
	return &OverviewHandler{
		getStatsHandler: getStatsHandler,
	}
}

// GetStats 获取系统统计信息
//
// @Summary      获取系统概览统计
// @Description  获取用户、角色、权限、菜单等统计信息
// @Tags         系统概览 (Overview)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.DataResponse[stats.StatsDTO] "统计信息"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/overview/stats [get]
func (h *OverviewHandler) GetStats(c *gin.Context) {
	result, err := h.getStatsHandler.Handle(c.Request.Context(), stats.GetStatsQuery{
		RecentLogsLimit: 5,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "success", result)
}
