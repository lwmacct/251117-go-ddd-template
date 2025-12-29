package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"
)

// AuditMiddleware 创建审计日志中间件。
// 基于 operation registry 决策是否记录审计日志，未注册或无需审计的操作静默跳过。
func AuditMiddleware(handler *auditlog.CreateHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过 operation registry 查找操作
		op := operation.ByMethodAndPath(
			operation.HTTPMethod(c.Request.Method),
			c.Request.URL.Path,
		)

		// 未注册或无需审计的操作静默跳过
		if !op.Valid() || !op.NeedsAudit() {
			c.Next()
			return
		}

		// 提取用户信息
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 未认证请求跳过审计
		if userID == nil || username == nil {
			c.Next()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			c.Next()
			return
		}

		uname, ok := username.(string)
		if !ok {
			c.Next()
			return
		}

		// 读取请求体用于记录详情
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 执行请求
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// 根据响应状态码确定审计状态
		status := "success"
		if c.Writer.Status() >= 400 {
			status = "failure"
		}

		// 从路径提取资源 ID
		resourceID := extractResourceID(op.Path(), c.Request.URL.Path)

		// 创建审计日志命令
		cmd := auditlog.CreateCommand{
			UserID:      uid,
			Username:    uname,
			Action:      op.AuditAction(),      // 语义化操作标识：setting.update
			Resource:    string(op.AuditCat()), // 资源分类：setting
			ResourceID:  resourceID,
			IPAddress:   c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
			Details:     formatDetails(c.Request.Method, requestBody, c.Writer.Status(), duration),
			Status:      status,
			RequestID:   GetRequestID(c),
			OperationID: string(op), // 操作标识：admin.settings.update
		}

		// 异步保存审计日志，避免阻塞响应
		asyncCtx := context.WithoutCancel(c.Request.Context())
		go func() {
			if err := handler.Handle(asyncCtx, cmd); err != nil {
				slog.Error("failed to create audit log", "error", err)
			}
		}()
	}
}

// extractResourceID 从实际路径中提取资源 ID。
// pattern: /api/admin/users/:id
// actual:  /api/admin/users/123
// 返回: "123"
func extractResourceID(pattern, actual string) string {
	patternSegs := splitPathSegments(pattern)
	actualSegs := splitPathSegments(actual)

	if len(patternSegs) != len(actualSegs) {
		return ""
	}

	// 查找第一个路径参数对应的值
	for i, seg := range patternSegs {
		if len(seg) > 0 && seg[0] == ':' {
			return actualSegs[i]
		}
	}

	return ""
}

// splitPathSegments 将路径分割为段。
func splitPathSegments(path string) []string {
	if len(path) == 0 {
		return nil
	}

	// 移除开头的斜杠
	if path[0] == '/' {
		path = path[1:]
	}

	// 移除结尾的斜杠
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if len(path) == 0 {
		return nil
	}

	// 分割路径
	var segments []string
	start := 0
	for i := range len(path) {
		if path[i] == '/' {
			if i > start {
				segments = append(segments, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		segments = append(segments, path[start:])
	}

	return segments
}

// formatDetails 格式化请求详情为 JSON。
func formatDetails(method string, requestBody []byte, statusCode int, duration time.Duration) string {
	details := make(map[string]any)
	details["method"] = method
	details["status_code"] = statusCode
	details["duration_ms"] = duration.Milliseconds()

	// 解析请求体为 JSON（限制 10KB）
	if len(requestBody) > 0 && len(requestBody) < 10000 {
		var bodyJSON map[string]any
		if err := json.Unmarshal(requestBody, &bodyJSON); err == nil {
			// 移除敏感字段
			delete(bodyJSON, "password")
			delete(bodyJSON, "old_password")
			delete(bodyJSON, "new_password")
			details["request_body"] = bodyJSON
		}
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Sprintf(`{"method": "%s", "status_code": %d}`, method, statusCode)
	}

	return string(detailsJSON)
}
