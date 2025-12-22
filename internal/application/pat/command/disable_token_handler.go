package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

// DisableTokenHandler 禁用 Token 命令处理器
type DisableTokenHandler struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
}

// NewDisableTokenHandler 创建 DisableTokenHandler 实例
func NewDisableTokenHandler(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
) *DisableTokenHandler {
	return &DisableTokenHandler{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
	}
}

// Handle 处理禁用 Token 命令
func (h *DisableTokenHandler) Handle(ctx context.Context, cmd DisableTokenCommand) error {
	token, err := h.patQueryRepo.FindByID(ctx, cmd.TokenID)
	if err != nil || token == nil {
		return errors.New("token not found")
	}

	if token.UserID != cmd.UserID {
		return errors.New("token does not belong to this user")
	}

	if err := h.patCommandRepo.Disable(ctx, cmd.TokenID); err != nil {
		return fmt.Errorf("failed to disable token: %w", err)
	}

	return nil
}
