package operation

import "strings"

// ============================================================================
// Resource 资源标识符
// ============================================================================

// Resource 资源标识符。
// 格式：{type}/{id} 或 {type}/*
//
// 示例：
//   - user/123     具体用户资源
//   - user/*       所有用户资源
//   - user/self    当前用户自身（特殊标记）
//   - *            所有资源
type Resource string

// 预定义资源常量。
const (
	ResourceAll  Resource = "*"    // 所有资源
	ResourceSelf Resource = "self" // 当前用户自身（用于 ID 部分）
)

// NewResource 创建资源标识符。
func NewResource(resourceType, id string) Resource {
	if resourceType == "*" {
		return ResourceAll
	}
	return Resource(resourceType + "/" + id)
}

// Type 返回资源类型。
// "user/123" → "user"
// "*" → "*"
func (r Resource) Type() string {
	s := string(r)
	if s == "*" {
		return "*"
	}
	idx := strings.Index(s, "/")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// ID 返回资源 ID。
// "user/123" → "123"
// "user/*" → "*"
// "*" → "*"
func (r Resource) ID() string {
	s := string(r)
	if s == "*" {
		return "*"
	}
	idx := strings.Index(s, "/")
	if idx == -1 {
		return "*"
	}
	return s[idx+1:]
}

// IsWildcard 报告资源是否为通配符（全局或类型通配）。
func (r Resource) IsWildcard() bool {
	return r == ResourceAll || r.ID() == "*"
}

// IsSelf 报告资源是否为 self 标记（当前用户自身）。
func (r Resource) IsSelf() bool {
	return r.ID() == string(ResourceSelf)
}

// String 返回资源标识符字符串。
func (r Resource) String() string {
	return string(r)
}
