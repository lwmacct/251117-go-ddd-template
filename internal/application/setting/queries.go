package setting

// GetQuery 获取设置查询
type GetQuery struct {
	Key string
}

// ListQuery 获取设置列表查询
type ListQuery struct {
	Category string // 可选: 按类别过滤
}

// ListSchemaQuery 获取设置 Schema 查询
// Schema 返回按 Category → Group → Settings 层级组织的数据
type ListSchemaQuery struct{}
