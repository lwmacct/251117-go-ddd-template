package setting

import (
	"context"
	"fmt"
	"sort"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// UserListCategoriesHandler 获取用户可见的分类列表处理器
// 只返回包含 scope="user" 设置的分类，用于懒加载场景
type UserListCategoriesHandler struct {
	settingQueryRepo  setting.QueryRepository
	categoryQueryRepo setting.SettingCategoryQueryRepository
}

// NewUserListCategoriesHandler 创建 UserListCategoriesHandler 实例
func NewUserListCategoriesHandler(
	settingQueryRepo setting.QueryRepository,
	categoryQueryRepo setting.SettingCategoryQueryRepository,
) *UserListCategoriesHandler {
	return &UserListCategoriesHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
	}
}

// Handle 处理获取用户可见分类列表查询
// 返回包含 scope="user" 设置的分类元信息（不含 settings 数据）
func (h *UserListCategoriesHandler) Handle(ctx context.Context, _ UserListCategoriesQuery) ([]CategoryMetaDTO, error) {
	// 1. 查询所有 user scope 的设置
	defs, err := h.settingQueryRepo.FindByScope(ctx, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to find user scope settings: %w", err)
	}

	// 2. 提取唯一的 CategoryID
	categoryIDs := extractUniqueCategoryIDs(defs)
	if len(categoryIDs) == 0 {
		return []CategoryMetaDTO{}, nil
	}

	// 3. 查询分类元信息
	categories, err := h.categoryQueryRepo.FindByIDs(ctx, categoryIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories by IDs: %w", err)
	}

	// 4. 转换为 DTO 并按 Order 排序
	return toCategoryMetaDTOs(categories), nil
}

// extractUniqueCategoryIDs 从设置列表中提取唯一的 CategoryID
func extractUniqueCategoryIDs(settings []*setting.Setting) []uint {
	seen := make(map[uint]struct{})
	var ids []uint
	for _, s := range settings {
		if _, ok := seen[s.CategoryID]; !ok {
			seen[s.CategoryID] = struct{}{}
			ids = append(ids, s.CategoryID)
		}
	}
	return ids
}

// toCategoryMetaDTOs 将分类实体转换为元信息 DTO
func toCategoryMetaDTOs(categories []*setting.SettingCategory) []CategoryMetaDTO {
	result := make([]CategoryMetaDTO, 0, len(categories))
	for _, cat := range categories {
		result = append(result, CategoryMetaDTO{
			Category: cat.Key,
			Label:    cat.Label,
			Icon:     cat.Icon,
			Order:    cat.Order,
		})
	}
	// 按 Order 排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Order < result[j].Order
	})
	return result
}
