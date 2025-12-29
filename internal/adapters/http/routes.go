package http

import (
	"github.com/gin-gonic/gin"
	op "github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"
)

// RouteBinding 路由绑定，将 Operation 与 Handler 关联
type RouteBinding struct {
	Op      op.Operation
	Handler gin.HandlerFunc
}

// AllRouteBindings 返回所有路由绑定
// 绑定顺序决定路由注册顺序，对于同路径前缀的路由需注意顺序
func (deps *RouterDependencies) AllRouteBindings() []RouteBinding {
	return []RouteBinding{
		// ==================== Auth 域（公开） ====================
		{op.AuthRegister, deps.AuthHandler.Register},
		{op.AuthLogin, deps.AuthHandler.Login},
		{op.AuthLogin2FA, deps.AuthHandler.Login2FA},
		{op.AuthRefresh, deps.AuthHandler.RefreshToken},
		{op.AuthCaptcha, deps.CaptchaHandler.GetCaptcha},

		// ==================== Auth 域 - 2FA ====================
		{op.Auth2FASetup, deps.TwoFAHandler.Setup},
		{op.Auth2FAVerify, deps.TwoFAHandler.VerifyAndEnable},
		{op.Auth2FADisable, deps.TwoFAHandler.Disable},
		{op.Auth2FAStatus, deps.TwoFAHandler.GetStatus},

		// ==================== Admin 域 - 用户管理 ====================
		{op.AdminUsersCreate, deps.AdminUserHandler.CreateUser},
		{op.AdminUsersBatchCreate, deps.AdminUserHandler.BatchCreateUsers},
		{op.AdminUsersList, deps.AdminUserHandler.ListUsers},
		{op.AdminUsersGet, deps.AdminUserHandler.GetUser},
		{op.AdminUsersUpdate, deps.AdminUserHandler.UpdateUser},
		{op.AdminUsersDelete, deps.AdminUserHandler.DeleteUser},
		{op.AdminUsersAssignRoles, deps.AdminUserHandler.AssignRoles},

		// ==================== Admin 域 - 角色管理 ====================
		{op.AdminRolesCreate, deps.RoleHandler.CreateRole},
		{op.AdminRolesList, deps.RoleHandler.ListRoles},
		{op.AdminRolesGet, deps.RoleHandler.GetRole},
		{op.AdminRolesUpdate, deps.RoleHandler.UpdateRole},
		{op.AdminRolesDelete, deps.RoleHandler.DeleteRole},
		{op.AdminRolesSetPermissions, deps.RoleHandler.SetPermissions},

		// ==================== Admin 域 - 权限管理 ====================
		{op.AdminPermissionsList, deps.RoleHandler.ListPermissions},

		// ==================== Admin 域 - 审计日志 ====================
		// 注意：actions 路由必须在 :id 路由之前
		{op.AdminAuditLogsActions, deps.AuditLogHandler.GetActions},
		{op.AdminAuditLogsList, deps.AuditLogHandler.ListLogs},
		{op.AdminAuditLogsGet, deps.AuditLogHandler.GetLog},

		// ==================== Admin 域 - 菜单管理 ====================
		// 注意：reorder 路由必须在 :id 路由之前
		{op.AdminMenusReorder, deps.MenuHandler.Reorder},
		{op.AdminMenusCreate, deps.MenuHandler.Create},
		{op.AdminMenusList, deps.MenuHandler.List},
		{op.AdminMenusGet, deps.MenuHandler.Get},
		{op.AdminMenusUpdate, deps.MenuHandler.Update},
		{op.AdminMenusDelete, deps.MenuHandler.Delete},

		// ==================== Admin 域 - 系统概览 ====================
		{op.AdminOverviewStats, deps.OverviewHandler.GetStats},

		// ==================== Admin 域 - 配置分类（必须在 :key 之前） ====================
		{op.AdminSettingCategoriesList, deps.SettingHandler.GetCategories},
		{op.AdminSettingCategoriesGet, deps.SettingHandler.GetCategory},
		{op.AdminSettingCategoriesCreate, deps.SettingHandler.CreateCategory},
		{op.AdminSettingCategoriesUpdate, deps.SettingHandler.UpdateCategory},
		{op.AdminSettingCategoriesDelete, deps.SettingHandler.DeleteCategory},

		// ==================== Admin 域 - 系统配置 ====================
		// 注意：batch 路由必须在 :key 路由之前
		{op.AdminSettingsBatchUpdate, deps.SettingHandler.BatchUpdateSettings},
		{op.AdminSettingsCreate, deps.SettingHandler.CreateSetting},
		{op.AdminSettingsList, deps.SettingHandler.GetSettings},
		{op.AdminSettingsGet, deps.SettingHandler.GetSetting},
		{op.AdminSettingsUpdate, deps.SettingHandler.UpdateSetting},
		{op.AdminSettingsDelete, deps.SettingHandler.DeleteSetting},

		// ==================== Admin 域 - 缓存管理 ====================
		{op.AdminCacheInfo, deps.CacheHandler.Info},
		{op.AdminCacheScanKeys, deps.CacheHandler.ScanKeys},
		{op.AdminCacheGetKey, deps.CacheHandler.GetKey},
		{op.AdminCacheDeleteKey, deps.CacheHandler.DeleteKey},
		{op.AdminCacheDeletePattern, deps.CacheHandler.DeleteByPattern},

		// ==================== User 域 - 个人资料 ====================
		{op.UserProfileGet, deps.UserProfileHandler.GetProfile},
		{op.UserProfileUpdate, deps.UserProfileHandler.UpdateProfile},
		{op.UserPasswordUpdate, deps.UserProfileHandler.ChangePassword},
		{op.UserAccountDelete, deps.UserProfileHandler.DeleteAccount},

		// ==================== User 域 - 访问令牌 ====================
		{op.UserTokensCreate, deps.PATHandler.CreateToken},
		{op.UserTokensList, deps.PATHandler.ListTokens},
		{op.UserTokensGet, deps.PATHandler.GetToken},
		{op.UserTokensDelete, deps.PATHandler.DeleteToken},
		{op.UserTokensDisable, deps.PATHandler.DisableToken},
		{op.UserTokensEnable, deps.PATHandler.EnableToken},

		// ==================== User 域 - 用户配置 ====================
		// 注意：categories 和 batch 路由必须在 :key 路由之前
		{op.UserSettingsCategoriesList, deps.UserSettingHandler.ListUserSettingCategories},
		{op.UserSettingsBatchSet, deps.UserSettingHandler.BatchSetUserSettings},
		{op.UserSettingsList, deps.UserSettingHandler.GetUserSettings},
		{op.UserSettingsGet, deps.UserSettingHandler.GetUserSetting},
		{op.UserSettingsSet, deps.UserSettingHandler.SetUserSetting},
		{op.UserSettingsReset, deps.UserSettingHandler.ResetUserSetting},
	}
}
