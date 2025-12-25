import { createRouter, createWebHashHistory } from "vue-router";
import type { RouteLocationNormalized, NavigationGuardNext } from "vue-router";
import { adminRoutes } from "./admin";
import { authRoutes } from "./auth";
import { userRoutes } from "./user";
import { accessToken } from "@/utils/auth";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: "/auth/login",
    },
    authRoutes,
    adminRoutes,
    userRoutes,
  ],
});

/**
 * 路由守卫：认证检查
 * 拦截所有需要认证的路由，未登录用户跳转到登录页
 */
router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
  const token = accessToken.value;
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);

  // 需要认证但没有 token
  if (requiresAuth && !token) {
    next({
      path: "/auth/login",
      query: { redirect: to.fullPath }, // 保存目标路由，登录后可跳转回来
    });
    return;
  }

  // 已登录用户访问登录/注册页，重定向到管理后台
  if (token && (to.path === "/auth/login" || to.path === "/auth/register")) {
    next("/admin/overview");
    return;
  }

  next();
});

export default router;
