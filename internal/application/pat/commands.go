package pat

import (
	"time"
)

// CreateCommand 创建 Token 命令
type CreateCommand struct {
	UserID      uint
	Name        string
	Permissions []string
	ExpiresAt   *time.Time
	IPWhitelist []string
	Description string
}

// DeleteCommand 删除 Token 命令
type DeleteCommand struct {
	UserID  uint
	TokenID uint
}

// DisableCommand 禁用 Token 命令
type DisableCommand struct {
	UserID  uint
	TokenID uint
}

// EnableCommand 启用 Token 命令
type EnableCommand struct {
	UserID  uint
	TokenID uint
}
