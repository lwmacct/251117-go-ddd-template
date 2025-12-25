package pat

import (
	"time"
)

// CreateTokenCommand 创建 Token 命令
type CreateTokenCommand struct {
	UserID      uint
	Name        string
	Permissions []string
	ExpiresAt   *time.Time
	IPWhitelist []string
	Description string
}
