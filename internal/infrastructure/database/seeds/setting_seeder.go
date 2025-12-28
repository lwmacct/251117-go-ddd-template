package seeds

import (
	"context"
	"fmt"
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

	// 1. 查询所有 Category ID
	categoryIDs, err := s.loadCategoryIDs(db)
	if err != nil {
		return fmt.Errorf("load category IDs: %w", err)
	}

	// 验证必需的 Category 存在
	requiredCategories := []string{"general", "security", "notification", "backup"}
	for _, key := range requiredCategories {
		if _, ok := categoryIDs[key]; !ok {
			return fmt.Errorf("required category not found: %s (run SettingCategorySeeder first)", key)
		}
	}

	// 2. 构建配置定义（使用 CategoryID）
	definitions := s.buildDefinitions(categoryIDs)

	// 3. 批量插入/更新
	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"ui_config", "group", "order", "label", "scope", "view_permission", "edit_permission",
		}), // 更新 UI 元数据和权限字段，不覆盖用户修改的默认值
	}).Create(&definitions)
	if result.Error != nil {
		return result.Error
	}

	slog.Info("Seeded setting definitions", "attempted", len(definitions), "inserted", result.RowsAffected)
	return nil
}

// Name 返回 seeder 名称
func (s *SettingSeeder) Name() string {
	return "SettingSeeder"
}

// loadCategoryIDs 从数据库加载 Category key -> ID 映射
func (s *SettingSeeder) loadCategoryIDs(db *gorm.DB) (map[string]uint, error) {
	var categories []_persistence.SettingCategoryModel
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}

	result := make(map[string]uint, len(categories))
	for _, cat := range categories {
		result[cat.Key] = cat.ID
	}
	return result, nil
}

