package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	appUserDTO "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	userCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/user/command"
	userQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/user/query"
)

// AdminUserHandler handles admin user management operations
type AdminUserHandler struct {
	createUserHandler      *userCommand.CreateUserHandler
	updateUserHandler      *userCommand.UpdateUserHandler
	deleteUserHandler      *userCommand.DeleteUserHandler
	assignRolesHandler     *userCommand.AssignRolesHandler
	batchCreateUserHandler *userCommand.BatchCreateUsersHandler
	getUserHandler         *userQuery.GetUserHandler
	listUsersHandler       *userQuery.ListUsersHandler
}

// NewAdminUserHandler creates a new AdminUserHandler instance
func NewAdminUserHandler(
	createUserHandler *userCommand.CreateUserHandler,
	updateUserHandler *userCommand.UpdateUserHandler,
	deleteUserHandler *userCommand.DeleteUserHandler,
	assignRolesHandler *userCommand.AssignRolesHandler,
	batchCreateUserHandler *userCommand.BatchCreateUsersHandler,
	getUserHandler *userQuery.GetUserHandler,
	listUsersHandler *userQuery.ListUsersHandler,
) *AdminUserHandler {
	return &AdminUserHandler{
		createUserHandler:      createUserHandler,
		updateUserHandler:      updateUserHandler,
		deleteUserHandler:      deleteUserHandler,
		assignRolesHandler:     assignRolesHandler,
		batchCreateUserHandler: batchCreateUserHandler,
		getUserHandler:         getUserHandler,
		listUsersHandler:       listUsersHandler,
	}
}

// CreateUser creates a new user (admin only)
//
// @Summary      创建用户
// @Description  管理员创建新用户账号，可同时分配角色
// @Tags         管理员 - 用户管理 (Admin - User Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body appUserDTO.CreateUserDTO true "用户信息"
// @Success      201 {object} response.Response{data=object{message=string,user=appUserDTO.UserWithRolesResponse}} "用户创建成功"
// @Failure      400 {object} response.ErrorResponse "参数错误或用户名/邮箱已存在"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/users [post]
// @x-permission {"scope":"admin:users:create"}
func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	var dto appUserDTO.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createUserHandler.Handle(c.Request.Context(), userCommand.CreateUserCommand{
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		FullName: dto.FullName,
		RoleIDs:  dto.RoleIDs,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	createdUser, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    result.UserID,
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "user created successfully", gin.H{
		"user": createdUser,
	})
}

// ListUsers lists all users with pagination (admin only)
//
// @Summary      获取用户列表
// @Description  分页获取所有用户列表（包含角色信息）
// @Tags         管理员 - 用户管理 (Admin - User Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "页码" default(1) minimum(1)
// @Param        limit query int false "每页数量" default(20) minimum(1) maximum(100)
// @Param        search query string false "搜索关键词（用户名或邮箱）"
// @Success      200 {object} response.ListResponse{data=[]appUserDTO.UserWithRolesResponse,meta=response.PaginationMeta} "用户列表"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/users [get]
// @x-permission {"scope":"admin:users:read"}
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	result, err := h.listUsersHandler.Handle(c.Request.Context(), userQuery.ListUsersQuery{
		Offset: offset,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), page, limit)
	response.List(c, "success", result.Users, meta)
}

// GetUser gets a user by ID (admin only)
//
// @Summary      获取用户详情
// @Description  根据用户ID获取用户详细信息（包含角色信息）
// @Tags         管理员 - 用户管理 (Admin - User Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID" minimum(1)
// @Success      200 {object} response.Response{data=appUserDTO.UserWithRolesResponse} "用户详情"
// @Failure      400 {object} response.ErrorResponse "无效的用户ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "用户不存在"
// @Router       /api/admin/users/{id} [get]
// @x-permission {"scope":"admin:users:read"}
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	userResp, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    uint(id),
		WithRoles: true,
	})
	if err != nil {
		response.NotFound(c, "user")
		return
	}

	response.OK(c, "success", userResp)
}

