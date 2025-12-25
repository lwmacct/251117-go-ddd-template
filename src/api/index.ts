/**
 * API 统一导出
 *
 * 本项目 API 层结构：
 * - types.ts: 前端专用类型（状态、扩展、枚举）
 * - helpers.ts: 响应提取辅助函数
 * - auth/client.ts: axios 实例 + API 实例
 * - generated/: OpenAPI Generator 自动生成的代码（勿手动修改）
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
  adminMenuApi,
  adminRoleApi,
  adminSettingsApi,
  adminUserApi,
  authApi,
  auth2faApi,
  overviewApi,
  systemApi,
  userTokensApi,
  userProfileApi,
} from "./auth/client";

// ============== 认证相关 ==============
export * from "./auth/platformAuth";

// ============== 错误处理 ==============
export * from "./errors";
