import type { RouteRecordRaw } from "vue-router";

/**
 * Admin 管理后台路由配置
 *
 * 架构说明：
 * - 模块化路由设计，每个业务模块管理自己的路由
 * - 采用懒加载提升性能
 * - 与后端 DDD 架构思想保持一致
 */

export const adminRoutes: RouteRecordRaw = {
  path: "/admin",
  name: "AdminLayout",
  component: () => import("@/layout/AdminLayout.vue"),
  redirect: "/admin/overview",
  meta: {
    title: "管理后台",
    requiresAuth: true,
  },
  children: [
    {
      path: "overview",
      name: "AdminOverview",
      component: () => import("@/pages/admin/overview/index.vue"),
      meta: {
        title: "数据概览",
        icon: "mdi-speedometer",
      },
    },
    {
      path: "roles",
      name: "AdminRoles",
      component: () => import("@/pages/admin/roles/index.vue"),
      meta: {
        title: "角色管理",
        icon: "mdi-account-group",
      },
    },
    {
      path: "users",
      name: "AdminUsers",
      component: () => import("@/pages/admin/users/index.vue"),
      meta: {
        title: "用户管理",
        icon: "mdi-account",
      },
    },
    {
      path: "settings",
      name: "AdminSettings",
      component: () => import("@/pages/admin/settings/index.vue"),
      meta: {
        title: "系统设置",
        icon: "mdi-cog",
      },
    },
    {
      path: "menus",
      name: "AdminMenus",
      component: () => import("@/pages/admin/menus/index.vue"),
      meta: {
        title: "菜单管理",
        icon: "mdi-menu",
      },
    },
    {
      path: "auditlogs",
      name: "AdminAuditLogs",
      component: () => import("@/pages/admin/auditlogs/index.vue"),
      meta: {
        title: "审计日志",
        icon: "mdi-file-document-outline",
      },
    },
  ],
};
