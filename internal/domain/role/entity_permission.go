package role

import "time"

// Permission represents a permission in the RBAC system
// Uses three-part format: domain:resource:action
type Permission struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
	Domain      string     `json:"domain"`
	Resource    string     `json:"resource"`
	Action      string     `json:"action"`
	Description string     `json:"description"`
	Code        string     `json:"code"`
}

// PermissionCode returns the full permission code in format "domain:resource:action"
func (p *Permission) PermissionCode() string {
	return p.Code
}

// IsValid 检查权限是否有效（所有字段都不为空）
func (p *Permission) IsValid() bool {
	return p.Domain != "" && p.Resource != "" && p.Action != "" && p.Code != ""
}

// GetComponents 获取权限的三个组成部分
//
//nolint:nonamedreturns // named returns for self-documenting API
func (p *Permission) GetComponents() (domain, resource, action string) {
	return p.Domain, p.Resource, p.Action
}

// Matches 检查权限是否匹配指定模式（支持通配符 *）
// 例如: "user:*:*" 匹配所有 user 域的权限
func (p *Permission) Matches(pattern string) bool {
	if pattern == "*" || pattern == "*:*:*" {
		return true
	}
	if pattern == p.Code {
		return true
	}

	// 简单通配符匹配
	parts := splitPermissionCode(pattern)
	if len(parts) != 3 {
		return false
	}

	return matchPart(parts[0], p.Domain) &&
		matchPart(parts[1], p.Resource) &&
		matchPart(parts[2], p.Action)
}

// BuildCode 根据 Domain/Resource/Action 构建权限代码
func (p *Permission) BuildCode() string {
	return p.Domain + ":" + p.Resource + ":" + p.Action
}

// splitPermissionCode 分割权限代码
func splitPermissionCode(code string) []string {
	result := make([]string, 0, 3)
	start := 0
	for i, c := range code {
		if c == ':' {
			result = append(result, code[start:i])
			start = i + 1
		}
	}
	result = append(result, code[start:])
	return result
}

// matchPart 检查模式部分是否匹配
func matchPart(pattern, value string) bool {
	return pattern == "*" || pattern == value
}
