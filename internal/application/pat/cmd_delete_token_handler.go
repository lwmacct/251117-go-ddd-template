package pat

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

// DeleteTokenHandler 删除 Token 命令处理器
type DeleteTokenHandler struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
}

// NewDeleteTokenHandler 创建 DeleteTokenHandler 实例
func NewDeleteTokenHandler(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
) *DeleteTokenHandler {
	return &DeleteTokenHandler{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
	}
}

// Handle 处理删除 Token 命令
func (h *DeleteTokenHandler) Handle(ctx context.Context, cmd DeleteTokenCommand) error {
	token, err := h.patQueryRepo.FindByID(ctx, cmd.TokenID)
	if err != nil || token == nil {
		return errors.New("token not found")
	}

	if token.UserID != cmd.UserID {
		return errors.New("token does not belong to this user")
	}

	if err := h.patCommandRepo.Delete(ctx, cmd.TokenID); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}
