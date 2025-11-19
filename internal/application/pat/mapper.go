// Package pat 提供领域模型到应用层 DTO 的映射函数
package pat

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

// ToTokenResponse 将领域模型 PersonalAccessToken 转换为应用层 TokenResponse DTO
func ToTokenResponse(token *pat.PersonalAccessToken, plainToken string) *TokenResponse {
	if token == nil {
		return nil
	}

	return &TokenResponse{
		ID:          token.ID,
		Name:        token.Name,
		Token:       plainToken, // 仅在创建时传入，其他时候为空字符串
		Permissions: token.Permissions,
		ExpiresAt:   token.ExpiresAt,
		LastUsedAt:  token.LastUsedAt,
		CreatedAt:   token.CreatedAt,
	}
}

// ToTokenListResponse 将领域模型 TokenListItem 数组转换为应用层 TokenListResponse DTO
func ToTokenListResponse(items []*pat.TokenListItem) *TokenListResponse {
	responses := make([]*TokenResponse, len(items))
	for i, item := range items {
		responses[i] = &TokenResponse{
			ID:          item.ID,
			Name:        item.Name,
			Permissions: item.Permissions,
			ExpiresAt:   item.ExpiresAt,
			LastUsedAt:  item.LastUsedAt,
			CreatedAt:   item.CreatedAt,
		}
	}

	return &TokenListResponse{
		Tokens: responses,
		Total:  int64(len(responses)),
	}
}

// ToTokenInfoResponse 将领域实体转换为 TokenInfoResponse（不包含 token）
func ToTokenInfoResponse(token *pat.PersonalAccessToken) *TokenInfoResponse {
	if token == nil {
		return nil
	}

	return &TokenInfoResponse{
		ID:          token.ID,
		Name:        token.Name,
		Permissions: token.Permissions,
		ExpiresAt:   token.ExpiresAt,
		LastUsedAt:  token.LastUsedAt,
		CreatedAt:   token.CreatedAt,
	}
}
