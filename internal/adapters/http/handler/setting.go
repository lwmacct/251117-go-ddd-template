package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/setting"
)

// SettingHandler handles setting management operations (DDD+CQRS Use Case Pattern)
type SettingHandler struct {
	// Command Handlers
	createHandler      *setting.CreateHandler
	updateHandler      *setting.UpdateHandler
	deleteHandler      *setting.DeleteHandler
	batchUpdateHandler *setting.BatchUpdateHandler

	// Query Handlers
	getHandler        *setting.GetHandler
	listHandler       *setting.ListHandler
	listSchemaHandler *setting.ListSchemaHandler
}

// NewSettingHandler creates a new SettingHandler instance
func NewSettingHandler(
	createHandler *setting.CreateHandler,
	updateHandler *setting.UpdateHandler,
	deleteHandler *setting.DeleteHandler,
	batchUpdateHandler *setting.BatchUpdateHandler,
	getHandler *setting.GetHandler,
	listHandler *setting.ListHandler,
	listSchemaHandler *setting.ListSchemaHandler,
) *SettingHandler {
	return &SettingHandler{
		createHandler:      createHandler,
		updateHandler:      updateHandler,
		deleteHandler:      deleteHandler,
		batchUpdateHandler: batchUpdateHandler,
		getHandler:         getHandler,
		listHandler:        listHandler,
		listSchemaHandler:  listSchemaHandler,
	}
}

