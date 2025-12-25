package events

import (
	"strconv"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
)

// ============================================================================
// 用户事件
// ============================================================================

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
	event.BaseEvent

	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// NewUserCreatedEvent 创建用户创建事件
func NewUserCreatedEvent(userID uint, username, email string) *UserCreatedEvent {
	return &UserCreatedEvent{
		BaseEvent: event.NewBaseEvent("user.created", "user", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
		Username:  username,
		Email:     email,
	}
}

// UserUpdatedEvent 用户更新事件
type UserUpdatedEvent struct {
	event.BaseEvent

	UserID uint `json:"user_id"`
}

// NewUserUpdatedEvent 创建用户更新事件
func NewUserUpdatedEvent(userID uint) *UserUpdatedEvent {
	return &UserUpdatedEvent{
		BaseEvent: event.NewBaseEvent("user.updated", "user", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
	}
}

// UserDeletedEvent 用户删除事件
type UserDeletedEvent struct {
	event.BaseEvent

	UserID uint `json:"user_id"`
}

// NewUserDeletedEvent 创建用户删除事件
func NewUserDeletedEvent(userID uint) *UserDeletedEvent {
	return &UserDeletedEvent{
		BaseEvent: event.NewBaseEvent("user.deleted", "user", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
	}
}

// UserRoleAssignedEvent 用户角色分配事件
type UserRoleAssignedEvent struct {
	event.BaseEvent

	UserID  uint   `json:"user_id"`
	RoleIDs []uint `json:"role_ids"`
}

// NewUserRoleAssignedEvent 创建用户角色分配事件
func NewUserRoleAssignedEvent(userID uint, roleIDs []uint) *UserRoleAssignedEvent {
	return &UserRoleAssignedEvent{
		BaseEvent: event.NewBaseEvent("user.role_assigned", "user", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
		RoleIDs:   roleIDs,
	}
}

// UserStatusChangedEvent 用户状态变更事件
type UserStatusChangedEvent struct {
	event.BaseEvent

	UserID    uint   `json:"user_id"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
}

// NewUserStatusChangedEvent 创建用户状态变更事件
func NewUserStatusChangedEvent(userID uint, oldStatus, newStatus string) *UserStatusChangedEvent {
	return &UserStatusChangedEvent{
		BaseEvent: event.NewBaseEvent("user.status_changed", "user", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}
}

// PasswordChangedEvent 密码变更事件
type PasswordChangedEvent struct {
	event.BaseEvent

	UserID uint `json:"user_id"`
}

// NewPasswordChangedEvent 创建密码变更事件
func NewPasswordChangedEvent(userID uint) *PasswordChangedEvent {
	return &PasswordChangedEvent{
		BaseEvent: event.NewBaseEvent("user.password_changed", "user", strconv.FormatUint(uint64(userID), 10)),
		UserID:    userID,
	}
}
