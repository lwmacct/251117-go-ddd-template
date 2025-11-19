import type { RouteRecordRaw } from "vue-router";

/**
 * User 用户中心路由配置
 *
 * 架构说明：
 * - 对应后端 /api/user/* API 命名空间
 * - 用户自我管理相关功能（个人资料、安全设置、访问令牌）
 * - 与 AdminRoutes 保持架构对称性
 */

export const userRoutes: RouteRecordRaw = {
  path: "/user",
  name: "UserLayout",
  component: () => import("@/layout/UserLayout.vue"),
  redirect: "/user/profile",
  meta: {
    title: "用户中心",
    requiresAuth: true,
  },
  children: [
    {
      path: "profile",
      name: "UserProfile",
      component: () => import("@/pages/user/profile/index.vue"),
      meta: {
        title: "个人资料",
        icon: "mdi-account-circle",
      },
    },
    {
      path: "security",
      name: "UserSecurity",
      component: () => import("@/pages/user/security/index.vue"),
      meta: {
        title: "安全设置",
        icon: "mdi-shield-lock",
      },
    },
    {
      path: "tokens",
      name: "UserTokens",
      component: () => import("@/pages/user/tokens/index.vue"),
      meta: {
        title: "访问令牌",
        icon: "mdi-key-variant",
      },
    },
  ],
};
