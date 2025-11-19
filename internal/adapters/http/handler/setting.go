package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

type SettingHandler struct {
	settingCommandRepo setting.CommandRepository
	settingQueryRepo   setting.QueryRepository
}

func NewSettingHandler(settingCommandRepo setting.CommandRepository, settingQueryRepo setting.QueryRepository) *SettingHandler {
	return &SettingHandler{
		settingCommandRepo: settingCommandRepo,
		settingQueryRepo:   settingQueryRepo,
	}
}

// GetSettings 获取配置列表
func (h *SettingHandler) GetSettings(c *gin.Context) {
	category := c.Query("category")

	var settings []*setting.Setting
	var err error

	if category != "" {
		settings, err = h.settingQueryRepo.FindByCategory(c.Request.Context(), category)
	} else {
		settings, err = h.settingQueryRepo.FindAll(c.Request.Context())
	}

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, settings)
}

// GetSetting 获取单个配置
func (h *SettingHandler) GetSetting(c *gin.Context) {
	key := c.Param("key")

	s, err := h.settingQueryRepo.FindByKey(c.Request.Context(), key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if s == nil {
		response.NotFound(c, "setting")
		return
	}

	response.OK(c, s)
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

	// 检查 Key 是否已存在
	existing, err := h.settingQueryRepo.FindByKey(c.Request.Context(), req.Key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if existing != nil {
		response.BadRequest(c, "配置键已存在")
		return
	}

	// 默认值类型
	if req.ValueType == "" {
		req.ValueType = setting.ValueTypeString
	}

	s := &setting.Setting{
		Key:       req.Key,
		Value:     req.Value,
		Category:  req.Category,
		ValueType: req.ValueType,
		Label:     req.Label,
	}

	if err := h.settingCommandRepo.Create(c.Request.Context(), s); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, s)
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

	s, err := h.settingQueryRepo.FindByKey(c.Request.Context(), key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if s == nil {
		response.NotFound(c, "setting")
		return
	}

	s.Value = req.Value
	if req.ValueType != "" {
		s.ValueType = req.ValueType
	}
	if req.Label != "" {
		s.Label = req.Label
	}

	if err := h.settingCommandRepo.Update(c.Request.Context(), s); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, s)
}

// DeleteSetting 删除配置
func (h *SettingHandler) DeleteSetting(c *gin.Context) {
	key := c.Param("key")

	s, err := h.settingQueryRepo.FindByKey(c.Request.Context(), key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if s == nil {
		response.NotFound(c, "setting")
		return
	}

	if err := h.settingCommandRepo.Delete(c.Request.Context(), s.ID); err != nil {
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

	// 构建配置列表
	settings := make([]*setting.Setting, 0, len(req.Settings))
	for _, s := range req.Settings {
		// 查找现有配置以获取完整信息
		existing, err := h.settingQueryRepo.FindByKey(c.Request.Context(), s.Key)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		if existing != nil {
			// 更新现有配置的值
			existing.Value = s.Value
			settings = append(settings, existing)
		} else {
			// 如果配置不存在,跳过(或者可以返回错误)
			response.BadRequest(c, "配置键 "+s.Key+" 不存在")
			return
		}
	}

	// 批量更新
	if err := h.settingCommandRepo.BatchUpsert(c.Request.Context(), settings); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "批量更新成功"})
}
