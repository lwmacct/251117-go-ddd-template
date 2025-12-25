package twofa

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
)

// VerifyEnableHandler 验证并启用 2FA 命令处理器
type VerifyEnableHandler struct {
	twofaService twofa.Service
}

// NewVerifyEnableHandler 创建验证并启用 2FA 命令处理器
func NewVerifyEnableHandler(twofaService twofa.Service) *VerifyEnableHandler {
	return &VerifyEnableHandler{
		twofaService: twofaService,
	}
}

// Handle 处理验证并启用 2FA 命令
func (h *VerifyEnableHandler) Handle(ctx context.Context, cmd VerifyEnableCommand) (*VerifyEnableResultDTO, error) {
	recoveryCodes, err := h.twofaService.VerifyAndEnable(ctx, cmd.UserID, cmd.Code)
	if err != nil {
		return nil, err
	}

	return &VerifyEnableResultDTO{
		RecoveryCodes: recoveryCodes,
	}, nil
}
