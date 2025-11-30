// Package command 定义用户模块的写操作命令。
package command

// BatchCreateUsersCommand 批量创建用户命令
type BatchCreateUsersCommand struct {
	Users []BatchUserItem
}

// BatchUserItem 批量创建中的单个用户数据
type BatchUserItem struct {
	Username string
	Email    string
	Password string
	FullName string
	Status   string // active, inactive
	RoleIDs  []uint
}

// BatchCreateUsersResult 批量创建用户结果
type BatchCreateUsersResult struct {
	Total   int                    // 总数
	Success int                    // 成功数
	Failed  int                    // 失败数
	Errors  []BatchCreateUserError // 失败详情
}

// BatchCreateUserError 单个用户创建失败的错误信息
type BatchCreateUserError struct {
	Index    int    // 在批量列表中的索引（从 0 开始）
	Username string // 用户名
	Email    string // 邮箱
	Error    string // 错误原因
}
