/**
 * 用户相关类型定义
 */

/** 用户信息 */
export interface User {
  id: number;
  username: string;
  email: string;
  full_name?: string;
  status: string;
  created_at?: string;
  updated_at?: string;
}
