/**
 * Admin 用户管理相关类型定义
 */

import type { Role } from './role';

/** Admin 用户信息（包含角色） */
export interface AdminUser {
  id: number;
  username: string;
  email: string;
  full_name?: string;
  avatar?: string;
  bio?: string;
  status: 'active' | 'inactive' | 'banned';
  roles?: Role[];
  created_at?: string;
  updated_at?: string;
}

/** 创建用户请求 */
export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  full_name?: string;
  status?: 'active' | 'inactive';
}

/** 更新用户请求 */
export interface UpdateUserRequest {
  email?: string;
  full_name?: string;
  avatar?: string;
  bio?: string;
  status?: 'active' | 'inactive' | 'banned';
}

/** 分配角色请求 */
export interface AssignRolesRequest {
  role_ids: number[];
}
