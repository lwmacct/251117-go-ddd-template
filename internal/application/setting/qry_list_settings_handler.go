package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// ListSettingsHandler 获取设置列表查询处理器
type ListSettingsHandler struct {
	settingQueryRepo setting.QueryRepository
}

// NewListSettingsHandler 创建 ListSettingsHandler 实例
func NewListSettingsHandler(settingQueryRepo setting.QueryRepository) *ListSettingsHandler {
	return &ListSettingsHandler{
		settingQueryRepo: settingQueryRepo,
	}
}

// Handle 处理获取设置列表查询
func (h *ListSettingsHandler) Handle(ctx context.Context, query ListSettingsQuery) ([]*SettingDTO, error) {
	var settings []*setting.Setting
	var err error

	if query.Category != "" {
		settings, err = h.settingQueryRepo.FindByCategory(ctx, query.Category)
	} else {
		settings, err = h.settingQueryRepo.FindAll(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch settings: %w", err)
	}

	// 转换为 DTO
	settingResponses := make([]*SettingDTO, 0, len(settings))
	for _, s := range settings {
		settingResponses = append(settingResponses, ToSettingDTO(s))
	}

	return settingResponses, nil
}
