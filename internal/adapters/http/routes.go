package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/permission"
)

// RouteBinding 路由绑定，将 Operation 与 Handler 关联
type RouteBinding struct {
	Op      permission.Operation
	Handler gin.HandlerFunc
}

// AllRouteBindings 返回所有路由绑定
// 绑定顺序决定路由注册顺序，对于同路径前缀的路由需注意顺序
func (deps *RouterDependencies) AllRouteBindings() []RouteBinding {
	return []RouteBinding{
		// ==================== Public 域（公开） ====================
		{permission.PublicAuthRegister, deps.AuthHandler.Register},
		{permission.PublicAuthLogin, deps.AuthHandler.Login},
		{permission.PublicAuthLogin2FA, deps.AuthHandler.Login2FA},
		{permission.PublicAuthRefresh, deps.AuthHandler.RefreshToken},
		{permission.PublicAuthCaptcha, deps.CaptchaHandler.GetCaptcha},

		// ==================== Self 域 - 2FA ====================
		{permission.Self2FASetup, deps.TwoFAHandler.Setup},
		{permission.Self2FAVerify, deps.TwoFAHandler.VerifyAndEnable},
		{permission.Self2FADisable, deps.TwoFAHandler.Disable},
		{permission.Self2FAStatus, deps.TwoFAHandler.GetStatus},

		// ==================== Sys 域 - 用户管理 ====================
		{permission.SysUsersCreate, deps.AdminUserHandler.CreateUser},
		{permission.SysUsersBatchCreate, deps.AdminUserHandler.BatchCreateUsers},
		{permission.SysUsersList, deps.AdminUserHandler.ListUsers},
		{permission.SysUsersGet, deps.AdminUserHandler.GetUser},
		{permission.SysUsersUpdate, deps.AdminUserHandler.UpdateUser},
		{permission.SysUsersDelete, deps.AdminUserHandler.DeleteUser},
		{permission.SysUsersAssignRoles, deps.AdminUserHandler.AssignRoles},

		// ==================== Sys 域 - 角色管理 ====================
		{permission.SysRolesCreate, deps.RoleHandler.CreateRole},
		{permission.SysRolesList, deps.RoleHandler.ListRoles},
		{permission.SysRolesGet, deps.RoleHandler.GetRole},
		{permission.SysRolesUpdate, deps.RoleHandler.UpdateRole},
		{permission.SysRolesDelete, deps.RoleHandler.DeleteRole},
		{permission.SysRolesSetPermissions, deps.RoleHandler.SetPermissions},

		// ==================== Sys 域 - 操作列表 ====================
		{permission.SysOperationsList, deps.OperationHandler.ListOperations},

		// ==================== Sys 域 - 审计日志 ====================
		// 注意：actions 路由必须在 :id 路由之前
		{permission.SysAuditLogsActions, deps.AuditLogHandler.GetActions},
		{permission.SysAuditLogsList, deps.AuditLogHandler.ListLogs},
		{permission.SysAuditLogsGet, deps.AuditLogHandler.GetLog},

		// ==================== Sys 域 - 系统概览 ====================
		{permission.SysOverviewStats, deps.OverviewHandler.GetStats},

		// ==================== Sys 域 - 配置分类（必须在 :key 之前） ====================
		{permission.SysSettingCategoriesList, deps.SettingHandler.GetCategories},
		{permission.SysSettingCategoriesGet, deps.SettingHandler.GetCategory},
		{permission.SysSettingCategoriesCreate, deps.SettingHandler.CreateCategory},
		{permission.SysSettingCategoriesUpdate, deps.SettingHandler.UpdateCategory},
		{permission.SysSettingCategoriesDelete, deps.SettingHandler.DeleteCategory},

		// ==================== Sys 域 - 系统配置 ====================
		// 注意：batch 路由必须在 :key 路由之前
		{permission.SysSettingsBatchUpdate, deps.SettingHandler.BatchUpdateSettings},
		{permission.SysSettingsCreate, deps.SettingHandler.CreateSetting},
		{permission.SysSettingsList, deps.SettingHandler.GetSettings},
		{permission.SysSettingsGet, deps.SettingHandler.GetSetting},
		{permission.SysSettingsUpdate, deps.SettingHandler.UpdateSetting},
		{permission.SysSettingsDelete, deps.SettingHandler.DeleteSetting},

		// ==================== Sys 域 - 缓存管理 ====================
		{permission.SysCacheInfo, deps.CacheHandler.Info},
		{permission.SysCacheScanKeys, deps.CacheHandler.ScanKeys},
		{permission.SysCacheGetKey, deps.CacheHandler.GetKey},
		{permission.SysCacheDeleteKey, deps.CacheHandler.DeleteKey},
		{permission.SysCacheDeletePattern, deps.CacheHandler.DeleteByPattern},

		// ==================== Self 域 - 个人资料 ====================
		{permission.SelfProfileGet, deps.UserProfileHandler.GetProfile},
		{permission.SelfProfileUpdate, deps.UserProfileHandler.UpdateProfile},
		{permission.SelfPasswordUpdate, deps.UserProfileHandler.ChangePassword},
		{permission.SelfAccountDelete, deps.UserProfileHandler.DeleteAccount},

		// ==================== Self 域 - 访问令牌 ====================
		// 注意：scopes 必须在 :id 路由之前
		{permission.SelfTokensCreate, deps.PATHandler.CreateToken},
		{permission.SelfTokensList, deps.PATHandler.ListTokens},
		{permission.SelfTokensScopes, deps.PATHandler.ListScopes},
		{permission.SelfTokensGet, deps.PATHandler.GetToken},
		{permission.SelfTokensDelete, deps.PATHandler.DeleteToken},
		{permission.SelfTokensDisable, deps.PATHandler.DisableToken},
		{permission.SelfTokensEnable, deps.PATHandler.EnableToken},

		// ==================== Self 域 - 用户配置 ====================
		// 注意：categories 和 batch 路由必须在 :key 路由之前
		{permission.SelfSettingsCategoriesList, deps.UserSettingHandler.ListUserSettingCategories},
		{permission.SelfSettingsBatchSet, deps.UserSettingHandler.BatchSetUserSettings},
		{permission.SelfSettingsList, deps.UserSettingHandler.GetUserSettings},
		{permission.SelfSettingsGet, deps.UserSettingHandler.GetUserSetting},
		{permission.SelfSettingsSet, deps.UserSettingHandler.SetUserSetting},
		{permission.SelfSettingsReset, deps.UserSettingHandler.ResetUserSetting},

		// ==================== Self 域 - 用户组织/团队 ====================
		{permission.SelfOrgsList, deps.UserOrgHandler.ListMyOrganizations},
		{permission.SelfTeamsList, deps.UserOrgHandler.ListMyTeams},

		// ==================== Sys 域 - 组织管理 ====================
		{permission.SysOrgsCreate, deps.OrgHandler.Create},
		{permission.SysOrgsList, deps.OrgHandler.List},
		{permission.SysOrgsGet, deps.OrgHandler.Get},
		{permission.SysOrgsUpdate, deps.OrgHandler.Update},
		{permission.SysOrgsDelete, deps.OrgHandler.Delete},

		// ==================== Org 域 - 组织成员管理 ====================
		{permission.OrgMembersList, deps.OrgMemberHandler.List},
		{permission.OrgMembersAdd, deps.OrgMemberHandler.Add},
		{permission.OrgMembersRemove, deps.OrgMemberHandler.Remove},
		{permission.OrgMembersUpdateRole, deps.OrgMemberHandler.UpdateRole},

		// ==================== Org 域 - 团队管理 ====================
		{permission.OrgTeamsCreate, deps.TeamHandler.Create},
		{permission.OrgTeamsList, deps.TeamHandler.List},
		{permission.OrgTeamsGet, deps.TeamHandler.Get},
		{permission.OrgTeamsUpdate, deps.TeamHandler.Update},
		{permission.OrgTeamsDelete, deps.TeamHandler.Delete},

		// ==================== Org 域 - 团队成员管理 ====================
		{permission.OrgTeamMembersList, deps.TeamMemberHandler.List},
		{permission.OrgTeamMembersAdd, deps.TeamMemberHandler.Add},
		{permission.OrgTeamMembersRemove, deps.TeamMemberHandler.Remove},

		// ==================== Sys 域 - 产品管理 ====================
		{permission.SysProductsCreate, deps.ProductHandler.Create},
		{permission.SysProductsList, deps.ProductHandler.List},
		{permission.SysProductsGet, deps.ProductHandler.Get},
		{permission.SysProductsUpdate, deps.ProductHandler.Update},
		{permission.SysProductsDelete, deps.ProductHandler.Delete},
	}
}
