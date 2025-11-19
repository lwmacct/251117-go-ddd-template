// Package command 定义 PAT 命令处理器
package command

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

// RevokeTokenHandler 撤销 Token 命令处理器
type RevokeTokenHandler struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
}

// NewRevokeTokenHandler 创建 RevokeTokenHandler 实例
func NewRevokeTokenHandler(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
) *RevokeTokenHandler {
	return &RevokeTokenHandler{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
	}
}

// Handle 处理撤销 Token 命令
func (h *RevokeTokenHandler) Handle(ctx context.Context, cmd RevokeTokenCommand) error {
	// 1. 验证 Token 是否存在且属于该用户
	token, err := h.patQueryRepo.FindByID(ctx, cmd.TokenID)
	if err != nil || token == nil {
		return fmt.Errorf("token not found")
	}

	if token.UserID != cmd.UserID {
		return fmt.Errorf("token does not belong to this user")
	}

	// 2. 删除 Token
	if err := h.patCommandRepo.Delete(ctx, cmd.TokenID); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}
