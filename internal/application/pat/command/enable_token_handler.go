// Package command 定义 PAT 命令处理器
package command

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

// EnableTokenHandler 启用 Token 命令处理器
type EnableTokenHandler struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
}

// NewEnableTokenHandler 创建 EnableTokenHandler 实例
func NewEnableTokenHandler(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
) *EnableTokenHandler {
	return &EnableTokenHandler{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
	}
}

// Handle 处理启用 Token 命令
func (h *EnableTokenHandler) Handle(ctx context.Context, cmd EnableTokenCommand) error {
	token, err := h.patQueryRepo.FindByID(ctx, cmd.TokenID)
	if err != nil || token == nil {
		return fmt.Errorf("token not found")
	}

	if token.UserID != cmd.UserID {
		return fmt.Errorf("token does not belong to this user")
	}

	if token.IsExpired() {
		return fmt.Errorf("token is expired and cannot be enabled")
	}

	if err := h.patCommandRepo.Enable(ctx, cmd.TokenID); err != nil {
		return fmt.Errorf("failed to enable token: %w", err)
	}

	return nil
}
