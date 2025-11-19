/**
 * 角色和权限相关类型定义
 */

/** 权限定义 */
export interface Permission {
  id: number;
  domain: string;
  resource: string;
  action: string;
  code: string; // "domain:resource:action"
  description?: string;
  created_at?: string;
}

/** 角色定义 */
export interface Role {
  id: number;
  name: string;
  display_name: string;
  description?: string;
  is_system: boolean;
  permissions?: Permission[];
  created_at?: string;
  updated_at?: string;
}

/** 创建角色请求 */
export interface CreateRoleRequest {
  name: string;
  display_name: string;
  description?: string;
}

/** 更新角色请求 */
export interface UpdateRoleRequest {
  display_name?: string;
  description?: string;
}

/** 设置权限请求 */
export interface SetPermissionsRequest {
  permission_ids: number[];
}

/** 权限树节点（用于 UI 展示） */
export interface PermissionTreeNode {
  key: string; // domain 或 domain:resource
  label: string;
  children?: PermissionTreeNode[];
  permission?: Permission; // 叶子节点的权限对象
}
