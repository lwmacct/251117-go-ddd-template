package query

import (
	"context"
	"errors"

	patApp "github.com/lwmacct/251117-go-ddd-template/internal/application/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

// GetTokenHandler 获取 Token 查询处理器
type GetTokenHandler struct {
	patQueryRepo pat.QueryRepository
}

// NewGetTokenHandler 创建 GetTokenHandler 实例
func NewGetTokenHandler(patQueryRepo pat.QueryRepository) *GetTokenHandler {
	return &GetTokenHandler{
		patQueryRepo: patQueryRepo,
	}
}

// Handle 处理获取 Token 查询
func (h *GetTokenHandler) Handle(ctx context.Context, query GetTokenQuery) (*patApp.TokenInfoResponse, error) {
	// 1. 查询 Token
	token, err := h.patQueryRepo.FindByID(ctx, query.TokenID)
	if err != nil || token == nil {
		return nil, errors.New("token not found")
	}

	// 2. 验证所有权
	if token.UserID != query.UserID {
		return nil, errors.New("token does not belong to this user")
	}

	return patApp.ToTokenInfoResponse(token), nil
}
