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
  details?: string;
  status: string;
  created_at: string;
}

/** 审计日志查询参数 */
export interface AuditLogQueryParams {
  page?: number;
  limit?: number;
  user_id?: number;
  action?: string;
  resource?: string;
  status?: string;
  start_date?: string;
  end_date?: string;
}
