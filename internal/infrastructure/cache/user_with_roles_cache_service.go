package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	appuser "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

const userWithRolesCacheTTL = 5 * time.Minute

// userWithRolesCacheDTO 用户实体缓存数据结构。
// 包含完整用户信息、角色和权限，避免直接序列化 Domain 实体。
type userWithRolesCacheDTO struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at,omitempty"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	FullName  string         `json:"full_name"`
	Avatar    string         `json:"avatar"`
	Bio       string         `json:"bio"`
	Status    string         `json:"status"`
	Roles     []roleCacheDTO `json:"roles"`
}

// roleCacheDTO 角色缓存数据结构。
type roleCacheDTO struct {
	ID          uint                 `json:"id"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Name        string               `json:"name"`
	DisplayName string               `json:"display_name"`
	Description string               `json:"description"`
	IsSystem    bool                 `json:"is_system"`
	Permissions []permissionCacheDTO `json:"permissions"`
}

// permissionCacheDTO 权限缓存数据结构。
// 新 RBAC 模型：存储 Operation + Resource Pattern。
type permissionCacheDTO struct {
	OperationPattern string `json:"operation_pattern"`
	ResourcePattern  string `json:"resource_pattern"`
}

// userWithRolesCacheService 用户实体缓存服务的 Redis 实现。
//
// 提供用户实体的缓存操作：
//   - Key 格式：{prefix}user:entity:{userID}
//   - TTL：5 分钟
//   - RedisJSON 存储
type userWithRolesCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewUserWithRolesCacheService 创建用户实体缓存服务。
func NewUserWithRolesCacheService(client *redis.Client, keyPrefix string) appuser.UserWithRolesCacheService {
	return &userWithRolesCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// GetUserWithRoles 获取缓存的用户实体（使用 RedisJSON）。
// 缓存未命中返回 nil, nil。
func (s *userWithRolesCacheService) GetUserWithRoles(ctx context.Context, userID uint) (*user.User, error) {
	data, err := s.client.JSONGet(ctx, s.buildKey(userID), "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil //nolint:nilnil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	var wrapper []userWithRolesCacheDTO
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, s.buildKey(userID))
		return nil, nil //nolint:nilnil,nilerr // corrupted cache treated as miss
	}

	if len(wrapper) == 0 {
		return nil, nil //nolint:nilnil // empty wrapper
	}

	return s.toEntity(&wrapper[0]), nil
}

// SetUserWithRoles 设置用户实体缓存（使用 RedisJSON）。
func (s *userWithRolesCacheService) SetUserWithRoles(ctx context.Context, u *user.User) error {
	dto := s.toDTO(u)

	// 使用 Pipeline 执行 JSON.SET + EXPIRE
	key := s.buildKey(u.ID)
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", dto)
	pipe.Expire(ctx, key, userWithRolesCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set user entity cache: %w", err)
	}

	return nil
}

// InvalidateUser 失效单个用户缓存。
func (s *userWithRolesCacheService) InvalidateUser(ctx context.Context, userID uint) error {
	return s.client.Del(ctx, s.buildKey(userID)).Err()
}

func (s *userWithRolesCacheService) buildKey(userID uint) string {
	return fmt.Sprintf("%suser:entity:%d", s.keyPrefix, userID)
}

// toDTO 将 Domain 实体转换为缓存 DTO。
func (s *userWithRolesCacheService) toDTO(u *user.User) *userWithRolesCacheDTO {
	dto := &userWithRolesCacheDTO{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Status:    u.Status,
		Roles:     make([]roleCacheDTO, 0, len(u.Roles)),
	}

	for _, r := range u.Roles {
		roleDTO := roleCacheDTO{
			ID:          r.ID,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
			Name:        r.Name,
			DisplayName: r.DisplayName,
			Description: r.Description,
			IsSystem:    r.IsSystem,
			Permissions: make([]permissionCacheDTO, 0, len(r.Permissions)),
		}

		for _, p := range r.Permissions {
			roleDTO.Permissions = append(roleDTO.Permissions, permissionCacheDTO{
				OperationPattern: p.OperationPattern,
				ResourcePattern:  p.ResourcePattern,
			})
		}

		dto.Roles = append(dto.Roles, roleDTO)
	}

	return dto
}

// toEntity 将缓存 DTO 转换为 Domain 实体。
func (s *userWithRolesCacheService) toEntity(dto *userWithRolesCacheDTO) *user.User {
	u := &user.User{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		DeletedAt: dto.DeletedAt,
		Username:  dto.Username,
		Email:     dto.Email,
		FullName:  dto.FullName,
		Avatar:    dto.Avatar,
		Bio:       dto.Bio,
		Status:    dto.Status,
		Roles:     make([]role.Role, 0, len(dto.Roles)),
	}

	for _, r := range dto.Roles {
		roleEntity := role.Role{
			ID:          r.ID,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
			Name:        r.Name,
			DisplayName: r.DisplayName,
			Description: r.Description,
			IsSystem:    r.IsSystem,
			Permissions: make([]role.Permission, 0, len(r.Permissions)),
		}

		for _, p := range r.Permissions {
			roleEntity.Permissions = append(roleEntity.Permissions, role.Permission{
				OperationPattern: p.OperationPattern,
				ResourcePattern:  p.ResourcePattern,
			})
		}

		u.Roles = append(u.Roles, roleEntity)
	}

	return u
}

var _ appuser.UserWithRolesCacheService = (*userWithRolesCacheService)(nil)
