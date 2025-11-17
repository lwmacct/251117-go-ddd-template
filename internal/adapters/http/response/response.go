package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 错误响应（符合 RFC 7807 Problem Details）
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Code    string      `json:"code"`              // 业务错误码
	Message string      `json:"message"`           // 错误消息
	Details interface{} `json:"details,omitempty"` // 额外详情（如验证错误列表）
}

// ListResponse 列表响应（带分页信息）
type ListResponse struct {
	Data interface{}     `json:"data"`
	Meta *PaginationMeta `json:"meta,omitempty"`
}

// PaginationMeta 分页元数据
type PaginationMeta struct {
	Total      int  `json:"total"`                 // 总记录数
	Page       int  `json:"page"`                  // 当前页码
	PerPage    int  `json:"per_page"`              // 每页数量
	TotalPages int  `json:"total_pages,omitempty"` // 总页数
	HasMore    bool `json:"has_more,omitempty"`    // 是否有下一页
}

// JSON 通用 JSON 响应（成功时直接返回数据）
func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// OK 200 成功响应（返回单个资源）
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// Created 201 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// NoContent 204 无内容响应（删除、更新成功等）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// List 200 列表响应（带可选分页）
func List(c *gin.Context, data interface{}, meta *PaginationMeta) {
	c.JSON(http.StatusOK, ListResponse{
		Data: data,
		Meta: meta,
	})
}

// Error 通用错误响应
func Error(c *gin.Context, statusCode int, code string, message string, details ...interface{}) {
	errDetail := ErrorDetail{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		errDetail.Details = details[0]
	}

	c.JSON(statusCode, ErrorResponse{
		Error: errDetail,
	})
}

// BadRequest 400 请求错误
func BadRequest(c *gin.Context, message string, details ...interface{}) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, details...)
}

// ValidationError 400 验证错误
func ValidationError(c *gin.Context, details interface{}) {
	Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
}

// Unauthorized 401 未认证
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden 403 无权限
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Access forbidden"
	}
	Error(c, http.StatusForbidden, "FORBIDDEN", message)
}

// NotFound 404 资源不存在
func NotFound(c *gin.Context, resource string) {
	message := "Resource not found"
	if resource != "" {
		message = resource + " not found"
	}
	Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict 409 资源冲突
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = "Resource conflict"
	}
	Error(c, http.StatusConflict, "CONFLICT", message)
}

// TooManyRequests 429 请求过多
func TooManyRequests(c *gin.Context) {
	Error(c, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "Rate limit exceeded")
}

// InternalError 500 服务器错误
func InternalError(c *gin.Context, details ...interface{}) {
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", details...)
}

// ServiceUnavailable 503 服务不可用
func ServiceUnavailable(c *gin.Context, message string) {
	if message == "" {
		message = "Service temporarily unavailable"
	}
	Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message)
}

// NewPaginationMeta 创建分页元数据
func NewPaginationMeta(total, page, perPage int) *PaginationMeta {
	totalPages := (total + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}

	return &PaginationMeta{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}
}