// UpdateUser updates a user (admin only)
//
// @Summary      更新用户信息
// @Description  管理员更新用户的基本信息和状态
// @Tags         管理员 - 用户管理 (Admin - User Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID" minimum(1)
// @Param        request body appUserDTO.UpdateUserDTO true "更新信息"
// @Success      200 {object} response.Response{data=object{message=string,user=appUserDTO.UserWithRolesResponse}} "用户更新成功"
// @Failure      400 {object} response.ErrorResponse "无效的用户ID或参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "用户不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/users/{id} [put]
// @x-permission {"scope":"admin:users:update"}
func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	var dto appUserDTO.UpdateUserDTO
	if err = c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	_, err = h.updateUserHandler.Handle(c.Request.Context(), userCommand.UpdateUserCommand{
		UserID:   uint(id),
		FullName: dto.FullName,
		Avatar:   dto.Avatar,
		Bio:      dto.Bio,
		Status:   dto.Status,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	updatedUser, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    uint(id),
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "user updated successfully", gin.H{
		"user": updatedUser,
	})
}

// DeleteUser deletes a user (admin only)
//
// @Summary      删除用户
// @Description  管理员删除指定用户（物理删除或软删除）
// @Tags         管理员 - 用户管理 (Admin - User Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID" minimum(1)
// @Success      200 {object} response.Response{data=object{message=string}} "用户删除成功"
// @Failure      400 {object} response.ErrorResponse "无效的用户ID"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "用户不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/users/{id} [delete]
// @x-permission {"scope":"admin:users:delete"}
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	if err := h.deleteUserHandler.Handle(c.Request.Context(), userCommand.DeleteUserCommand{
		UserID: uint(id),
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "user deleted successfully", nil)
}

// AssignRoles assigns roles to a user (admin only)
//
// @Summary      分配用户角色
// @Description  管理员为指定用户分配角色（会覆盖现有角色）
// @Tags         管理员 - 用户管理 (Admin - User Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID" minimum(1)
// @Param        request body appUserDTO.AssignRolesDTO true "角色ID列表"
// @Success      200 {object} response.Response{data=object{message=string}} "角色分配成功"
// @Failure      400 {object} response.ErrorResponse "无效的用户ID或参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "用户不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/users/{id}/roles [put]
// @x-permission {"scope":"admin:users:update"}
func (h *AdminUserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	var req appUserDTO.AssignRolesDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err = h.assignRolesHandler.Handle(c.Request.Context(), userCommand.AssignRolesCommand{
		UserID:  uint(id),
		RoleIDs: req.RoleIDs,
	}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 获取更新后的用户信息（包含角色）
	updatedUser, err := h.getUserHandler.Handle(c.Request.Context(), userQuery.GetUserQuery{
		UserID:    uint(id),
		WithRoles: true,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "roles assigned successfully", updatedUser)
}

// BatchCreateUsers creates multiple users at once (admin only)
//
// @Summary      批量创建用户
// @Description  管理员从 CSV 等来源批量创建用户，支持部分失败（单个失败不影响其他用户）
// @Tags         管理员 - 用户管理 (Admin - User Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body appUserDTO.BatchCreateUserDTO true "用户列表（最多 100 个）"
// @Success      200 {object} response.Response{data=appUserDTO.BatchCreateUserResponse} "批量创建结果"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/users/batch [post]
// @x-permission {"scope":"admin:users:create"}
func (h *AdminUserHandler) BatchCreateUsers(c *gin.Context) {
	var dto appUserDTO.BatchCreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 构建 Command
	users := make([]userCommand.BatchUserItem, len(dto.Users))
	for i, u := range dto.Users {
		users[i] = userCommand.BatchUserItem{
			Username: u.Username,
			Email:    u.Email,
			Password: u.Password,
			FullName: u.FullName,
			Status:   u.Status,
			RoleIDs:  u.RoleIDs,
		}
	}

	result, err := h.batchCreateUserHandler.Handle(c.Request.Context(), userCommand.BatchCreateUsersCommand{
		Users: users,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 构建响应
	resp := appUserDTO.BatchCreateUserResponse{
		Total:   result.Total,
		Success: result.Success,
		Failed:  result.Failed,
		Errors:  make([]appUserDTO.BatchCreateErrorDetail, len(result.Errors)),
	}
	for i, e := range result.Errors {
		resp.Errors[i] = appUserDTO.BatchCreateErrorDetail{
			Index:    e.Index,
			Username: e.Username,
			Email:    e.Email,
			Error:    e.Error,
		}
	}

	response.OK(c, "batch import completed", resp)
}
