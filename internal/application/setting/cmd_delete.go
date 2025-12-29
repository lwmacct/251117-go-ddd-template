package setting

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// DeleteHandler 删除配置命令处理器
type DeleteHandler struct {
	commandRepo setting.CommandRepository
	queryRepo   setting.QueryRepository
	schemaCache SchemaCacheService
}

// NewDeleteHandler 创建 DeleteHandler 实例
func NewDeleteHandler(
	commandRepo setting.CommandRepository,
	queryRepo setting.QueryRepository,
	schemaCache SchemaCacheService,
) *DeleteHandler {
	return &DeleteHandler{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		schemaCache: schemaCache,
	}
}

// Handle 处理删除配置命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 1. 查询配置定义
	def, err := h.queryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return errors.New("setting not found")
	}

	// 2. 删除配置定义
	if err := h.commandRepo.Delete(ctx, cmd.Key); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}

	// 3. 失效 Schema 缓存
	if h.schemaCache != nil {
		if err := h.schemaCache.DeleteAdminSchemaAll(ctx); err != nil {
			slog.Warn("admin schema cache invalidation failed", "key", cmd.Key, "err", err)
		}
	}

	return nil
}
