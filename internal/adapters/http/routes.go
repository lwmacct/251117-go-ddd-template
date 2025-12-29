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

		// ==================== Sys 域 - 用户管理 ====================
		{op.SysUsersCreate, deps.AdminUserHandler.CreateUser},
		{op.SysUsersBatchCreate, deps.AdminUserHandler.BatchCreateUsers},
		{op.SysUsersList, deps.AdminUserHandler.ListUsers},
		{op.SysUsersGet, deps.AdminUserHandler.GetUser},
		{op.SysUsersUpdate, deps.AdminUserHandler.UpdateUser},
		{op.SysUsersDelete, deps.AdminUserHandler.DeleteUser},
		{op.SysUsersAssignRoles, deps.AdminUserHandler.AssignRoles},

		// ==================== Sys 域 - 角色管理 ====================
		{op.SysRolesCreate, deps.RoleHandler.CreateRole},
		{op.SysRolesList, deps.RoleHandler.ListRoles},
		{op.SysRolesGet, deps.RoleHandler.GetRole},
		{op.SysRolesUpdate, deps.RoleHandler.UpdateRole},
		{op.SysRolesDelete, deps.RoleHandler.DeleteRole},
		{op.SysRolesSetPermissions, deps.RoleHandler.SetPermissions},

		// ==================== Sys 域 - 审计日志 ====================
		// 注意：actions 路由必须在 :id 路由之前
		{op.SysAuditLogsActions, deps.AuditLogHandler.GetActions},
		{op.SysAuditLogsList, deps.AuditLogHandler.ListLogs},
		{op.SysAuditLogsGet, deps.AuditLogHandler.GetLog},

		// ==================== Sys 域 - 菜单管理 ====================
		// 注意：reorder 路由必须在 :id 路由之前
		{op.SysMenusReorder, deps.MenuHandler.Reorder},
		{op.SysMenusCreate, deps.MenuHandler.Create},
		{op.SysMenusList, deps.MenuHandler.List},
		{op.SysMenusGet, deps.MenuHandler.Get},
		{op.SysMenusUpdate, deps.MenuHandler.Update},
		{op.SysMenusDelete, deps.MenuHandler.Delete},

		// ==================== Sys 域 - 系统概览 ====================
		{op.SysOverviewStats, deps.OverviewHandler.GetStats},

		// ==================== Sys 域 - 配置分类（必须在 :key 之前） ====================
		{op.SysSettingCategoriesList, deps.SettingHandler.GetCategories},
		{op.SysSettingCategoriesGet, deps.SettingHandler.GetCategory},
		{op.SysSettingCategoriesCreate, deps.SettingHandler.CreateCategory},
		{op.SysSettingCategoriesUpdate, deps.SettingHandler.UpdateCategory},
		{op.SysSettingCategoriesDelete, deps.SettingHandler.DeleteCategory},

		// ==================== Sys 域 - 系统配置 ====================
		// 注意：batch 路由必须在 :key 路由之前
		{op.SysSettingsBatchUpdate, deps.SettingHandler.BatchUpdateSettings},
		{op.SysSettingsCreate, deps.SettingHandler.CreateSetting},
		{op.SysSettingsList, deps.SettingHandler.GetSettings},
		{op.SysSettingsGet, deps.SettingHandler.GetSetting},
		{op.SysSettingsUpdate, deps.SettingHandler.UpdateSetting},
		{op.SysSettingsDelete, deps.SettingHandler.DeleteSetting},

		// ==================== Sys 域 - 缓存管理 ====================
		{op.SysCacheInfo, deps.CacheHandler.Info},
		{op.SysCacheScanKeys, deps.CacheHandler.ScanKeys},
		{op.SysCacheGetKey, deps.CacheHandler.GetKey},
		{op.SysCacheDeleteKey, deps.CacheHandler.DeleteKey},
		{op.SysCacheDeletePattern, deps.CacheHandler.DeleteByPattern},

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
