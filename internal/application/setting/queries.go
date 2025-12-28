package setting

// ==================== Category Queries ====================

// GetCategoryQuery 获取单个配置分类查询
type GetCategoryQuery struct {
	ID uint
}

// ListCategoriesQuery 获取配置分类列表查询
type ListCategoriesQuery struct{}

// ==================== Setting Queries ====================

// GetQuery 获取配置查询
type GetQuery struct {
	Key string
}

// ListQuery 获取配置列表查询
type ListQuery struct {
	CategoryID uint // 可选: 按类别 ID 过滤
}

// ListSchemaQuery 获取配置 Schema 查询（系统配置）
// Schema 返回按 Category → Group → Settings 层级组织的数据
type ListSchemaQuery struct{}

// ==================== UserSetting Queries ====================

// UserGetQuery 获取用户配置查询（合并默认值）
type UserGetQuery struct {
	UserID uint
	Key    string
}

// UserListQuery 获取用户配置列表查询（合并默认值）
type UserListQuery struct {
	UserID     uint
	CategoryID uint // 可选: 按类别 ID 过滤
}

// UserListSchemaQuery 获取用户配置 Schema 查询（带合并值）
type UserListSchemaQuery struct {
	UserID uint
}
