import type { RouteRecordRaw } from "vue-router";

/**
 * 认证相关路由配置
 *
 * 架构说明：
 * - 模块化路由设计，认证模块管理自己的路由
 * - 采用懒加载提升性能
 * - 与后端 DDD 架构思想保持一致
 */

export const authRoutes: RouteRecordRaw = {
  path: "/auth",
  name: "Auth",
  redirect: "/auth/login",
  meta: {
    title: "认证",
    requiresAuth: false,
  },
  children: [
    {
      path: "login",
      name: "Login",
      component: () => import("@/pages/auth/login/index.vue"),
      meta: {
        title: "登录",
      },
    },
    {
      path: "register",
      name: "Register",
      component: () => import("@/pages/auth/register/index.vue"),
      meta: {
        title: "注册",
      },
    },
  ],
};
