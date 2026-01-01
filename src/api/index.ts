/**
 * API 统一导出
 *
 * 本项目 API 层结构：
 * - types.ts: 前端专用类型（状态、扩展、枚举）
 * - helpers.ts: 响应提取辅助函数
 * - client.ts: axios 实例 + API 实例
 * - auth.ts: 认证业务封装
 *
 * 业务 DTO 直接从 generated/models 导出，无别名层
 */

// ============== 生成的模型（业务 DTO） ==============
export * from "@models";

// ============== 前端专用类型 ==============
export * from "./types";

// ============== 辅助函数 ==============
export * from "./helpers";

// ============== API 实例 ==============
export {
  apiClient,
  adminAuditLogApi,
  adminOrgApi,
  adminProductApi,
  adminRoleApi,
  adminSettingCategoriesApi,
  adminSettingsApi,
  adminUserApi,
  authApi,
  auth2faApi,
  orgMemberApi,
  orgTeamApi,
  orgTeamMemberApi,
  overviewApi,
  systemApi,
  userOrgApi,
  userSettingsApi,
  userTokensApi,
  userProfileApi,
} from "./client";

// ============== 认证服务 ==============
export { AuthService } from "./auth";

// ============== 错误处理 ==============
export * from "./errors";
