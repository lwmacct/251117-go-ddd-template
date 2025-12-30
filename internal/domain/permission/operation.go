package permission

import "strings"

// Operation 统一操作标识符（URN 格式）。
//
// 格式：{scope}:{type}:{action}
//
// 常用 scope：
//   - public: 公开操作（无需认证）
//   - sys:    系统管理操作
//   - self:   用户自服务操作
//
// 示例：
//   - public:auth:login
//   - sys:users:create
//   - self:profile:update
type Operation string

// String 返回操作标识符字符串。
func (o Operation) String() string { return string(o) }

// Scope 返回操作的 scope。
//
//	Operation("sys:users:create").Scope() // "sys"
func (o Operation) Scope() string {
	return parseURN(string(o)).Scope
}

// Type 返回操作的 type。
//
//	Operation("sys:users:create").Type() // "users"
func (o Operation) Type() string {
	return parseURN(string(o)).Type
}

// Identifier 返回操作的 identifier（action）。
//
//	Operation("sys:users:create").Identifier() // "create"
func (o Operation) Identifier() string {
	return parseURN(string(o)).Identifier
}

// IsPublic 报告操作是否公开（无需权限检查）。
//
// scope 为 "public" 的操作无需认证和权限检查。
func (o Operation) IsPublic() bool {
	return o.Scope() == "public"
}

// ============================================================================
// URN 解析
// ============================================================================

// urnParts 包含 URN 的解析结果。
type urnParts struct {
	Scope      string   // Scope，如 "sys.admin" 或 "org.acme"
	ScopeParts []string // Scope 层级，如 ["sys", "admin"]
	Type       string   // 类型/模块，如 "users"
	Identifier string   // 标识符，如 "create" 或 "123"
}

// parseURN 解析 URN 字符串。
//
// 解析规则：
//   - "*" → {Scope: "*", Type: "*", Identifier: "*"}
//   - "scope:type:id" → {Scope: scope, Type: type, Identifier: id}
//   - "scope:type" → {Scope: scope, Type: type, Identifier: "*"}
//   - "scope" → {Scope: scope, Type: "*", Identifier: "*"}
func parseURN(s string) urnParts {
	// 超级通配符
	if s == "*" {
		return urnParts{
			Scope:      "*",
			ScopeParts: []string{"*"},
			Type:       "*",
			Identifier: "*",
		}
	}

	parts := strings.SplitN(s, ":", 3)

	scope := "*"
	typ := "*"
	identifier := "*"

	switch len(parts) {
	case 3:
		scope = parts[0]
		typ = parts[1]
		identifier = parts[2]
	case 2:
		scope = parts[0]
		typ = parts[1]
	case 1:
		scope = parts[0]
	}

	// 解析 scope 层级
	var scopeParts []string
	if scope == "*" {
		scopeParts = []string{"*"}
	} else {
		scopeParts = strings.Split(scope, ".")
	}

	return urnParts{
		Scope:      scope,
		ScopeParts: scopeParts,
		Type:       typ,
		Identifier: identifier,
	}
}

// NewURN 创建 URN 字符串。
func NewURN(scope, typ, identifier string) string {
	return scope + ":" + typ + ":" + identifier
}
