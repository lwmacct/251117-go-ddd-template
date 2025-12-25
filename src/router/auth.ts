import type { RouteRecordRaw } from "vue-router";

/**
 * 认证相关路由配置
 *
 * 架构说明：
 * - 使用嵌套路由实现 login/register 无缝切换
 * - 父路由渲染共享布局（品牌区固定）
 * - 子路由在表单区内切换
 */

export const authRoutes: RouteRecordRaw = {
  path: "/auth",
  name: "Auth",
  component: () => import("@/layout/AuthLayout.vue"),
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
