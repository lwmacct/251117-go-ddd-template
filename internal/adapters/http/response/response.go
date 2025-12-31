// Package response 提供统一的 HTTP 响应格式。
//
// 本包定义了前后端约定的响应结构，确保 API 响应的一致性：
//
// 成功响应格式：
//
//	{
//	  "code": 200,
//	  "message": "操作成功",
//	  "data": { ... }
//	}
//
// 错误响应格式：
//
//	{
//	  "code": 400,
//	  "message": "验证失败",
//	  "error": { ... }
//	}
//
// 列表响应格式：
//
//	{
//	  "code": 200,
//	  "message": "操作成功",
//	  "data": [...],
//	  "meta": { "total": 100, "page": 1, "per_page": 10 }
//	}
//
// 使用示例：
//
//	response.OK(c, response.MsgSuccess, user)
//	response.Created(c, response.MsgCreated, nil)
//	response.BadRequest(c, "无效输入")
//	response.List(c, response.MsgSuccess, users, response.NewPaginationMeta(total, page, limit))
package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// 响应结构定义
// ============================================================================

// UnifiedResponse 统一响应结构
// 前端期望格式：{ code: number, message: string, data?: any, error?: any }
type UnifiedResponse struct {
	Code    int    `json:"code"`            // HTTP 状态码（数字）
	Message string `json:"message"`         // 消息描述
	Data    any    `json:"data,omitempty"`  // 成功时的数据
	Error   any    `json:"error,omitempty"` // 失败时的错误详情
}

// Response 通用响应结构（用于 Swagger 文档 - 已废弃，请使用泛型版本）
//
// Deprecated: 使用 DataResponse[T] 替代
type Response struct {
	Message string `json:"message,omitempty"` // 消息
	Data    any    `json:"data,omitempty"`    // 数据
}

// DataResponse 泛型数据响应（用于 Swagger 文档，消除 allOf 嵌套）
//
//	@Description	统一数据响应格式
type DataResponse[T any] struct {
	Code    int    `json:"code"`            // HTTP 状态码
	Message string `json:"message"`         // 消息描述
	Data    T      `json:"data,omitempty"`  // 响应数据
	Error   any    `json:"error,omitempty"` // 错误详情（仅失败时）
}

// MessageResponse 纯消息响应（无数据）
//
//	@Description	纯消息响应格式
type MessageResponse struct {
	Code    int    `json:"code"`    // HTTP 状态码
	Message string `json:"message"` // 消息描述
}

// ErrorResponse 错误响应结构（用于 Swagger 文档）
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Code    string `json:"code"`              // 业务错误码（小写下划线）
	Message string `json:"message"`           // 错误消息
	Details any    `json:"details,omitempty"` // 额外详情（如验证错误列表）
}

// ListResponse 列表响应（带分页信息 - 已废弃，请使用泛型版本）
//
// Deprecated: 使用 PagedResponse[T] 替代
type ListResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    any             `json:"data"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

// PagedResponse 泛型分页响应（用于 Swagger 文档，消除 allOf 嵌套）
//
//	@Description	分页列表响应格式
type PagedResponse[T any] struct {
	Code    int             `json:"code"`           // HTTP 状态码
	Message string          `json:"message"`        // 消息描述
	Data    []T             `json:"data"`           // 数据列表
	Meta    *PaginationMeta `json:"meta,omitempty"` // 分页信息
}

// PaginationMeta 分页元数据
type PaginationMeta struct {
	Total      int    `json:"total"`                 // 总记录数
	Page       int    `json:"page"`                  // 当前页码
	PerPage    int    `json:"per_page"`              // 每页数量
	TotalPages int    `json:"total_pages,omitempty"` // 总页数
	HasMore    bool   `json:"has_more,omitempty"`    // 是否有下一页
	Warning    string `json:"warning,omitempty"`     // 页码越界警告
}

// ============================================================================
// 成功响应函数
// ============================================================================

// Success 统一成功响应
// 返回格式：{ code: 200, message: "...", data: {...} }
func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, UnifiedResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// OK 200 成功响应
func OK(c *gin.Context, message string, data any) {
	Success(c, http.StatusOK, message, data)
}

// Created 201 创建成功响应
func Created(c *gin.Context, message string, data any) {
	Success(c, http.StatusCreated, message, data)
}

// NoContent 204 无内容响应
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// List 200 列表响应（带分页）
func List(c *gin.Context, message string, data any, meta *PaginationMeta) {
	c.JSON(http.StatusOK, ListResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ============================================================================
// 错误响应函数
// ============================================================================

// Failure 统一错误响应
// 返回格式：{ code: 400, message: "...", error: {...} }
func Failure(c *gin.Context, statusCode int, message string, errorDetails ...any) {
	resp := UnifiedResponse{
		Code:    statusCode,
		Message: message,
	}

	if len(errorDetails) > 0 {
		resp.Error = errorDetails[0]
	}

	c.JSON(statusCode, resp)
}

// BadRequest 400 请求错误
func BadRequest(c *gin.Context, message string, details ...any) {
	Failure(c, http.StatusBadRequest, message, details...)
}

// ValidationError 400 验证错误
func ValidationError(c *gin.Context, details any) {
	Failure(c, http.StatusBadRequest, MsgValidationFailed, details)
}

// Unauthorized 401 未认证
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = MsgAuthenticationRequired
	}
	Failure(c, http.StatusUnauthorized, message)
}

// Forbidden 403 无权限
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = MsgAccessForbidden
	}
	Failure(c, http.StatusForbidden, message)
}

// NotFound 404 资源不存在
func NotFound(c *gin.Context, resource string) {
	message := MsgResourceNotFound
	if resource != "" {
		message = resource + " not found"
	}
	Failure(c, http.StatusNotFound, message)
}

// NotFoundMessage 404 资源不存在（自定义消息）
// 用于传递完整的错误消息（如来自 Domain 层的 err.Error()）
func NotFoundMessage(c *gin.Context, message string) {
	if message == "" {
		message = MsgResourceNotFound
	}
	Failure(c, http.StatusNotFound, message)
}

// Conflict 409 资源冲突
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = MsgResourceConflict
	}
	Failure(c, http.StatusConflict, message)
}

// TooManyRequests 429 请求过多
func TooManyRequests(c *gin.Context) {
	Failure(c, http.StatusTooManyRequests, MsgRateLimitExceeded)
}

// InternalError 500 服务器错误
func InternalError(c *gin.Context, details ...any) {
	Failure(c, http.StatusInternalServerError, MsgInternalError, details...)
}

// ServiceUnavailable 503 服务不可用
func ServiceUnavailable(c *gin.Context, message string) {
	if message == "" {
		message = MsgServiceUnavailable
	}
	Failure(c, http.StatusServiceUnavailable, message)
}

// ============================================================================
// 工具函数
// ============================================================================

// NewPaginationMeta 创建分页元数据
func NewPaginationMeta(total, page, perPage int) *PaginationMeta {
	totalPages := max((total+perPage-1)/perPage, 1)

	meta := &PaginationMeta{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}

	if page > totalPages {
		meta.Warning = fmt.Sprintf("page %d exceeds total_pages %d", page, totalPages)
	}

	return meta
}
