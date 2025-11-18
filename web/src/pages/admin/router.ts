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
  component: () => import("./index.vue"),
  redirect: "/admin/overview",
  meta: {
    title: "管理后台",
    requiresAuth: false,
  },
  children: [
    {
      path: "overview",
      name: "AdminOverview",
      component: () => import("./views/overview/index.vue"),
      meta: {
        title: "数据概览",
        icon: "mdi-speedometer",
      },
    },
    {
      path: "roles",
      name: "AdminRoles",
      component: () => import("./views/roles/index.vue"),
      meta: {
        title: "角色管理",
        icon: "mdi-account-group",
      },
    },
    {
      path: "users",
      name: "AdminUsers",
      component: () => import("./views/users/index.vue"),
      meta: {
        title: "用户管理",
        icon: "mdi-account",
      },
    },
    {
      path: "settings",
      name: "AdminSettings",
      component: () => import("./views/settings/index.vue"),
      meta: {
        title: "系统设置",
        icon: "mdi-cog",
      },
    },
    {
      path: "menus",
      name: "AdminMenus",
      component: () => import("./views/menus/index.vue"),
      meta: {
        title: "菜单管理",
        icon: "mdi-menu",
      },
    },
  ],
};