// buildDefinitions 构建配置定义列表
func (s *SettingSeeder) buildDefinitions(categoryIDs map[string]uint) []_persistence.SettingModel {
	return []_persistence.SettingModel{
		// ==================== General 常规设置 ====================
		// basic 分组 - 基础信息（系统级，管理员可见可编辑）
		{
			Key: "general.site_name", DefaultValue: datatypes.JSON(`""`), CategoryID: categoryIDs["general"], Group: "basic",
			Scope: "system", ViewPermission: "*:settings:read", EditPermission: "admin:settings:update",
			ValueType: "string", Label: "站点名称", Order: 10,
			UIConfig: datatypes.JSON(`{"input_type":"text","hint":"显示在浏览器标签和页面标题中"}`),
		},
		{
			Key: "general.site_url", DefaultValue: datatypes.JSON(`""`), CategoryID: categoryIDs["general"], Group: "basic",
			Scope: "system", ViewPermission: "*:settings:read", EditPermission: "admin:settings:update",
			ValueType: "string", Label: "站点 URL", Order: 20,
			UIConfig: datatypes.JSON(`{"input_type":"url","hint":"站点完整 URL，如 https://example.com"}`),
		},
		{
			Key: "general.admin_email", DefaultValue: datatypes.JSON(`""`), CategoryID: categoryIDs["general"], Group: "basic",
			Scope: "system", ViewPermission: "admin:settings:read", EditPermission: "admin:settings:update",
			ValueType: "string", Label: "管理员邮箱", Order: 30,
			UIConfig: datatypes.JSON(`{"input_type":"email","hint":"用于接收系统通知和报警邮件"}`),
		},
		// locale 分组 - 本地化设置（用户级，用户可覆盖）
		{
			Key: "general.timezone", DefaultValue: datatypes.JSON(`"Asia/Shanghai"`), CategoryID: categoryIDs["general"], Group: "locale",
			Scope: "user", ViewPermission: "user:settings:read", EditPermission: "user:settings:update",
			ValueType: "string", Label: "时区", Order: 40,
			UIConfig: datatypes.JSON(`{"input_type":"select","options":[{"value":"Asia/Shanghai","label":"中国标准时间 (UTC+8)"},{"value":"Asia/Tokyo","label":"日本标准时间 (UTC+9)"},{"value":"America/New_York","label":"美国东部时间 (UTC-5)"},{"value":"Europe/London","label":"格林威治时间 (UTC+0)"},{"value":"UTC","label":"协调世界时 (UTC)"}]}`),
		},
		{
			Key: "general.language", DefaultValue: datatypes.JSON(`"zh-CN"`), CategoryID: categoryIDs["general"], Group: "locale",
			Scope: "user", ViewPermission: "user:settings:read", EditPermission: "user:settings:update",
			ValueType: "string", Label: "语言", Order: 50,
			UIConfig: datatypes.JSON(`{"input_type":"select","options":[{"value":"zh-CN","label":"简体中文"},{"value":"zh-TW","label":"繁體中文"},{"value":"en-US","label":"English (US)"},{"value":"ja-JP","label":"日本語"}]}`),
		},
		// appearance 分组 - 外观设置（用户级，用户可覆盖）
		{
			Key: "general.theme", DefaultValue: datatypes.JSON(`"light"`), CategoryID: categoryIDs["general"], Group: "appearance",
			Scope: "user", ViewPermission: "user:settings:read", EditPermission: "user:settings:update",
			ValueType: "string", Label: "默认主题", Order: 60,
			UIConfig: datatypes.JSON(`{"input_type":"select","hint":"新用户默认使用的主题","options":[{"value":"light","label":"浅色模式"},{"value":"dark","label":"深色模式"},{"value":"system","label":"跟随系统"}]}`),
		},

		// ==================== Security 安全设置（系统级，管理员专用）====================
		// password 分组 - 密码策略
		{
			Key: "security.password_min_length", DefaultValue: datatypes.JSON(`8`), CategoryID: categoryIDs["security"], Group: "password",
			Scope: "system", ViewPermission: "admin:settings:read", EditPermission: "admin:settings:update",
			ValueType: "number", Label: "密码最小长度", Order: 10,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"建议至少 8 位","validation":{"min":6,"max":32}}`),
		},
		{
			Key: "security.max_login_attempts", DefaultValue: datatypes.JSON(`5`), CategoryID: categoryIDs["security"], Group: "password",
			Scope: "system", ViewPermission: "admin:settings:read", EditPermission: "admin:settings:update",
			ValueType: "number", Label: "最大登录尝试次数", Order: 20,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"超过后账户将被临时锁定","validation":{"min":3,"max":10}}`),
		},
		// session 分组 - 会话管理
		{
			Key: "security.session_timeout", DefaultValue: datatypes.JSON(`30`), CategoryID: categoryIDs["security"], Group: "session",
			Scope: "system", ViewPermission: "admin:settings:read", EditPermission: "admin:settings:update",
			ValueType: "number", Label: "会话超时时间（分钟）", Order: 30,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"用户无操作后自动登出的时间","validation":{"min":5,"max":1440}}`),
		},
		// advanced 分组 - 高级安全
		{
			Key: "security.enable_twofa", DefaultValue: datatypes.JSON(`false`), CategoryID: categoryIDs["security"], Group: "advanced",
			Scope: "system", ViewPermission: "admin:settings:read", EditPermission: "admin:settings:update",
			ValueType: "boolean", Label: "强制启用两步验证", Order: 40,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"启用后所有用户必须配置两步验证才能登录"}`),
		},

		// ==================== Notification 通知设置 ====================
		// general 分组 - 通知总开关（系统级）
		{
			Key: "notification.enable_notifications", DefaultValue: datatypes.JSON(`true`), CategoryID: categoryIDs["notification"], Group: "general",
			Scope: "system", ViewPermission: "admin:settings:read", EditPermission: "admin:settings:update",
			ValueType: "boolean", Label: "启用系统通知", Order: 10,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"关闭后所有通知渠道将停止发送"}`),
		},
		// email 分组 - 邮件通知（用户级，用户可开关）
		{
			Key: "notification.enable_email", DefaultValue: datatypes.JSON(`true`), CategoryID: categoryIDs["notification"], Group: "email",
			Scope: "user", ViewPermission: "user:settings:read", EditPermission: "user:settings:update",
			ValueType: "boolean", Label: "启用邮件通知", Order: 20,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"通过邮件发送系统通知","depends_on":{"key":"notification.enable_notifications","value":true}}`),
		},
		// sms 分组 - 短信通知（用户级，用户可开关）
		{
			Key: "notification.enable_sms", DefaultValue: datatypes.JSON(`false`), CategoryID: categoryIDs["notification"], Group: "sms",
			Scope: "user", ViewPermission: "user:settings:read", EditPermission: "user:settings:update",
			ValueType: "boolean", Label: "启用短信通知", Order: 30,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"通过短信发送重要通知（需配置短信服务商）","depends_on":{"key":"notification.enable_notifications","value":true}}`),
		},

		// ==================== Backup 备份设置（系统级，管理员专用）====================
		// general 分组 - 备份总开关
		{
			Key: "backup.enable_backup", DefaultValue: datatypes.JSON(`false`), CategoryID: categoryIDs["backup"], Group: "general",
			Scope: "system", ViewPermission: "admin:settings:read", EditPermission: "admin:settings:update",
			ValueType: "boolean", Label: "启用自动备份", Order: 10,
			UIConfig: datatypes.JSON(`{"input_type":"switch","hint":"开启数据自动备份功能"}`),
		},
		// schedule 分组 - 备份计划
		{
			Key: "backup.backup_frequency", DefaultValue: datatypes.JSON(`24`), CategoryID: categoryIDs["backup"], Group: "schedule",
			Scope: "system", ViewPermission: "admin:settings:read", EditPermission: "admin:settings:update",
			ValueType: "number", Label: "备份频率（小时）", Order: 20,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"每隔多少小时执行一次备份","validation":{"and":[{">=":[{"var":"value"},1]},{"<=":[{"var":"value"},168]}]},"depends_on":{"key":"backup.enable_backup","value":true}}`),
		},
		{
			Key: "backup.retention_days", DefaultValue: datatypes.JSON(`30`), CategoryID: categoryIDs["backup"], Group: "schedule",
			Scope: "system", ViewPermission: "admin:settings:read", EditPermission: "admin:settings:update",
			ValueType: "number", Label: "备份保留天数", Order: 30,
			UIConfig: datatypes.JSON(`{"input_type":"number","hint":"超过保留期的备份将被自动删除","validation":{"min":7,"max":365},"depends_on":{"key":"backup.enable_backup","value":true}}`),
		},
	}
}
