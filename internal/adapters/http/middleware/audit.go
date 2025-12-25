package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/auditlog"
)

// AuditMiddleware creates a middleware that logs user actions (CQRS: 只写操作)
func AuditMiddleware(handler *auditlog.CreateLogHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip audit for GET requests (read-only operations)
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		// Skip audit for health check and other non-sensitive endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/api/auth/login" || c.Request.URL.Path == "/api/auth/register" {
			c.Next()
			return
		}

		// Extract user information from context
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// If no user context, skip audit (unauthenticated request)
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

		// Read request body for details
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process the request
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// Determine status based on response code
		status := "success"
		if c.Writer.Status() >= 400 {
			status = "failure"
		}

		// Extract resource information from path
		resource, resourceID := extractResourceInfo(c.Request.URL.Path)

		// Create audit log command
		cmd := auditlog.CreateLogCommand{
			UserID:     uid,
			Username:   uname,
			Action:     fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
			Resource:   resource,
			ResourceID: resourceID,
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Details:    formatDetails(c.Request.Method, requestBody, c.Writer.Status(), duration),
			Status:     status,
		}

		// Save audit log asynchronously to avoid blocking the response
		// 使用 WithoutCancel 保留 context 中的值（trace ID 等），但不会被请求取消影响
		asyncCtx := context.WithoutCancel(c.Request.Context())
		go func() {
			if err := handler.Handle(asyncCtx, cmd); err != nil {
				// Log error but don't fail the request
				slog.Error("failed to create audit log", "error", err)
			}
		}()
	}
}

// extractResourceInfo extracts resource type and ID from URL path
//
//nolint:nonamedreturns // named returns for self-documenting API
func extractResourceInfo(path string) (resource, resourceID string) {
	// Simple extraction logic
	// Examples:
	// /api/admin/users/123 -> resource: "users", resourceID: "123"
	// /api/user/me -> resource: "profile", resourceID: "me"
	// /api/admin/roles/5/permissions -> resource: "roles", resourceID: "5"

	if len(path) == 0 {
		return "", ""
	}

	// Remove leading slash
	if path[0] == '/' {
		path = path[1:]
	}

	// Split path into segments
	segments := splitPath(path)

	// Find resource and ID
	for i, segment := range segments {
		// Common resource identifiers
		if segment == "users" || segment == "roles" || segment == "permissions" || segment == "audit-logs" {
			resource = segment
			if i+1 < len(segments) && !isAction(segments[i+1]) {
				resourceID = segments[i+1]
			}
			break
		}
		if segment == "me" {
			resource = "profile"
			resourceID = "me"
			break
		}
	}

	if resource == "" {
		resource = "unknown"
	}

	return resource, resourceID
}

// splitPath splits a path into segments
func splitPath(path string) []string {
	var segments []string
	var current string

	for _, char := range path {
		if char == '/' {
			if current != "" {
				segments = append(segments, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		segments = append(segments, current)
	}

	return segments
}

// isAction checks if a segment is an action (not a resource ID)
func isAction(segment string) bool {
	actions := []string{"permissions", "roles", "status", "password", "email"}
	return slices.Contains(actions, segment)
}

// formatDetails formats request details for audit log
func formatDetails(method string, requestBody []byte, statusCode int, duration time.Duration) string {
	details := make(map[string]any)
	details["method"] = method
	details["status_code"] = statusCode
	details["duration_ms"] = duration.Milliseconds()

	// Try to parse request body as JSON
	if len(requestBody) > 0 && len(requestBody) < 10000 { // Limit to 10KB
		var bodyJSON map[string]any
		if err := json.Unmarshal(requestBody, &bodyJSON); err == nil {
			// Remove sensitive fields
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
