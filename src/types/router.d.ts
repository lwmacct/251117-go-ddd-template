/**
 * Vue Router 类型扩展
 * 扩展路由元信息接口，支持菜单和导航功能
 */

import "vue-router";

declare module "vue-router" {
  interface RouteMeta {
    /** 页面标题 */
    title?: string;
    /** 页面图标 (MDI 格式) */
    icon?: string;
    /** 页面描述 */
    description?: string;
    /** 页面分类 */
    category?: string;
    /** 排序优先级 */
    priority?: number;
    /** 是否在菜单中显示 */
    showInMenu?: boolean;
    /** 是否需要认证 */
    requireAuth?: boolean;
  }
}
