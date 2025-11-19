// Package command 定义 PAT 命令处理器
package command

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
)

// CreateTokenHandler 创建 Token 命令处理器
type CreateTokenHandler struct {
	patCommandRepo pat.CommandRepository
	patQueryRepo   pat.QueryRepository
	tokenGenerator *auth.TokenGenerator
}

// NewCreateTokenHandler 创建 CreateTokenHandler 实例
func NewCreateTokenHandler(
	patCommandRepo pat.CommandRepository,
	patQueryRepo pat.QueryRepository,
	tokenGenerator *auth.TokenGenerator,
) *CreateTokenHandler {
	return &CreateTokenHandler{
		patCommandRepo: patCommandRepo,
		patQueryRepo:   patQueryRepo,
		tokenGenerator: tokenGenerator,
	}
}

// Handle 处理创建 Token 命令
func (h *CreateTokenHandler) Handle(ctx context.Context, cmd CreateTokenCommand) (*CreateTokenResult, error) {
	// 1. 生成 Token
	plainToken, hashedToken, _, err := h.tokenGenerator.GeneratePAT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 2. 解析过期时间
	var expiresAt *time.Time
	if cmd.ExpiresAt != nil && *cmd.ExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, *cmd.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("invalid expiration date format: %w", err)
		}
		expiresAt = &parsedTime
	}

	// 3. 创建 PAT 实体
	patEntity := &pat.PersonalAccessToken{
		UserID:      cmd.UserID,
		Name:        cmd.Name,
		Token:       hashedToken,
		Permissions: cmd.Permissions,
		ExpiresAt:   expiresAt,
		LastUsedAt:  nil,
	}

	// 4. 保存 Token
	if err := h.patCommandRepo.Create(ctx, patEntity); err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &CreateTokenResult{
		ID:          patEntity.ID,
		Token:       plainToken, // 返回明文 token（仅此一次）
		Name:        patEntity.Name,
		Permissions: patEntity.Permissions,
	}, nil
}
