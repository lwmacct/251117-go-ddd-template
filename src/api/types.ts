/**
 * 前端专用类型定义
 *
 * 仅包含：
 * - 前端状态类型（不在后端 DTO 中）
 * - 必要的类型扩展（如树形结构）
 * - 前端枚举/联合类型
 *
 * 业务 DTO 直接从 @/api/generated/models 导入
 */

import type { MenuMenuDTO, AuthLoginResponseDTO } from "@models";

// ============== 前端状态类型 ==============

/**
 * 登录结果类型（用于 composable 返回，前端内部状态）
 * 表示登录操作的结果，包括 2FA 中间状态
 */
export interface LoginResult {
  success: boolean;
  message?: string;
  requiresTwoFactor: boolean;
  sessionToken?: string;
}

/**
 * 分页状态类型（用于组件内部的分页状态管理）
 */
export interface PaginationState {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

// ============== 扩展类型 ==============

/**
 * 菜单类型（扩展 children 字段用于树形结构）
 * 后端返回平面数据，前端根据 parent_id 构建树
 */
export interface Menu extends MenuMenuDTO {
  children?: Menu[];
}

/**
 * 登录 API 响应类型（联合类型）
 * 后端可能返回：正常登录响应 或 需要 2FA 的响应
 */
export interface LoginApiResponse extends AuthLoginResponseDTO {
  /** 是否需要 2FA（当需要时为 true） */
  requires_2fa?: boolean;
  /** 2FA 会话令牌（当需要 2FA 时返回） */
  session_token?: string;
}

// ============== 前端枚举类型 ==============

/**
 * 用户状态类型
 */
export type UserStatus = "active" | "inactive" | "banned";

/**
 * 审计日志操作类型（空字符串表示全部）
 */
export type AuditAction = "create" | "update" | "delete" | "login" | "logout" | "other" | "";

/**
 * 审计日志状态类型（空字符串表示全部）
 */
export type AuditStatus = "success" | "failure" | "";
