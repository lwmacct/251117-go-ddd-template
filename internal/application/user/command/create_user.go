// Package command 定义用户模块的写操作命令。
//
// 本包实现 CQRS 模式中的 Command 端，包含：
//   - CreateUserCommand: 创建用户（支持同时分配角色）
//   - UpdateUserCommand: 更新用户资料
//   - DeleteUserCommand: 删除用户（软删除）
//   - AssignRolesCommand: 分配/替换用户角色
//   - ChangePasswordCommand: 修改密码
//
// 每个 Command 对应一个 Handler，Handler 负责：
//   - 参数验证
//   - 业务规则检查
//   - 调用 Repository 执行持久化
//   - 返回操作结果
package command

// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
	Username string
	Email    string
	Password string
	FullName string
	RoleIDs  []uint // 可选：创建时分配角色
}
