package setting

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// ToSettingDTO 将领域实体转换为 DTO
func ToSettingDTO(setting *setting.Setting) *SettingDTO {
	if setting == nil {
		return nil
	}

	return &SettingDTO{
		ID:        setting.ID,
		Key:       setting.Key,
		Value:     setting.Value,
		Category:  setting.Category,
		ValueType: setting.ValueType,
		Label:     setting.Label,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}
