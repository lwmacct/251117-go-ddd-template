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

		// ==================== Admin 域 - 用户管理 ====================
		{permission.AdminUsersCreate, deps.AdminUserHandler.CreateUser},
		{permission.AdminUsersBatchCreate, deps.AdminUserHandler.BatchCreateUsers},
		{permission.AdminUsersList, deps.AdminUserHandler.ListUsers},
		{permission.AdminUsersGet, deps.AdminUserHandler.GetUser},
		{permission.AdminUsersUpdate, deps.AdminUserHandler.UpdateUser},
		{permission.AdminUsersDelete, deps.AdminUserHandler.DeleteUser},
		{permission.AdminUsersAssignRoles, deps.AdminUserHandler.AssignRoles},

		// ==================== Admin 域 - 角色管理 ====================
		{permission.AdminRolesCreate, deps.RoleHandler.CreateRole},
		{permission.AdminRolesList, deps.RoleHandler.ListRoles},
		{permission.AdminRolesGet, deps.RoleHandler.GetRole},
		{permission.AdminRolesUpdate, deps.RoleHandler.UpdateRole},
		{permission.AdminRolesDelete, deps.RoleHandler.DeleteRole},
		{permission.AdminRolesSetPermissions, deps.RoleHandler.SetPermissions},

		// ==================== Admin 域 - 操作列表 ====================
		{permission.AdminOperationsList, deps.OperationHandler.ListOperations},

		// ==================== Admin 域 - 审计日志 ====================
		// 注意：actions 路由必须在 :id 路由之前
		{permission.AdminAuditActions, deps.AuditHandler.GetActions},
		{permission.AdminAuditList, deps.AuditHandler.ListLogs},
		{permission.AdminAuditGet, deps.AuditHandler.GetLog},

		// ==================== Admin 域 - 系统概览 ====================
		{permission.AdminOverviewStats, deps.OverviewHandler.GetStats},

		// ==================== Admin 域 - 配置分类（必须在 :key 之前） ====================
		{permission.AdminSettingCategoriesList, deps.SettingHandler.GetCategories},
		{permission.AdminSettingCategoriesGet, deps.SettingHandler.GetCategory},
		{permission.AdminSettingCategoriesCreate, deps.SettingHandler.CreateCategory},
		{permission.AdminSettingCategoriesUpdate, deps.SettingHandler.UpdateCategory},
		{permission.AdminSettingCategoriesDelete, deps.SettingHandler.DeleteCategory},

		// ==================== Admin 域 - 系统配置 ====================
		// 注意：batch 路由必须在 :key 路由之前
		{permission.AdminSettingsBatchUpdate, deps.SettingHandler.BatchUpdateSettings},
		{permission.AdminSettingsCreate, deps.SettingHandler.CreateSetting},
		{permission.AdminSettingsList, deps.SettingHandler.GetSettings},
		{permission.AdminSettingsGet, deps.SettingHandler.GetSetting},
		{permission.AdminSettingsUpdate, deps.SettingHandler.UpdateSetting},
		{permission.AdminSettingsDelete, deps.SettingHandler.DeleteSetting},

		// ==================== Admin 域 - 缓存管理 ====================
		{permission.AdminCacheInfo, deps.CacheHandler.Info},
		{permission.AdminCacheScanKeys, deps.CacheHandler.ScanKeys},
		{permission.AdminCacheGetKey, deps.CacheHandler.GetKey},
		{permission.AdminCacheDeleteKey, deps.CacheHandler.DeleteKey},
		{permission.AdminCacheDeletePattern, deps.CacheHandler.DeleteByPattern},

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

		// ==================== Admin 域 - 组织管理 ====================
		{permission.AdminOrgsCreate, deps.OrgHandler.Create},
		{permission.AdminOrgsList, deps.OrgHandler.List},
		{permission.AdminOrgsGet, deps.OrgHandler.Get},
		{permission.AdminOrgsUpdate, deps.OrgHandler.Update},
		{permission.AdminOrgsDelete, deps.OrgHandler.Delete},

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

		// ==================== Org 域 - 团队任务管理 ====================
		{permission.OrgTasksCreate, deps.TaskHandler.Create},
		{permission.OrgTasksList, deps.TaskHandler.List},
		{permission.OrgTasksGet, deps.TaskHandler.Get},
		{permission.OrgTasksUpdate, deps.TaskHandler.Update},
		{permission.OrgTasksDelete, deps.TaskHandler.Delete},

		// ==================== Admin 域 - 产品管理 ====================
		{permission.AdminProductsCreate, deps.ProductHandler.Create},
		{permission.AdminProductsList, deps.ProductHandler.List},
		{permission.AdminProductsGet, deps.ProductHandler.Get},
		{permission.AdminProductsUpdate, deps.ProductHandler.Update},
		{permission.AdminProductsDelete, deps.ProductHandler.Delete},
	}
}
