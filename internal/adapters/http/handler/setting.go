package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	settingCommand "github.com/lwmacct/251117-go-ddd-template/internal/application/setting/command"
	settingQuery "github.com/lwmacct/251117-go-ddd-template/internal/application/setting/query"
)

// SettingHandler handles setting management operations (DDD+CQRS Use Case Pattern)
type SettingHandler struct {
	// Command Handlers
	createSettingHandler       *settingCommand.CreateSettingHandler
	updateSettingHandler       *settingCommand.UpdateSettingHandler
	deleteSettingHandler       *settingCommand.DeleteSettingHandler
	batchUpdateSettingsHandler *settingCommand.BatchUpdateSettingsHandler

	// Query Handlers
	getSettingHandler   *settingQuery.GetSettingHandler
	listSettingsHandler *settingQuery.ListSettingsHandler
}

// NewSettingHandler creates a new SettingHandler instance
func NewSettingHandler(
	createSettingHandler *settingCommand.CreateSettingHandler,
	updateSettingHandler *settingCommand.UpdateSettingHandler,
	deleteSettingHandler *settingCommand.DeleteSettingHandler,
	batchUpdateSettingsHandler *settingCommand.BatchUpdateSettingsHandler,
	getSettingHandler *settingQuery.GetSettingHandler,
	listSettingsHandler *settingQuery.ListSettingsHandler,
) *SettingHandler {
	return &SettingHandler{
		createSettingHandler:       createSettingHandler,
		updateSettingHandler:       updateSettingHandler,
		deleteSettingHandler:       deleteSettingHandler,
		batchUpdateSettingsHandler: batchUpdateSettingsHandler,
		getSettingHandler:          getSettingHandler,
		listSettingsHandler:        listSettingsHandler,
	}
}

// GetSettings 获取配置列表
func (h *SettingHandler) GetSettings(c *gin.Context) {
	category := c.Query("category")

	// 调用 Use Case Handler
	settings, err := h.listSettingsHandler.Handle(c.Request.Context(), settingQuery.ListSettingsQuery{
		Category: category,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, settings)
}

// GetSetting 获取单个配置
func (h *SettingHandler) GetSetting(c *gin.Context) {
	key := c.Param("key")

	// 调用 Use Case Handler
	setting, err := h.getSettingHandler.Handle(c.Request.Context(), settingQuery.GetSettingQuery{
		Key: key,
	})

	if err != nil {
		response.NotFound(c, "setting")
		return
	}

	response.OK(c, setting)
}

// CreateSettingRequest 创建配置请求
type CreateSettingRequest struct {
	Key       string `json:"key" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Category  string `json:"category" binding:"required"`
	ValueType string `json:"value_type"`
	Label     string `json:"label"`
}

// CreateSetting 创建配置
func (h *SettingHandler) CreateSetting(c *gin.Context) {
	var req CreateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	result, err := h.createSettingHandler.Handle(c.Request.Context(), settingCommand.CreateSettingCommand{
		Key:       req.Key,
		Value:     req.Value,
		Category:  req.Category,
		ValueType: req.ValueType,
		Label:     req.Label,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// UpdateSettingRequest 更新配置请求
type UpdateSettingRequest struct {
	Value     string `json:"value" binding:"required"`
	ValueType string `json:"value_type"`
	Label     string `json:"label"`
}

// UpdateSetting 更新配置
func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	setting, err := h.updateSettingHandler.Handle(c.Request.Context(), settingCommand.UpdateSettingCommand{
		Key:       key,
		Value:     req.Value,
		ValueType: req.ValueType,
		Label:     req.Label,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, setting)
}

// DeleteSetting 删除配置
func (h *SettingHandler) DeleteSetting(c *gin.Context) {
	key := c.Param("key")

	// 调用 Use Case Handler
	err := h.deleteSettingHandler.Handle(c.Request.Context(), settingCommand.DeleteSettingCommand{
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
		Value string `json:"value" binding:"required"`
	} `json:"settings" binding:"required"`
}

// BatchUpdateSettings 批量更新配置
func (h *SettingHandler) BatchUpdateSettings(c *gin.Context) {
	var req BatchUpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换为 Command
	settings := make([]settingCommand.SettingItem, len(req.Settings))
	for i, s := range req.Settings {
		settings[i].Key = s.Key
		settings[i].Value = s.Value
	}

	// 调用 Use Case Handler
	err := h.batchUpdateSettingsHandler.Handle(c.Request.Context(), settingCommand.BatchUpdateSettingsCommand{
		Settings: settings,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "批量更新成功"})
}
