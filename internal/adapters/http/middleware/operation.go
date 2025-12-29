package middleware

import (
	"github.com/gin-gonic/gin"
)

// OperationIDKey 是 Operation ID 在 context 中的键名。
const OperationIDKey = "operation_id"

// SetOperationID 创建 Operation ID 中间件（直接设置）。
//
// 将给定的 Operation ID 存入 Gin context，供后续中间件和处理器使用。
// 这是声明式路由的首选方式，Operation ID 在路由注册时已知。
//
// Operation ID 格式：{domain}.{resource}.{action}
// 例如：admin.users.create, user.profile.get
func SetOperationID(operationID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(OperationIDKey, operationID)
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
