package pat

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// InternalCreateTokenResult Handler 内部返回类型（包含领域实体）
// 注意：这是内部类型，不应该直接序列化为 HTTP 响应
type InternalCreateTokenResult struct {
	Token      *pat.PersonalAccessToken
	PlainToken string
}

// CreateHandler 创建 Token 命令处理器
type CreateHandler struct {
	patCommandRepo pat.CommandRepository
	userQueryRepo  user.QueryRepository
	tokenGenerator auth.TokenGenerator
}

// NewCreateHandler 创建 CreateHandler 实例
func NewCreateHandler(
	patCommandRepo pat.CommandRepository,
	userQueryRepo user.QueryRepository,
	tokenGenerator auth.TokenGenerator,
) *CreateHandler {
	return &CreateHandler{
		patCommandRepo: patCommandRepo,
		userQueryRepo:  userQueryRepo,
		tokenGenerator: tokenGenerator,
	}
}

// Handle 处理创建 Token 命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*InternalCreateTokenResult, error) {
	if cmd.UserID == 0 {
		return nil, errors.New("user ID is required")
	}

	// 1. 拉取用户权限并校验
	u, err := h.userQueryRepo.GetByIDWithRoles(ctx, cmd.UserID)
	if err != nil || u == nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 新 RBAC 模型：权限为 Permission 对象数组
	userPerms := u.GetPermissions()
	if len(userPerms) == 0 {
		return nil, errors.New("user has no permissions")
	}

	// 转换为字符串格式供 PAT 存储
	userPermStrings := permissionsToStrings(userPerms)

	requestedPerms := cmd.Permissions
	if len(requestedPerms) == 0 {
		requestedPerms = userPermStrings // 默认继承全部权限
	}

	if err = validatePermissions(requestedPerms, userPermStrings); err != nil {
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

	return &InternalCreateTokenResult{
		Token:      patEntity,
		PlainToken: plainToken, // 返回明文 token（仅此一次）
	}, nil
}

// permissionsToStrings 将 Permission 对象转换为字符串格式
// 格式：operation_pattern|resource_pattern
func permissionsToStrings(perms []role.Permission) []string {
	result := make([]string, len(perms))
	for i, p := range perms {
		resPattern := p.ResourcePattern
		if resPattern == "" {
			resPattern = "*"
		}
		result[i] = p.OperationPattern + "|" + resPattern
	}
	return result
}

// validatePermissions ensures requested permissions are subset of user permissions.
// 使用模式匹配验证请求的权限是否被用户权限覆盖。
func validatePermissions(requested, userPerms []string) error {
	if len(requested) == 0 {
		return errors.New("at least one permission is required")
	}

	// 解析用户权限为 role.Permission 对象列表
	parsedUserPerms := make([]role.Permission, 0, len(userPerms))
	for _, perm := range userPerms {
		parts := strings.SplitN(perm, "|", 2)
		opPattern := parts[0]
		resPattern := "*"
		if len(parts) > 1 {
			resPattern = parts[1]
		}
		parsedUserPerms = append(parsedUserPerms, role.Permission{
			OperationPattern: opPattern,
			ResourcePattern:  resPattern,
		})
	}

	// 验证每个请求的权限是否被用户权限覆盖
	for _, reqPerm := range requested {
		// 解析请求的权限（格式: operation 或 operation|resource）
		parts := strings.SplitN(reqPerm, "|", 2)
		reqOp := parts[0]
		reqRes := "*"
		if len(parts) > 1 {
			reqRes = parts[1]
		}

		// 检查是否有任何用户权限能覆盖此请求
		matched := false
		for _, userPerm := range parsedUserPerms {
			if operation.MatchOperation(userPerm.OperationPattern, reqOp) &&
				operation.MatchResource(userPerm.ResourcePattern, reqRes) {
				matched = true
				break
			}
		}

		if !matched {
			return fmt.Errorf("permission '%s' is not granted to user", reqPerm)
		}
	}

	return nil
}
