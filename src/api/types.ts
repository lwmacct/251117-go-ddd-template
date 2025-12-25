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

import type { MenuMenuDTO, ResponsePaginationMeta, ResponseErrorDetail } from "@models";

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

// ============== API 响应类型 ==============

/**
 * 统一 API 响应结构（泛型包装）
 * 对应后端 UnifiedResponse
 */
export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data?: T;
  error?: ResponseErrorDetail;
}

/**
 * 列表响应结构（带分页）
 * 对应后端 ListResponse
 */
export interface ListApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
  meta?: ResponsePaginationMeta;
}

/**
 * 错误响应结构
 * 对应后端 ErrorResponse
 */
export interface ErrorResponse {
  code: number;
  message: string;
  error?: ResponseErrorDetail;
}
