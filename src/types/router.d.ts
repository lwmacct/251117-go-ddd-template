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
    /** 是否需要认证 */
    requiresAuth?: boolean;

    // === 角色权限控制 ===
    /** 允许访问的角色列表，空数组表示所有已登录用户可访问 */
    roles?: string[];

    // === 菜单配置 ===
    /** 是否在菜单中显示 */
    menuVisible?: boolean;
    /** 菜单排序（数值越小越靠前） */
    menuOrder?: number;
  }
}
