package seeds

import (
	"context"
	"log/slog"

	_persistence "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SettingSeeder 系统设置种子数据
type SettingSeeder struct{}

// Seed 执行系统设置种子数据填充
func (s *SettingSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	settings := []_persistence.SettingModel{
		// General 常规设置
		{Key: "general.site_name", Value: "", Category: "general", ValueType: "string", Label: "站点名称"},
		{Key: "general.site_url", Value: "", Category: "general", ValueType: "string", Label: "站点 URL"},
		{Key: "general.admin_email", Value: "", Category: "general", ValueType: "string", Label: "管理员邮箱"},
		{Key: "general.timezone", Value: "Asia/Shanghai", Category: "general", ValueType: "string", Label: "时区"},
		{Key: "general.language", Value: "zh-CN", Category: "general", ValueType: "string", Label: "语言"},
		{Key: "general.theme", Value: "light", Category: "general", ValueType: "string", Label: "主题"},
		// Security 安全设置
		{Key: "security.session_timeout", Value: "30", Category: "security", ValueType: "number", Label: "会话超时时间"},
		{Key: "security.password_min_length", Value: "8", Category: "security", ValueType: "number", Label: "密码最小长度"},
		{Key: "security.enable_twofa", Value: "false", Category: "security", ValueType: "boolean", Label: "强制启用两步验证"},
		{Key: "security.max_login_attempts", Value: "5", Category: "security", ValueType: "number", Label: "最大登录尝试次数"},
		// Notification 通知设置
		{Key: "notification.enable_notifications", Value: "true", Category: "notification", ValueType: "boolean", Label: "启用系统通知"},
		{Key: "notification.enable_email", Value: "true", Category: "notification", ValueType: "boolean", Label: "启用邮件通知"},
		{Key: "notification.enable_sms", Value: "false", Category: "notification", ValueType: "boolean", Label: "启用短信通知"},
		// Backup 备份设置
		{Key: "backup.enable_backup", Value: "false", Category: "backup", ValueType: "boolean", Label: "启用自动备份"},
		{Key: "backup.backup_frequency", Value: "24", Category: "backup", ValueType: "number", Label: "备份频率"},
		{Key: "backup.retention_days", Value: "30", Category: "backup", ValueType: "number", Label: "备份保留天数"},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoNothing: true, // 已存在则跳过，不覆盖用户修改
	}).Create(&settings)
	if result.Error != nil {
		return result.Error
	}

	slog.Info("Seeded system settings", "attempted", len(settings), "inserted", result.RowsAffected)
	return nil
}
