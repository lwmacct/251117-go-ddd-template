package setting

import (
	domainSetting "github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// ToSettingResponse 将领域实体转换为 DTO
func ToSettingResponse(setting *domainSetting.Setting) *SettingResponse {
	if setting == nil {
		return nil
	}

	return &SettingResponse{
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
