package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// ListSchemaHandler 获取设置 Schema 查询处理器
type ListSchemaHandler struct {
	settingQueryRepo  setting.QueryRepository
	categoryQueryRepo setting.SettingCategoryQueryRepository
}

// NewListSchemaHandler 创建 ListSchemaHandler 实例
func NewListSchemaHandler(
	settingQueryRepo setting.QueryRepository,
	categoryQueryRepo setting.SettingCategoryQueryRepository,
) *ListSchemaHandler {
	return &ListSchemaHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
	}
}

// Handle 处理获取设置 Schema 查询
// 返回按 Category → Group → Settings 层级组织的精简数据
//
// 支持 CategoryKey 过滤：
//   - 为空时返回全量系统设置（用于总配置页）
//   - 指定 Key 时只返回该分类（用于分散页面的懒加载）
func (h *ListSchemaHandler) Handle(ctx context.Context, query ListSchemaQuery) ([]SchemaCategoryDTO, error) {
	// 1. 根据 CategoryKey 决定查询范围
	settings, err := h.fetchSettings(ctx, query.CategoryKey)
	if err != nil {
		return nil, err
	}

	// 2. 查询所有分类元数据
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch setting categories: %w", err)
	}

	// 3. 使用共享构建器
	builder := NewSchemaBuilder(categories)
	return builder.Build(settings, nil, AdminSettingMapper), nil
}

// fetchSettings 根据 CategoryKey 获取设置列表
func (h *ListSchemaHandler) fetchSettings(ctx context.Context, categoryKey string) ([]*setting.Setting, error) {
	// 全量查询
	if categoryKey == "" {
		settings, err := h.settingQueryRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch settings: %w", err)
		}
		return settings, nil
	}

	// 按 Category Key 过滤
	category, err := h.categoryQueryRepo.FindByKey(ctx, categoryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find category by key %q: %w", categoryKey, err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found: %s", categoryKey)
	}

	settings, err := h.settingQueryRepo.FindByCategoryID(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch settings: %w", err)
	}
	return settings, nil
}
