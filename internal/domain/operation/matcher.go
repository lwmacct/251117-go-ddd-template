package operation

import "strings"

// ============================================================================
// Operation 匹配算法
// ============================================================================

// MatchOperation 匹配操作模式。
// 支持的模式：
//   - "sys:users.create"  精确匹配
//   - "sys:users.*"       模块通配（匹配 sys:users.* 所有动作）
//   - "sys:*"             域通配（匹配 sys 域所有操作）
//   - "*:*.read"          动作通配（匹配所有域所有模块的 read 动作）
//   - "*"                 超级通配（匹配所有操作）
func MatchOperation(pattern, operation string) bool {
	// 超级通配
	if pattern == "*" {
		return true
	}

	// 精确匹配
	if pattern == operation {
		return true
	}

	// 解析模式和操作
	patternParts := parseOperationParts(pattern)
	opParts := parseOperationParts(operation)

	// 逐部分匹配
	//nolint:gosec // G602 误报：固定大小数组 [3]string 与 range 3 循环配合使用
	for i := range 3 {
		if patternParts[i] != "*" && patternParts[i] != opParts[i] {
			return false
		}
	}
	return true
}

// parseOperationParts 解析操作标识符为三部分。
// "sys:users.create" → ["sys", "users", "create"]
// "sys:users.*" → ["sys", "users", "*"]
// "sys:*" → ["sys", "*", "*"]
// "*" → ["*", "*", "*"]
func parseOperationParts(op string) [3]string {
	// 处理超级通配
	if op == "*" {
		return [3]string{"*", "*", "*"}
	}

	// 按 : 分割域
	domain, rest, found := strings.Cut(op, ":")
	if !found {
		// 无域分隔符，视为无效格式
		return [3]string{op, "*", "*"}
	}

	// 处理域通配 "sys:*"
	if rest == "*" {
		return [3]string{domain, "*", "*"}
	}

	// 按 . 分割模块和动作
	module, action, _ := strings.Cut(rest, ".")

	return [3]string{domain, module, action}
}

// ============================================================================
// Resource 匹配算法
// ============================================================================

// MatchResource 匹配资源模式。
// 支持的模式：
//   - "user/123"   精确匹配
//   - "user/*"     类型通配（匹配所有 user 类型资源）
//   - "*"          超级通配（匹配所有资源）
//
// 特殊处理：
//   - "user/self" 模式在实际检查时需要替换为当前用户 ID
func MatchResource(pattern, resource string) bool {
	// 超级通配
	if pattern == "*" {
		return true
	}

	// 精确匹配
	if pattern == resource {
		return true
	}

	// 解析模式和资源
	patternType, patternID := parseResourceParts(pattern)
	resType, resID := parseResourceParts(resource)

	// 类型匹配
	if patternType != "*" && patternType != resType {
		return false
	}

	// ID 匹配（* 匹配任意）
	if patternID != "*" && patternID != resID {
		return false
	}

	return true
}

// parseResourceParts 解析资源标识符为类型和 ID。
// "user/123" → ("user", "123")
// "user/*" → ("user", "*")
// "*" → ("*", "*")
func parseResourceParts(res string) (string, string) {
	if res == "*" {
		return "*", "*"
	}

	idx := strings.Index(res, "/")
	if idx == -1 {
		return res, "*"
	}

	return res[:idx], res[idx+1:]
}

// ============================================================================
// Operation 扩展方法
// ============================================================================

// Domain 返回操作的域。
// "sys:users.create" → "sys"
func (o Operation) Domain() string {
	parts := parseOperationParts(string(o))
	return parts[0]
}

// Module 返回操作的模块。
// "sys:users.create" → "users"
func (o Operation) Module() string {
	parts := parseOperationParts(string(o))
	return parts[1]
}

// Action 返回操作的动作。
// "sys:users.create" → "create"
func (o Operation) Action() string {
	parts := parseOperationParts(string(o))
	return parts[2]
}

// Matches 检查操作是否匹配指定模式。
func (o Operation) Matches(pattern string) bool {
	return MatchOperation(pattern, string(o))
}
