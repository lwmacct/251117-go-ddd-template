package operation

import "github.com/lwmacct/251117-go-ddd-template/pkg/urn"

// ============================================================================
// Operation 匹配算法
// ============================================================================

// MatchOperation 匹配操作模式（委托给 URN 包）。
// 支持的模式：
//   - "sys:users:create"   精确匹配
//   - "sys:users:*"        模块通配（匹配 sys:users:* 所有动作）
//   - "sys:*:*"            域通配（匹配 sys 域所有操作）
//   - "*:*:create"         动作通配（匹配所有域所有模块的 create 动作）
//   - "*:*:*" 或 "*"       超级通配（匹配所有操作）
//   - "sys.*:*:*"          子域通配（匹配 sys 及其子域）
func MatchOperation(pattern, operation string) bool {
	return urn.MatchOperation(pattern, operation)
}

// ============================================================================
// Resource 匹配算法
// ============================================================================

// MatchResource 匹配资源模式（委托给 URN 包）。
// 支持的模式：
//   - "sys:user:123"       精确匹配
//   - "sys:user:*"         类型通配（匹配所有 sys:user 资源）
//   - "sys:*:*"            域通配（匹配 sys 域所有资源）
//   - "*:*:*" 或 "*"       超级通配（匹配所有资源）
//   - "org.acme.*:*:*"     子域通配
//
// 特殊处理：
//   - "self:user:@me" 模式在实际检查时需要使用 [urn.Resolver] 替换为当前用户 ID
func MatchResource(pattern, resource string) bool {
	return urn.MatchResource(pattern, resource)
}
