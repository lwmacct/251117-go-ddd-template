// Package command 定义 PAT 命令处理器
package command

import (
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
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

// CreateTokenResult 创建 Token 结果
type CreateTokenResult struct {
	Token      *pat.PersonalAccessToken
	PlainToken string
}
