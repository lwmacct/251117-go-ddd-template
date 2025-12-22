// Package command 定义个人访问令牌 (PAT) 的写操作命令。
//
// 本包处理 PAT 的生命周期管理：
//   - CreateTokenCommand: 创建新 PAT（返回明文 Token，仅此一次）
//   - DeleteTokenCommand: 删除 PAT
//   - DisableTokenCommand: 禁用 PAT
//   - EnableTokenCommand: 启用已禁用的 PAT
//
// 安全特性：
//   - PlainToken 仅在创建时返回，之后无法再次获取
//   - 支持过期时间和 IP 白名单配置
//   - 删除/禁用操作立即生效
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
