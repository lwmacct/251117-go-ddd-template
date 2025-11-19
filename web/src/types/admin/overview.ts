/**
 * 系统概览统计相关类型定义
 */

import type { AuditLog } from "./audit";

/** 系统统计信息 */
export interface SystemStats {
  total_users: number;
  active_users: number;
  inactive_users: number;
  banned_users: number;
  total_roles: number;
  total_permissions: number;
  total_menus: number;
  recent_audit_logs?: AuditLog[];
}
