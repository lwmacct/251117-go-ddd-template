package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// GetSettingHandler 获取设置查询处理器
type GetSettingHandler struct {
	settingQueryRepo setting.QueryRepository
}

// NewGetSettingHandler 创建 GetSettingHandler 实例
func NewGetSettingHandler(settingQueryRepo setting.QueryRepository) *GetSettingHandler {
	return &GetSettingHandler{
		settingQueryRepo: settingQueryRepo,
	}
}

// Handle 处理获取设置查询
func (h *GetSettingHandler) Handle(ctx context.Context, query GetSettingQuery) (*SettingDTO, error) {
	settingEntity, err := h.settingQueryRepo.FindByKey(ctx, query.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}

	if settingEntity == nil {
		return nil, errors.New("setting not found")
	}

	return ToSettingDTO(settingEntity), nil
}
