package permission

import "strings"

// MatchOperation 检查操作是否匹配模式。
//
// 匹配规则：
//   - 每个段（scope, type, identifier）独立匹配
//   - * 匹配该段的任意值
//   - Scope 段：.* 后缀匹配子 scope（sys.* 匹配 sys.admin）
//
// 示例：
//
//	MatchOperation("*:*:*", "sys:users:create")            // true - 超级通配符
//	MatchOperation("sys:*:*", "sys:users:create")          // true - scope 通配
//	MatchOperation("sys:users:*", "sys:users:create")      // true - type 通配
//	MatchOperation("sys:users:create", "sys:users:create") // true - 精确匹配
//	MatchOperation("sys.*:*:*", "sys.admin:config:update") // true - 子 scope 通配
func MatchOperation(pattern, operation string) bool {
	return match(pattern, operation)
}

// MatchResource 检查资源是否匹配模式。
//
// 匹配规则与 [MatchOperation] 相同。
//
// 注意：包含 @me 等变量的模式需要先使用 [Resolver] 解析。
func MatchResource(pattern, resource string) bool {
	return match(pattern, resource)
}

// match 执行 URN 模式匹配。
func match(pattern, target string) bool {
	// 超级通配符
	if pattern == "*" || pattern == "*:*:*" {
		return true
	}

	// 精确匹配
	if pattern == target {
		return true
	}

	p := parseURN(pattern)
	t := parseURN(target)

	// Scope 匹配
	if !matchScope(p.Scope, t.Scope) {
		return false
	}

	// Type 匹配
	if p.Type != "*" && p.Type != t.Type {
		return false
	}

	// Identifier 匹配
	if p.Identifier != "*" && p.Identifier != t.Identifier {
		return false
	}

	return true
}

// matchScope 匹配 scope，支持层级通配符。
//
// 规则：
//   - "*" 匹配任意 scope
//   - "sys" 精确匹配 "sys"
//   - "sys.*" 匹配 "sys"、"sys.admin"、"sys.readonly" 等
func matchScope(pattern, scope string) bool {
	// 通配符匹配所有
	if pattern == "*" {
		return true
	}

	// 层级通配符：sys.* 匹配 sys 及其子 scope
	if prefix, found := strings.CutSuffix(pattern, ".*"); found {
		return scope == prefix || strings.HasPrefix(scope, prefix+".")
	}

	// 精确匹配
	return pattern == scope
}
