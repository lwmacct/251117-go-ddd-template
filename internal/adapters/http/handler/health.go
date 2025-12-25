package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/health"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	checker health.Checker
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(checker health.Checker) *HealthHandler {
	return &HealthHandler{
		checker: checker,
	}
}

// Check 执行健康检查
//
// @Summary      健康检查
// @Description  检查系统服务健康状态（数据库、Redis）
// @Tags         系统 (System)
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response{status=string,checks=object{database=object,redis=object}} "服务健康"
// @Failure      503 {object} response.Response{status=string,checks=object{database=object,redis=object}} "服务降级"
// @Router       /api/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	report := h.checker.Check(ctx)

	if report.Status == health.StatusHealthy {
		response.OK(c, "healthy", report)
	} else {
		response.ServiceUnavailable(c, string(report.Status))
	}
}
