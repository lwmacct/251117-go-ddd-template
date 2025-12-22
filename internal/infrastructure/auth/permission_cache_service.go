package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"github.com/redis/go-redis/v9"
)

// UserPermissions 用户权限缓存数据结构
type UserPermissions struct {
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// PermissionCacheService 权限缓存服务
// 使用 Redis 缓存用户权限信息，提升高并发场景下的性能
// 缓存失效策略：5分钟 TTL + 权限变更时主动清除
type PermissionCacheService struct {
	redis         *redis.Client
	userQueryRepo user.QueryRepository
	keyPrefix     string        // Redis key 前缀
	cacheTTL      time.Duration // 缓存过期时间
}

// NewPermissionCacheService 创建权限缓存服务
func NewPermissionCacheService(
	redisClient *redis.Client,
	userQueryRepo user.QueryRepository,
	keyPrefix string,
) *PermissionCacheService {
	return &PermissionCacheService{
		redis:         redisClient,
		userQueryRepo: userQueryRepo,
		keyPrefix:     keyPrefix,
		cacheTTL:      5 * time.Minute, // 5分钟过期
	}
}

// GetUserPermissions 获取用户权限（先查缓存，未命中则查数据库）
func (s *PermissionCacheService) GetUserPermissions(ctx context.Context, userID uint) ([]string, []string, error) {
	// 1. 尝试从 Redis 读取缓存
	cached, err := s.getFromCache(ctx, userID)
	if err == nil && cached != nil {
		return cached.Roles, cached.Permissions, nil
	}

	// 2. 缓存未命中，查询数据库
	u, err := s.userQueryRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	roles := u.GetRoleNames()
	permissions := u.GetPermissionCodes()

	// 3. 写入 Redis 缓存（异步写入，不阻塞请求）
	go func(cacheUserID uint, cacheRoles, cachePerms []string) { //nolint:contextcheck // 异步缓存写入使用独立超时
		cacheCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := s.setToCache(cacheCtx, cacheUserID, cacheRoles, cachePerms); err != nil {
			slog.Warn("Failed to cache user permissions",
				"user_id", cacheUserID,
				"error", err,
			)
		}
	}(userID, roles, permissions)

	return roles, permissions, nil
}

// InvalidateUser 清除指定用户的权限缓存
// 用于用户角色变更、用户状态变更等场景
func (s *PermissionCacheService) InvalidateUser(ctx context.Context, userID uint) error {
	key := s.getCacheKey(userID)

	if err := s.redis.Del(ctx, key).Err(); err != nil {
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
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, s.getCacheKey(userID))
	}

	if len(keys) > 0 {
		if err := s.redis.Del(ctx, keys...).Err(); err != nil {
			slog.Warn("Failed to invalidate users cache",
				"role_id", roleID,
				"user_count", len(keys),
				"error", err,
			)
			return err
		}

		slog.Info("Users permissions cache invalidated",
			"role_id", roleID,
			"user_count", len(keys),
		)
	}

	return nil
}

// InvalidateAllUsers 清除所有用户权限缓存（慎用，仅用于全局权限系统变更）
func (s *PermissionCacheService) InvalidateAllUsers(ctx context.Context) error {
	pattern := s.keyPrefix + "user:perms:*"

	// 使用 SCAN 命令遍历所有匹配的 key
	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}

	// 批量删除
	if len(keys) > 0 {
		if err := s.redis.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}

		slog.Info("All users permissions cache invalidated",
			"count", len(keys),
		)
	}

	return nil
}

// getFromCache 从 Redis 读取缓存
func (s *PermissionCacheService) getFromCache(ctx context.Context, userID uint) (*UserPermissions, error) {
	key := s.getCacheKey(userID)

	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss returns nil, nil
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var perms UserPermissions
	if err := json.Unmarshal(data, &perms); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.redis.Del(ctx, key)
		return nil, nil //nolint:nilerr,nilnil // intentionally treat corrupted cache as cache miss
	}

	return &perms, nil
}

// setToCache 写入 Redis 缓存
func (s *PermissionCacheService) setToCache(ctx context.Context, userID uint, roles, permissions []string) error {
	key := s.getCacheKey(userID)

	perms := UserPermissions{
		Roles:       roles,
		Permissions: permissions,
	}

	data, err := json.Marshal(perms)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %w", err)
	}

	if err := s.redis.Set(ctx, key, data, s.cacheTTL).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

// getCacheKey 生成 Redis 缓存 key
func (s *PermissionCacheService) getCacheKey(userID uint) string {
	return fmt.Sprintf("%suser:perms:%d", s.keyPrefix, userID)
}
