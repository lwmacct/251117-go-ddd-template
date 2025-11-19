package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	auditlogQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog/query"
)

// AuditLogHandler handles audit log operations (DDD+CQRS Use Case Pattern)
type AuditLogHandler struct {
	// Query Handlers
	listLogsHandler *auditlogQuery.ListLogsHandler
	getLogHandler   *auditlogQuery.GetLogHandler
}

// NewAuditLogHandler creates a new AuditLogHandler instance
func NewAuditLogHandler(
	listLogsHandler *auditlogQuery.ListLogsHandler,
	getLogHandler *auditlogQuery.GetLogHandler,
) *AuditLogHandler {
	return &AuditLogHandler{
		listLogsHandler: listLogsHandler,
		getLogHandler:   getLogHandler,
	}
}

// ListLogs lists audit logs with filtering
func (h *AuditLogHandler) ListLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	query := auditlogQuery.ListLogsQuery{
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": result.Logs,
		"pagination": gin.H{
			"page":  result.Page,
			"limit": result.Limit,
			"total": result.Total,
		},
	})
}

// GetLog gets an audit log by ID
func (h *AuditLogHandler) GetLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	// 调用 Use Case Handler
	log, err := h.getLogHandler.Handle(c.Request.Context(), auditlogQuery.GetLogQuery{
		LogID: uint(id),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audit log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}
