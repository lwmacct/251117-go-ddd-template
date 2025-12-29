package operation

// ============================================================================
// 派生查询函数
// ============================================================================

// PermissionDefinition 权限定义（替代 permission.Definition）
type PermissionDefinition struct {
	Code        string `json:"code"`        // 权限代码，如 admin:users:create
	Domain      string `json:"domain"`      // 域，如 admin
	Resource    string `json:"resource"`    // 资源，如 users
	Action      string `json:"action"`      // 操作，如 create
	Description string `json:"description"` // 英文描述
}

// AllPermissions 返回所有权限定义，供 Seeder 使用。
// 从注册表派生，仅返回有权限定义的操作。
func AllPermissions() []PermissionDefinition {
	seen := make(map[string]bool)
	perms := make([]PermissionDefinition, 0, len(operationRegistry))

	for _, meta := range operationRegistry {
		if meta.Permission == "" {
			continue
		}
		if seen[meta.Permission] {
			continue
		}
		seen[meta.Permission] = true

		// 解析权限代码：domain:resource:action
		parts := splitPermission(meta.Permission)
		perms = append(perms, PermissionDefinition{
			Code:        meta.Permission,
			Domain:      parts[0],
			Resource:    parts[1],
			Action:      parts[2],
			Description: meta.Description,
		})
	}

	return perms
}

// splitPermission 解析权限代码为 [domain, resource, action]
func splitPermission(code string) [3]string {
	var result [3]string
	idx := 0
	start := 0
	for i, c := range code {
		if c == ':' {
			if idx < 3 {
				result[idx] = code[start:i]
				idx++
				start = i + 1
			}
		}
	}
	if idx < 3 && start < len(code) {
		result[idx] = code[start:]
	}
	return result
}

// AuditActionDefinition 审计操作定义，供前端动态选项使用。
type AuditActionDefinition struct {
	Action      string         `json:"action"`       // 审计操作标识，如 user.create
	Operation   AuditOperation `json:"operation"`    // 操作类型，如 create
	Category    string         `json:"category"`     // 分类，如 user
	Label       string         `json:"label"`        // 中文标签
	Description string         `json:"description"`  // 英文描述
	OperationID string         `json:"operation_id"` // API 操作标识
}

// AllAuditActions 返回所有审计操作定义。
// 从注册表派生，仅返回有审计定义的操作。
func AllAuditActions() []AuditActionDefinition {
	actions := make([]AuditActionDefinition, 0, len(operationRegistry))

	for op, meta := range operationRegistry {
		if meta.AuditAction == "" {
			continue
		}
		actions = append(actions, AuditActionDefinition{
			Action:      meta.AuditAction,
			Operation:   meta.AuditOperation,
			Category:    meta.AuditCategory,
			Label:       meta.Label,
			Description: meta.Description,
			OperationID: string(op),
		})
	}

	return actions
}

// CategoryOption 分类选项（用于前端下拉框）
type CategoryOption struct {
	Value string `json:"value"` // 分类值
	Label string `json:"label"` // 显示标签
}

// AllAuditCategories 返回所有审计分类选项。
func AllAuditCategories() []CategoryOption {
	seen := make(map[string]bool)
	categories := make([]CategoryOption, 0, 16) // 分类数量有限

	categoryLabels := map[string]string{
		"auth":         "认证",
		"user":         "用户",
		"role":         "角色",
		"menu":         "菜单",
		"setting":      "配置",
		"cache":        "缓存",
		"profile":      "个人资料",
		"token":        "访问令牌",
		"user_setting": "用户配置",
	}

	for _, meta := range operationRegistry {
		if meta.AuditCategory == "" || seen[meta.AuditCategory] {
			continue
		}
		seen[meta.AuditCategory] = true

		label := categoryLabels[meta.AuditCategory]
		if label == "" {
			label = meta.AuditCategory
		}
		categories = append(categories, CategoryOption{
			Value: meta.AuditCategory,
			Label: label,
		})
	}

	return categories
}

// OperationTypeOption 操作类型选项
type OperationTypeOption struct {
	Value string `json:"value"` // 操作类型值
	Label string `json:"label"` // 显示标签
}

// AllAuditOperations 返回所有审计操作类型选项。
func AllAuditOperations() []OperationTypeOption {
	return []OperationTypeOption{
		{Value: string(AuditOpCreate), Label: "创建"},
		{Value: string(AuditOpUpdate), Label: "更新"},
		{Value: string(AuditOpDelete), Label: "删除"},
		{Value: string(AuditOpAccess), Label: "访问"},
		{Value: string(AuditOpAuthenticate), Label: "认证"},
	}
}

// ByOperationID 通过操作标识符查找操作。
// 如果未找到返回空 Operation。
func ByOperationID(id string) Operation {
	op := Operation(id)
	if op.Valid() {
		return op
	}
	return ""
}

// All 返回所有操作。
func All() []Operation {
	ops := make([]Operation, 0, len(operationRegistry))
	for op := range operationRegistry {
		ops = append(ops, op)
	}
	return ops
}
