/**
 * 系统概览统计相关类型定义
 */

/** 概览页展示的审计日志摘要 */
export interface OverviewAuditLog {
  id: number;
  user_id: number;
  username: string;
  action: string;
  resource: string;
  status: string;
  created_at: string;
}

/** 系统统计信息 */
export interface SystemStats {
  total_users: number;
  active_users: number;
  inactive_users: number;
  banned_users: number;
  total_roles: number;
  total_permissions: number;
  total_menus: number;
  recent_audit_logs?: OverviewAuditLog[];
}
