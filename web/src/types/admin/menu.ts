/**
 * 菜单管理相关类型定义
 */

/** 菜单项 */
export interface Menu {
  id: number;
  title: string;
  path: string;
  icon?: string;
  parent_id?: number;
  order: number;
  visible: boolean;
  children?: Menu[];
  created_at?: string;
  updated_at?: string;
}

/** 创建菜单请求 */
export interface CreateMenuRequest {
  title: string;
  path: string;
  icon?: string;
  parent_id?: number;
  order?: number;
  visible?: boolean;
}

/** 更新菜单请求 */
export interface UpdateMenuRequest {
  title?: string;
  path?: string;
  icon?: string;
  parent_id?: number;
  order?: number;
  visible?: boolean;
}

/** 批量更新排序请求 */
export interface ReorderMenusRequest {
  menus: Array<{
    id: number;
    order: number;
    parent_id?: number;
  }>;
}
