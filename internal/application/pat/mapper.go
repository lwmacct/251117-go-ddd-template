package pat

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

// ToTokenResponse 将领域模型 PersonalAccessToken 转换为应用层 TokenResponse DTO
func ToTokenResponse(token *pat.PersonalAccessToken) *TokenResponse {
	if token == nil {
		return nil
	}

	return &TokenResponse{
		ID:          token.ID,
		UserID:      token.UserID,
		Name:        token.Name,
		TokenPrefix: token.TokenPrefix,
		Permissions: token.Permissions,
		IPWhitelist: token.IPWhitelist,
		Status:      token.Status,
		ExpiresAt:   token.ExpiresAt,
		LastUsedAt:  token.LastUsedAt,
		CreatedAt:   token.CreatedAt,
		UpdatedAt:   token.UpdatedAt,
	}
}

// ToCreateTokenResponse 将领域模型 PersonalAccessToken 转换为创建响应 DTO（携带一次性明文 token）
func ToCreateTokenResponse(token *pat.PersonalAccessToken, plainToken string) *CreateTokenResponse {
	if token == nil {
		return nil
	}

	return &CreateTokenResponse{
		PlainToken: plainToken,
		Token:      ToTokenResponse(token),
	}
}

// ToTokenListResponse 将领域模型 TokenListItem 数组转换为应用层 TokenListResponse DTO
func ToTokenListResponse(items []*pat.TokenListItem) *TokenListResponse {
	responses := make([]*TokenResponse, len(items))
	for i, item := range items {
		responses[i] = &TokenResponse{
			ID:          item.ID,
			Name:        item.Name,
			TokenPrefix: item.TokenPrefix,
			Permissions: item.Permissions,
			Status:      item.Status,
			ExpiresAt:   item.ExpiresAt,
			LastUsedAt:  item.LastUsedAt,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.CreatedAt,
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

	return ToTokenResponse(token)
}
