package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// CreateTokenHandler 创建 Token 命令处理器
type CreateTokenHandler struct {
	patCommandRepo pat.CommandRepository
	userQueryRepo  user.QueryRepository
	tokenGenerator auth.TokenGenerator
}

// NewCreateTokenHandler 创建 CreateTokenHandler 实例
func NewCreateTokenHandler(
	patCommandRepo pat.CommandRepository,
	userQueryRepo user.QueryRepository,
	tokenGenerator auth.TokenGenerator,
) *CreateTokenHandler {
	return &CreateTokenHandler{
		patCommandRepo: patCommandRepo,
		userQueryRepo:  userQueryRepo,
		tokenGenerator: tokenGenerator,
	}
}

// Handle 处理创建 Token 命令
func (h *CreateTokenHandler) Handle(ctx context.Context, cmd CreateTokenCommand) (*CreateTokenResult, error) {
	if cmd.UserID == 0 {
		return nil, errors.New("user ID is required")
	}

	// 1. 拉取用户权限并校验
	u, err := h.userQueryRepo.GetByIDWithRoles(ctx, cmd.UserID)
	if err != nil || u == nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	userPerms := u.GetPermissionCodes()
	if len(userPerms) == 0 {
		return nil, errors.New("user has no permissions")
	}

	requestedPerms := cmd.Permissions
	if len(requestedPerms) == 0 {
		requestedPerms = userPerms // 默认继承全部权限
	}

	if err = validatePermissions(requestedPerms, userPerms); err != nil {
		return nil, err
	}

	// 1. 生成 Token
	plainToken, hashedToken, prefix, err := h.tokenGenerator.GeneratePAT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 2. 创建 PAT 实体
	patEntity := &pat.PersonalAccessToken{
		UserID:      cmd.UserID,
		Name:        cmd.Name,
		Token:       hashedToken,
		TokenPrefix: prefix,
		Permissions: requestedPerms,
		ExpiresAt:   cmd.ExpiresAt,
		LastUsedAt:  nil,
		Status:      "active",
		IPWhitelist: cmd.IPWhitelist,
		Description: cmd.Description,
	}

	// 3. 保存 Token
	if err := h.patCommandRepo.Create(ctx, patEntity); err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &CreateTokenResult{
		Token:      patEntity,
		PlainToken: plainToken, // 返回明文 token（仅此一次）
	}, nil
}

// validatePermissions ensures requested permissions are subset of user permissions.
func validatePermissions(requested, userPerms []string) error {
	if len(requested) == 0 {
		return errors.New("at least one permission is required")
	}

	permSet := make(map[string]struct{}, len(userPerms))
	for _, perm := range userPerms {
		permSet[perm] = struct{}{}
	}

	for _, perm := range requested {
		if _, ok := permSet[perm]; !ok {
			return fmt.Errorf("permission '%s' is not granted to user", perm)
		}
	}

	return nil
}
