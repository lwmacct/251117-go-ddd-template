package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// BatchUpdateSettingsHandler 批量更新设置命令处理器
type BatchUpdateSettingsHandler struct {
	settingCommandRepo setting.CommandRepository
	settingQueryRepo   setting.QueryRepository
}

// NewBatchUpdateSettingsHandler 创建 BatchUpdateSettingsHandler 实例
func NewBatchUpdateSettingsHandler(
	settingCommandRepo setting.CommandRepository,
	settingQueryRepo setting.QueryRepository,
) *BatchUpdateSettingsHandler {
	return &BatchUpdateSettingsHandler{
		settingCommandRepo: settingCommandRepo,
		settingQueryRepo:   settingQueryRepo,
	}
}

// Handle 处理批量更新设置命令
func (h *BatchUpdateSettingsHandler) Handle(ctx context.Context, cmd BatchUpdateSettingsCommand) error {
	// 1. 验证所有配置是否存在并构建更新列表
	settings := make([]*setting.Setting, 0, len(cmd.Settings))
	for _, item := range cmd.Settings {
		existing, err := h.settingQueryRepo.FindByKey(ctx, item.Key)
		if err != nil {
			return fmt.Errorf("failed to find setting %s: %w", item.Key, err)
		}
		if existing == nil {
			return fmt.Errorf("setting key %s does not exist", item.Key)
		}

		// 更新值
		existing.Value = item.Value
		settings = append(settings, existing)
	}

	// 2. 批量更新
	if err := h.settingCommandRepo.BatchUpsert(ctx, settings); err != nil {
		return fmt.Errorf("failed to batch update settings: %w", err)
	}

	return nil
}
