package permission

// Resource 资源标识符（URN 格式）。
//
// 格式：{scope}:{type}:{id}
//
// 示例：
//   - sys:user:123       系统用户 ID 123
//   - self:user:@me      当前用户（运行时替换）
//   - org.acme:user:*    org.acme 组织所有用户
//   - *:*:*              所有资源
//
// 特殊标识符：
//   - * 通配符，匹配任意值
//   - @me 当前用户 ID（运行时替换）
//   - @org 当前组织 ID（运行时替换）
type Resource string

// 预定义资源常量。
const (
	// ResourceAll 匹配所有资源。
	ResourceAll Resource = "*:*:*"
)

// NewResource 创建资源标识符。
func NewResource(scope, resourceType, id string) Resource {
	return Resource(NewURN(scope, resourceType, id))
}

// String 返回资源标识符字符串。
func (r Resource) String() string {
	return string(r)
}

// Scope 返回资源的 scope。
func (r Resource) Scope() string {
	return parseURN(string(r)).Scope
}

// Type 返回资源的 type。
func (r Resource) Type() string {
	return parseURN(string(r)).Type
}

// Identifier 返回资源的 identifier（ID）。
func (r Resource) Identifier() string {
	return parseURN(string(r)).Identifier
}

// IsWildcard 报告资源是否为通配符。
func (r Resource) IsWildcard() bool {
	s := string(r)
	return s == "*" || s == "*:*:*"
}
