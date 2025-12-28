package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/registry"
)

// OperationIDKey 是 Operation ID 在 context 中的键名。
const OperationIDKey = "operation_id"

// OperationID 创建 Operation ID 中间件。
//
// 根据请求的 HTTP 方法和路径，从 API 注册表查找对应的 Operation ID，
// 并将其存入 Gin context，供后续中间件和处理器使用。
//
// Operation ID 格式：{domain}.{resource}.{action}
// 例如：admin.users.create, user.profile.get
func OperationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.FullPath() // 获取路由模式而非实际路径

		// 从注册表查找端点
		if ep := registry.ByPath(method, path); ep != nil {
			c.Set(OperationIDKey, ep.OperationID)
		}

		c.Next()
	}
}

// GetOperationID 从 Gin context 获取 Operation ID。
// 如果不存在返回空字符串。
func GetOperationID(c *gin.Context) string {
	if id, ok := c.Get(OperationIDKey); ok {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}
