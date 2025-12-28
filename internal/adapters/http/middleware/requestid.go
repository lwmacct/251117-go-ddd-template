package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

// RequestIDKey 是 Request ID 在 context 中的键名。
const RequestIDKey = "request_id"

// RequestIDHeader 是 Request ID 响应头名称。
const RequestIDHeader = "X-Request-ID"

// RequestID 创建 Request ID 中间件。
//
// 优先从 OpenTelemetry Span 提取 Trace ID 作为 Request ID，
// 如果 OTel 未启用，则生成 UUID 作为 fallback。
//
// Request ID 会被：
//   - 存入 Gin context（供后续中间件和处理器使用）
//   - 设置到 X-Request-ID 响应头（供客户端追踪）
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestID string

		// 优先从 OTel Span 获取 Trace ID
		span := trace.SpanFromContext(c.Request.Context())
		if span.SpanContext().IsValid() {
			requestID = span.SpanContext().TraceID().String()
		} else {
			// Fallback: 生成 UUID
			requestID = uuid.New().String()
		}

		// 存入 Gin context
		c.Set(RequestIDKey, requestID)

		// 设置响应头（便于客户端追踪）
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID 从 Gin context 获取 Request ID。
// 如果不存在返回空字符串。
func GetRequestID(c *gin.Context) string {
	if id, ok := c.Get(RequestIDKey); ok {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}
