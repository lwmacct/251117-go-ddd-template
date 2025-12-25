package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog"
)

// ListAuditLogsQuery 审计日志列表查询参数
type ListAuditLogsQuery struct {
	response.PaginationQueryDTO

	// UserID 按用户 ID 过滤
	UserID *uint `form:"user_id" json:"user_id" binding:"omitempty,gt=0"`
	// Action 操作类型过滤
	Action string `form:"action" json:"action" binding:"omitempty,oneof=create update delete login logout" enums:"create,update,delete,login,logout"`
	// Resource 资源类型过滤
	Resource string `form:"resource" json:"resource" binding:"omitempty,oneof=user role menu setting" enums:"user,role,menu,setting"`
	// Status 状态过滤
	Status string `form:"status" json:"status" binding:"omitempty,oneof=success failure" enums:"success,failure"`
	// StartDate 开始时间（RFC3339 格式）
	StartDate string `form:"start_date" json:"start_date" binding:"omitempty"`
	// EndDate 结束时间（RFC3339 格式）
	EndDate string `form:"end_date" json:"end_date" binding:"omitempty"`
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListAuditLogsQuery) ToQuery() auditlog.ListLogsQuery {
	result := auditlog.ListLogsQuery{
		Page:     q.GetPage(),
		Limit:    q.GetLimit(),
		UserID:   q.UserID,
		Action:   q.Action,
		Resource: q.Resource,
		Status:   q.Status,
	}

	// 解析时间字符串
	if q.StartDate != "" {
		if t, err := time.Parse(time.RFC3339, q.StartDate); err == nil {
			result.StartDate = &t
		}
	}
	if q.EndDate != "" {
		if t, err := time.Parse(time.RFC3339, q.EndDate); err == nil {
			result.EndDate = &t
		}
	}

	return result
}

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
// @Param        params query handler.ListAuditLogsQuery false "查询参数"
// @Success      200 {object} response.PagedResponse[auditlog.AuditLogDTO] "审计日志列表"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/auditlogs [get]
// @x-permission {"scope":"admin:audit_logs:read"}
func (h *AuditLogHandler) ListLogs(c *gin.Context) {
	var q ListAuditLogsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listLogsHandler.Handle(c.Request.Context(), q.ToQuery())
	if err != nil {
		response.InternalError(c, "failed to list audit logs")
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), q.GetPage(), q.GetLimit())
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

	log, err := h.getLogHandler.Handle(c.Request.Context(), auditlog.GetLogQuery{
		LogID: uint(id),
	})

	if err != nil {
		response.NotFound(c, "audit log")
		return
	}

	response.OK(c, "success", log)
}
