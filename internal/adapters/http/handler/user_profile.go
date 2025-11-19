package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	appUserDTO "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	userCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/user/command"
	userQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/user/query"
)

// UserProfileHandler handles user profile operations
type UserProfileHandler struct {
	getUserHandler        *userQuery.GetUserHandler
	updateUserHandler     *userCommand.UpdateUserHandler
	changePasswordHandler *userCommand.ChangePasswordHandler
	deleteUserHandler     *userCommand.DeleteUserHandler
}

// NewUserProfileHandler creates a new UserProfileHandler instance
func NewUserProfileHandler(
	getUserHandler *userQuery.GetUserHandler,
	updateUserHandler *userCommand.UpdateUserHandler,
	changePasswordHandler *userCommand.ChangePasswordHandler,
	deleteUserHandler *userCommand.DeleteUserHandler,
) *UserProfileHandler {
	return &UserProfileHandler{
		getUserHandler:        getUserHandler,
		updateUserHandler:     updateUserHandler,
		changePasswordHandler: changePasswordHandler,
		deleteUserHandler:     deleteUserHandler,
	}
}

// GetProfile gets the current user's profile
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

	u, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    uid,
		WithRoles: true,
	})
	if err != nil {
		response.NotFound(c, "user")
		return
	}

	response.OK(c, u)
}

// UpdateProfile updates the current user's profile
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

	var req struct {
		FullName *string `json:"full_name" binding:"omitempty,max=100"`
		Avatar   *string `json:"avatar" binding:"omitempty,max=255"`
		Bio      *string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if _, err := h.updateUserHandler.Handle(c.Request.Context(), userCommand.UpdateUserCommand{
		UserID:   uid,
		FullName: req.FullName,
		Avatar:   req.Avatar,
		Bio:      req.Bio,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	updatedUser, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    uid,
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"message": "profile updated successfully",
		"user":    updatedUser,
	})
}

// ChangePassword changes the current user's password
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

	var req appUserDTO.ChangePasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.changePasswordHandler.Handle(c.Request.Context(), userCommand.ChangePasswordCommand{
		UserID:      uid,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "password changed successfully"})
}

// DeleteAccount deletes the current user's account
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

	if err := h.deleteUserHandler.Handle(c.Request.Context(), userCommand.DeleteUserCommand{
		UserID: uid,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "account deleted successfully"})
}
