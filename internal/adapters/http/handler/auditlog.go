package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// AuditLogHandler handles audit log operations
type AuditLogHandler struct {
	repo auditlog.Repository
}

// NewAuditLogHandler creates a new AuditLogHandler instance
func NewAuditLogHandler(repo auditlog.Repository) *AuditLogHandler {
	return &AuditLogHandler{
		repo: repo,
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

	filter := auditlog.FilterOptions{
		Page:  page,
		Limit: limit,
	}

	// Parse filters
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(userID)
			filter.UserID = &uid
		}
	}

	if action := c.Query("action"); action != "" {
		filter.Action = action
	}

	if resource := c.Query("resource"); resource != "" {
		filter.Resource = resource
	}

	if status := c.Query("status"); status != "" {
		filter.Status = status
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	logs, total, err := h.repo.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
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

	log, err := h.repo.FindByID(c.Request.Context(), uint(id))
	if err != nil || log == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audit log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}
