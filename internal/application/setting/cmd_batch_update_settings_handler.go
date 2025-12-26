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
	if len(cmd.Settings) == 0 {
		return nil
	}

	// 1. 提取所有 keys 并构建 key -> value 映射
	keys := make([]string, len(cmd.Settings))
	keyValueMap := make(map[string]string, len(cmd.Settings))
	for i, item := range cmd.Settings {
		keys[i] = item.Key
		keyValueMap[item.Key] = item.Value
	}

	// 2. 批量查询所有现有配置（1 次查询，解决 N+1 问题）
	existingSettings, err := h.settingQueryRepo.FindByKeys(ctx, keys)
	if err != nil {
		return fmt.Errorf("failed to find settings: %w", err)
	}

	// 3. 构建 key -> setting 映射，用于验证存在性
	existingMap := make(map[string]*setting.Setting, len(existingSettings))
	for _, s := range existingSettings {
		existingMap[s.Key] = s
	}

	// 4. 验证所有 key 存在并更新值
	settings := make([]*setting.Setting, 0, len(cmd.Settings))
	for _, key := range keys {
		existing, ok := existingMap[key]
		if !ok {
			return fmt.Errorf("setting key %s does not exist", key)
		}
		existing.Value = keyValueMap[key]
		settings = append(settings, existing)
	}

	// 5. 批量更新（1 次写入）
	if err := h.settingCommandRepo.BatchUpsert(ctx, settings); err != nil {
		return fmt.Errorf("failed to batch update settings: %w", err)
	}

	return nil
}
