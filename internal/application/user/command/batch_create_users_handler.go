package command

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// BatchCreateUsersHandler 批量创建用户命令处理器
type BatchCreateUsersHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
	authService     domainAuth.Service
}

// NewBatchCreateUsersHandler 创建批量用户命令处理器
func NewBatchCreateUsersHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	authService domainAuth.Service,
) *BatchCreateUsersHandler {
	return &BatchCreateUsersHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
		authService:     authService,
	}
}

// Handle 处理批量创建用户命令
// 采用部分失败策略：单个用户失败不影响其他用户的创建
func (h *BatchCreateUsersHandler) Handle(ctx context.Context, cmd BatchCreateUsersCommand) (*BatchCreateUsersResult, error) {
	result := &BatchCreateUsersResult{
		Total:  len(cmd.Users),
		Errors: make([]BatchCreateUserError, 0),
	}

	for i, item := range cmd.Users {
		err := h.createSingleUser(ctx, item)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, BatchCreateUserError{
				Index:    i,
				Username: item.Username,
				Email:    item.Email,
				Error:    err.Error(),
			})
		} else {
			result.Success++
		}
	}

	return result, nil
}

// createSingleUser 创建单个用户（内部方法）
func (h *BatchCreateUsersHandler) createSingleUser(ctx context.Context, item BatchUserItem) error {
	// 1. 基本验证
	if err := h.validateUserItem(item); err != nil {
		return err
	}

	// 2. 验证密码策略
	if err := h.authService.ValidatePasswordPolicy(ctx, item.Password); err != nil {
		return fmt.Errorf("密码不符合策略: %w", err)
	}

	// 3. 检查用户名是否已存在
	exists, err := h.userQueryRepo.ExistsByUsername(ctx, item.Username)
	if err != nil {
		return fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 4. 检查邮箱是否已存在
	exists, err = h.userQueryRepo.ExistsByEmail(ctx, item.Email)
	if err != nil {
		return fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return errors.New("邮箱已存在")
	}

	// 5. 生成密码哈希
	hashedPassword, err := h.authService.GeneratePasswordHash(ctx, item.Password)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 6. 确定用户状态
	status := item.Status
	if status == "" {
		status = "active"
	}
	if status != "active" && status != "inactive" {
		return fmt.Errorf("无效的状态值: %s", status)
	}

	// 7. 创建用户实体
	newUser := &user.User{
		Username: item.Username,
		Email:    item.Email,
		Password: hashedPassword,
		FullName: item.FullName,
		Status:   status,
	}

	// 8. 保存用户
	if err := h.userCommandRepo.Create(ctx, newUser); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	// 9. 分配角色（如果提供）
	if len(item.RoleIDs) > 0 {
		if err := h.userCommandRepo.AssignRoles(ctx, newUser.ID, item.RoleIDs); err != nil {
			// 用户已创建，但角色分配失败，记录警告但不回滚
			return fmt.Errorf("用户已创建但角色分配失败: %w", err)
		}
	}

	return nil
}

// validateUserItem 验证单个用户数据
func (h *BatchCreateUsersHandler) validateUserItem(item BatchUserItem) error {
	// 用户名验证
	if strings.TrimSpace(item.Username) == "" {
		return errors.New("用户名不能为空")
	}
	if len(item.Username) < 3 || len(item.Username) > 50 {
		return errors.New("用户名长度必须在 3-50 个字符之间")
	}

	// 邮箱验证
	if strings.TrimSpace(item.Email) == "" {
		return errors.New("邮箱不能为空")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(item.Email) {
		return errors.New("邮箱格式无效")
	}

	// 密码验证
	if strings.TrimSpace(item.Password) == "" {
		return errors.New("密码不能为空")
	}

	return nil
}
