package twofa

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
)

// GetStatusHandler 获取 2FA 状态查询处理器
type GetStatusHandler struct {
	twofaService twofa.Service
}

// NewGetStatusHandler 创建获取 2FA 状态查询处理器
func NewGetStatusHandler(twofaService twofa.Service) *GetStatusHandler {
	return &GetStatusHandler{
		twofaService: twofaService,
	}
}

// Handle 处理获取 2FA 状态查询
func (h *GetStatusHandler) Handle(ctx context.Context, query GetStatusQuery) (*GetStatusResultDTO, error) {
	enabled, count, err := h.twofaService.GetStatus(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	return &GetStatusResultDTO{
		Enabled:            enabled,
		RecoveryCodesCount: count,
	}, nil
}
