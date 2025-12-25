package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 返回一个基于 slog 的 Gin 日志中间件
//
// 记录每个 HTTP 请求的详细信息，包括：
// - 请求方法、路径、状态码
// - 响应时间、客户端 IP、User-Agent
// - 错误信息（如果有）
//
// 日志级别根据状态码自动选择：
// - 2xx: INFO
// - 3xx: INFO
// - 4xx: WARN
// - 5xx: ERROR
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算请求处理时间
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 构建日志字段
		fields := []any{
			"method", method,
			"path", path,
			"status", statusCode,
			"latency", latency.String(),
			"ip", clientIP,
			"user_agent", userAgent,
		}

		// 如果有查询参数，添加到日志
		if query != "" {
			fields = append(fields, "query", query)
		}

		// 如果有错误，添加到日志
		if len(c.Errors) > 0 {
			fields = append(fields, "errors", c.Errors.String())
		}

		// 根据状态码选择日志级别
		msg := "HTTP Request"
		switch {
		case statusCode >= 500:
			slog.Error(msg, fields...)
		case statusCode >= 400:
			slog.Warn(msg, fields...)
		default:
			slog.Info(msg, fields...)
		}
	}
}

// LoggerSkipPaths 返回一个跳过指定路径的日志中间件
//
// 使用场景：
// - 跳过健康检查端点 /health
// - 跳过静态资源请求
// - 跳过其他高频低价值的请求
//
// 示例：
//
//	r.Use(middleware.LoggerSkipPaths("/health", "/metrics"))
func LoggerSkipPaths(skipPaths ...string) gin.HandlerFunc {
	skipPathsMap := make(map[string]bool, len(skipPaths))
	for _, path := range skipPaths {
		skipPathsMap[path] = true
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果路径在跳过列表中，直接继续处理
		if skipPathsMap[path] {
			c.Next()
			return
		}

		// 否则执行日志记录
		start := time.Now()
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		fields := []any{
			"method", method,
			"path", path,
			"status", statusCode,
			"latency", latency.String(),
			"ip", clientIP,
			"user_agent", userAgent,
		}

		if query != "" {
			fields = append(fields, "query", query)
		}

		if len(c.Errors) > 0 {
			fields = append(fields, "errors", c.Errors.String())
		}

		msg := "HTTP Request"
		switch {
		case statusCode >= 500:
			slog.Error(msg, fields...)
		case statusCode >= 400:
			slog.Warn(msg, fields...)
		default:
			slog.Info(msg, fields...)
		}
	}
}
