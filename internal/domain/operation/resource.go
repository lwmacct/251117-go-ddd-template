package operation

import "github.com/lwmacct/251117-go-ddd-template/pkg/urn"

// ============================================================================
// Resource 资源标识符
// ============================================================================

// Resource 资源标识符。
// 格式：{scope}:{type}:{id}
//
// 示例：
//   - sys:user:123      系统用户 123
//   - self:user:@me     当前用户自身
//   - org.acme:user:*   org.acme 组织所有用户
//   - *:*:*             所有资源
//
// 特殊标识符：
//   - * 通配符，匹配任意值
//   - @me 当前用户 ID（运行时替换）
//   - @org 当前组织 ID（运行时替换）
//
// 本类型基于 [urn.URN] 格式。
type Resource string

// 预定义资源常量。
const (
	ResourceAll      Resource = "*:*:*"         // 所有资源
	ResourceSelfUser Resource = "self:user:@me" // 当前用户自身
)

// NewResource 创建资源标识符。
func NewResource(scope, resourceType, id string) Resource {
	return Resource(urn.New(scope, resourceType, id))
}

// URN 返回底层 URN 类型。
func (r Resource) URN() urn.URN {
	return urn.URN(r)
}

// Scope 返回资源的 scope。
func (r Resource) Scope() string {
	return r.URN().Scope()
}

// Type 返回资源的 type。
func (r Resource) Type() string {
	return r.URN().Type()
}

// Identifier 返回资源的 identifier（ID）。
func (r Resource) Identifier() string {
	return r.URN().Identifier()
}

// Match 检查资源是否匹配指定模式。
func (r Resource) Match(pattern string) bool {
	return r.URN().Match(urn.URN(pattern))
}

// IsWildcard 报告资源是否为通配符。
func (r Resource) IsWildcard() bool {
	return r.URN().IsWildcard()
}

// String 返回资源字符串。
func (r Resource) String() string {
	return string(r)
}
