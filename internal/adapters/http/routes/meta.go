package routes

// HTTPMethod HTTP 请求方法。
type HTTPMethod string

// HTTP 方法常量。
const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
	PATCH  HTTPMethod = "PATCH"
)

// routeMeta 路由元数据。
type routeMeta struct {
	// HTTP 路由
	Method HTTPMethod // HTTP 方法
	Path   string     // 路由路径（Gin 格式），如 /api/admin/users/:id

	// 中间件配置
	ReadOnly bool // 只读操作（对于团队操作，使用 TeamContextOptional 而非 TeamContext）

	// 审计配置
	Audit bool // 是否启用审计（审计详情从 Operation 派生）

	// Swagger 注解字段
	Tags        string //	@Tags，如	Admin - Users
	Summary     string //	@Summary，如	创建用户
	Description string //	@Description（可选）
}
