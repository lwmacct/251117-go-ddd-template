package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog"
)

// AuditLogHandler handles audit log operations (DDD+CQRS Use Case Pattern)
type AuditLogHandler struct {
	// Query Handlers
	listLogsHandler *auditlog.ListLogsHandler
	getLogHandler   *auditlog.GetLogHandler
}

// NewAuditLogHandler creates a new AuditLogHandler instance
func NewAuditLogHandler(
	listLogsHandler *auditlog.ListLogsHandler,
	getLogHandler *auditlog.GetLogHandler,
) *AuditLogHandler {
	return &AuditLogHandler{
		listLogsHandler: listLogsHandler,
		getLogHandler:   getLogHandler,
	}
}

// ListLogs lists audit logs with filtering
//
// @Summary      获取审计日志列表
// @Description  分页获取审计日志，支持按用户、操作、资源、状态、时间范围筛选
// @Tags         管理员 - 审计日志 (Admin - Audit Log)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "页码" default(1) minimum(1)
// @Param        limit query int false "每页数量" default(20) minimum(1) maximum(100)
// @Param        user_id query int false "用户ID"
// @Param        action query string false "操作类型" Enums(create, update, delete, login, logout)
// @Param        resource query string false "资源类型" Enums(user, role, menu, setting)
// @Param        status query string false "状态" Enums(success, failure)
// @Param        start_date query string false "开始时间(RFC3339格式)" example:"2024-01-01T00:00:00Z"
// @Param        end_date query string false "结束时间(RFC3339格式)" example:"2024-12-31T23:59:59Z"
// @Success      200 {object} response.PagedResponse[auditlog.AuditLogDTO] "审计日志列表"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/auditlogs [get]
// @x-permission {"scope":"admin:audit_logs:read"}
func (h *AuditLogHandler) ListLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	query := auditlog.ListLogsQuery{
		Page:  page,
		Limit: limit,
	}

	// Parse filters
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(userID)
			query.UserID = &uid
		}
	}

	if action := c.Query("action"); action != "" {
		query.Action = action
	}

	if resource := c.Query("resource"); resource != "" {
		query.Resource = resource
	}

	if status := c.Query("status"); status != "" {
		query.Status = status
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			query.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			query.EndDate = &endDate
		}
	}

	// 调用 Use Case Handler
	result, err := h.listLogsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		response.InternalError(c, "failed to list audit logs")
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), page, limit)
	response.List(c, "success", result.Logs, meta)
}

// GetLog gets an audit log by ID
//
// @Summary      获取审计日志详情
// @Description  根据日志ID获取审计日志详细信息
// @Tags         管理员 - 审计日志 (Admin - Audit Log)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "日志ID" minimum(1)
// @Success      200 {object} response.DataResponse[auditlog.AuditLogDTO] "日志详情"
// @Failure      400 {object} response.ErrorResponse "无效的日志ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "日志不存在"
// @Router       /api/admin/auditlogs/{id} [get]
// @x-permission {"scope":"admin:audit_logs:read"}
func (h *AuditLogHandler) GetLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid log ID")
		return
	}

	// 调用 Use Case Handler
	log, err := h.getLogHandler.Handle(c.Request.Context(), auditlog.GetLogQuery{
		LogID: uint(id),
	})

	if err != nil {
		response.NotFound(c, "audit log")
		return
	}

	response.OK(c, "success", log)
}
