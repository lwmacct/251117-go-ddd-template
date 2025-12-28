package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// PermissionCacheService 权限缓存服务（Cache-Aside 模式）
//
// 本服务封装 Cache-Aside 逻辑，供中间件使用：
//   - GetUserPermissions: 先查缓存，未命中则查数据库并回写缓存
//   - InvalidateUsersWithRole: 按角色批量失效（需查询数据库获取用户列表）
//
// 性能优化：
// 查询数据库时，会同时写入用户实体缓存（UserWithRolesCacheService），
// 供后续 Handler 通过 Repository 缓存装饰器直接命中，避免重复查询。
//
// 底层缓存操作委托给 [cache.PermissionCacheService] 接口实现。
type PermissionCacheService struct {
	cache         cache.PermissionCacheService
	userCache     cache.UserWithRolesCacheService
	userQueryRepo user.QueryRepository
}

// NewPermissionCacheService 创建权限缓存服务
func NewPermissionCacheService(
	cacheService cache.PermissionCacheService,
	userCacheService cache.UserWithRolesCacheService,
	userQueryRepo user.QueryRepository,
) *PermissionCacheService {
	return &PermissionCacheService{
		cache:         cacheService,
		userCache:     userCacheService,
		userQueryRepo: userQueryRepo,
	}
}

// GetUserPermissions 获取用户权限（Cache-Aside 模式）
//
// 执行流程：
//  1. 尝试从缓存读取
//  2. 缓存未命中，查询数据库
//  3. 同步写入权限缓存
//  4. 同步写入用户实体缓存（供后续 Handler 通过 Repository 装饰器命中）
func (s *PermissionCacheService) GetUserPermissions(ctx context.Context, userID uint) ([]string, []string, error) {
	// 1. 尝试从缓存读取
	roles, permissions, err := s.cache.GetUserPermissions(ctx, userID)
	if err == nil && (roles != nil || permissions != nil) {
		return roles, permissions, nil
	}

	// 2. 缓存未命中，查询数据库
	u, err := s.userQueryRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	roles = u.GetRoleNames()
	permissions = u.GetPermissionCodes()

	// 3. 同步写入权限缓存（Redis 写入 < 1ms，延迟可忽略）
	if err := s.cache.SetUserPermissions(ctx, userID, roles, permissions); err != nil {
		slog.Warn("Failed to cache user permissions",
			"user_id", userID,
			"error", err,
		)
	}

	// 4. 同步写入用户实体缓存（关键！供后续 Handler 通过 Repository 装饰器命中）
	if err := s.userCache.SetUserWithRoles(ctx, u); err != nil {
		slog.Warn("Failed to cache user entity",
			"user_id", userID,
			"error", err,
		)
	}

	return roles, permissions, nil
}

// InvalidateUser 清除指定用户的权限缓存
// 用于用户角色变更、用户状态变更等场景
func (s *PermissionCacheService) InvalidateUser(ctx context.Context, userID uint) error {
	if err := s.cache.InvalidateUser(ctx, userID); err != nil {
		slog.Warn("Failed to invalidate user cache",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	slog.Debug("User permissions cache invalidated",
		"user_id", userID,
	)

	return nil
}

// InvalidateUsersWithRole 清除拥有指定角色的所有用户缓存
// 用于角色权限变更场景（需要使所有拥有该角色的用户缓存失效）
func (s *PermissionCacheService) InvalidateUsersWithRole(ctx context.Context, roleID uint) error {
	// 查询所有拥有该角色的用户
	userIDs, err := s.userQueryRepo.GetUserIDsByRole(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get users with role %d: %w", roleID, err)
	}

	// 批量清除缓存
	if len(userIDs) > 0 {
		if err := s.cache.InvalidateUsers(ctx, userIDs); err != nil {
			slog.Warn("Failed to invalidate users cache",
				"role_id", roleID,
				"user_count", len(userIDs),
				"error", err,
			)
			return err
		}

		slog.Info("Users permissions cache invalidated",
			"role_id", roleID,
			"user_count", len(userIDs),
		)
	}

	return nil
}

// InvalidateAllUsers 清除所有用户权限缓存（慎用，仅用于全局权限系统变更）
func (s *PermissionCacheService) InvalidateAllUsers(ctx context.Context) error {
	if err := s.cache.InvalidateAll(ctx); err != nil {
		return fmt.Errorf("failed to invalidate all users cache: %w", err)
	}

	slog.Info("All users permissions cache invalidated")
	return nil
}
