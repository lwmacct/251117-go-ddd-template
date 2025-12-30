package operation

// ============================================================================
// 派生查询函数
// ============================================================================

// OperationDefinition 操作定义，供前端权限配置使用。
type OperationDefinition struct {
	Code        string `json:"code"`        // 操作代码，如 sys:users:create
	Scope       string `json:"scope"`       // Scope，如 sys
	Type        string `json:"type"`        // 类型，如 users
	Identifier  string `json:"identifier"`  // 标识符，如 create
	Label       string `json:"label"`       // 中文标签
	Description string `json:"description"` // 英文描述
	Group       string `json:"group"`       // Swagger 分组
}

// AllOperationDefinitions 返回所有操作定义，供前端权限配置使用。
// 仅返回非公开操作（需要权限检查的操作）。
func AllOperationDefinitions() []OperationDefinition {
	ops := make([]OperationDefinition, 0, len(operationRegistry))

	for op, meta := range operationRegistry {
		// 跳过公开操作（通过 scope 判断）
		if IsPublic(op) {
			continue
		}

		ops = append(ops, OperationDefinition{
			Code:        string(op),
			Scope:       op.Scope(),
			Type:        op.Type(),
			Identifier:  op.Identifier(),
			Label:       meta.Label,
			Description: meta.Description,
			Group:       meta.Group,
		})
	}

	return ops
}

// AuditActionDefinition 审计操作定义，供前端动态选项使用。
type AuditActionDefinition struct {
	Action      string         `json:"action"`       // 审计操作标识，如 user.create
	Operation   AuditOperation `json:"operation"`    // 操作类型，如 create
	Category    AuditCategory  `json:"category"`     // 分类，如 user
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
	seen := make(map[AuditCategory]bool)
	categories := make([]CategoryOption, 0, 16)

	for _, meta := range operationRegistry {
		if meta.AuditCategory == "" || seen[meta.AuditCategory] {
			continue
		}
		seen[meta.AuditCategory] = true
		categories = append(categories, CategoryOption{
			Value: string(meta.AuditCategory),
			Label: meta.AuditCategory.Label(),
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
	ops := []AuditOperation{
		AuditOpCreate, AuditOpUpdate, AuditOpDelete,
		AuditOpAccess, AuditOpAuthenticate,
	}
	result := make([]OperationTypeOption, len(ops))
	for i, op := range ops {
		result[i] = OperationTypeOption{
			Value: string(op),
			Label: op.Label(),
		}
	}
	return result
}

// ByOperationID 通过操作标识符查找操作。
// 如果未找到返回空 Operation。
func ByOperationID(id string) Operation {
	op := Operation(id)
	if Valid(op) {
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

// ByMethodAndPath 通过 HTTP 方法和路径查找操作。
// 支持路径参数匹配：/api/system/users/:id 匹配 /api/system/users/123
// 如果未找到返回空 Operation。
func ByMethodAndPath(method HTTPMethod, path string) Operation {
	for op, meta := range operationRegistry {
		if meta.Method != method {
			continue
		}
		if matchPath(meta.Path, path) {
			return op
		}
	}
	return ""
}

// matchPath 检查实际路径是否匹配模式路径。
// pattern: /api/system/users/:id（Gin 风格路径参数）
// actual:  /api/system/users/123
func matchPath(pattern, actual string) bool {
	patternSegs := splitPathSegments(pattern)
	actualSegs := splitPathSegments(actual)

	if len(patternSegs) != len(actualSegs) {
		return false
	}

	for i, seg := range patternSegs {
		// :param 匹配任意非空值
		if len(seg) > 0 && seg[0] == ':' {
			continue
		}
		if seg != actualSegs[i] {
			return false
		}
	}

	return true
}

// splitPathSegments 将路径分割为段。
func splitPathSegments(path string) []string {
	if len(path) == 0 {
		return nil
	}

	// 移除开头的斜杠
	if path[0] == '/' {
		path = path[1:]
	}

	// 移除结尾的斜杠
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if len(path) == 0 {
		return nil
	}

	// 分割路径
	var segments []string
	start := 0
	for i := range len(path) {
		if path[i] == '/' {
			if i > start {
				segments = append(segments, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		segments = append(segments, path[start:])
	}

	return segments
}
