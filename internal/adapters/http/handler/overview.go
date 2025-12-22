package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	statsQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/stats/query"
)

// OverviewHandler 系统概览处理器
type OverviewHandler struct {
	getStatsHandler *statsQuery.GetStatsHandler
}

// NewOverviewHandler 创建 OverviewHandler 实例
func NewOverviewHandler(getStatsHandler *statsQuery.GetStatsHandler) *OverviewHandler {
	return &OverviewHandler{
		getStatsHandler: getStatsHandler,
	}
}

// GetStats 获取系统统计信息
// @Summary 获取系统概览统计
// @Description 获取用户、角色、权限、菜单等统计信息
// @Tags Overview
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=statsQuery.GetStatsResult}
// @Failure 500 {object} response.Response
// @Router /overview/stats [get]
func (h *OverviewHandler) GetStats(c *gin.Context) {
	result, err := h.getStatsHandler.Handle(context.Background(), statsQuery.GetStatsQuery{
		RecentLogsLimit: 5,
	})
	if err != nil {
		response.InternalError(c, "failed to get system stats", err.Error())
		return
	}

	response.OK(c, "success", result)
}
