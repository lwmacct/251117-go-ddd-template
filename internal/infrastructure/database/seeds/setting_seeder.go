package seeds

import (
	"context"
	"log/slog"

	_persistence "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SettingSeeder 系统设置种子数据
type SettingSeeder struct{}

// Seed 执行系统设置种子数据填充
func (s *SettingSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	settings := []_persistence.SettingModel{
		// ==================== General 常规设置 ====================
		// basic 分组 - 基础信息
		{
			Key: "general.site_name", Value: "", Category: "general", Group: "basic",
			ValueType: "string", Label: "站点名称", Order: 10,
			UIConfig: datatypes.JSON(`{"input_type":"text","hint":"显示在浏览器标签和页面标题中"}`),
		},
		{
			Key: "general.site_url", Value: "", Category: "general", Group: "basic",
			ValueType: "string", Label: "站点 URL", Order: 20,
			UIConfig: datatypes.JSON(`{"input_type":"url","hint":"站点完整 URL，如 https://example.com"}`),
		},
		{
			Key: "general.admin_email", Value: "", Category: "general", Group: "basic",
			ValueType: "string", Label: "管理员邮箱", Order: 30,
			UIConfig: datatypes.JSON(`{"input_type":"email","hint":"用于接收系统通知和报警邮件"}`),
		},
		// locale 分组 - 本地化设置
		{
			Key: "general.timezone", Value: "Asia/Shanghai", Category: "general", Group: "locale",
			ValueType: "string", Label: "时区", Order: 40,
			UIConfig: datatypes.JSON(`{"input_type":"select","options":[{"value":"Asia/Shanghai","label":"中国标准时间 (UTC+8)"},{"value":"Asia/Tokyo","label":"日本标准时间 (UTC+9)"},{"value":"America/New_York","label":"美国东部时间 (UTC-5)"},{"value":"Europe/London","label":"格林威治时间 (UTC+0)"},{"value":"UTC","label":"协调世界时 (UTC)"}]}`),
		},
		{
			Key: "general.language", Value: "zh-CN", Category: "general", Group: "locale",
			ValueType: "string", Label: "语言", Order: 50,
			UIConfig: datatypes.JSON(`{"input_type":"select","options":[{"value":"zh-CN","label":"简体中文"},{"value":"zh-TW","label":"繁體中文"},{"value":"en-US","label":"English (US)"},{"value":"ja-JP","label":"日本語"}]}`),
		},
		// appearance 分组 - 外观设置
		{
			Key: "general.theme", Value: "light", Category: "general", Group: "appearance",
			ValueType: "string", Label: "默认主题", Order: 60,
			UIConfig: datatypes.JSON(`{"input_type":"select","hint":"新用户默认使用的主题","options":[{"value":"light","label":"浅色模式"},{"value":"dark","label":"深色模式"},{"value":"system","label":"跟随系统"}]}`),
		},

		// ==================== Security 安全设置 ====================
		// password 分组 - 密码策略
		{
			Key: "security.password_min_length", Value: "8", Category: "security", Group: "password",
			ValueType: "number", Label: "密码最小长度", Order: 10,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"建议至少 8 位","validation":{"min":6,"max":32}}`),
		},
		{
			Key: "security.max_login_attempts", Value: "5", Category: "security", Group: "password",
			ValueType: "number", Label: "最大登录尝试次数", Order: 20,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"超过后账户将被临时锁定","validation":{"min":3,"max":10}}`),
		},
		// session 分组 - 会话管理
		{
			Key: "security.session_timeout", Value: "30", Category: "security", Group: "session",
			ValueType: "number", Label: "会话超时时间（分钟）", Order: 30,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"用户无操作后自动登出的时间","validation":{"min":5,"max":1440}}`),
		},
		// advanced 分组 - 高级安全
		{
			Key: "security.enable_twofa", Value: "false", Category: "security", Group: "advanced",
			ValueType: "boolean", Label: "强制启用两步验证", Order: 40,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"启用后所有用户必须配置两步验证才能登录"}`),
		},

		// ==================== Notification 通知设置 ====================
		// general 分组 - 通知总开关
		{
			Key: "notification.enable_notifications", Value: "true", Category: "notification", Group: "general",
			ValueType: "boolean", Label: "启用系统通知", Order: 10,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"关闭后所有通知渠道将停止发送"}`),
		},
		// email 分组 - 邮件通知
		{
			Key: "notification.enable_email", Value: "true", Category: "notification", Group: "email",
			ValueType: "boolean", Label: "启用邮件通知", Order: 20,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"通过邮件发送系统通知","depends_on":{"key":"notification.enable_notifications","value":true}}`),
		},
		// sms 分组 - 短信通知
		{
			Key: "notification.enable_sms", Value: "false", Category: "notification", Group: "sms",
			ValueType: "boolean", Label: "启用短信通知", Order: 30,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"通过短信发送重要通知（需配置短信服务商）","depends_on":{"key":"notification.enable_notifications","value":true}}`),
		},

		// ==================== Backup 备份设置 ====================
		// general 分组 - 备份总开关
		{
			Key: "backup.enable_backup", Value: "false", Category: "backup", Group: "general",
			ValueType: "boolean", Label: "启用自动备份", Order: 10,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"开启数据自动备份功能"}`),
		},
		// schedule 分组 - 备份计划
		{
			Key: "backup.backup_frequency", Value: "24", Category: "backup", Group: "schedule",
			ValueType: "number", Label: "备份频率（小时）", Order: 20,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"每隔多少小时执行一次备份","validation":{"and":[{">=":[{"var":"value"},1]},{"<=":[{"var":"value"},168]}]},"depends_on":{"key":"backup.enable_backup","value":true}}`),
		},
		{
			Key: "backup.retention_days", Value: "30", Category: "backup", Group: "schedule",
			ValueType: "number", Label: "备份保留天数", Order: 30,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"超过保留期的备份将被自动删除","validation":{"min":7,"max":365},"depends_on":{"key":"backup.enable_backup","value":true}}`),
		},
	}

	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"ui_config", "group", "order",
		}), // 更新 UI 元数据字段
	}).Create(&settings)
	if result.Error != nil {
		return result.Error
	}

	slog.Info("Seeded system settings", "attempted", len(settings), "inserted", result.RowsAffected)
	return nil
}
