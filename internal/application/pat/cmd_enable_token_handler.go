package pat

import (
	"context"
	"errors"
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
		return errors.New("token not found")
	}

	if token.UserID != cmd.UserID {
		return errors.New("token does not belong to this user")
	}

	if token.IsExpired() {
		return errors.New("token is expired and cannot be enabled")
	}

	if err := h.patCommandRepo.Enable(ctx, cmd.TokenID); err != nil {
		return fmt.Errorf("failed to enable token: %w", err)
	}

	return nil
}