// GetSettings 获取配置列表
//
// @Summary      获取系统配置列表
// @Description  获取所有系统配置，可按类别筛选
// @Tags         管理员 - 系统配置 (Admin - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category query string false "配置类别" example:"system"
// @Success      200 {object} response.DataResponse[[]setting.SettingDTO] "配置列表"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/settings [get]
// @x-permission {"scope":"admin:settings:read"}
func (h *SettingHandler) GetSettings(c *gin.Context) {
	category := c.Query("category")

	// 调用 Use Case Handler
	settings, err := h.listHandler.Handle(c.Request.Context(), setting.ListQuery{
		Category: category,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "success", settings)
}

// GetSetting 获取单个配置
//
// @Summary      获取单个配置
// @Description  根据配置键获取配置详情
// @Tags         管理员 - 系统配置 (Admin - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        key path string true "配置键" example:"site_name"
// @Success      200 {object} response.DataResponse[setting.SettingDTO] "配置详情"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "配置不存在"
// @Router       /api/admin/settings/{key} [get]
// @x-permission {"scope":"admin:settings:read"}
func (h *SettingHandler) GetSetting(c *gin.Context) {
	key := c.Param("key")

	// 调用 Use Case Handler
	setting, err := h.getHandler.Handle(c.Request.Context(), setting.GetQuery{
		Key: key,
	})

	if err != nil {
		response.NotFound(c, "setting")
		return
	}

	response.OK(c, "success", setting)
}

// CreateSettingRequest 创建配置请求
type CreateSettingRequest struct {
	Key          string `json:"key" binding:"required" example:"site_name"`
	DefaultValue any    `json:"default_value" binding:"required"`
	Category     string `json:"category" binding:"required" example:"general"`
	Group        string `json:"group" example:"basic"`
	ValueType    string `json:"value_type" example:"string"`
	Label        string `json:"label" example:"网站名称"`
	UIConfig     string `json:"ui_config" example:"{}"`
	Order        int    `json:"order" example:"0"`
}

// CreateSetting 创建配置
//
// @Summary      创建配置
// @Description  管理员创建新的系统配置项
// @Tags         管理员 - 系统配置 (Admin - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateSettingRequest true "配置信息"
// @Success      201 {object} response.DataResponse[setting.SettingDTO] "配置创建成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/settings [post]
// @x-permission {"scope":"admin:settings:create"}
func (h *SettingHandler) CreateSetting(c *gin.Context) {
	var req CreateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用 Use Case Handler 创建配置
	_, err := h.createHandler.Handle(c.Request.Context(), setting.CreateCommand{
		Key:          req.Key,
		DefaultValue: req.DefaultValue,
		Category:     req.Category,
		Group:        req.Group,
		ValueType:    req.ValueType,
		Label:        req.Label,
		UIConfig:     req.UIConfig,
		Order:        req.Order,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 查询完整的配置信息返回
	settingDTO, err := h.getHandler.Handle(c.Request.Context(), setting.GetQuery{
		Key: req.Key,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "setting created successfully", settingDTO)
}

// UpdateSettingRequest 更新配置请求
type UpdateSettingRequest struct {
	DefaultValue any    `json:"default_value"`
	Label        string `json:"label" example:"更新后的标签"`
	UIConfig     string `json:"ui_config"`
	Order        int    `json:"order"`
}

// UpdateSetting 更新配置
//
// @Summary      更新配置
// @Description  管理员更新指定配置项的值和标签
// @Tags         管理员 - 系统配置 (Admin - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        key path string true "配置键" example:"site_name"
// @Param        request body UpdateSettingRequest true "更新信息"
// @Success      200 {object} response.DataResponse[setting.SettingDTO] "配置更新成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "配置不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/settings/{key} [put]
// @x-permission {"scope":"admin:settings:update"}
func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	settingDTO, err := h.updateHandler.Handle(c.Request.Context(), setting.UpdateCommand{
		Key:          key,
		DefaultValue: req.DefaultValue,
		Label:        req.Label,
		UIConfig:     req.UIConfig,
		Order:        req.Order,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "setting updated successfully", settingDTO)
}

// DeleteSetting 删除配置
//
// @Summary      删除配置
// @Description  管理员删除指定的系统配置项
// @Tags         管理员 - 系统配置 (Admin - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        key path string true "配置键" example:"site_name"
// @Success      204 "配置删除成功"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      404 {object} response.ErrorResponse "配置不存在"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/settings/{key} [delete]
// @x-permission {"scope":"admin:settings:delete"}
func (h *SettingHandler) DeleteSetting(c *gin.Context) {
	key := c.Param("key")

	// 调用 Use Case Handler
	err := h.deleteHandler.Handle(c.Request.Context(), setting.DeleteCommand{
		Key: key,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}

// BatchUpdateSettingsRequest 批量更新配置请求
type BatchUpdateSettingsRequest struct {
	Settings []struct {
		Key   string `json:"key" binding:"required"`
		Value any    `json:"value"` // JSONB 原生值
	} `json:"settings" binding:"required,min=1"` // 至少需要一个设置项
}

// BatchUpdateSettings 批量更新配置
//
// @Summary      批量更新配置
// @Description  管理员批量更新多个系统配置项的值
// @Tags         管理员 - 系统配置 (Admin - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body BatchUpdateSettingsRequest true "配置列表"
// @Success      200 {object} response.MessageResponse "批量更新成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/settings/batch [post]
// @x-permission {"scope":"admin:settings:update"}
func (h *SettingHandler) BatchUpdateSettings(c *gin.Context) {
	var req BatchUpdateSettingsRequest
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

	// 调用 Use Case Handler
	err := h.batchUpdateHandler.Handle(c.Request.Context(), setting.BatchUpdateCommand{
		Settings: items,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "批量更新成功", nil)
}

// GetSettingsSchema 获取配置 Schema
//
// @Summary      获取配置 Schema
// @Description  获取按 Category → Group → Settings 层级组织的配置数据，用于前端动态渲染设置页面
// @Tags         管理员 - 系统配置 (Admin - Settings)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.DataResponse[[]setting.SchemaCategoryDTO] "配置 Schema"
// @Failure      401 {object} response.ErrorResponse "未授权"
// @Failure      403 {object} response.ErrorResponse "权限不足"
// @Failure      500 {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/admin/settings/schema [get]
// @x-permission {"scope":"admin:settings:read"}
func (h *SettingHandler) GetSettingsSchema(c *gin.Context) {
	// 调用 Use Case Handler
	schema, err := h.listSchemaHandler.Handle(c.Request.Context(), setting.ListSchemaQuery{})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "success", schema)
}
