import { createRouter, createWebHashHistory } from "vue-router";
import type { RouteLocationNormalized, NavigationGuardNext } from "vue-router";
import { adminRoutes } from "./admin";
import { authRoutes } from "./auth";
import { userRoutes } from "./user";
import { accessToken } from "@/utils/auth";
import { useAuthStore } from "@/stores/auth";

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
    // 403 无权限页面
    {
      path: "/403",
      name: "Forbidden",
      component: () => import("@/pages/error/403.vue"),
      meta: {
        title: "无权访问",
      },
    },
    // 404 页面（可选）
    {
      path: "/:pathMatch(.*)*",
      name: "NotFound",
      component: () => import("@/pages/error/404.vue"),
      meta: {
        title: "页面不存在",
      },
    },
  ],
});

/**
 * 路由守卫：认证和角色检查
 * 1. 未登录用户访问需认证页面 → 跳转登录页
 * 2. 已登录用户无权限访问 → 跳转 403 页面
 * 3. 已登录用户访问登录页 → 跳转管理后台
 */
router.beforeEach(async (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
  const token = accessToken.value;
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);

  // 1. 需要认证但没有 token
  if (requiresAuth && !token) {
    next({
      path: "/auth/login",
      query: { redirect: to.fullPath },
    });
    return;
  }

  // 2. 角色检查（仅对需要认证的路由）
  if (token && requiresAuth) {
    // 获取路由 meta 中的 roles 配置
    const routeRoles = to.meta?.roles as string[] | undefined;

    // 如果路由配置了 roles 且不为空数组，进行角色检查
    if (routeRoles && routeRoles.length > 0) {
      const authStore = useAuthStore();

      // 确保用户信息已加载
      if (!authStore.currentUser) {
        await authStore.initAuth();
      }

      // 提取用户角色名称
      const userRoles = authStore.currentUser?.roles?.map((r) => r.name).filter((n): n is string => !!n) || [];

      // 检查用户是否有权限
      const hasPermission = routeRoles.some((r) => userRoles.includes(r));

      if (!hasPermission) {
        next({ name: "Forbidden" });
        return;
      }
    }
  }

  // 3. 已登录用户访问登录/注册页，重定向到管理后台
  if (token && (to.path === "/auth/login" || to.path === "/auth/register")) {
    next("/admin/overview");
    return;
  }

  next();
});

export default router;
