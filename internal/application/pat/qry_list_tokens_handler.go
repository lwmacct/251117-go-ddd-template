package pat

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

// ListTokensHandler 获取 Token 列表查询处理器
type ListTokensHandler struct {
	patQueryRepo pat.QueryRepository
}

// NewListTokensHandler 创建 ListTokensHandler 实例
func NewListTokensHandler(patQueryRepo pat.QueryRepository) *ListTokensHandler {
	return &ListTokensHandler{
		patQueryRepo: patQueryRepo,
	}
}

// Handle 处理获取 Token 列表查询
func (h *ListTokensHandler) Handle(ctx context.Context, query ListTokensQuery) ([]*TokenInfoDTO, error) {
	tokens, err := h.patQueryRepo.ListByUser(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tokens: %w", err)
	}

	// 转换为 DTO
	tokenResponses := make([]*TokenInfoDTO, 0, len(tokens))
	for _, token := range tokens {
		tokenResponses = append(tokenResponses, ToTokenInfoDTO(token))
	}

	return tokenResponses, nil
}
