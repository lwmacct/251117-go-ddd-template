package role

// GetQuery 获取角色查询
type GetQuery struct {
	RoleID uint
}

// ListQuery 列出角色查询
type ListQuery struct {
	Page  int
	Limit int
}

// ListPermissionsQuery 列出权限查询
type ListPermissionsQuery struct {
	Page  int
	Limit int
}
