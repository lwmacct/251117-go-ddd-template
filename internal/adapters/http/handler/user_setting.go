package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
)

// UserSettingHandler handles user setting operations (DDD+CQRS Use Case Pattern)
type UserSettingHandler struct {
	// Command Handlers
	setHandler      *setting.UserSetHandler
	batchSetHandler *setting.UserBatchSetHandler
	resetHandler    *setting.UserResetHandler
	resetAllHandler *setting.UserResetAllHandler

	// Query Handlers
	getHandler        *setting.UserGetHandler
	listHandler       *setting.UserListHandler
	listSchemaHandler *setting.UserListSchemaHandler
}

// NewUserSettingHandler creates a new UserSettingHandler instance
func NewUserSettingHandler(
	setHandler *setting.UserSetHandler,
	batchSetHandler *setting.UserBatchSetHandler,
	resetHandler *setting.UserResetHandler,
	resetAllHandler *setting.UserResetAllHandler,
	getHandler *setting.UserGetHandler,
	listHandler *setting.UserListHandler,
	listSchemaHandler *setting.UserListSchemaHandler,
) *UserSettingHandler {
	return &UserSettingHandler{
		setHandler:        setHandler,
		batchSetHandler:   batchSetHandler,
		resetHandler:      resetHandler,
		resetAllHandler:   resetAllHandler,
		getHandler:        getHandler,
		listHandler:       listHandler,
		listSchemaHandler: listSchemaHandler,
	}
}

// GetUserSettings 获取用户配置列表
//
// @Summary      获取用户配置列表
// @Description  获取当前用户的所有配置（合并系统默认值）
// @Tags         用户 - 个人配置 (User - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @x-permission {"scope":"user:settings:read"}
// @Param        category_id query int false "配置类别 ID"
// @Success      200 {object} response.DataResponse[[]setting.UserSettingDTO] "配置列表"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/user/settings [get]
func (h *UserSettingHandler) GetUserSettings(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var categoryID uint
	if id := c.Query("category_id"); id != "" {
		parsed, _ := strconv.ParseUint(id, 10, 64)
		categoryID = uint(parsed)
	}

	settings, err := h.listHandler.Handle(c.Request.Context(), setting.UserListQuery{
		UserID:     userID,
		CategoryID: categoryID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "success", settings)
}

// GetUserSetting 获取单个用户配置
//
// @Summary      获取单个用户配置
// @Description  根据配置键获取用户配置（合并系统默认值）
// @Tags         用户 - 个人配置 (User - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @x-permission {"scope":"user:settings:read"}
// @Param        key path string true "配置键" example:"theme"
// @Success      200 {object} response.DataResponse[setting.UserSettingDTO] "配置详情"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      404 {object} response.ErrorResponse "配置不存在"
// @Router       /api/user/settings/{key} [get]
func (h *UserSettingHandler) GetUserSetting(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	key := c.Param("key")

	settingDTO, err := h.getHandler.Handle(c.Request.Context(), setting.UserGetQuery{
		UserID: userID,
		Key:    key,
	})
	if err != nil {
		response.NotFound(c, "setting")
		return
	}

	response.OK(c, "success", settingDTO)
}

// SetUserSettingRequest 设置用户配置请求
type SetUserSettingRequest struct {
	Value any `json:"value"` // JSONB 原生值
}

// SetUserSetting 设置用户配置
//
// @Summary      设置用户配置
// @Description  用户设置自定义配置值
// @Tags         用户 - 个人配置 (User - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @x-permission {"scope":"user:settings:update"}
// @Param        key path string true "配置键" example:"theme"
// @Param        request body SetUserSettingRequest true "配置值"
// @Success      200 {object} response.DataResponse[setting.UserSettingDTO] "设置成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      404 {object} response.ErrorResponse "配置不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/user/settings/{key} [put]
func (h *UserSettingHandler) SetUserSetting(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	key := c.Param("key")

	var req SetUserSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	settingDTO, err := h.setHandler.Handle(c.Request.Context(), setting.UserSetCommand{
		UserID: userID,
		Key:    key,
		Value:  req.Value,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "设置成功", settingDTO)
}

// ResetUserSetting 重置用户配置
//
// @Summary      重置用户配置
// @Description  删除用户自定义配置，恢复为系统默认值
// @Tags         用户 - 个人配置 (User - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @x-permission {"scope":"user:settings:update"}
// @Param        key path string true "配置键" example:"theme"
// @Success      204 "重置成功"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/user/settings/{key} [delete]
func (h *UserSettingHandler) ResetUserSetting(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	key := c.Param("key")

	err := h.resetHandler.Handle(c.Request.Context(), setting.UserResetCommand{
		UserID: userID,
		Key:    key,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}

// BatchSetUserSettingsRequest 批量设置用户配置请求
type BatchSetUserSettingsRequest struct {
	Settings []struct {
		Key   string `json:"key" binding:"required"`
		Value any    `json:"value"` // JSONB 原生值
	} `json:"settings" binding:"required,min=1"`
}

// BatchSetUserSettings 批量设置用户配置
//
// @Summary      批量设置用户配置
// @Description  用户批量设置多个自定义配置值
// @Tags         用户 - 个人配置 (User - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @x-permission {"scope":"user:settings:update"}
// @Param        request body BatchSetUserSettingsRequest true "配置列表"
// @Success      200 {object} response.MessageResponse "批量设置成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/user/settings/batch [post]
func (h *UserSettingHandler) BatchSetUserSettings(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req BatchSetUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换为 Command
	items := make([]setting.SettingItemCommand, len(req.Settings))
	for i, s := range req.Settings {
		items[i] = setting.SettingItemCommand{
			Key:   s.Key,
			Value: s.Value,
		}
	}

	err := h.batchSetHandler.Handle(c.Request.Context(), setting.UserBatchSetCommand{
		UserID:   userID,
		Settings: items,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "批量设置成功", nil)
}

// GetUserSettingsSchema 获取用户配置 Schema
//
// @Summary      获取用户配置 Schema
// @Description  获取按 Category → Group → Settings 层级组织的配置数据，包含用户自定义值
// @Tags         用户 - 个人配置 (User - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @x-permission {"scope":"user:settings:read"}
// @Success      200 {object} response.DataResponse[[]setting.UserSchemaCategoryDTO] "配置 Schema"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/user/settings/schema [get]
func (h *UserSettingHandler) GetUserSettingsSchema(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	schema, err := h.listSchemaHandler.Handle(c.Request.Context(), setting.UserListSchemaQuery{
		UserID: userID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "success", schema)
}
