package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/user"
)

// UserProfileHandler handles user profile operations
type UserProfileHandler struct {
	getUserHandler        *user.GetUserHandler
	updateUserHandler     *user.UpdateUserHandler
	changePasswordHandler *user.ChangePasswordHandler
	deleteUserHandler     *user.DeleteUserHandler
}

// NewUserProfileHandler creates a new UserProfileHandler instance
func NewUserProfileHandler(
	getUserHandler *user.GetUserHandler,
	updateUserHandler *user.UpdateUserHandler,
	changePasswordHandler *user.ChangePasswordHandler,
	deleteUserHandler *user.DeleteUserHandler,
) *UserProfileHandler {
	return &UserProfileHandler{
		getUserHandler:        getUserHandler,
		updateUserHandler:     updateUserHandler,
		changePasswordHandler: changePasswordHandler,
		deleteUserHandler:     deleteUserHandler,
	}
}

// GetProfile gets the current user's profile
//
// @Summary      获取个人资料
// @Description  获取当前登录用户的个人资料和角色信息
// @Tags         用户 - 个人资料 (User - Profile)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.DataResponse[user.UserWithRolesDTO] "个人资料"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      404 {object} response.ErrorResponse "用户不存在"
// @Router       /api/user/profile [get]
// @x-permission {"scope":"user:profile:read"}
func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.InternalError(c, "invalid user ID")
		return
	}

	u, err := h.getUserHandler.Handle(c.Request.Context(), user.GetUserQuery{
		UserID:    uid,
		WithRoles: true,
	})
	if err != nil {
		response.NotFound(c, "user")
		return
	}

	response.OK(c, "success", u)
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	FullName *string `json:"full_name" binding:"omitempty,max=100" example:"张三"`
	Avatar   *string `json:"avatar" binding:"omitempty,max=255" example:"https://example.com/avatar.jpg"`
	Bio      *string `json:"bio" example:"这是我的个人简介"`
}

// UpdateProfile updates the current user's profile
//
// @Summary      更新个人资料
// @Description  用户更新自己的姓名、头像和个人简介
// @Tags         用户 - 个人资料 (User - Profile)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateProfileRequest true "更新信息"
// @Success      200 {object} response.DataResponse[user.UserWithRolesDTO] "资料更新成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/user/profile [put]
// @x-permission {"scope":"user:profile:update"}
func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.InternalError(c, "invalid user ID")
		return
	}

	var req UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if _, err := h.updateUserHandler.Handle(c.Request.Context(), user.UpdateUserCommand{
		UserID:   uid,
		FullName: req.FullName,
		Avatar:   req.Avatar,
		Bio:      req.Bio,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	updatedUser, err := h.getUserHandler.Handle(c.Request.Context(), user.GetUserQuery{
		UserID:    uid,
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "profile updated successfully", updatedUser)
}

// ChangePassword changes the current user's password
//
// @Summary      修改密码
// @Description  用户修改自己的登录密码
// @Tags         用户 - 个人资料 (User - Profile)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body user.ChangePasswordDTO true "密码信息"
// @Success      200 {object} response.MessageResponse "密码修改成功"
// @Failure      400 {object} response.ErrorResponse "参数错误或旧密码不正确"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Router       /api/user/password [put]
// @x-permission {"scope":"user:password:update"}
func (h *UserProfileHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.InternalError(c, "invalid user ID")
		return
	}

	var req user.ChangePasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.changePasswordHandler.Handle(c.Request.Context(), user.ChangePasswordCommand{
		UserID:      uid,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "password changed successfully", nil)
}

// DeleteAccount deletes the current user's account
//
// @Summary      删除账号
// @Description  用户删除自己的账号（不可恢复）
// @Tags         用户 - 个人资料 (User - Profile)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.MessageResponse "账号删除成功"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/user/account [delete]
// @x-permission {"scope":"user:profile:delete"}
func (h *UserProfileHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.InternalError(c, "invalid user ID")
		return
	}

	if err := h.deleteUserHandler.Handle(c.Request.Context(), user.DeleteUserCommand{
		UserID: uid,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "account deleted successfully", nil)
}
