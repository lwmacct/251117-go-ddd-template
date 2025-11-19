/**
 * 审计日志相关类型定义
 */

/** 审计日志 */
export interface AuditLog {
  id: number;
  user_id: number;
  username?: string;
  action: string;
  resource: string;
  resource_id?: number;
  ip_address?: string;
  user_agent?: string;
  details?: Record<string, any>;
  created_at: string;
}
